package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var (
	ErrChatFinishChatWithSameOne = errors.New("not chat user with other one")
	ErrTodayIsFinished           = errors.New("today has finished")
	ErrTimeHasNotBeenReached     = errors.New("time has not been reached")
)

func CreateOrUpdateTask(task *db.Task) error {
	task.CreatedAt = time.Now()
	// 没有就创建没有就更新
	upTask := db.Task{}
	utils.CopyStructFields(&upTask, task)
	return db.DB.MysqlDB.DefaultGormDB().Table("task").Where("task_id=?", task.Id).
		Assign(upTask).FirstOrCreate(&task).Error
}

func GetTaskById(taskId int32) (resultdata *db.Task, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("task").Where("task_id=?", taskId).
		Take(&resultdata).Error
	return
}

// 获取开启的任务
func GetOpenTaskById(taskId int32) (resultdata *db.Task, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("task").Where("task_id=? and status = ?", taskId, constant.TaskStatusOpen).
		Take(&resultdata).Error
	return
}

func HasUserTaskRecord(userTaskId string) bool {
	var count int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_task").
		Where("id=?", userTaskId).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func GetTaskList(classify string) (resultdata []db.Task, err error) {
	if classify == "" {
		err = db.DB.MysqlDB.DefaultGormDB().Table("task").Where("status=0").
			Find(&resultdata).Error
	} else {
		err = db.DB.MysqlDB.DefaultGormDB().Table("task").Where("classify=? and status=0", classify).
			Find(&resultdata).Error
	}
	return
}

func GetUserTaskCount(userId string, taskId int32) int64 {
	var count int64
	db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("user_id=? and task_id=?", userId, taskId).
		Count(&count)
	return count
}

func DeleteTask(taskId int64) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("task").Delete("task_id=?", taskId).Error
	return err
}

func InsertIntoUserTask(userTask *db.UserTask, eventEx string) error {
	userTask.StartTime = time.Now()
	task, err := GetOpenTaskById(userTask.TaskID)
	if err != nil {
		return err
	}
	return db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		err := tx.Table("user_task").Create(userTask).Error
		if err != nil {
			return err
		}
		if userTask.Progress >= task.CompletionCount {
			return UserTaskFinished(userTask.ID, userTask.UserID, task.Reward, eventEx, task.Name, tx)
		}
		return nil
	})
}

// 获取用户领取的任务 0 代表所有状态
func GetUserClaimTask(userId string, status int) (resultdata []db.UserTask, err error) {
	if status == 0 {
		err = db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("user_id=?", userId).Preload("Task").
			Find(&resultdata).Error
	} else {
		err = db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("user_id=? and status = ?", userId, status).Preload("Task").
			Find(&resultdata).Error
	}
	return
}

func GetUserTaskByID(userTaskId string) (resultdata db.UserTask, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("id=?", userTaskId).
		Take(&resultdata).Error
	return
}

// 获取开启状态的任务
func GetOpenUserTaskByID(userTaskId string) (resultdata db.UserTask, err error) {
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("id=?", userTaskId).
		Take(&resultdata).Error
	return
}

func ClaimUserTaskSuccess(userTaskId, userId string, task *db.Task) error {
	return UserTaskFinished(userTaskId, userId, task.Reward, "", task.Name, db.DB.MysqlDB.DefaultGormDB())
}

// 更新用户任务状态为已完成并插入奖励表
func UserTaskFinished(userTaskId string, userId string, taskReward uint64, eventEx string, taskInfo string, tx *gorm.DB) error {
	// 插入奖励表
	taskEventLogId := fmt.Sprintf("%s:%s", constant.RewardTypeTask, userTaskId)
	if eventEx != "" {
		taskEventLogId += ":" + eventEx
	}
	rewardEventLogs := &db.RewardEventLogs{
		ID:         taskEventLogId,
		UserID:     userId,
		RewardType: constant.RewardTypeTask,
		Reward:     taskReward,
		CreatedAt:  time.Now(),
		UserJSON:   fmt.Sprintf(`{"type":"%s", "info":"%s"}`, constant.RewardTypeTask, taskInfo),
	}
	err := tx.Table("reward_event_logs").Create(rewardEventLogs).Error
	if err != nil {
		return err
	}
	// 更新用户任务状态为已完成并关联 evelogs_id（如果查看详情页的话，可以直接查询外键查看同步时间等）
	return tx.Table("user_task").Where("id=?", userTaskId).
		Updates(map[string]interface{}{
			"status":               constant.UserTaskStatusFinished,
			"reward_event_logs_id": taskEventLogId,
		}).Error
}

// 更新用户任务状态为已领取
func UpdateUserTaskSync(userTaskId string) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("id=?", userTaskId).
		Update("status", constant.UserTaskStatusClaimed).Error
}

func IsFinishUserTask(userTaskId string) (bool, error) {
	// 获取 user_task user_task_id 为 userTaskId 的记录数量，判断是否大于 0
	var count int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("id=? and status = ?", userTaskId, constant.UserTaskStatusFinished).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func FinishFollowOfficialTwitterTask(userId string) error {
	return InsertIntoUserTask(&db.UserTask{
		ID:       fmt.Sprintf("%s:%d", userId, constant.TaskIDFollowOfficialTwitter),
		UserID:   userId,
		TaskID:   constant.TaskIDFollowOfficialTwitter,
		Progress: 1,
	}, "")
}

func IsFinishFollowOfficialTwitterTask(userId string) (bool, error) {
	return IsFinishUserTask(fmt.Sprintf("%s:%d", userId, constant.TaskIDFollowOfficialTwitter))
}

func FinishBindTwitterTask(userId string) error {
	return InsertIntoUserTask(&db.UserTask{
		ID:       fmt.Sprintf("%s:%d", userId, constant.TaskIDBindTwitter),
		UserID:   userId,
		TaskID:   constant.TaskIDBindTwitter,
		Progress: 1,
	}, "")
}

func IsFinishBindTwitterTask(userId string) (bool, error) {
	return IsFinishUserTask(fmt.Sprintf("%s:%d", userId, constant.TaskIDBindTwitter))
}

func FinishJoinOfficialSpaceTask(userId string) error {
	return InsertIntoUserTask(&db.UserTask{
		ID:        fmt.Sprintf("%s:%d", userId, constant.TaskIdJoinOfficialSpace),
		UserID:    userId,
		TaskID:    constant.TaskIdJoinOfficialSpace,
		Status:    constant.UserTaskStatusDoing,
		StartTime: time.Now(),
	}, "")
}

func FinishInviteBindTwitterTask(userId string, formUser string) error {
	userTaskId := fmt.Sprintf("%s:%d:%s", userId, constant.TaskIDInviteBindTwitter, formUser)
	userTaskId = utils.Md5(userTaskId)
	return InsertIntoUserTask(&db.UserTask{
		ID:       userTaskId,
		UserID:   userId,
		TaskID:   constant.TaskIDInviteBindTwitter,
		Progress: 1,
		Ex:       formUser,
	}, "")
}

func FinishInviteUploadNftHeadTask(userId string, formUser string) error {
	userTaskId := fmt.Sprintf("%s:%d:%s", userId, constant.TaskIdInviteUploadNftHead, formUser)
	userTaskId = utils.Md5(userTaskId)
	return InsertIntoUserTask(&db.UserTask{
		ID:       userTaskId,
		UserID:   userId,
		TaskID:   constant.TaskIdInviteUploadNftHead,
		Progress: 1,
		Ex:       formUser,
	}, "")
}
func FinishUploadNftHeadTask(userId string) error {
	// 先去创建任务，再自动去领取奖励
	userTaskId := fmt.Sprintf("%s:%d", userId, constant.TaskIdUploadNftHead)
	userTask := &db.UserTask{
		ID:       userTaskId,
		UserID:   userId,
		TaskID:   constant.TaskIdUploadNftHead,
		Progress: 1,
		Status:   constant.UserTaskStatusFinished,
	}
	return db.DB.MysqlDB.DefaultGormDB().
		Table("user_task").
		Where("id= ?", userTaskId).
		Assign(db.UserTask{Progress: 1, Status: constant.UserTaskStatusFinished}).
		FirstOrCreate(userTask).Error
}

func CloseUploadNftHeadTask(userId string) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.UserTask{}).
		Where("user_id=? and task_id=?", userId, constant.TaskIdUploadNftHead).
		Updates(map[string]interface{}{
			"progress": 0,
			"status":   constant.UserTaskStatusNoStart,
		}).Error
}

// 完成携带官方NFT任务
func FinishOfficialNFTHeadTask(userId string) error {
	// 先去创建任务，再自动去领取奖励
	userTaskId := fmt.Sprintf("%s:%d", userId, constant.TaskIdUploadOfficialNftHead)
	userTask := &db.UserTask{
		ID:       userTaskId,
		UserID:   userId,
		TaskID:   constant.TaskIdUploadOfficialNftHead,
		Progress: 1,
		Status:   constant.UserTaskStatusFinished,
	}
	return db.DB.MysqlDB.DefaultGormDB().
		Table("user_task").
		Where("id= ?", userTaskId).
		Assign(db.UserTask{Progress: 1, Status: constant.UserTaskStatusFinished}).
		FirstOrCreate(userTask).Error
}

func FinishInviteFollowOfficialTwitterTask(userId string, formUser string) error {
	userTaskId := fmt.Sprintf("%s:%d:%s", userId, constant.TaskIDInviteFollowOfficialTwitter, formUser)
	userTaskId = utils.Md5(userTaskId)
	return InsertIntoUserTask(&db.UserTask{
		ID:       userTaskId,
		UserID:   userId,
		TaskID:   constant.TaskIDInviteFollowOfficialTwitter,
		Progress: 1,
		Ex:       formUser,
	}, "")
}

// func IsFinishInviteFollowOfficialTwitterTask(userId string) bool {
// 	return IsFinishUserTask(fmt.Sprintf("%s:%d", userId, constant.TaskIDInviteFollowOfficialTwitter))
// }

// 是否完成携带NFT与新地址聊天任务
func IsFinishDailyChatNFTHeadWithNewUserTask(userId string, chatUser string) (bool, error) {
	userTaskId := fmt.Sprintf("%s:%d", userId, constant.TaskIDNFTHeadDailyChatWithNewUser)
	userTask, err := GetUserTaskByID(userTaskId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if userTask.Ex == "" {
		return false, nil
	}
	if userTask.Ex == chatUser {
		return true, nil
	}
	today := time.Now().Format("20060102")
	return userTask.StartTime.Format("20060102") == today, nil
}

// 完成携带NFT与新地址聊天任务
func FinishDailyChatNFTHeadWithNewUserTask(userId string, chatUser string) error {
	// 查看是否有记录，没有则插入，有则更新
	userTaskId := fmt.Sprintf("%s:%d", userId, constant.TaskIDNFTHeadDailyChatWithNewUser)
	userTask, err := GetUserTaskByID(userTaskId)
	today := time.Now().Format("20060102")
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return InsertIntoUserTask(&db.UserTask{
			ID:       userTaskId,
			UserID:   userId,
			TaskID:   constant.TaskIDNFTHeadDailyChatWithNewUser,
			Progress: 1,
			Ex:       chatUser,
			Status:   constant.UserTaskStatusFinished,
		}, today)
	}
	if err != nil {
		return err
	}
	remark := userTask.Ex
	// 对比 chatUser 是否是新地址
	if remark == chatUser {
		return ErrChatFinishChatWithSameOne
	}
	return ClaimDailyTaskReward(userId, userTask.TaskID, false, chatUser)
}

// 完成携带官方NFT与新地址聊天任务
func FinishOfficialNFTHeadDailyChatWithNewUserTask(userId string, chatUser string) error {
	// 查看是否有记录，没有则插入，有则更新
	userTaskId := fmt.Sprintf("%s:%d", userId, constant.TaskIDOfficialNFTHeadDailyChatWithNewUser)
	userTask, err := GetUserTaskByID(userTaskId)
	today := time.Now().Format("20060102")
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return InsertIntoUserTask(&db.UserTask{
			ID:       userTaskId,
			UserID:   userId,
			TaskID:   constant.TaskIDOfficialNFTHeadDailyChatWithNewUser,
			Progress: 1,
			Ex:       chatUser,
			Status:   constant.UserTaskStatusFinished,
		}, today)
	}
	if err != nil {
		return err
	}
	remark := userTask.Ex
	// 对比 chatUser 是否是新地址
	if remark == chatUser {
		return fmt.Errorf("not chat user with other one")
	}
	return ClaimDailyTaskReward(userId, userTask.TaskID, false, chatUser)
}

// 是否完成携带官方NFT与新地址聊天任务
func IsFinishOfficialNFTHeadDailyChatWithNewUserTask(userId string, chatUser string) (bool, error) {
	userTaskId := fmt.Sprintf("%s:%d", userId, constant.TaskIDOfficialNFTHeadDailyChatWithNewUser)
	userTask, err := GetUserTaskByID(userTaskId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if userTask.Ex == "" {
		return false, nil
	}
	if userTask.Ex == chatUser {
		return true, nil
	}
	today := time.Now().Format("20060102")
	return userTask.StartTime.Format("20060102") == today, nil
}

// 领取每日任务奖励
func ClaimDailyTaskReward(userId string, taskId int32, resetState bool, usetTaskEx string) error {
	userTaskId := fmt.Sprintf("%s:%d", userId, taskId)
	userTask, err := GetUserTaskByID(userTaskId)
	if err != nil {
		return err
	}
	task, err := GetOpenTaskById(taskId)
	if err != nil {
		return err
	}
	// 判断是否完成任务
	if userTask.Progress < task.CompletionCount {
		return fmt.Errorf("task not finished")
	}
	// 判断今天是否完成过任务
	today := time.Now().Format("20060102")
	if userTask.StartTime.Format("20060102") == today {
		return ErrTodayIsFinished
	}
	// 插入奖励，更改完成状态和startTime
	return db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		eventLogId := fmt.Sprintf("%s:%s:%s", constant.RewardTypeTask, userTaskId, today)
		rewardInfo := fmt.Sprintf(`{"type":"%s", "info":"%s"}`, constant.RewardTypeTask, task.Name)
		rewardEventLogs := &db.RewardEventLogs{
			ID:         eventLogId,
			UserID:     userId,
			RewardType: constant.RewardTypeTask,
			Reward:     task.Reward,
			CreatedAt:  time.Now(),
			UserJSON:   rewardInfo,
		}
		err := tx.Table("reward_event_logs").Create(rewardEventLogs).Error
		if err != nil {
			return err
		}
		progress := task.CompletionCount
		if resetState {
			progress = 0
		}
		// 更新一下领取时间
		return tx.Table("user_task").Where("id=?", userTaskId).
			Updates(map[string]interface{}{
				"status":     constant.UserTaskStatusFinished,
				"start_time": time.Now(),
				"progress":   progress,
				"ex":         usetTaskEx,
			}).Error
	})
}

// 是否完成每日签到任务
func IsFinishUserDailyCheckInTask(userId string) (bool, error) {
	// 如果没有 userTask 记录，插入一条
	userTaskId := fmt.Sprintf("%s:%d", userId, constant.TaskIDDailyCheckIn)
	userTask, err := GetUserTaskByID(userTaskId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	today := time.Now().Format("20060102")
	return userTask.StartTime.Format("20060102") == today, nil
}

// 完成每日签到任务
func FinishUserDailyCheckInTask(userId string) error {
	// 如果没有 userTask 记录，插入一条
	userTaskId := fmt.Sprintf("%s:%d", userId, constant.TaskIDDailyCheckIn)
	today := time.Now().Format("20060102")
	if !HasUserTaskRecord(userTaskId) {
		return InsertIntoUserTask(&db.UserTask{
			ID:       userTaskId,
			UserID:   userId,
			TaskID:   constant.TaskIDDailyCheckIn,
			Progress: 1,
			Status:   constant.UserTaskStatusFinished,
		}, today)
	}
	// 如果有记录就去领取奖励
	return ClaimDailyTaskReward(userId, constant.TaskIDDailyCheckIn, false, "")
}

func CloseOfficialNFTHeadTask(userId string) error {
	return db.DB.MysqlDB.DefaultGormDB().
		Model(&db.UserTask{}).
		Where("user_id=? and task_id=?", userId, constant.TaskIdUploadOfficialNftHead).
		Updates(map[string]interface{}{
			"progress": 0,
			"status":   constant.UserTaskStatusNoStart,
		}).Error
}

func GetTaskLastClaimTime(userId string, taskId int) int {
	userTaskId := fmt.Sprintf("%s:%d", userId, taskId)
	// 先去 user_task 表中查找是否有该任务的记录
	userTask, err := GetUserTaskByID(userTaskId)
	// 如果有记录查看它的领取时间是否超过了30天
	if err == nil {
		startTime := userTask.StartTime
		// 距离上次领取的时间超过了几天
		return int(time.Since(startTime).Hours() / 24) //天
	}
	return 0
}

func FinishCreateSpaceTask(userId string) error {
	return InsertIntoUserTask(&db.UserTask{
		ID:       fmt.Sprintf("%s:%d", userId, constant.TaskIdCreateSapce),
		UserID:   userId,
		TaskID:   constant.TaskIdCreateSapce,
		Progress: 1,
		Status:   constant.UserTaskStatusDoing,
	}, "")
}

func CancelCreateSpaceTask(userId string) error {
	return DeleteUserTaskById(fmt.Sprintf("%s:%d", userId, constant.TaskIdCreateSapce))
}

func DeleteUserTaskById(userTaskId string) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("id=?", userTaskId).Delete(&db.UserTask{}).Error
}

// 领取 开启全网推送的任务
func ClaimOpenWholeNetworkPushTaskByTx(userId string) error {
	return InsertIntoUserTask(&db.UserTask{
		ID:        fmt.Sprintf("%s:%d", userId, constant.TaskIdOpenWholeNetworkPush),
		UserID:    userId,
		Progress:  0,
		TaskID:    constant.TaskIdOpenWholeNetworkPush,
		StartTime: time.Now(),
		Status:    constant.UserTaskStatusDoing,
	}, "")
}

func ClaimTimeProgressTaskReward(userId string, taskId int32) error {
	task, err := GetOpenTaskById(taskId)
	if err != nil {
		return err
	}
	dayCount := GetTaskLastClaimTime(userId, int(taskId))
	if int32(dayCount) < task.CompletionCount {
		return ErrTimeHasNotBeenReached
	}
	return db.DB.MysqlDB.DefaultGormDB().Transaction(func(tx *gorm.DB) error {
		userTaskId := fmt.Sprintf("%s:%d", userId, taskId)
		day := time.Now().Format("20060102")
		eventLogId := fmt.Sprintf("%s:%s:%s", constant.RewardTypeTask, userTaskId, day)
		rewardInfo := fmt.Sprintf(`{"type":"%s", "info":"%s"}`, constant.RewardTypeTask, task.Name)
		rewardEventLogs := &db.RewardEventLogs{
			ID:         eventLogId,
			UserID:     userId,
			RewardType: constant.RewardTypeTask,
			Reward:     task.Reward,
			CreatedAt:  time.Now(),
			UserJSON:   rewardInfo,
		}
		err := tx.Table("reward_event_logs").Create(rewardEventLogs).Error
		if err != nil {
			return err
		}
		// 更新一下领取时间
		return tx.Table("user_task").Where("id=?", userTaskId).
			Updates(map[string]interface{}{
				"status":     constant.UserTaskStatusDoing,
				"start_time": time.Now(),
			}).Error
	})
}

func CancelClaimJoinOfficialSpaceTask(userId string) error {
	return DeleteUserTaskById(fmt.Sprintf("%s:%d", userId, constant.TaskIdJoinOfficialSpace))
}

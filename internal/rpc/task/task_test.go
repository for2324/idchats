package task

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	pbTask "Open_IM/pkg/proto/task"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetTaskList(t *testing.T) {
	s := &TaskServer{}
	optUserId := "0xcbd033ea3c05dc9504610061c86c7ae191c5c913"
	_, err := s.GetUserTaskList(context.Background(), &pbTask.GetUserTaskListReq{
		UserId: optUserId,
	})
	if err != nil {
		t.Errorf("GetUserTaskList() error = %v", err)
		return
	}
}

func TestCheckInUser(t *testing.T) {
	s := &TaskServer{}
	userId := "0xcbd033ea3c05dc9504610061c86c7ae191c5c913"
	res, err := s.DailyCheckIn(context.Background(), &pbTask.DailyCheckInReq{
		UserId: userId,
	})
	if err != nil {
		t.Errorf("DailyCheckIn() error = %v", err)
		return
	}
	if res.CommonResp.ErrCode != 0 {
		t.Errorf("DailyCheckIn() error = %v", res.CommonResp.ErrMsg)
		return
	}
	checkIn, err := db.DB.GetUserIsCheckIn(userId)
	if err != nil {
		t.Errorf("GetUserIsCheckIn() error = %v", err)
		return
	}
	if !checkIn {
		t.Errorf("GetUserIsCheckIn() error = %v", err)
		return
	}
}

// 完成每日携带NFT与新地址聊天任务
func TestFinishDailyChatNFTHeadWithNewUserTask(t *testing.T) {
	s := &TaskServer{}
	t.Run("查看是否 0x2 聊天过", func(t *testing.T) {
		resp, err := s.IsFinishDailyChatNFTHeadWithNewUserTask(context.TODO(), &pbTask.IsFinishDailyChatNFTHeadWithNewUserTaskReq{
			UserId:   "0x1",
			ChatUser: "0x2",
		})
		assert.Nil(t, err)
		assert.Equal(t, false, resp.IsFinish)
	})
	t.Run("完成和 0x2 聊天任务", func(t *testing.T) {
		resp, err := s.FinishDailyChatNFTHeadWithNewUserTask(context.TODO(), &pbTask.FinishDailyChatNFTHeadWithNewUserTaskReq{
			UserId:   "0x1",
			ChatUser: "0x2",
		})
		assert.Nil(t, err)
		assert.EqualValues(t, 0, resp.CommonResp.ErrCode)
		resp2, err := s.IsFinishDailyChatNFTHeadWithNewUserTask(context.TODO(), &pbTask.IsFinishDailyChatNFTHeadWithNewUserTaskReq{
			UserId:   "0x1",
			ChatUser: "0x2",
		})
		assert.Nil(t, err)
		assert.Equal(t, true, resp2.IsFinish)
		// 重复完成当天聊天任务
		resp3, err := s.FinishDailyChatNFTHeadWithNewUserTask(context.TODO(), &pbTask.FinishDailyChatNFTHeadWithNewUserTaskReq{
			UserId:   "0x1",
			ChatUser: "0x3",
		})
		assert.Nil(t, err)
		assert.Equal(t, resp3.CommonResp.ErrMsg, imdb.ErrTodayIsFinished.Error())
	})
	t.Run("改变完成聊天任务的时间，查看是否重复跟 0x2 聊天完成任务", func(t *testing.T) {
		userTaskId := fmt.Sprintf("%s:%d", "0x1", constant.TaskIDNFTHeadDailyChatWithNewUser)
		db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("id=? ", userTaskId).Updates(map[string]interface{}{
			"start_time": time.Now().Add(-24 * time.Hour),
		})
		db.DB.MysqlDB.DefaultGormDB().Where("user_id=? ", "0x1").Delete(&db.RewardEventLogs{})
		resp, err := s.IsFinishDailyChatNFTHeadWithNewUserTask(context.TODO(), &pbTask.IsFinishDailyChatNFTHeadWithNewUserTaskReq{
			UserId:   "0x1",
			ChatUser: "0x2",
		})
		assert.Nil(t, err)
		assert.Equal(t, resp.IsFinish, true)
		resp2, err := s.FinishDailyChatNFTHeadWithNewUserTask(context.TODO(), &pbTask.FinishDailyChatNFTHeadWithNewUserTaskReq{
			UserId:   "0x1",
			ChatUser: "0x2",
		})
		assert.Nil(t, err)
		assert.Equal(t, resp2.CommonResp.ErrMsg, imdb.ErrChatFinishChatWithSameOne.Error())
	})
	t.Run("改变完成聊天任务的时间，与新用户 0x3 聊天完成任务", func(t *testing.T) {
		resp, err := s.IsFinishDailyChatNFTHeadWithNewUserTask(context.TODO(), &pbTask.IsFinishDailyChatNFTHeadWithNewUserTaskReq{
			UserId:   "0x1",
			ChatUser: "0x3",
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp.CommonResp.ErrCode, 0)
		assert.Equal(t, resp.IsFinish, false)
		resp2, err := s.FinishDailyChatNFTHeadWithNewUserTask(context.TODO(), &pbTask.FinishDailyChatNFTHeadWithNewUserTaskReq{
			UserId:   "0x1",
			ChatUser: "0x3",
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp2.CommonResp.ErrCode, 0)
	})
}

// 完成上传NFT头像任务
func TestFinishUploadNftHeadTask(t *testing.T) {
	s := &TaskServer{}

	t.Run("测试完成和重复完成上传NFT头像任务", func(t *testing.T) {
		resp, err := s.FinishUploadNftHeadTask(context.TODO(), &pbTask.FinishUploadNftHeadTaskReq{
			UserId: "0x1",
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp.CommonResp.ErrCode, 0)

		resp2, err := s.FinishUploadNftHeadTask(context.TODO(), &pbTask.FinishUploadNftHeadTaskReq{
			UserId: "0x1",
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp2.CommonResp.ErrCode, 0)
	})
	t.Run("领取和重复领取上传NFT头像任务", func(t *testing.T) {
		resp, err := s.ClaimTaskRewards(context.TODO(), &pbTask.ClaimTaskRewardsReq{
			UserId: "0x1",
			TaskId: constant.TaskIdUploadNftHead,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp.CommonResp.ErrMsg, imdb.ErrTodayIsFinished.Error())
	})

	t.Run("改变时间再次领取上传NFT头像任务", func(t *testing.T) {
		userTaskId := fmt.Sprintf("%s:%d", "0x1", constant.TaskIdUploadNftHead)
		db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("id=? ", userTaskId).Updates(map[string]interface{}{
			"start_time": time.Now().Add(-24 * time.Hour),
		})
		db.DB.MysqlDB.DefaultGormDB().Where("user_id=? ", "0x1").Delete(&db.RewardEventLogs{})
		resp, err := s.ClaimTaskRewards(context.TODO(), &pbTask.ClaimTaskRewardsReq{
			UserId: "0x1",
			TaskId: constant.TaskIdUploadNftHead,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp.CommonResp.ErrCode, 0)
	})
}

// 签到
func TestDailyCheckIn(t *testing.T) {
	s := &TaskServer{}
	t.Run("测试签到", func(t *testing.T) {
		resp, err := s.DailyCheckIn(context.TODO(), &pbTask.DailyCheckInReq{
			UserId: "0x1",
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp.CommonResp.ErrCode, 0)
	})
	t.Run("测试重复签到", func(t *testing.T) {
		resp, err := s.DailyCheckIn(context.TODO(), &pbTask.DailyCheckInReq{
			UserId: "0x1",
		})
		assert.Nil(t, err)
		assert.Equal(t, resp.CommonResp.ErrMsg, ErrTodayIsCheckInEd.Error())
	})
	t.Run("改变时间再次签到", func(t *testing.T) {
		userTaskId := fmt.Sprintf("%s:%d", "0x1", constant.TaskIDDailyCheckIn)
		db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("id=? ", userTaskId).Updates(map[string]interface{}{
			"start_time": time.Now().Add(-24 * time.Hour),
		})
		db.DB.MysqlDB.DefaultGormDB().Where("user_id=? ", "0x1").Delete(&db.RewardEventLogs{})
		resp, err := s.DailyCheckIn(context.TODO(), &pbTask.DailyCheckInReq{
			UserId: "0x1",
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp.CommonResp.ErrCode, 0)
	})
}

// 完成订阅任务
func TestFinishJoinOfficialSpaceTask(t *testing.T) {
	s := &TaskServer{}
	resp, err := s.FinishJoinOfficialSpaceTask(context.TODO(), &pbTask.FinishJoinOfficialSpaceTaskReq{
		UserId: "0x1",
	})
	assert.Nil(t, err)
	assert.EqualValues(t, resp.CommonResp.ErrCode, 0)
	t.Run("领取订阅任务奖励", func(t *testing.T) {
		resp, err := s.ClaimTaskRewards(context.TODO(), &pbTask.ClaimTaskRewardsReq{
			UserId: "0x1",
			TaskId: constant.TaskIdJoinOfficialSpace,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp.CommonResp.ErrMsg, imdb.ErrTimeHasNotBeenReached.Error())
	})

	t.Run("改变时间再次领取订阅任务", func(t *testing.T) {
		userTaskId := fmt.Sprintf("%s:%d", "0x1", constant.TaskIdJoinOfficialSpace)
		db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("id=? ", userTaskId).Updates(map[string]interface{}{
			"start_time": time.Now().Add(-24 * time.Hour),
		})
		db.DB.MysqlDB.DefaultGormDB().Where("user_id=? ", "0x1").Delete(&db.RewardEventLogs{})
		resp, err := s.ClaimTaskRewards(context.TODO(), &pbTask.ClaimTaskRewardsReq{
			UserId: "0x1",
			TaskId: constant.TaskIdJoinOfficialSpace,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp.CommonResp.ErrMsg, imdb.ErrTimeHasNotBeenReached.Error())

		db.DB.MysqlDB.DefaultGormDB().Table("user_task").Where("id=? ", userTaskId).Updates(map[string]interface{}{
			"start_time": time.Now().Add(-30 * 24 * time.Hour),
		})
		db.DB.MysqlDB.DefaultGormDB().Where("user_id=? ", "0x1").Delete(&db.RewardEventLogs{})
		resp2, err := s.ClaimTaskRewards(context.TODO(), &pbTask.ClaimTaskRewardsReq{
			UserId: "0x1",
			TaskId: constant.TaskIdJoinOfficialSpace,
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp2.CommonResp.ErrCode, 0)
	})
}

// 领取邀请绑定推特任务
func TestFinishInviteBindTwitterTask(t *testing.T) {
	s := &TaskServer{}
	t.Run("领取邀请绑定推特任务", func(t *testing.T) {
		resp, err := s.FinishInviteBindTwitterTask(context.TODO(), &pbTask.FinishInviteBindTwitterTaskReq{
			UserId:     "0x1",
			FormUserId: "0x2",
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp.CommonResp.ErrCode, 0)
		resp2, err := s.FinishInviteBindTwitterTask(context.TODO(), &pbTask.FinishInviteBindTwitterTaskReq{
			UserId:     "0x1",
			FormUserId: "0x2",
		})
		assert.Nil(t, err)
		assert.Equal(t, "Error 1062: Duplicate entry '5af57038deae4744a795f9db5c6f285f' for key 'PRIMARY'", resp2.CommonResp.ErrMsg)

		resp3, _ := s.GetUserTaskList(context.TODO(), &pbTask.GetUserTaskListReq{
			UserId:   "0x1",
			Classify: "invite",
		})
		for _, v := range resp3.Data {
			if v.TaskId == constant.TaskIDInviteBindTwitter {
				assert.EqualValues(t, v.Status, constant.UserTaskStatusFinished)
				assert.EqualValues(t, v.Progress, 1)
			}
		}
	})
	t.Run("邀请人数状态", func(t *testing.T) {
		resp, err := s.FinishInviteBindTwitterTask(context.TODO(), &pbTask.FinishInviteBindTwitterTaskReq{
			UserId:     "0x1",
			FormUserId: "0x3",
		})
		assert.Nil(t, err)
		assert.EqualValues(t, resp.CommonResp.ErrCode, 0)
		resp2, _ := s.GetUserTaskList(context.TODO(), &pbTask.GetUserTaskListReq{
			UserId:   "0x1",
			Classify: "invite",
		})
		for _, v := range resp2.Data {
			if v.TaskId == constant.TaskIDInviteBindTwitter {
				assert.EqualValues(t, v.Status, constant.UserTaskStatusFinished)
				assert.EqualValues(t, v.Progress, 2)
			}
		}
	})
}

func finishBindingTwitterTask(userId, inviteCode string) error {
	if err := imdb.FinishBindTwitterTask(userId); err == nil {
		// 帮助邀请人领取邀请关联任务
		if inviteCode != "" {
			if err := imdb.FinishInviteBindTwitterTask(inviteCode, userId); err != nil {
				log.Info("", "FinishInviteBindTwitterTask err", err, inviteCode)
			} else {
				log.Info("", "FinishInviteBindTwitterTask success", inviteCode)
			}
		} else {
			if yaoqingren, err := imdb.GetRegisterInfo(userId); err == nil && yaoqingren.InvitationCode != "" {
				if err := imdb.FinishInviteBindTwitterTask(yaoqingren.InvitationCode, userId); err != nil {
					log.Info("", "FinishInviteBindTwitterTask err", err, yaoqingren.InvitationCode)
				} else {
					log.Info("", "FinishInviteBindTwitterTask success", yaoqingren.InvitationCode)
				}
			}
		}
	}
	return nil
}

func TestSimulateBindingTwitter(t *testing.T) {
	simulateBindingTwitter("0x1", "knight", "0x2")
}

func TestSimulateFollowTwitter(t *testing.T) {
	simulateFollowTwitter("0x1", "0x2")
}

func simulateBindingTwitter(userId, twitterName, inviteCode string) error {
	userBind, err := imdb.HasBindTwitter(userId)
	if err != nil {
		return err
	}
	if userBind {
		return fmt.Errorf("禁止重复绑定推特用户")
	}
	//判断是否有其数据， 如果没有这个数据的情况 做插入的操作， 如果存在这个数据就做删除的操作。
	var userthird db.UserThird
	err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").Where("user_id=? ", userId).First(&userthird).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userthird.UserId = userId
		userthird.Twitter = twitterName
		err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").Create(&userthird).Error
	} else if err == nil {
		userthird.UserId = userId
		userthird.Twitter = twitterName
		err = db.DB.MysqlDB.DefaultGormDB().Table("user_third").Where("user_id=?", userId).
			Updates(&userthird).Error
	}
	if err != nil {
		return err
	}
	err = imdb.UpdateUserInfo(db.User{
		UserID: userId,
		Ex:     fmt.Sprintf(`{"twitter":"%s"}`, twitterName),
	})
	if err != nil {
		return err
	}
	if err := imdb.FinishBindTwitterTask(userId); err == nil {
		// 帮助邀请人领取邀请关联任务
		if inviteCode != "" {
			if err := imdb.FinishInviteBindTwitterTask(inviteCode, userId); err != nil {
				log.Info("", "FinishInviteBindTwitterTask err", err, inviteCode)
			} else {
				log.Info("", "FinishInviteBindTwitterTask success", inviteCode)
			}
		} else {
			if yaoqingren, err := imdb.GetRegisterInfo(userId); err == nil && yaoqingren.InvitationCode != "" {
				if err := imdb.FinishInviteBindTwitterTask(yaoqingren.InvitationCode, userId); err != nil {
					log.Info("", "FinishInviteBindTwitterTask err", err, yaoqingren.InvitationCode)
				} else {
					log.Info("", "FinishInviteBindTwitterTask success", yaoqingren.InvitationCode)
				}
			}
		}
	}
	return nil
}

func simulateFollowTwitter(userId, inviteCode string) error {
	if finish, err := imdb.IsFinishFollowOfficialTwitterTask(userId); err != nil {
		return err
	} else if finish {
		return fmt.Errorf("禁止重复关注官方推特")
	}
	//检查是否绑定的twitter
	twitterNameValue, err := imdb.GetUserTwitter(userId)
	if err != nil {
		return err
	}
	if twitterNameValue == "" {
		return fmt.Errorf("your not  bind twitte")
	}
	//已经有关注 旧需要去新增事件
	if err := imdb.FinishFollowOfficialTwitterTask(userId); err == nil {
		// 帮助邀请人领取邀请关联任务
		if inviteCode != "" {
			if err := imdb.FinishInviteFollowOfficialTwitterTask(inviteCode, userId); err != nil {
				log.Info("", "FinishInviteFollowOfficialTwitterTask err", err, inviteCode)
			} else {
				log.Info("", "FinishInviteFollowOfficialTwitterTask success", inviteCode)
			}
		} else {
			if yaoqingren, err := imdb.GetRegisterInfo(userId); err == nil && yaoqingren.InvitationCode != "" {
				if err := imdb.FinishInviteFollowOfficialTwitterTask(yaoqingren.InvitationCode, userId); err != nil {
					log.Info("", "FinishInviteFollowOfficialTwitterTask err", err, yaoqingren.InvitationCode)
				} else {
					log.Info("", "FinishInviteFollowOfficialTwitterTask success", yaoqingren.InvitationCode)
				}
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}
func TestCloseUploadNftHeadTask(t *testing.T) {
	s := &TaskServer{}
	resp, err := s.CloseUploadNftHeadTask(context.TODO(), &pbTask.CloseUploadNftHeadTaskReq{
		UserId: "0xcbd033ea3c05dc9504610061c86c7ae191c5c913",
	})
	assert.Nil(t, err)
	assert.EqualValues(t, resp.CommonResp.ErrCode, 0)
}

func TestMain(m *testing.M) {
	// Open database connection
	// 清空历史数据，要确保连接到的数据库是本地的
	db.DB.MysqlDB.DefaultGormDB().Where("user_id = ?", "0x1").Delete(&db.UserTask{})
	db.DB.MysqlDB.DefaultGormDB().Where("user_id = ?", "0x1").Delete(&db.RewardEventLogs{})

	imdb.UserRegister(db.User{
		UserID:   "0x1",
		Nickname: "0x1",
	})
	m.Run()
}

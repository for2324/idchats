package task

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbTask "Open_IM/pkg/proto/task"
	pbWeb3 "Open_IM/pkg/proto/web3pub"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
)

var (
	ErrTodayIsCheckInEd = errors.New("Signed in today")
)

type TaskServer struct {
	rpcPort          int
	rpcRegisterName  string
	etcdSchema       string
	etcdAddr         []string
	rpcWeb3PubClient pbWeb3.Web3PubClient
}

func NewTaskServer(port int) *TaskServer {
	log.NewPrivateLog(constant.LogFileName)
	return &TaskServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImTask,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *TaskServer) Run() {
	log.NewInfo("0", "TaskServer run...")

	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)

	//listener network
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + s.rpcRegisterName)
	}
	log.NewInfo("0", "listen ok ", address)
	defer listener.Close()
	//grpc server
	var grpcOpts []grpc.ServerOption
	if config.Config.Prometheus.Enable {
		promePkg.NewGrpcRequestCounter()
		promePkg.NewGrpcRequestFailedCounter()
		promePkg.NewGrpcRequestSuccessCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	// web3 client
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImWeb3Js, "")
	if etcdConn == nil {
		panic("task Run server getcdv3.GetDefaultConn == nil")
	}
	s.rpcWeb3PubClient = pbWeb3.NewWeb3PubClient(etcdConn)

	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()
	//User friend related services register to etcd
	pbTask.RegisterTaskServiceServer(srv, s)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema,
		strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(), s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register friend module  rpc to etcd err"))
	}

	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error(), listener)
		return
	}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *pbTask.CreateTaskReq) (*pbTask.CreateTaskResp, error) {
	for _, v := range req.TaskList {
		taskInfo := db.Task{}
		utils.CopyStructFields(&taskInfo, v)
		err := imdb.CreateOrUpdateTask(&taskInfo)
		if err != nil {
			log.NewError(req.OperationID, "InsertIntoTask failed, ", err.Error(), taskInfo)
			return &pbTask.CreateTaskResp{CommonResp: &pbTask.CommonResp{
				ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg,
			}}, nil
		}
		log.NewInfo(req.OperationID, "InsertIntoTask success, ", taskInfo)
	}
	return &pbTask.CreateTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

func (s *TaskServer) GetTaskList(ctx context.Context, req *pbTask.GetTaskListReq) (*pbTask.GetTaskListResp, error) {
	var (
		resp pbTask.GetTaskListResp
	)
	tasks, err := imdb.GetTaskList(req.Classify)
	if err != nil {
		log.NewError(req.OperationID, "GetTaskList failed, ", err.Error(), tasks)
		resp.CommonResp = &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg,
		}
		return &resp, nil
	}
	for _, v := range tasks {
		var node pbTask.Task
		utils.CopyStructFields(&node, &v)
		resp.Data = append(resp.Data, &node)
	}
	resp.CommonResp = &pbTask.CommonResp{}
	return &resp, nil
}
func (s *TaskServer) GetUserClaimTaskList(ctx context.Context, req *pbTask.GetUserClaimTaskListReq) (*pbTask.GetUserClaimTaskListResp, error) {
	var (
		resp pbTask.GetUserClaimTaskListResp
	)
	tasks, err := imdb.GetUserClaimTask(req.UserId, int(req.Status))
	if err != nil {
		log.NewError(req.OperationID, "GetTaskList failed, ", err.Error(), tasks)
		resp.CommonResp = &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg,
		}
		return &resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetUserClaimTask: ", tasks)
	for _, v := range tasks {
		var node pbTask.UserTask
		utils.CopyStructFields(&node, &v)
		resp.Data = append(resp.Data, &node)
	}
	resp.CommonResp = &pbTask.CommonResp{}
	return &resp, nil
}

func (s *TaskServer) GetUserTaskList(ctx context.Context, req *pbTask.GetUserTaskListReq) (*pbTask.GetUserTaskListResp, error) {
	var (
		resp pbTask.GetUserTaskListResp
	)
	tasks, err := imdb.GetUserClaimTask(req.UserId, 0)
	if err != nil {
		log.NewError(req.OperationID, "GetTaskList failed, ", err.Error(), tasks)
		resp.CommonResp = &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg,
		}
		return &resp, nil
	}
	userTaskIDMap := make(map[string]db.UserTask)
	for _, v := range tasks {
		userTaskIDMap[v.ID] = v
	}
	// get all task
	allTasks, err := imdb.GetTaskList(req.Classify)
	if err != nil {
		log.NewError(req.OperationID, "GetTaskList failed, ", err.Error(), tasks)
		resp.CommonResp = &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg,
		}
		return &resp, nil
	}
	for _, v := range allTasks {
		task := pbTask.Task{}
		utils.CopyStructFields(&task, &v)
		task.Id = v.Id
		node := pbTask.UserTask{
			UserId:   req.UserId,
			TaskId:   v.Id,
			Id:       fmt.Sprintf("%s:%d", req.UserId, v.Id),
			Task:     &task,
			Status:   constant.UserTaskStatusNoStart,
			Progress: 0,
		}
		if v.Classify == "invite" {
			node.Progress = int32(imdb.GetUserTaskCount(req.UserId, v.Id))
			if node.Progress >= v.CompletionCount {
				node.Status = constant.UserTaskStatusFinished
			}
			resp.Data = append(resp.Data, &node)
			continue
		}
		// 判断是否有领取记录，如果没有的话就是未领取
		if userTask, ok := userTaskIDMap[node.Id]; ok {
			switch userTask.Status {
			// 已领取直接返回
			case constant.UserTaskStatusClaimed:
				node.Status = int32(userTask.Status)
				node.Progress = userTask.Progress
			// 已完成 判断是否同步
			case constant.UserTaskStatusFinished:
				node.Status = int32(userTask.Status)
				node.Progress = userTask.Progress
				if v.Type == constant.TaskTypeDaily {
					// 判断是否是今天完成的任务（对比 startTime）
					if userTask.StartTime.Format("20060102") == time.Now().Format("20060102") {
						node.Status = constant.UserTaskStatusFinished
						node.Progress = userTask.Progress
					} else {
						// 正常每日任务，自动完成，手动领取
						node.Status = constant.UserTaskStatusDoing
						node.Progress = userTask.Progress
						// 手动完成领取并自动领取
						if s.CheckTaskIsExpire(req.OperationID, req.UserId, v.Id) {
							node.Status = constant.UserTaskStatusNoStart
							node.Progress = 0
						}
					}
				}
			case constant.UserTaskStatusDoing:
				node.Status = int32(userTask.Status)
				node.Progress = userTask.Progress
				// 如果是时间累计型任务，让用户直接领取
				if task.Type == constant.TaskTypeTimeProgress {
					day := int32(time.Since(userTask.StartTime).Hours() / 24)
					node.Progress = day
					if node.Progress >= task.CompletionCount {
						node.Progress = task.CompletionCount
					}
				}
			default:

			}
		}
		resp.Data = append(resp.Data, &node)
	}
	resp.CommonResp = &pbTask.CommonResp{}
	return &resp, nil
}

// 检查任务是否过期
func (s *TaskServer) CheckTaskIsExpire(OperationID, userId string, taskId int32) bool {
	// 如果是签到任务、每日互动任务，得手动完成，自动领取
	if taskId == constant.TaskIDNFTHeadDailyChatWithNewUser || taskId == constant.TaskIDOfficialNFTHeadDailyChatWithNewUser || taskId == constant.TaskIDDailyCheckIn {
		return true
	}
	// 上传头像任务，检查头像是否有效
	if taskId == constant.TaskIdUploadNftHead {
		resp, err := s.rpcWeb3PubClient.CheckIsHaveNftRecvID(context.Background(), &pbWeb3.CheckIsHaveNftRecvIDReq{
			OperatorID: OperationID,
			UserId:     userId,
		})
		if err != nil {
			log.NewError(OperationID, "CheckIsHaveNftRecvID failed, ", err.Error())
			return true
		}
		if !resp.HaveNft {
			imdb.CloseUploadNftHeadTask(userId)
		}
		return !resp.HaveNft
	}
	if taskId == constant.TaskIdUploadOfficialNftHead {
		resp, err := s.rpcWeb3PubClient.CheckIsHaveGuanFangNftRecvID(context.Background(), &pbWeb3.CheckIsHaveGuanFangNftRecvIDReq{
			OperatorID: OperationID,
			UserId:     userId,
		})
		if err != nil {
			log.NewError(OperationID, "CheckIsHaveNftRecvID failed, ", err.Error())
			return true
		}
		if !resp.HaveNft {
			imdb.CloseOfficialNFTHeadTask(userId)
		}
		return !resp.HaveNft
	}
	return false
}

func (s *TaskServer) DailyCheckIn(ctx context.Context, req *pbTask.DailyCheckInReq) (*pbTask.DailyCheckInResp, error) {
	// 判断是否已经签到
	isCheckIn, err := imdb.IsFinishUserDailyCheckInTask(req.UserId)
	if err != nil {
		log.NewError(req.OperationID, "IsFinishUserDailyCheckInTask failed, ", req.UserId, err.Error())
		return &pbTask.DailyCheckInResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	if isCheckIn {
		if err := db.DB.SetUserCheckIn(req.UserId); err != nil {
			log.NewError(req.OperationID, "SetUserCheckIn failed, ", req.UserId, err.Error())
		}
		return &pbTask.DailyCheckInResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: ErrTodayIsCheckInEd.Error(),
		}}, nil
	}
	// 签到
	if err := imdb.FinishUserDailyCheckInTask(req.UserId); err != nil {
		log.NewError(req.OperationID, "UserDailySign failed, ", req.UserId, err.Error())
		return &pbTask.DailyCheckInResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	if err := db.DB.SetUserCheckIn(req.UserId); err != nil {
		log.NewError(req.OperationID, "SetUserCheckIn failed 2, ", req.UserId, err.Error())
	}
	return &pbTask.DailyCheckInResp{CommonResp: &pbTask.CommonResp{}}, nil
}

func (s *TaskServer) ClaimTaskRewards(ctx context.Context, req *pbTask.ClaimTaskRewardsReq) (*pbTask.ClaimTaskRewardsResp, error) {
	// get task
	task, err := imdb.GetTaskById(req.TaskId)
	if err != nil {
		log.NewError(req.OperationID, "GetTaskById failed, ", err.Error(), req.TaskId)
		return &pbTask.ClaimTaskRewardsResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg,
		}}, nil
	}
	// 创建空间,检查是否 nft >= 100
	if task.Id == constant.TaskIdCreateSapce {
		groupInfo, err := imdb.GetOneGroupInfoByUserID(req.UserId)
		if err != nil {
			log.NewError(req.OperationID, "GetOneGroupInfoByUserID failed, ", req.UserId, err.Error())
			return &pbTask.ClaimTaskRewardsResp{CommonResp: &pbTask.CommonResp{
				ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
			}}, nil
		}
		count, err := imdb.GetGroupHaveNftMemberCount(groupInfo.GroupID)
		if err != nil {
			log.NewError(req.OperationID, "GetOneGroupInfoByUserID failed, ", req.UserId, err.Error())
			return &pbTask.ClaimTaskRewardsResp{CommonResp: &pbTask.CommonResp{
				ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
			}}, nil
		}
		log.NewDebug(req.OperationID, "GetGroupHaveNftMemberCount", count)
		if count < 100 {
			return &pbTask.ClaimTaskRewardsResp{CommonResp: &pbTask.CommonResp{
				ErrCode: constant.ErrDB.ErrCode, ErrMsg: "nft count < 100",
			}}, nil
		}
	}
	// 日常任务
	if task.Type == constant.TaskTypeDaily {
		resetState := false
		// 互动任务重置状态
		if task.Id == constant.TaskIDNFTHeadDailyChatWithNewUser || task.Id == constant.TaskIDOfficialNFTHeadDailyChatWithNewUser {
			resetState = true
		}
		err := imdb.ClaimDailyTaskReward(req.UserId, req.TaskId, resetState, "")
		if err != nil {
			log.NewError(req.OperationID, "ClaimDailyTaskReward failed, ", err.Error(), req.UserId, req.TaskId)
			return &pbTask.ClaimTaskRewardsResp{CommonResp: &pbTask.CommonResp{
				ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
			}}, nil
		}
		log.NewInfo(req.OperationID, "ClaimDailyTaskReward success, ", req.UserId, req.TaskId)
	}
	// 时间累计任务
	if task.Type == constant.TaskTypeTimeProgress {
		if err := imdb.ClaimTimeProgressTaskReward(req.UserId, req.TaskId); err != nil {
			log.NewError(req.OperationID, "ClaimTimeProgressTaskReward failed, ", req.UserId, req.TaskId, err.Error())
			return &pbTask.ClaimTaskRewardsResp{CommonResp: &pbTask.CommonResp{
				ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
			}}, nil
		}
		log.NewError(req.OperationID, "ClaimTimeProgressTaskReward success, ", req.UserId)
	}
	return &pbTask.ClaimTaskRewardsResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 完成携带NFT与新地址聊天任务
func (s *TaskServer) FinishDailyChatNFTHeadWithNewUserTask(ctx context.Context, req *pbTask.FinishDailyChatNFTHeadWithNewUserTaskReq) (*pbTask.FinishDailyChatNFTHeadWithNewUserTaskResp, error) {
	if err := imdb.FinishDailyChatNFTHeadWithNewUserTask(req.UserId, req.ChatUser); err != nil {
		log.NewError(req.OperationID, "FinishDailyChatNFTHeadWithNewUserTask failed, ", req.UserId, req.ChatUser, err.Error())
		return &pbTask.FinishDailyChatNFTHeadWithNewUserTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.FinishDailyChatNFTHeadWithNewUserTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 是否完成携带NFT与新地址聊天任务
func (s *TaskServer) IsFinishDailyChatNFTHeadWithNewUserTask(ctx context.Context, req *pbTask.IsFinishDailyChatNFTHeadWithNewUserTaskReq) (*pbTask.IsFinishDailyChatNFTHeadWithNewUserTaskResp, error) {
	isFinish, err := imdb.IsFinishDailyChatNFTHeadWithNewUserTask(req.UserId, req.ChatUser)
	if err != nil {
		log.NewError(req.OperationID, "IsFinishDailyChatNFTHeadWithNewUserTask failed, ", req.UserId, req.ChatUser, err.Error())
		return &pbTask.IsFinishDailyChatNFTHeadWithNewUserTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.IsFinishDailyChatNFTHeadWithNewUserTaskResp{CommonResp: &pbTask.CommonResp{}, IsFinish: isFinish}, nil
}

// 是否完成携带官方NFT与新地址聊天任务
func (s *TaskServer) IsFinishOfficialNFTHeadDailyChatWithNewUserTask(ctx context.Context, req *pbTask.IsFinishOfficialNFTHeadDailyChatWithNewUserTaskReq) (*pbTask.IsFinishOfficialNFTHeadDailyChatWithNewUserTaskResp, error) {
	isFinish, err := imdb.IsFinishOfficialNFTHeadDailyChatWithNewUserTask(req.UserId, req.ChatUser)
	if err != nil {
		log.NewError(req.OperationID, "IsFinishOfficialNFTHeadDailyChatWithNewUserTask failed, ", req.UserId, req.ChatUser, err.Error())
		return &pbTask.IsFinishOfficialNFTHeadDailyChatWithNewUserTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.IsFinishOfficialNFTHeadDailyChatWithNewUserTaskResp{CommonResp: &pbTask.CommonResp{}, IsFinish: isFinish}, nil
}

// 完成携带NFT与新地址聊天任务
func (s *TaskServer) FinishOfficialNFTHeadDailyChatWithNewUserTask(ctx context.Context, req *pbTask.FinishOfficialNFTHeadDailyChatWithNewUserTaskReq) (*pbTask.FinishOfficialNFTHeadDailyChatWithNewUserTaskResp, error) {
	if err := imdb.FinishOfficialNFTHeadDailyChatWithNewUserTask(req.UserId, req.ChatUser); err != nil {
		log.NewError(req.OperationID, "FinishOfficialNFTHeadDailyChatWithNewUserTask failed, ", req.UserId, req.ChatUser, err.Error())
		return &pbTask.FinishOfficialNFTHeadDailyChatWithNewUserTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.FinishOfficialNFTHeadDailyChatWithNewUserTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 完成上传NFT头像任务
func (s *TaskServer) FinishUploadNftHeadTask(ctx context.Context, req *pbTask.FinishUploadNftHeadTaskReq) (*pbTask.FinishUploadNftHeadTaskResp, error) {
	if err := imdb.FinishUploadNftHeadTask(req.UserId); err != nil {
		log.NewError(req.OperationID, "FinishUploadNftHeadTask failed, ", req.UserId, err.Error())
		return &pbTask.FinishUploadNftHeadTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	// 自动发放奖励
	imdb.ClaimDailyTaskReward(req.UserId, constant.TaskIdUploadNftHead, false, "")
	// 帮助邀请人领取邀请上传头像任务
	if yaoqingren, err := imdb.GetRegisterInfo(req.UserId); err == nil && yaoqingren.InvitationCode != "" {
		if err := imdb.FinishInviteUploadNftHeadTask(yaoqingren.InvitationCode, req.UserId); err != nil {
			log.Info(req.OperationID, "FinishInviteUploadNftHeadTask err", err, yaoqingren.InvitationCode)
		} else {
			log.Info(req.OperationID, "FinishInviteUploadNftHeadTask success", yaoqingren.InvitationCode)
		}
	}
	return &pbTask.FinishUploadNftHeadTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 是否完成上传NFT头像任务
func (s *TaskServer) IsFinishUploadNftHeadTask(ctx context.Context, req *pbTask.IsFinishUploadNftHeadTaskReq) (*pbTask.IsFinishUploadNftHeadTaskResp, error) {
	userTaskId := fmt.Sprintf("%s:%d", req.UserId, constant.TaskIdUploadNftHead)
	isFinish, err := imdb.IsFinishUserTask(userTaskId)
	if err != nil {
		log.NewError(req.OperationID, "IsFinishUploadNftHeadTask failed, ", req.UserId, err.Error())
		return &pbTask.IsFinishUploadNftHeadTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.IsFinishUploadNftHeadTaskResp{CommonResp: &pbTask.CommonResp{}, IsFinish: isFinish}, nil
}

// 完成官方NFT头像任务
func (s *TaskServer) FinishOfficialNFTHeadTask(ctx context.Context, req *pbTask.FinishOfficialNFTHeadTaskReq) (*pbTask.FinishOfficialNFTHeadTaskResp, error) {
	if err := imdb.FinishOfficialNFTHeadTask(req.UserId); err != nil {
		log.NewError(req.OperationID, "FinishOfficialNFTHeadTask failed, ", req.UserId, err.Error())
		return &pbTask.FinishOfficialNFTHeadTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	// 自动发放奖励
	imdb.ClaimDailyTaskReward(req.UserId, constant.TaskIdUploadOfficialNftHead, false, "")
	return &pbTask.FinishOfficialNFTHeadTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

func (s *TaskServer) IsFinishOfficialNFTHeadTask(ctx context.Context, req *pbTask.IsFinishOfficialNFTHeadTaskReq) (*pbTask.IsFinishOfficialNFTHeadTaskResp, error) {
	userTaskId := fmt.Sprintf("%s:%d", req.UserId, constant.TaskIdUploadOfficialNftHead)
	isFinish, err := imdb.IsFinishUserTask(userTaskId)
	if err != nil {
		log.NewError(req.OperationID, "IsFinishOfficialNFTHeadTask failed, ", req.UserId, err.Error())
		return &pbTask.IsFinishOfficialNFTHeadTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.IsFinishOfficialNFTHeadTaskResp{CommonResp: &pbTask.CommonResp{}, IsFinish: isFinish}, nil
}

// 完成创建空间任务
func (s *TaskServer) FinishCreateSpaceTask(ctx context.Context, req *pbTask.FinishCreateSpaceTaskReq) (*pbTask.FinishCreateSpaceTaskResp, error) {
	if err := imdb.FinishCreateSpaceTask(req.UserId); err != nil {
		log.NewError(req.OperationID, "FinishCreateSpaceTask failed, ", req.UserId, err.Error())
		return &pbTask.FinishCreateSpaceTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.FinishCreateSpaceTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 完成订阅官方空间
func (s *TaskServer) FinishJoinOfficialSpaceTask(ctx context.Context, req *pbTask.FinishJoinOfficialSpaceTaskReq) (*pbTask.FinishJoinOfficialSpaceTaskResp, error) {
	if err := imdb.FinishJoinOfficialSpaceTask(req.UserId); err != nil {
		log.NewError(req.OperationID, "FinishJoinOfficialSpaceTask failed, ", req.UserId, err.Error())
		return &pbTask.FinishJoinOfficialSpaceTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.FinishJoinOfficialSpaceTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 完成关注官方推特
func (s *TaskServer) FinishFollowOfficialTwitterTask(ctx context.Context, req *pbTask.FinishFollowOfficialTwitterTaskReq) (*pbTask.FinishFollowOfficialTwitterTaskResp, error) {
	if err := imdb.FinishFollowOfficialTwitterTask(req.UserId); err != nil {
		log.NewError(req.OperationID, "FinishFollowOfficialTwitterTask failed, ", req.UserId, err.Error())
		return &pbTask.FinishFollowOfficialTwitterTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	// 帮助邀请人领取邀请关联任务
	if yaoqingren, err := imdb.GetRegisterInfo(req.UserId); err == nil && yaoqingren.InvitationCode != "" {
		if err := imdb.FinishInviteFollowOfficialTwitterTask(yaoqingren.InvitationCode, req.UserId); err != nil {
			log.NewError(req.OperationID, "FinishInviteFollowOfficialTwitterTask err", err, yaoqingren.InvitationCode)
		} else {
			log.Info(req.OperationID, "FinishInviteFollowOfficialTwitterTask success", yaoqingren.InvitationCode)
		}
	}
	return &pbTask.FinishFollowOfficialTwitterTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 完成绑定推特任务
func (s *TaskServer) FinishBindTwitterTask(ctx context.Context, req *pbTask.FinishBindTwitterTaskReq) (*pbTask.FinishBindTwitterTaskResp, error) {
	if err := imdb.FinishBindTwitterTask(req.UserId); err != nil {
		log.NewError(req.OperationID, "FinishBindTwitterTask failed, ", req.UserId, err.Error())
		return &pbTask.FinishBindTwitterTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	// 帮助邀请人领取邀请关联任务
	if yaoqingren, err := imdb.GetRegisterInfo(req.UserId); err == nil && yaoqingren.InvitationCode != "" {
		if err := imdb.FinishInviteBindTwitterTask(yaoqingren.InvitationCode, req.UserId); err != nil {
			log.NewError(req.OperationID, "FinishInviteBindTwitterTask err", err, yaoqingren.InvitationCode)
		} else {
			log.Info(req.OperationID, "FinishInviteBindTwitterTask success", yaoqingren.InvitationCode)
		}
	}
	return &pbTask.FinishBindTwitterTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 完成邀请绑定推特任务
func (s *TaskServer) FinishInviteBindTwitterTask(ctx context.Context, req *pbTask.FinishInviteBindTwitterTaskReq) (*pbTask.FinishInviteBindTwitterTaskResp, error) {
	if err := imdb.FinishInviteBindTwitterTask(req.UserId, req.FormUserId); err != nil {
		log.NewError(req.OperationID, "FinishInviteBindTwitterTask failed, ", req.UserId, err.Error())
		return &pbTask.FinishInviteBindTwitterTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.FinishInviteBindTwitterTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 完成邀请绑定头像
func (s *TaskServer) FinishInviteUploadNftHeadTask(ctx context.Context, req *pbTask.FinishInviteUploadNftHeadTaskReq) (*pbTask.FinishInviteUploadNftHeadTaskResp, error) {
	if err := imdb.FinishInviteUploadNftHeadTask(req.UserId, req.FormUserId); err != nil {
		log.NewError(req.OperationID, "FinishInviteUploadNftHeadTask failed, ", req.UserId, err.Error())
		return &pbTask.FinishInviteUploadNftHeadTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.FinishInviteUploadNftHeadTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 完成邀请关注官方推特
func (s *TaskServer) FinishInviteFollowOfficialTwitterTask(ctx context.Context, req *pbTask.FinishInviteFollowOfficialTwitterTaskReq) (*pbTask.FinishInviteFollowOfficialTwitterTaskResp, error) {
	if err := imdb.FinishInviteFollowOfficialTwitterTask(req.UserId, req.FormUserId); err != nil {
		log.NewError(req.OperationID, "FinishInviteFollowOfficialTwitterTask failed, ", req.UserId, err.Error())
		return &pbTask.FinishInviteFollowOfficialTwitterTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.FinishInviteFollowOfficialTwitterTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 取消官方NFT头像任务
func (s *TaskServer) CloseOfficialNFTHeadTask(ctx context.Context, req *pbTask.CloseOfficialNFTHeadTaskReq) (*pbTask.CloseOfficialNFTHeadTaskResp, error) {
	if err := imdb.CloseOfficialNFTHeadTask(req.UserId); err != nil {
		log.NewError(req.OperationID, "CloseOfficialNFTHeadTask failed, ", req.UserId, err.Error())
		return &pbTask.CloseOfficialNFTHeadTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.CloseOfficialNFTHeadTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 取消官方上传NFT头像任务
func (s *TaskServer) CloseUploadNftHeadTask(ctx context.Context, req *pbTask.CloseUploadNftHeadTaskReq) (*pbTask.CloseUploadNftHeadTaskResp, error) {
	if err := imdb.CloseUploadNftHeadTask(req.UserId); err != nil {
		log.NewError(req.OperationID, "CloseUploadNftHeadTask failed, ", req.UserId, err.Error())
		return &pbTask.CloseUploadNftHeadTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.CloseUploadNftHeadTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 取消创建空间任务
func (s *TaskServer) CancelCreateSpaceTask(ctx context.Context, req *pbTask.CancelCreateSpaceTaskReq) (*pbTask.CancelCreateSpaceTaskResp, error) {
	if err := imdb.CancelCreateSpaceTask(req.UserId); err != nil {
		log.NewError(req.OperationID, "CancelCreateSpaceTask failed, ", req.UserId, err.Error())
		return &pbTask.CancelCreateSpaceTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.CancelCreateSpaceTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

// 取消加入官方空间任务
func (s *TaskServer) CancelClaimJoinOfficialSpaceTask(ctx context.Context, req *pbTask.CancelClaimJoinOfficialSpaceTaskReq) (*pbTask.CancelClaimJoinOfficialSpaceTaskResp, error) {
	if err := imdb.CancelClaimJoinOfficialSpaceTask(req.UserId); err != nil {
		log.NewError(req.OperationID, "CancelClaimJoinOfficialSpaceTask failed, ", req.UserId, err.Error())
		return &pbTask.CancelClaimJoinOfficialSpaceTaskResp{CommonResp: &pbTask.CommonResp{
			ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error(),
		}}, nil
	}
	return &pbTask.CancelClaimJoinOfficialSpaceTaskResp{CommonResp: &pbTask.CommonResp{}}, nil
}

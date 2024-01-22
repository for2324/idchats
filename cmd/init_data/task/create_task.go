package task

import (
	rpcTask "Open_IM/internal/rpc/task"
	pbTask "Open_IM/pkg/proto/task"

	"context"
	"fmt"
)

func CreateTaskList() {
	fmt.Println("CreateTaskList start ... ")
	// claimConditions
	// 1：头像设置NFT
	// 2：头像设置为 Biubiu NFT
	// 3：陌生人互动
	// 4：开启全网推送
	// 5：订阅官方空间
	// 6：创建空间
	// 7：NFT人数大于100
	// 8：关联推特
	// 9：关注官方推特
	// 10：邀请关联推特
	// 11：邀请关注 Biubiu
	// 12：签到
	// 13：邀请上传nft
	taskList := []*pbTask.Task{
		{
			Id:              4,
			Name:            "陌生地址互动",
			Head:            "",
			Type:            "Daily",
			EventType:       "chat_5",
			Classify:        "daily",
			Desc:            "头像设置为NFT，与陌生地址完成一次聊天互动",
			Reward:          50000,
			CompletionCount: 1,
			Status:          0,
			ClaimConditions: "1,3",
			StartTime:       0,
			EndTime:         0,
		},
		{
			Id:              7,
			Name:            "陌生地址互动",
			Head:            "",
			Type:            "Daily",
			EventType:       "chat_2",
			Classify:        "daily",
			Desc:            "头像设置为IDChats NFT，与陌生地址完成一次聊天互动",
			Reward:          100000,
			CompletionCount: 1,
			Status:          0,
			ClaimConditions: "2,3",
			StartTime:       0,
			EndTime:         0,
		},
		{
			Id:              1,
			Name:            "关联推特",
			Head:            "",
			Type:            "CountProgress",
			EventType:       "phone",
			Classify:        "follow",
			Desc:            "关联个人推特",
			Reward:          25000,
			CompletionCount: 1,
			Status:          0,
			ClaimConditions: "8",
			StartTime:       0,
			EndTime:         0,
		},
		{
			Id:              3,
			Name:            "关注官方推特",
			Head:            "",
			Type:            "CountProgress",
			EventType:       "follow_twitter",
			Classify:        "follow",
			Desc:            "关注官方推特@IDChatsBD",
			Reward:          50000,
			CompletionCount: 1,
			Status:          0,
			ClaimConditions: "9",
			StartTime:       0,
			EndTime:         0,
		},
		{
			Id:              8,
			Name:            "邀请关联推特",
			Head:            "",
			Type:            "CountProgress",
			EventType:       "invite_friends",
			Classify:        "invite",
			Desc:            "邀请新地址关联推特",
			Reward:          50000,
			CompletionCount: 1,
			Status:          0,
			ClaimConditions: "10",
			StartTime:       0,
			EndTime:         0,
		},
		{
			Id:              9,
			Name:            "邀请关注@IDChatsBD",
			Head:            "",
			Type:            "CountProgress",
			EventType:       "invite_twitter",
			Classify:        "invite",
			Desc:            "邀请新地址关注@IDChatsBD",
			Reward:          50000,
			CompletionCount: 1,
			Status:          0,
			ClaimConditions: "11",
			StartTime:       0,
			EndTime:         0,
		},
		{
			Id:              11,
			Name:            "每日签到",
			Head:            "",
			Type:            "Daily",
			EventType:       "sign",
			Classify:        "sign",
			Desc:            "每日签到",
			Reward:          1000,
			CompletionCount: 1,
			Status:          0,
			ClaimConditions: "12",
			StartTime:       0,
			EndTime:         0,
		},
		{
			Id:              12,
			Name:            "全网推送",
			Head:            "",
			Type:            "TimeProgress",
			EventType:       "open_full_push",
			Classify:        "space",
			Desc:            "打开全网推送，累计30天",
			Reward:          10000,
			CompletionCount: 30,
			Status:          0,
			ClaimConditions: "4",
			StartTime:       0,
			EndTime:         0,
		},
		{
			Id:              13,
			Name:            "订阅空间",
			Head:            "",
			Type:            "TimeProgress",
			EventType:       "sub_space",
			Classify:        "space",
			Desc:            "加入官方空间，累计30天",
			Reward:          10000,
			CompletionCount: 30,
			Status:          1,
			ClaimConditions: "5",
			StartTime:       0,
			EndTime:         0,
		},
		{
			Id:              14,
			Name:            "空间福利",
			Head:            "",
			Type:            "TimeProgress",
			EventType:       "create_space",
			Classify:        "space",
			Desc:            "空间头像为NFT且人数>100,累计30天",
			Reward:          250000,
			CompletionCount: 30,
			Status:          1,
			ClaimConditions: "6,7",
			StartTime:       0,
			EndTime:         0,
		},

		{
			Id:              15,
			Name:            "上传官方NFT头像",
			Head:            "",
			Type:            "Daily",
			EventType:       "official_nft_avatar",
			Classify:        "daily",
			Desc:            "头像设置为IDChats NFT",
			Reward:          5000,
			CompletionCount: 1,
			Status:          0,
			ClaimConditions: "2",
			StartTime:       0,
			EndTime:         0,
		},
		{
			Id:              16,
			Name:            "上传任意NFT头像",
			Head:            "",
			Type:            "Daily",
			EventType:       "upload_nft_head",
			Classify:        "daily",
			Desc:            "头像设置为任意NFT",
			Reward:          2000,
			CompletionCount: 1,
			Status:          0,
			ClaimConditions: "1",
			StartTime:       0,
			EndTime:         0,
		},
		{
			Id:              17,
			Name:            "邀请上传任意NFT头像",
			Head:            "",
			Type:            "CountProgress",
			EventType:       "invite_upload_nft_head",
			Classify:        "invite",
			Desc:            "邀请头像设置为任意NFT",
			Reward:          10000,
			CompletionCount: 1,
			Status:          0,
			ClaimConditions: "13",
			StartTime:       0,
			EndTime:         0,
		},
	}
	// 创建任务
	s := &rpcTask.TaskServer{}
	_, err := s.CreateTask(context.Background(), &pbTask.CreateTaskReq{
		TaskList: taskList,
	})
	fmt.Println("CreateTaskList end ... ")
	if err != nil {
		fmt.Errorf("CreateTask() error = %v", err)
		return
	}
}

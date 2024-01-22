package base_info

type GetUserTaskListReq struct {
	OperationID string `json:"operationid"`
	Classify    string `json:"classify"`
}
type Task struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Head            string `json:"head"`
	Type            string `json:"type"`
	EventType       string `json:"eventType"`
	Classify        string `json:"classify"`
	Desc            string `json:"desc"`
	Reward          int    `json:"reward"`
	CompletionCount int    `json:"completionCount"`
	Status          int    `json:"status"`
	ClaimConditions string `json:"claimConditions"`
	StartTime       int64  `json:"startTime"`
	EndTime         int64  `json:"endTime"`
}
type UserTask struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	TaskId    int64  `json:"taskId"`
	Status    int8   `json:"status"`
	Progress  int32  `json:"progress"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Task      *Task  `json:"task"`
}

type GetUserTaskListResp struct {
	CommResp
	Data []*UserTask
}

type ClaimTaskRewardsReq struct {
	OperationID string `json:"operationid"`
	TaskId      int    `json:"taskId"`
}
type ClaimTaskRewardsResp struct {
	CommResp
}

type DailyCheckInReq struct {
	OperationID string `json:"operationid"`
}
type DailyCheckInResp struct {
	CommResp
}

type DailyIsCheckInReq struct {
	OperationID string `json:"operationid"`
}
type DailyIsCheckInResp struct {
	CommResp
}

type CheckIsHaveNftRecvIDReq struct {
	OperationID string `json:"operationid"`
}
type CheckIsHaveNftRecvIDResp struct {
	OperationID string `json:"operationid"`
}
type ObtainUnmetConditionsResp struct {
	CommResp
	Data bool `json:"data"`
}

type CheckIsHaveOfficialNftRecvIDReq struct {
	OperationID string `json:"operationid"`
}
type CheckIsHaveOfficialNftRecvIDResp struct {
	CommResp
	Data bool `json:"data"`
}
type CheckIsFollowSystemTwitterReq struct {
	OperationID string `json:"operationid"`
}
type CheckIsFollowSystemTwitterResp struct {
	CommResp
	Data bool `json:"data"`
}

// type CreateTask struct {
// 	Id              string `json:"id"` // 任务id
// 	Name            string `json:"name"`
// 	Head            string `json:"head"`
// 	Type            string `json:"type"`
// 	Classify        string `json:"classify"`
// 	Desc            string `json:"desc"`
// 	Reward          int    `json:"reward"`
// 	EventType       string `json:"eventType"` // 事件类型
// 	CompletionCount int    `json:"completionCount"`
// 	ClaimConditions string `json:"claimConditions"`
// 	Status          int8   `json:"status"`
// 	StartTime       int64  `json:"startTime"`
// 	EndTime         int64  `json:"endTime"`
// }

type CreateTaskReq struct {
	OperationID string  `json:"operationID"`
	TaskList    []*Task `json:"taskList"`
}

type CreateTaskResp struct {
	CommResp
}

type GetTaskListReq struct {
	OperationID string `json:"operationid"`
	Classify    string `json:"classify"`
}
type GetTaskListResp struct {
	CommResp
	Data []*Task
}

type GetUserClaimTaskListReq struct {
	OperationID string `json:"operationid"`
	Status      int    `json:"status"`
}
type GetUserClaimTaskListResp struct {
	CommResp
	Data []*UserTask
}

type ClaimTaskReq struct {
	OperationID string `json:"operationid"`
	TaskId      int    `json:"taskId"`
}
type ClaimTaskResp struct {
	CommResp
}


package base_info

import "Open_IM/pkg/common/db"

type PostStartGameReq struct {
	OperationID string  `json:"operationID"`
	GameID      string  `json:"gameID"`
	OperatorID  string  `json:"operatorID"`
	UserID      string  `json:"userID"`
	Status      int32   `json:"status"` //1为注册开始游戏 2 为注册结束游戏
	StartTime   int64   `json:"startTime"`
	EndTime     int64   `json:"endTime"`
	Score       float64 `json:"score"`
}

type PostStartGameResp struct {
	CommResp
	StartTime int64 `json:"data"`
}
type PostGameRankListReq struct {
	OperationID string
	GameID      string
}

type PostGameRankListResp struct {
	CommResp
	RankLinkInfoInfo RankLinkInfo `json:"data"`
}
type RankLinkInfo struct {
	UserScore     []*UserGameScore `json:"userScore"`
	UserSelfScore *UserGameScore   `json:"userSelfScore"`
}
type UserGameScore struct {
	UserID             string `json:"userID"`
	Nickname           string `json:"nickname"`
	FaceURL            string `json:"faceURL"`
	Reward             int64  `json:"reward"`
	Score              int64  `json:"score"`
	RankIndex          int32  `json:"rankIndex"`
	TokenContractChain string `json:"tokenContractChain"`
}

type PostGameListResp struct {
	CommResp
	GameConfigList []*ApiGameConfig `json:"data"`
}
type ApiGameConfig struct {
	db.GameConfig
	NextTime int64 `json:"nextTime"`
}

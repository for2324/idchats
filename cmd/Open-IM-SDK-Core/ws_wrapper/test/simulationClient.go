package test

import (
	"encoding/json"
	"net/url"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/ws_wrapper/test/client"
	"open_im_sdk/ws_wrapper/ws_local_server"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var totalConnNum int
var lock sync.Mutex
var TotalSendMsgNum int

func StartSimulationJSClient(token, api, jssdkURL, userID string) {
	// 模拟登录 认证 ws连接初始化
	user := client.NewIMClient(token, userID, api, jssdkURL, 5)
	// var err error
	// user.Token, err = user.GetToken()
	// if err != nil {
	// 	log.NewError("", "generate token failed", userID, api, err.Error())
	// 	user.Token, err = user.GetToken()
	// 	if err != nil {
	// 		log.NewError("", "generate token failed", userID, api, err.Error())
	// 		return
	// 	}
	// }
	v := url.Values{}
	v.Set("sendID", userID)
	v.Set("token", user.Token)
	v.Set("platformID", utils.IntToString(5))
	c, _, err := websocket.DefaultDialer.Dial(jssdkURL+"?"+v.Encode(), nil)
	if err != nil {
		log.NewInfo("", "dial:", err.Error(), "userID", userID)
		return
	}
	lock.Lock()
	totalConnNum += 1
	log.NewInfo("", "connect success", userID, "total conn num", totalConnNum)
	lock.Unlock()
	user.Conn = c
	// user.WsLogout()
	user.WsLogin()
	time.Sleep(time.Second * 2)

	// 模拟登录同步
	// go user.GetSelfUserInfo()
	// go user.GetAllConversationList()
	// go user.GetBlackList()
	// go user.GetFriendList()
	// go user.GetRecvFriendApplicationList()
	// go user.GetRecvGroupApplicationList()
	// go user.GetSendFriendApplicationList()
	// go user.GetSendGroupApplicationList()

	go func() {
		time.Sleep(time.Second * 5)
		user.GetFollowFriendApplicationList()
	}()

	// 模拟监听回调
	for {
		resp := ws_local_server.EventData{}
		_, message, err := c.ReadMessage()
		if err != nil {
			log.NewError("", "read:", err, "error an connet failed", userID)
			return
		}
		// log.Printf("recv: %s", message)
		_ = json.Unmarshal(message, &resp)
		log.NewInfo(resp.Event, resp.Data)
		// if resp.Event == "CreateTextMessage" {
		// 	msg := sdk_struct.MsgStruct{}
		// 	msg.InitStruct()
		// 	_ = json.Unmarshal([]byte(resp.Data), &msg)
		// 	type Data struct {
		// 		RecvID          string `json:"recvID"`
		// 		GroupID         string `json:"groupID"`
		// 		OfflinePushInfo string `json:"offlinePushInfo"`
		// 		Message         string `json:"message"`
		// 	}
		// 	offlinePushBytes, _ := json.Marshal(server_api_params.OfflinePushInfo{Title: "push offline"})
		// 	messageBytes, _ := json.Marshal(msg)
		// 	data := Data{RecvID: userID, OfflinePushInfo: string(offlinePushBytes), Message: string(messageBytes)}
		// 	err = user.SendMsg(userID, data)
		// 	//fmt.Println(msg)
		// 	lock.Lock()
		// 	TotalSendMsgNum += 1
		// 	lock.Unlock()
		// }
	}

	// 模拟给随机用户发消息
	// go func() {
	// 	for {
	// 		err = user.CreateTextMessage(userID)
	// 		if err != nil {
	// 			log.NewError("", err, i, userID)
	// 		}
	// 		time.Sleep(time.Second * 2)
	// 	}
	// }()

	// // 模拟获取登陆状态
	// go func() {
	// 	for {
	// 		if err = user.GetLoginStatus(); err != nil {
	// 			log.NewError("", err, i, userID)
	// 		}
	// 		time.Sleep(time.Second * 10)
	// 	}
	// }()
}

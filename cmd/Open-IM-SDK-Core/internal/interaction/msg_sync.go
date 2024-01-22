package interaction

import (
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

type UserInterface interface {
	UserScoreInfoUpdatedNotification(msg *server_api_params.MsgData, operationID string)
	UserAnnouncementNotification(msg *server_api_params.MsgData, operationID string)
	UserOrderPayNotification(msg *server_api_params.MsgData, operationID string)
}

type SeqPair struct {
	BeginSeq uint32
	EndSeq   uint32
}

type MsgSync struct {
	*db.DataBase
	*Ws
	LoginUserID        string
	conversationCh     chan common.Cmd2Value
	PushMsgAndMaxSeqCh chan common.Cmd2Value

	selfMsgSync *SelfMsgSync
	//selfMsgSyncLatestModel *SelfMsgSyncLatestModel
	//superGroupMsgSync *SuperGroupMsgSync
	isSyncFinished            bool
	readDiffusionGroupMsgSync *ReadDiffusionGroupMsgSync
	User                      UserInterface
}

func (m *MsgSync) compareSeq() {
	operationID := utils.OperationIDGenerator()
	m.selfMsgSync.compareSeq(operationID)
	m.readDiffusionGroupMsgSync.compareSeq(operationID)
}

func (m *MsgSync) doMaxSeq(cmd common.Cmd2Value) {
	operationID := cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync).OperationID
	if !m.isSyncFinished {
		m.readDiffusionGroupMsgSync.TriggerCmdNewMsgCome(nil, operationID, constant.MsgSyncBegin)
	}
	m.readDiffusionGroupMsgSync.doMaxSeq(cmd)
	m.selfMsgSync.doMaxSeq(cmd)
	if !m.isSyncFinished {
		m.readDiffusionGroupMsgSync.TriggerCmdNewMsgCome(nil, operationID, constant.MsgSyncEnd)
	}
	m.isSyncFinished = true
}

func (m *MsgSync) doPushMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	switch msg.SessionType {
	case constant.SuperGroupChatType:
		m.readDiffusionGroupMsgSync.doPushMsg(cmd)
	case constant.NotificationOnlinePushType:
		m.HandlerOnlinePush(cmd)
	default:
		m.selfMsgSync.doPushMsg(cmd)
	}
}

func (m *MsgSync) Work(cmd common.Cmd2Value) {
	switch cmd.Cmd {
	case constant.CmdPushMsg:
		m.doPushMsg(cmd)
	case constant.CmdMaxSeq:
		m.doMaxSeq(cmd)
	default:
		log.Error("", "cmd failed ", cmd.Cmd)
	}
}

func (m *MsgSync) GetCh() chan common.Cmd2Value {
	return m.PushMsgAndMaxSeqCh
}

func NewMsgSync(dataBase *db.DataBase, ws *Ws, user UserInterface, loginUserID string, ch chan common.Cmd2Value, pushMsgAndMaxSeqCh chan common.Cmd2Value, joinedSuperGroupCh chan common.Cmd2Value) *MsgSync {
	p := &MsgSync{
		Ws:                 ws,
		conversationCh:     ch,
		User:               user,
		DataBase:           dataBase,
		LoginUserID:        loginUserID,
		PushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh,
	}
	//	p.superGroupMsgSync = NewSuperGroupMsgSync(dataBase, ws, loginUserID, ch, joinedSuperGroupCh)
	p.selfMsgSync = NewSelfMsgSync(dataBase, ws, loginUserID, ch)
	p.readDiffusionGroupMsgSync = NewReadDiffusionGroupMsgSync(dataBase, ws, loginUserID, ch, joinedSuperGroupCh)
	//	p.selfMsgSync = NewSelfMsgSyncLatestModel(dataBase, ws, loginUserID, ch)
	p.compareSeq()
	go common.DoListener(p)
	return p
}

func (m *MsgSync) HandlerOnlinePush(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	operationID := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).OperationID
	log.Debug(operationID, utils.GetSelfFuncName(), " args ", " msg ", msg)
	if msg.ContentType == constant.UserScoreInfoUpdatedNotification {
		m.User.UserScoreInfoUpdatedNotification(msg, operationID)
	} else if msg.ContentType == constant.UserAnnouncement {
		m.User.UserAnnouncementNotification(msg, operationID)
	} else if msg.ContentType == constant.UserOrderPayNotification {
		m.User.UserOrderPayNotification(msg, operationID)
	}
}

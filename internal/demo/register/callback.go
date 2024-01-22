package register

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbGroup "Open_IM/pkg/proto/group"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

func CallBackOnUserLogin(userId string) {
	err := AutoCreateCommunity(userId)
	if err != nil {
		log.NewInfo("CallBackOnUserLogin AutoCreateCommunity", err.Error())
	}
}

func AutoCreateCommunity(userID string) error {
	_, err := imdb.GetGroupInfoByGroupID(userID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	req := pbGroup.CreateCommunityReq{
		OperationID: "CallBackOnUserLogin AutoCreateCommunity",
		OpUserID:    userID,
		OwnerUserID: userID,
		GroupInfo: &sdk_ws.GroupInfo{
			GroupID:     userID,
			GroupName:   userID,
			OwnerUserID: userID,
			GroupType:   constant.WorkingGroup,
			FaceURL:     "",
		},
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema,
		strings.Join(config.Config.Etcd.EtcdAddr, ","),
		config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		return errors.New("etcdConn is nil")
	}
	client := pbGroup.NewGroupClient(etcdConn)
	RpcResp, err := client.CreateCommunity(context.Background(), &req)
	if err != nil {
		return err
	}
	if RpcResp.ErrCode != 0 {
		return errors.New(RpcResp.ErrMsg)
	}
	return nil
}

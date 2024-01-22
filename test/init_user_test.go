package test

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"fmt"
	"testing"
	"time"
)

func TestInitUser(t *testing.T) {
	fmt.Println(config.Config.InitUser.UserId)
	for _, userId := range config.Config.InitUser.UserId {

		_, err := im_mysql_model.GetUserByUserID(userId)
		if err != nil {
		} else {
			continue
		}
		var initUser db.User
		initUser.UserID = userId
		initUser.Nickname = userId
		initUser.AppMangerLevel = constant.AppAdmin
		initUser.CreateTime = time.Now()
		initUser.Birth = time.Now()
		err = im_mysql_model.UserRegister(initUser)
		if err != nil {
			fmt.Println("InitUser insert error ", err.Error(), initUser)
		} else {
			fmt.Println("InitUser insert ", initUser)
		}
	}
}

func TestRpcCreateUserGroup(t *testing.T) {
}

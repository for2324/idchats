package ens

import (
	api "Open_IM/pkg/base_info"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"encoding/json"
	"testing"
)

func TestGetMyAppointmentList(t *testing.T) {
	var ensList []api.AppointmentUserInfo
	{
		ensList, _ = imdb.MyAppointmentList("0x56d9003b84762b0c8c9d703874f8f3ef", 0, 10)
		for i := range ensList {
			ensList[i].CreateTime = ensList[i].CreatedAt.Unix()
		}
		bytes, _ := json.Marshal(ensList)
		log.NewInfo("MyAppointmentList", string(bytes))

	}
	{
		ensList, _ = imdb.AppointmentList(0, 10)
		for i := range ensList {
			ensList[i].CreateTime = ensList[i].CreatedAt.Unix()
		}
		bytes, _ := json.Marshal(ensList)
		log.NewInfo("AppointmentList", string(bytes))
	}
}

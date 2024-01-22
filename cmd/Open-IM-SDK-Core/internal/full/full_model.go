package full

import (
	"fmt"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (u *Full) GetGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error) {
	fmt.Println("aaaaaaaaaa", utils.RunFuncName(), "GetGroupInfoByGroupID")
	g1, err := u.SuperGroup.GetGroupInfoFromLocal2Svr(groupID)
	if err == nil {
		return g1, nil
	}
	g2, err := u.group.GetGroupInfoFromLocal2Svr(groupID)
	return g2, err
}

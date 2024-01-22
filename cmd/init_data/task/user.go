package task

import (
	"Open_IM/internal/rpc/web3pub"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	pbWeb3 "Open_IM/pkg/proto/web3pub"
	utils2 "Open_IM/pkg/utils"

	"context"
	"encoding/json"
	"strings"

	"fmt"
)

func InitUserHeadInfo() {
	fmt.Println("InitUserInfo start ... ")

	users, err := imdb.GetTokenAllUser()
	if err != nil {
		fmt.Println("GetAllUser err ... ", err)
		return
	}
	fmt.Println("get user Count", len(users))
	for i, user := range users {
		fmt.Println("start resolve", i, "[", len(users), "]", user.UserID)
		// 完成绑定官方头像的任务
		if checkIsHaveGuanFangNftRecvID(user.UserID, user.TokenId, "0x23AC6898a81b4a1144d26c6f1b580B72d33a860d", "5") {
			fmt.Println("checkIsHaveNftRecvID true ... ", user.UserID)
			if err := imdb.FinishOfficialNFTHeadTask(user.UserID); err != nil {
				fmt.Println("FinishOfficialNFTHeadTask fail ... ", err)
			}
		} else {
			fmt.Println("checkIsHaveNftRecvID false ... ", user.UserID, user.TokenContractChain, user.TokenId)
		}
		fmt.Println("end resolve", i, user.UserID)
	}
	fmt.Println("InitUserInfo end ... ")
}

func InitTaskCheckIsFollowTwitter() {
	list, err := imdb.GetBindTwitterUserList()
	rpcServer := web3pub.NewWeb3PubServer(0)
	if err == nil {
		fmt.Println("GetBindTwitterUserList", len(list))
		for _, user := range list {
			if err := imdb.FinishBindTwitterTask(user.UserId); err == nil {
				// 帮助邀请人领取邀请关联任务
				if yaoqingren, err := imdb.GetRegisterInfo(user.UserId); err == nil && yaoqingren.InvitationCode != "" {
					if err := imdb.FinishInviteBindTwitterTask(yaoqingren.InvitationCode, user.UserId); err != nil {
						fmt.Println("FinishInviteBindTwitterTask err", err, yaoqingren.InvitationCode)
					} else {
						fmt.Println("FinishInviteBindTwitterTask success", yaoqingren.InvitationCode)
					}
				}
			}
			// 完成关注官方推特任务
			resp, err := rpcServer.CheckIsFollowSystemTwitter(
				context.Background(),
				&pbWeb3.CheckUserIsFollowSystemTwitterReq{
					UserId: user.UserId,
				},
			)
			if err == nil && resp.CommonResp.ErrCode == 0 {
				fmt.Println("CheckIsFollowSystemTwitter success ... ", user.UserId)
			}
		}
	}
}

func checkIsHaveGuanFangNftRecvID(userId string, tokenId string, contractAddress string, ChainID string) bool {
	if tokenId != "" {
		PostCheckData, _ := json.Marshal(&api.RequestTokenIdReq{
			ChainID:         ChainID,
			TokenID:         tokenId,
			ContractAddress: contractAddress,
		})
		resultByte, err := utils2.HttpPost(config.Config.EnsPostCheck.Url+"/graph/tokenOwnerAddress",
			"", map[string]string{"Content-Type": "application/json", "chainId": "1"}, PostCheckData)
		if err == nil {
			var resultData api.RequestTokenIdResp
			json.Unmarshal(resultByte, &resultData)
			if strings.EqualFold(resultData.TokenOwnerAddress, userId) {
				return true
			}
		}
	}
	return false
}

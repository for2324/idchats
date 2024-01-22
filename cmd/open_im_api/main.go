package main

import (
	_ "Open_IM/cmd/open_im_api/docs"
	"Open_IM/internal/api/announcement"
	apiAuth "Open_IM/internal/api/auth"
	"Open_IM/internal/api/brc20"
	"Open_IM/internal/api/business"
	clientInit "Open_IM/internal/api/client_init"
	"Open_IM/internal/api/conversation"
	"Open_IM/internal/api/ens"
	"Open_IM/internal/api/friend"
	"Open_IM/internal/api/game"
	"Open_IM/internal/api/group"
	"Open_IM/internal/api/manage"
	apiChat "Open_IM/internal/api/msg"
	"Open_IM/internal/api/nft"
	"Open_IM/internal/api/office"
	"Open_IM/internal/api/order"
	"Open_IM/internal/api/organization"
	"Open_IM/internal/api/robot"
	"Open_IM/internal/api/task"
	apiThird "Open_IM/internal/api/third"
	"Open_IM/internal/api/user"
	"Open_IM/internal/api/userthird"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/juju/ratelimit"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	//"syscall"
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
)

// @title			biubiuid
// @version		1.0
// @description	biubiuid 的API服务器文档, 文档中所有请求都有一个operationID字段用于链路追踪
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath		/
func main() {
	log.NewPrivateLog(constant.LogFileName)
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("../logs/api.log")
	gin.DefaultWriter = io.MultiWriter(f)
	//	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(utils.CorsHandler())
	log.Info("load config: ", config.Config)
	if !config.Config.IsPublicEnv {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	if config.Config.Prometheus.Enable {
		promePkg.NewApiRequestCounter()
		promePkg.NewApiRequestFailedCounter()
		promePkg.NewApiRequestSuccessCounter()
		r.Use(promePkg.PromeTheusMiddleware)
		r.GET("/metrics", promePkg.PrometheusHandler())
	}
	thirdVerify := r.Group("/thirdSns") //第三方社交平台
	{
		thirdVerify.POST("/postSign", userthird.ThirdBindPostFunc)
		thirdVerify.POST("/verifyingSign", userthird.VerifyThirdPlatformSign)
	}
	businessService := r.Group("/business")
	{
		businessService.POST("/get_business_list", business.GetBusinessList)
		businessService.POST("/add_business_api_key", business.AddBusinessApiKey)
		businessService.POST("/update_business_api_key", business.UpdateBusinessApiKey)
		businessService.POST("/get_user_trade_score", business.GetUserBusinessTrade)
		businessService.POST("/claim_self_reward", business.ClaimSelfReward)
		businessService.POST("/get_total_volume_sub_personal", business.GetSubPersonalTotalVolume)
	}

	robotGroup := r.Group("/robot") //第三方社交平台
	{
		robotGroup.POST("/create_robot", robot.CreateRobot)
		robotGroup.POST("/delegate_call", robot.DelegateCall)
		robotGroup.POST("/get_robot", robot.GetRobot)
		robotGroup.POST("/change_ord_status", robot.ChangeOrderStatus)
		robotGroup.POST("/token_price", robot.TokenPrice)
		robotGroup.POST("/export_wallet", robot.ExportsWallet)

		robotGroup.POST("/v2/create_robot", robot.CreateRobotV2)
		robotGroup.POST("/v2/delegate_call", robot.DelegateCall)
		//robotGroup.POST("/v2/delegate_call", robot.DelegateCallV2)
		robotGroup.POST("/v2/get_robot", robot.GetRobot)
		robotGroup.POST("/v2/token_price", robot.TokenPrice)
		//检查被过期
		robotGroup.POST("/v2/check_expire_key", robot.CheckExpireKey)
		//重新签名延长日期
		robotGroup.POST("/v2/reloadKey", robot.ReLoadKey)
		robotGroup.POST("/v2/export_wallet", robot.ExportsWalletV2)
		robotGroup.POST("/v2/total_volume", robot.TradingVolume)
		robotGroup.POST("/v2/total_7days_volume", robot.TradingVolume7Days)
		robotGroup.POST("/v2/import_wallet", robot.ImportWalletV2)
	}

	// user routing group, which handles user registration and login services
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/update_user_info", user.UpdateUserInfo) //1
		userRouterGroup.POST("/update_user_head", user.UpdateUserHead) //1
		userRouterGroup.POST("/set_global_msg_recv_opt", user.SetGlobalRecvMessageOpt)
		userRouterGroup.POST("/set_show_balance", user.SetShowBalance)
		userRouterGroup.POST("/get_users_info", user.GetUsersPublicInfo)            //1
		userRouterGroup.POST("/get_self_user_info", user.GetSelfUserInfo)           //1
		userRouterGroup.POST("/get_users_online_status", user.GetUsersOnlineStatus) //1
		userRouterGroup.POST("/get_users_info_from_cache", user.GetUsersInfoFromCache)
		userRouterGroup.POST("/get_user_friend_from_cache", user.GetFriendIDListFromCache)
		userRouterGroup.POST("/get_user_followeachother_friend_from_cache", user.GetEachOtherFriendIdListFromCache)
		userRouterGroup.POST("/get_black_list_from_cache", user.GetBlackIDListFromCache)
		userRouterGroup.POST("/get_all_users_uid", manage.GetAllUsersUid) //1
		userRouterGroup.POST("/account_check", manage.AccountCheck)       //1
		userRouterGroup.POST("/get_users", user.GetUsers)
		userRouterGroup.POST("/get_third_status", user.GetThirdInfo)
		userRouterGroup.POST("/sendSms", user.BindUserTelephoneCode)
		userRouterGroup.POST("/updateTelephoneInfo", user.BindUserInfoTelephoneInfo)
		userRouterGroup.POST("/change_old_phone_number", user.ChangeOldPhoneNumber)
		userRouterGroup.POST("/delThirdPlatform", user.DeletePlatformInfo)
		userRouterGroup.POST("/show_third_platform", user.ShowPlatfomrInfo)
		//更新个人简介
		userRouterGroup.POST("/update_user_sign", user.UpdateUserSign)
		userRouterGroup.POST("/get_user_sign", user.GetUserSign)
		userRouterGroup.POST("/ask_user_question", user.AskUserQuestion)
		userRouterGroup.POST("/get_user_chat_token", user.GetUserChatToken)
		userRouterGroup.POST("/operator_user_chat_token", user.OperatorUserChatToken)
		userRouterGroup.POST("/user_chat_token_history", user.ChatTokenHistory)
		userRouterGroup.POST("/transfer_token_to_group", user.TransferTokenToGroup)
		userRouterGroup.POST("/share_user_profile", user.ShowUserProfile) //1
		userRouterGroup.POST("/bind_ens_domain", user.BindEnsDomain)      //绑定ens域名
		userRouterGroup.POST("/bind_web_domain", user.BindDnsDomain)      //绑定dns域名
		//userRouterGroup.POST("/get_email_verify_code", user.GetEmailVerifyCode) //申请邮箱验证码
		userRouterGroup.POST("/bind_user_email", user.BindUserEmail) //绑定dns域名
		userRouterGroup.POST("/bind_show_nft", user.BindShowNft)     //绑定展示的的nft
		userRouterGroup.POST("/get_show_nft", user.GetShowNft)       //获取绑定的nft

		userRouterGroup.POST("/get_user_setting_page", user.GetUserSettingPage)
		userRouterGroup.POST("/update_user_setting_page", user.UpdateUserSettingPage)
		userRouterGroup.POST("/check_dns_domain", user.CheckDnsDomain)

		userRouterGroup.POST("/get_user_score", user.GetUserScore)                  //获取用户积分
		userRouterGroup.POST("/get_reward_event_logs", user.GetRewardEventLogs)     //获取奖励事件列表
		userRouterGroup.POST("/withdraw_score", user.WithdrawScore)                 // 提现积分
		userRouterGroup.POST("/get_withdraw_score_logs", user.GetWithdrawScoreLogs) // 获取提现记录
	}
	nftRouterGroup := r.Group("/nft")
	{
		nftRouterGroup.POST("/get_show_balance", nft.GetShowBalance)        //获取绑定的nft
		nftRouterGroup.POST("/get_show_btc_balance", nft.GetBtcShowBalance) //获取绑定的nft
		nftRouterGroup.POST("/like_unlike_nft_id", nft.UpdateLikeShowNft)   //点赞 或者取消点赞
		nftRouterGroup.POST("/get_like_status_nft_id", nft.GetLikeShowNft)  //获取该nft是否点赞过 以及他的数量个数
	}
	announcedRouterGroup := r.Group("/announce")
	{
		announcedRouterGroup.POST("/set_global_announce", announcement.SetGlobalAnnouncementMessageOpt)
		announcedRouterGroup.POST("/createorupdate_user_announcement_draft", announcement.CreateOrUpdateUserAnnouncementDraft) //创建草稿
		announcedRouterGroup.POST("/delete_user_announcement_draft", announcement.DeleteUserAnnouncementDraft)                 //删除用户草稿信息内容
		announcedRouterGroup.POST("/get_user_announcement_draft", announcement.GetUserAnnouncementDraft)                       //获取用户草稿列表
		announcedRouterGroup.POST("/publish_user_announcement", announcement.PublishUserAnnouncement)                          //发布用户的草稿（推送信息）
		//announcedRouterGroup.DELETE("/delete_global_user_announcement_view", user.DeleteGlobalUserAnnouncementView)            //全网推送删除已经不要的数据
		announcedRouterGroup.POST("/delete_publish_user_announcement", announcement.DeletePublishUserAnnouncement)                                  //删除用户已经发送的推送
		announcedRouterGroup.POST("/like_unlike_announcement_article", announcement.LikeUnLikeAnnouncement)                                         //点赞 或者取消点赞
		announcedRouterGroup.POST("/get_publish_space_announcement_list", announcement.GetSpacePublishAnnouncementView)                             //获取空间公告列表
		announcedRouterGroup.POST("/get_publish_space_announcement_ido_list", announcement.GetSpacePublishAnnouncementViewWithIdo)                  //获取空间公告列表
		announcedRouterGroup.POST("/get_personal_publish_space_announcement_ido_list", announcement.GetSpacePersonalPublishAnnouncementViewWithIdo) //获取个人的推荐页面的文章

		announcedRouterGroup.POST("/pin_publish_space_announcement_ido_list", announcement.PinSpacePublishAnnouncementViewWithIdo) //置顶某个文章
		announcedRouterGroup.POST("/del_publish_space_announcement_ido_list", announcement.DelSpacePublishAnnouncementViewWithIdo) //删除空间公告列表

	}
	taskRouterGroup := r.Group("/task")
	{
		taskRouterGroup.POST("/get_task_list", task.GetTaskList)

		taskRouterGroup.POST("/get_user_claim_task_list", task.GetUserClaimTaskList)
		taskRouterGroup.POST("/get_user_task_list", task.GetUserTaskList)

		taskRouterGroup.POST("/claim_task_rewards", task.ClaimTaskRewards)

		taskRouterGroup.POST("/daily_check_in", task.DailyCheckIn)
		taskRouterGroup.POST("/daily_is_check_in", task.DailyIsCheckIn)
		taskRouterGroup.POST("/check_is_have_nft_recvid", task.CheckIsHaveNftRecvID)
		taskRouterGroup.POST("/check_is_have_official_nft_recvid", task.CheckIsHaveOfficialNftRecvID)
		taskRouterGroup.POST("/check_is_follow_system_twitter", task.CheckIsFollowSystemTwitter)

	}
	gameRouterGroup := r.Group("/gameapi")
	{
		gameRouterGroup.POST("/start_game", game.PostStartGame)
		gameRouterGroup.POST("/game_rank_list", game.GetGameRankList)
		gameRouterGroup.POST("/game_list", game.GetGameList)
	}
	//friend routing group
	friendRouterGroup := r.Group("/friend")
	{
		//	friendRouterGroup.POST("/get_friends_info", friend.GetFriendsInfo)
		friendRouterGroup.POST("/add_friend", friend.AddFriend) //1
		friendRouterGroup.POST("/follow_add_friend", friend.FollowAddFriend)
		friendRouterGroup.POST("/following_list", friend.FollowFriendList)
		friendRouterGroup.POST("/delete_friend", friend.DeleteFriend)                        //1
		friendRouterGroup.POST("/get_friend_apply_list", friend.GetFriendApplyList)          //1
		friendRouterGroup.POST("/get_self_friend_apply_list", friend.GetSelfFriendApplyList) //1
		friendRouterGroup.POST("/get_friend_list", friend.GetFriendList)                     //1
		friendRouterGroup.POST("/add_friend_response", friend.AddFriendResponse)             //1
		friendRouterGroup.POST("/set_friend_remark", friend.SetFriendRemark)                 //1

		friendRouterGroup.POST("/add_black", friend.AddBlack)          //1
		friendRouterGroup.POST("/get_black_list", friend.GetBlacklist) //1
		friendRouterGroup.POST("/remove_black", friend.RemoveBlack)    //1

		friendRouterGroup.POST("/import_friend", friend.ImportFriend) //1
		friendRouterGroup.POST("/is_friend", friend.IsFriend)         //1
		friendRouterGroup.POST("/is_follow_eachother_friend", friend.IsFollowEachOtherFriend)

		friendRouterGroup.POST("/get_user_followed_list", friend.GetUserFollowedList)
		friendRouterGroup.POST("/get_user_following_list", friend.GetUserFollowingList)

	}
	//group related routing group
	groupRouterGroup := r.Group("/group")
	{
		groupRouterGroup.POST("/create_community", group.CreateCommunity)
		groupRouterGroup.POST("/create_community_channel", group.CreateCommunityChannel)
		groupRouterGroup.POST("/search_community", group.SearchCommunity)
		//groupRouterGroup.POST("/join_community", group.JoinCommunity)

		//创建角色标签
		groupRouterGroup.POST("/create_role_tag", group.CreateRoleTag)
		////区块同步的内容
		//groupRouterGroup.POST("/commit_role_tag", group.CommitCreateRoleTag)
		//groupRouterGroup.POST("/mint_burn_role_tag", group.MinOrBurn1155ToUser)
		//群标签列表
		groupRouterGroup.POST("/get_community_role_tag", group.GetCommunityRoleTag)
		groupRouterGroup.POST("/get_community_role_tag_detail", group.GetCommunityRoleTagDetail)
		//频道关联标签 请从上面的的设置

		groupRouterGroup.POST("/get_community_channel_list", group.CommunityChannel)
		groupRouterGroup.POST("/group_channel_status", group.CommunityChannelStatus)
		groupRouterGroup.POST("/get_user_joined_group_list", group.GetUserJoinedGroupList)
		groupRouterGroup.POST("/get_group_all_history_message_list", group.GetGroupAllHistoryMessageList)

		groupRouterGroup.POST("/create_group", group.CreateGroup)                                   //1
		groupRouterGroup.POST("/set_group_info", group.SetGroupInfo)                                //1
		groupRouterGroup.POST("/join_group", group.JoinGroup)                                       //1
		groupRouterGroup.POST("/quit_group", group.QuitGroup)                                       //1
		groupRouterGroup.POST("/group_application_response", group.ApplicationGroupResponse)        //1
		groupRouterGroup.POST("/transfer_group", group.TransferGroupOwner)                          //1
		groupRouterGroup.POST("/get_recv_group_applicationList", group.GetRecvGroupApplicationList) //1
		groupRouterGroup.POST("/get_user_req_group_applicationList", group.GetUserReqGroupApplicationList)
		groupRouterGroup.POST("/get_groups_info", group.GetGroupsInfo) //1
		groupRouterGroup.POST("/kick_group", group.KickGroupMember)    //1
		//	groupRouterGroup.POST("/get_group_member_list", group.GetGroupMemberList)        //no use
		groupRouterGroup.POST("/get_group_all_member_list", group.GetGroupAllMemberList) //1
		groupRouterGroup.POST("/get_group_members_info", group.GetGroupMembersInfo)      //1
		groupRouterGroup.POST("/invite_user_to_group", group.InviteUserToGroup)          //1
		groupRouterGroup.POST("/get_joined_group_list", group.GetJoinedGroupList)
		groupRouterGroup.POST("/dismiss_group", group.DismissGroup) //
		groupRouterGroup.POST("/mute_group_member", group.MuteGroupMember)
		groupRouterGroup.POST("/cancel_mute_group_member", group.CancelMuteGroupMember) //MuteGroup
		groupRouterGroup.POST("/mute_group", group.MuteGroup)
		groupRouterGroup.POST("/cancel_mute_group", group.CancelMuteGroup)
		groupRouterGroup.POST("/set_group_member_nickname", group.SetGroupMemberNickname)
		groupRouterGroup.POST("/set_group_member_info", group.SetGroupMemberInfo)
		groupRouterGroup.POST("/get_group_abstract_info", group.GetGroupAbstractInfo)
		groupRouterGroup.POST("/create_sys_user_group", group.CreateSysUserGroup)
		groupRouterGroup.POST("/is_can_get_member_count_reword", group.IsCanGetMemberCountReword)

		groupRouterGroup.POST("/colligate_search", RateLimitMiddleware(time.Second/10, 100), group.ColligateSearch)
		groupRouterGroup.POST("/create_push_space_articel_order", group.CreatePushSpaceArticelOrder)

	}

	spaceRouterGroup := r.Group("/api_space")
	{
		spaceRouterGroup.POST("/get_hot_space", group.GetHotSpace)
		spaceRouterGroup.POST("/get_hot_space_banner", group.GetHotSpaceBanner)
		spaceRouterGroup.POST("/get_hot_space_banner_article", announcement.GetHotSpaceBannerArticle)
	}
	brc20RouterGroup := r.Group("/brc20")
	{
		brc20RouterGroup.POST("/brc20_all_tokens", brc20.GetBrc20Tokens)
		brc20RouterGroup.POST("/brc20_personal_tokens", brc20.GetPersonalBrc20Tokens)
		brc20RouterGroup.POST("/brc20_pledge", brc20.Brc20Pledge)
		brc20RouterGroup.POST("/brc20_transferable_list", brc20.GetBrc20TransferableList)
		brc20RouterGroup.GET("/brc20_pledge_pool_infoes", brc20.GetBrc20PledgePoolInfoes)
		brc20RouterGroup.POST("/brc20_personal_pledge_infoes", brc20.Brc20PledgePersonalInfoes)
		brc20RouterGroup.POST("/brc20_search", brc20.Brc20SearchScan)

	}
	superGroupRouterGroup := r.Group("/super_group")
	{
		superGroupRouterGroup.POST("/get_joined_group_list", group.GetJoinedSuperGroupList)
		superGroupRouterGroup.POST("/get_groups_info", group.GetSuperGroupsInfo)
	}
	//certificate
	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/user_register", apiAuth.UserRegister) //1
		//authRouterGroup.POST("/user_token", apiAuth.UserToken)       //1
		authRouterGroup.POST("/parse_token", apiAuth.ParseToken)   //1
		authRouterGroup.POST("/force_logout", apiAuth.ForceLogout) //1
	}
	//Third service
	thirdGroup := r.Group("/third")
	{
		thirdGroup.POST("/tencent_cloud_storage_credential", apiThird.TencentCloudStorageCredential)
		thirdGroup.POST("/ali_oss_credential", apiThird.AliOSSCredential)
		thirdGroup.POST("/minio_storage_credential", apiThird.MinioStorageCredential)
		thirdGroup.POST("/minio_upload", apiThird.MinioUploadFile)
		thirdGroup.POST("/upload_update_app", apiThird.UploadUpdateApp)
		thirdGroup.POST("/get_download_url", apiThird.GetDownloadURL)
		thirdGroup.POST("/get_rtc_invitation_info", apiThird.GetRTCInvitationInfo)
		thirdGroup.POST("/get_rtc_invitation_start_app", apiThird.GetRTCInvitationInfoStartApp)
		thirdGroup.POST("/fcm_update_token", apiThird.FcmUpdateToken)
		thirdGroup.POST("/aws_storage_credential", apiThird.AwsStorageCredential)
		thirdGroup.POST("/set_app_badge", apiThird.SetAppBadge)
	}
	//Message
	chatGroup := r.Group("/msg")
	{
		chatGroup.POST("/newest_seq", apiChat.GetSeq)
		chatGroup.POST("/send_msg", apiChat.SendMsg)
		chatGroup.POST("/pull_msg_by_seq", apiChat.PullMsgBySeqList)
		chatGroup.POST("/del_msg", apiChat.DelMsg)
		chatGroup.POST("/del_super_group_msg", apiChat.DelSuperGroupMsg)
		chatGroup.POST("/clear_msg", apiChat.ClearMsg)
		chatGroup.POST("/manage_send_msg", manage.ManagementSendMsg)
		chatGroup.POST("/batch_send_msg", manage.ManagementBatchSendMsg)
		chatGroup.POST("/check_msg_is_send_success", manage.CheckMsgIsSendSuccess)
		chatGroup.POST("/set_msg_min_seq", apiChat.SetMsgMinSeq)
		chatGroup.POST("/get_single_chat_history_message_list", apiChat.GetSingleChatHistoryMessageList)

	}
	//Conversation
	conversationGroup := r.Group("/conversation")
	{ //1
		conversationGroup.POST("/get_all_conversations", conversation.GetAllConversations)
		conversationGroup.POST("/get_conversation", conversation.GetConversation)
		conversationGroup.POST("/get_conversations", conversation.GetConversations)
		conversationGroup.POST("/set_conversation", conversation.SetConversation)
		conversationGroup.POST("/monitor_conversation", conversation.MonitorConversation)
		conversationGroup.POST("/batch_set_conversation", conversation.BatchSetConversations)
		conversationGroup.POST("/set_recv_msg_opt", conversation.SetRecvMsgOpt)
		conversationGroup.POST("/modify_conversation_field", conversation.ModifyConversationField)
	}
	// office
	officeGroup := r.Group("/office")
	{
		officeGroup.POST("/get_user_tags", office.GetUserTags)
		officeGroup.POST("/get_user_tag_by_id", office.GetUserTagByID)
		officeGroup.POST("/create_tag", office.CreateTag)
		officeGroup.POST("/delete_tag", office.DeleteTag)
		officeGroup.POST("/set_tag", office.SetTag)
		officeGroup.POST("/send_msg_to_tag", office.SendMsg2Tag)
		officeGroup.POST("/get_send_tag_log", office.GetTagSendLogs)

		officeGroup.POST("/create_one_work_moment", office.CreateOneWorkMoment)
		officeGroup.POST("/delete_one_work_moment", office.DeleteOneWorkMoment)
		officeGroup.POST("/like_one_work_moment", office.LikeOneWorkMoment)
		officeGroup.POST("/comment_one_work_moment", office.CommentOneWorkMoment)
		officeGroup.POST("/get_work_moment_by_id", office.GetWorkMomentByID)
		officeGroup.POST("/get_user_work_moments", office.GetUserWorkMoments)
		officeGroup.POST("/get_user_friend_work_moments", office.GetUserFriendWorkMoments)
		officeGroup.POST("/set_user_work_moments_level", office.SetUserWorkMomentsLevel)
		officeGroup.POST("/delete_comment", office.DeleteComment)
	}

	organizationGroup := r.Group("/organization")
	{
		organizationGroup.POST("/create_department", organization.CreateDepartment)
		organizationGroup.POST("/update_department", organization.UpdateDepartment)
		organizationGroup.POST("/get_sub_department", organization.GetSubDepartment)
		organizationGroup.POST("/delete_department", organization.DeleteDepartment)
		organizationGroup.POST("/get_all_department", organization.GetAllDepartment)
		organizationGroup.POST("/create_organization_user", organization.CreateOrganizationUser)
		organizationGroup.POST("/update_organization_user", organization.UpdateOrganizationUser)
		organizationGroup.POST("/delete_organization_user", organization.DeleteOrganizationUser)

		organizationGroup.POST("/create_department_member", organization.CreateDepartmentMember)
		organizationGroup.POST("/get_user_in_department", organization.GetUserInDepartment)
		organizationGroup.POST("/update_user_in_department", organization.UpdateUserInDepartment)

		organizationGroup.POST("/get_department_member", organization.GetDepartmentMember)
		organizationGroup.POST("/delete_user_in_department", organization.DeleteUserInDepartment)
		organizationGroup.POST("/get_user_in_organization", organization.GetUserInOrganization)
	}

	initGroup := r.Group("/init")
	{
		initGroup.POST("/set_client_config", clientInit.SetClientInitConfig)
		initGroup.POST("/get_client_config", clientInit.GetClientInitConfig)
	}

	ensGroup := r.Group("/ens")
	{
		ensGroup.POST("/appointment", ens.Appointment)
		ensGroup.POST("/appointment_list", ens.AppointmentList)
		ensGroup.POST("/create_register_ens_order", ens.CreateRegisterEnsOrder)
		ensGroup.POST("/get_ens_order_info", ens.GetEnsOrderInfo)
		ensGroup.POST("/has_appointment", ens.HasAppointment)
	}

	orderGroup := r.Group("/order")
	{
		orderGroup.POST("/get_support_coin_list", order.GetSupportCoinList)
		orderGroup.POST("/request_payment", order.RequestPayment)
		orderGroup.GET("/transactions/order/:id", order.GetOrderInfo)
		orderGroup.GET("/transactions/out_trade_no/:orderId", order.GetOutTradeOrderInfo)
		orderGroup.POST("/test_notify", order.TestNotify)
	}

	go apiThird.MinioInit()
	defaultPorts := config.Config.Api.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10002 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	fmt.Println("start api server, address: ", address, "OpenIM version: ", constant.CurrentVersion)
	err := r.Run(address)
	if err != nil {
		log.Error("", "api run failed ", address, err.Error())
		panic("api start failed " + err.Error())
	}
}

// 获取用户请求的真实ip
func realIPMiddleware(c *gin.Context) {
	forwardedFor := c.Request.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ", ")
		c.Request.RemoteAddr = ips[0]
	}
	c.Next()
}
func rateLimitMiddleware(ctx *gin.Context, countnumber int64) {
	key := fmt.Sprintf("ratelimit:%s", ctx.ClientIP())
	limit := countnumber         // 每秒钟最多访问 10 次
	expiration := time.Second    // 过期时间为 1 秒钟
	now := time.Now().UnixNano() // 当前时间戳（纳秒）
	pipe := db.DB.RDB.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-int64(expiration)))
	pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})
	pipe.ZCard(ctx, key)
	_, err := pipe.Exec(ctx)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	cmdList, err := pipe.Exec(ctx)
	if err == nil && len(cmdList) > 2 {
		count, _ := cmdList[2].(*redis.IntCmd).Result()
		if count > limit {
			ctx.AbortWithStatus(429)
			return
		}
	}
	ctx.Next()
}

// 整个接口限流量
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit...
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}

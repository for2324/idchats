package constant

const (
	GetSelfUserInfoRouter          = "/user/get_self_user_info"
	GetUsersInfoRouter             = "/user/get_users_info"
	UpdateSelfUserInfoRouter       = "/user/update_user_info"
	SetGlobalRecvMessageOptRouter  = "/user/set_global_msg_recv_opt"
	GetUsersInfoFromCacheRouter    = "/user/get_users_info_from_cache"
	PostFollowEachOtherFriendCache = "/user/get_user_followeachother_friend_from_cache"
	GetUserScoreInfoRouter         = "/user/get_user_score"

	AddFriendRouter                    = "/friend/add_friend"
	FollowAddFriendRouter              = "/friend/follow_add_friend"
	DeleteFriendRouter                 = "/friend/delete_friend"
	GetFriendApplicationListRouter     = "/friend/get_friend_apply_list"      //recv
	GetSelfFriendApplicationListRouter = "/friend/get_self_friend_apply_list" //send
	GetFriendListRouter                = "/friend/get_friend_list"
	GetFollowingList                   = "/friend/following_list"
	AddFriendResponse                  = "/friend/add_friend_response"
	SetFriendRemark                    = "/friend/set_friend_remark"
	PostIsFollowEachOtherFriend        = "/friend/is_follow_eachother_friend"

	GetGroupAllHistoryMessageList       = "/group/get_group_all_history_message_list"
	GetSinglerChatAllHistoryMessageList = "/msg/get_single_chat_history_message_list"

	AddBlackRouter     = "/friend/add_black"
	RemoveBlackRouter  = "/friend/remove_black"
	GetBlackListRouter = "/friend/get_black_list"

	SendMsgRouter          = "/msg/send_msg"
	PullUserMsgRouter      = "/chat/pull_msg"
	PullUserMsgBySeqRouter = "/msg/pull_msg_by_seq"
	NewestSeqRouter        = "/chat/newest_seq"

	//msg
	DeleteMsgRouter           = RouterMsg + "/del_msg"
	ClearMsgRouter            = RouterMsg + "/clear_msg"
	DeleteSuperGroupMsgRouter = RouterMsg + "/del_super_group_msg"

	TencentCloudStorageCredentialRouter = "/third/tencent_cloud_storage_credential"
	AliOSSCredentialRouter              = "/third/ali_oss_credential"
	MinioStorageCredentialRouter        = "/third/minio_storage_credential"
	AwsStorageCredentialRouter          = "/third/aws_storage_credential"

	//group
	CreateGroupRouter                 = RouterGroup + "/create_group"
	SetGroupInfoRouter                = RouterGroup + "/set_group_info"
	JoinGroupRouter                   = RouterGroup + "/join_group"
	SubscribeGroupChannel             = RouterGroup + "/subscribe_group_channel"
	UnSubscribeGroupChannel           = RouterGroup + "/un_subscribe_group_channel"
	QuitGroupRouter                   = RouterGroup + "/quit_group"
	GetGroupsInfoRouter               = RouterGroup + "/get_groups_info"
	GetGroupAllMemberListRouter       = RouterGroup + "/get_group_all_member_list"
	GetGroupMembersInfoRouter         = RouterGroup + "/get_group_members_info"
	InviteUserToGroupRouter           = RouterGroup + "/invite_user_to_group"
	GetJoinedGroupListRouter          = RouterGroup + "/get_joined_group_list"
	KickGroupMemberRouter             = RouterGroup + "/kick_group"
	TransferGroupRouter               = RouterGroup + "/transfer_group"
	GetRecvGroupApplicationListRouter = RouterGroup + "/get_recv_group_applicationList"
	GetSendGroupApplicationListRouter = RouterGroup + "/get_user_req_group_applicationList"
	AcceptGroupApplicationRouter      = RouterGroup + "/group_application_response"
	RefuseGroupApplicationRouter      = RouterGroup + "/group_application_response"
	DismissGroupRouter                = RouterGroup + "/dismiss_group"
	MuteGroupMemberRouter             = RouterGroup + "/mute_group_member"
	CancelMuteGroupMemberRouter       = RouterGroup + "/cancel_mute_group_member"
	MuteGroupRouter                   = RouterGroup + "/mute_group"
	CancelMuteGroupRouter             = RouterGroup + "/cancel_mute_group"
	SetGroupMemberNicknameRouter      = RouterGroup + "/set_group_member_nickname"
	SetGroupMemberInfoRouter          = RouterGroup + "/set_group_member_info"
	GetGroupAbstractInfoRouter        = RouterGroup + "/get_group_abstract_info"
	GetCommunityChannelInfoRouter     = RouterGroup + "/get_community_channel_list"
	MonitorGroupMessage               = RouterGroup + "/monitor_group_message"
	CommunityChannelStatus            = RouterGroup + "/group_channel_status"

	SetReceiveMessageOptRouter         = "/conversation/set_receive_message_opt"
	GetReceiveMessageOptRouter         = "/conversation/get_receive_message_opt"
	GetAllConversationMessageOptRouter = "/conversation/get_all_conversation_message_opt"
	SetConversationOptRouter           = ConversationGroup + "/set_conversation"
	MonitorConversationsRouter         = ConversationGroup + "/monitor_conversation"
	GetConversationsRouter             = ConversationGroup + "/get_conversations"
	GetAllConversationsRouter          = ConversationGroup + "/get_all_conversations"
	GetConversationRouter              = ConversationGroup + "/get_conversation"
	BatchSetConversationRouter         = ConversationGroup + "/batch_set_conversation"
	ModifyConversationFieldRouter      = ConversationGroup + "/modify_conversation_field"

	//organization
	GetSubDepartmentRouter    = RouterOrganization + "/get_sub_department"
	GetDepartmentMemberRouter = RouterOrganization + "/get_department_member"
	ParseTokenRouter          = RouterAuth + "/parse_token"

	//super_group
	GetJoinedSuperGroupListRouter = RouterSuperGroup + "/get_joined_group_list"
	GetSuperGroupsInfoRouter      = RouterSuperGroup + "/get_groups_info"

	//third
	FcmUpdateTokenRouter = RouterThird + "/fcm_update_token"
	SetAppBadgeRouter    = RouterThird + "/set_app_badge"
)
const (
	RouterGroup        = "/group"
	ConversationGroup  = "/conversation"
	RouterOrganization = "/organization"
	RouterAuth         = "/auth"
	RouterSuperGroup   = "/super_group"
	RouterMsg          = "/msg"
	RouterThird        = "/third"
)

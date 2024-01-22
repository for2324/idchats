package rocksCache

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"runtime/debug"
	"sort"
	"strconv"
	"time"

	"gorm.io/gorm"
)

const (
	userInfoCache             = "USER_INFO_CACHE:"
	userInfoLinkCache         = "USER_INFO_LINK_CACHE:"
	friendRelationCache       = "FRIEND_RELATION_CACHE:"
	FollowFriendCache         = "Follow_Friend_CACHE:" //用户相互关注的列表
	blackListCache            = "BLACK_LIST_CACHE:"
	groupCache                = "GROUP_CACHE:"
	groupInfoCache            = "GROUP_INFO_CACHE:"
	groupOwnerIDCache         = "GROUP_OWNER_ID:"
	joinedGroupListCache      = "JOINED_GROUP_LIST_CACHE:"
	groupMemberInfoCache      = "GROUP_MEMBER_INFO_CACHE:"
	groupAllMemberInfoCache   = "GROUP_ALL_MEMBER_INFO_CACHE:"
	allFriendInfoCache        = "ALL_FRIEND_INFO_CACHE:"
	allDepartmentCache        = "ALL_DEPARTMENT_CACHE:"
	allDepartmentMemberCache  = "ALL_DEPARTMENT_MEMBER_CACHE:"
	joinedSuperGroupListCache = "JOINED_SUPER_GROUP_LIST_CACHE:"
	groupMemberListHashCache  = "GROUP_MEMBER_LIST_HASH_CACHE:"
	groupMemberNumCache       = "GROUP_MEMBER_NUM_CACHE:"
	conversationCache         = "CONVERSATION_CACHE:"
	conversationIDListCache   = "CONVERSATION_ID_LIST_CACHE:"
	userEmailInfo             = "UserEmail:"
	userSpaceGroupInfo        = "USER_GROUP_INFO:"
	SystemOfficialNft         = "SYSTEM_OFFICIAL_NFT"
	TxtDomainInfo             = "TXT_DOMAIN_NAME"
)

func DelKeys() {
	fmt.Println("init to del old keys")
	for _, key := range []string{groupCache, friendRelationCache, blackListCache, userInfoCache, groupInfoCache, groupOwnerIDCache, joinedGroupListCache,
		groupMemberInfoCache, groupAllMemberInfoCache, allFriendInfoCache, TxtDomainInfo} {
		fName := utils.GetSelfFuncName()
		var cursor uint64
		var n int
		for {
			var keys []string
			var err error
			keys, cursor, err = db.DB.RDB.Scan(context.Background(), cursor, key+"*", 3000).Result()
			if err != nil {
				panic(err.Error())
			}
			n += len(keys)
			// for each for redis cluster
			for _, key := range keys {
				if err = db.DB.RDB.Del(context.Background(), key).Err(); err != nil {
					log.NewError("", fName, key, err.Error())
					err = db.DB.RDB.Del(context.Background(), key).Err()
					if err != nil {
						panic(err.Error())
					}
				}
			}
			if cursor == 0 {
				break
			}
		}
	}
}
func GetFollowEachOtherFromCache(userID string) ([]string, error) {
	getFriendIDList := func() (string, error) {
		friendIDList, err := imdb.GetFollowEachOtherUserId(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(friendIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	friendIDListStr, err := db.DB.Rc.Fetch(FollowFriendCache+userID, time.Second*30*60, getFriendIDList)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.Wrap(err, "")
	}
	var friendIDList []string
	err = json.Unmarshal([]byte(friendIDListStr), &friendIDList)
	return friendIDList, utils.Wrap(err, "")
}
func DelFollowFriendIDListFromCache(userID string) error {
	err := db.DB.Rc.TagAsDeleted(FollowFriendCache + userID)
	return err
}
func GetFriendIDListFromCache(userID string) ([]string, error) {
	getFriendIDList := func() (string, error) {
		friendIDList, err := imdb.GetFriendIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(friendIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	friendIDListStr, err := db.DB.Rc.Fetch(friendRelationCache+userID, time.Second*30*60, getFriendIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var friendIDList []string
	err = json.Unmarshal([]byte(friendIDListStr), &friendIDList)
	return friendIDList, utils.Wrap(err, "")
}

func DelFriendIDListFromCache(userID string) error {
	err := db.DB.Rc.TagAsDeleted(friendRelationCache + userID)
	return err
}

func GetBlackListFromCache(userID string) ([]string, error) {
	getBlackIDList := func() (string, error) {
		blackIDList, err := imdb.GetBlackIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(blackIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	blackIDListStr, err := db.DB.Rc.Fetch(blackListCache+userID, time.Second*30*60, getBlackIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var blackIDList []string
	err = json.Unmarshal([]byte(blackIDListStr), &blackIDList)
	return blackIDList, utils.Wrap(err, "")
}

func DelBlackIDListFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(blackListCache + userID)
}

func GetJoinedGroupIDListFromCache(userID string) ([]string, error) {
	getJoinedGroupIDList := func() (string, error) {
		joinedGroupList, err := imdb.GetJoinedGroupIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(joinedGroupList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	joinedGroupIDListStr, err := db.DB.Rc.Fetch(joinedGroupListCache+userID, time.Second*30*60, getJoinedGroupIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var joinedGroupList []string
	err = json.Unmarshal([]byte(joinedGroupIDListStr), &joinedGroupList)
	return joinedGroupList, utils.Wrap(err, "")
}

func DelJoinedGroupIDListFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(joinedGroupListCache + userID)
}

func GetGroupMemberIDListFromCache(groupID string) ([]string, error) {
	f := func() (string, error) {
		groupInfo, err := GetGroupInfoFromCache(groupID)
		if err != nil {
			return "", utils.Wrap(err, "GetGroupInfoFromCache failed")
		}
		var groupMemberIDList []string
		if groupInfo.GroupType == constant.SuperGroup {
			superGroup, err := db.DB.GetSuperGroup(groupID)
			if err != nil {
				return "", utils.Wrap(err, "")
			}
			groupMemberIDList = superGroup.MemberIDList
		} else {
			groupMemberIDList, err = imdb.GetGroupMemberIDListByGroupID(groupID)
			if err != nil {
				return "", utils.Wrap(err, "")
			}
		}
		bytes, err := json.Marshal(groupMemberIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	groupIDListStr, err := db.DB.Rc.Fetch(groupCache+groupID, time.Second*30*60, f)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var groupMemberIDList []string
	err = json.Unmarshal([]byte(groupIDListStr), &groupMemberIDList)
	return groupMemberIDList, utils.Wrap(err, "")
}

func DelGroupMemberIDListFromCache(groupID string) error {
	err := db.DB.Rc.TagAsDeleted(groupCache + groupID)
	return err
}
func GetUserBaseInfoFromCache(userID string) (*db.YeWuUser, error) {
	getUserInfo := func() (string, error) {
		userInfo, err := imdb.GetUserAndThirdPath(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		//userInfo.LinkTree, err = imdb.GetUserLinkTree(userID, "")
		//if err != nil {
		//	return "", utils.Wrap(err, "")
		//}
		bytes, err := json.Marshal(userInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	userInfoStr, err := db.DB.Rc.Fetch(userInfoCache+userID, time.Second*30*60, getUserInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	userInfo := &db.YeWuUser{}
	err = json.Unmarshal([]byte(userInfoStr), userInfo)
	return userInfo, utils.Wrap(err, "")
}

func GetUserBaseInfoFromCacheUserLink(userID string) ([]*db.UserLink, error) {
	getUserInfoLink := func() (string, error) {
		userInfoLink, err := imdb.GetUserLinkTree(userID, "")
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		//userInfo.LinkTree, err = imdb.GetUserLinkTree(userID, "")
		//if err != nil {
		//	return "", utils.Wrap(err, "")
		//}
		bytes, err := json.Marshal(userInfoLink)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	userInfoStr, err := db.DB.Rc.Fetch(userInfoLinkCache+userID, time.Second*30*60, getUserInfoLink)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	userInfoLink := make([]*db.UserLink, 0)
	err = json.Unmarshal([]byte(userInfoStr), &userInfoLink)
	return userInfoLink, utils.Wrap(err, "")
}
func DeleteUserBaseInfoFromCacheUserLink(userID string) error {
	return db.DB.Rc.TagAsDeleted(userInfoLinkCache + userID)

}
func GetUserInfoFromCacheByMerLin(userID string, chainID string) (*db.User, error) {
	getUserInfo := func() (string, error) {
		userInfo, err := imdb.GetUserByUserID(userID)
		if errors.Is(gorm.ErrRecordNotFound, err) {
			if !utils.CheckEthAddress(userID) {
				return "", errors.New("无效的以太坊地址")
			}
			imdb.UserRegister(db.User{
				UserID:           userID,
				Nickname:         userID,
				FaceURL:          "",
				Gender:           0,
				PhoneNumber:      "",
				Birth:            time.Now(),
				Email:            "",
				Ex:               "",
				CreateTime:       time.Now(),
				AppMangerLevel:   1,
				GlobalRecvMsgOpt: 0,
				Chainid:          5,
			})
			userInfo, err = imdb.GetUserByUserID(userID)
		} else if err != nil {
			fmt.Println(err.Error())
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(userInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	userInfoStr, err := db.DB.Rc.Fetch(userInfoCache+userID, time.Second*30*60, getUserInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	userInfo := &db.User{}
	err = json.Unmarshal([]byte(userInfoStr), userInfo)
	resultName, _ := imdb.GetUserEnsNameByChainID(userID, chainID)
	if resultName != "" {
		userInfo.Nickname = resultName
	}
	return userInfo, utils.Wrap(err, "")
}
func DelUserInfoFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(userInfoCache + userID)

}

func GetGroupMemberInfoFromCache(groupID, userID string) (*db.GroupMember, error) {
	getGroupMemberInfo := func() (string, error) {
		groupMemberInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserIDMerlin(groupID, userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupMemberInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	groupMemberInfoStr, err := db.DB.Rc.Fetch(groupMemberInfoCache+groupID+"-"+userID, time.Second*30*60, getGroupMemberInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	groupMember := &db.GroupMember{}
	err = json.Unmarshal([]byte(groupMemberInfoStr), groupMember)
	return groupMember, utils.Wrap(err, "")
}

func DelGroupMemberInfoFromCache(groupID, userID string) error {
	return db.DB.Rc.TagAsDeleted(groupMemberInfoCache + groupID + "-" + userID)
}

func GetGroupMembersInfoFromCache(count, offset int32, groupID string) ([]*db.GroupMember, error) {
	groupMemberIDList, err := GetGroupMemberIDListFromCache(groupID)
	if err != nil {
		return nil, err
	}
	if count < 0 || offset < 0 {
		return nil, nil
	}
	var groupMemberList []*db.GroupMember
	var start, stop int32
	start = offset
	stop = offset + count
	l := int32(len(groupMemberIDList))
	if start > stop {
		return nil, nil
	}
	if start >= l {
		return nil, nil
	}
	if count != 0 {
		if stop >= l {
			stop = l
		}
		groupMemberIDList = groupMemberIDList[start:stop]
	} else {
		if l < 1000 {
			stop = l
		} else {
			stop = 1000
		}
		groupMemberIDList = groupMemberIDList[start:stop]
	}
	//log.NewDebug("", utils.GetSelfFuncName(), "ID list: ", groupMemberIDList)
	for _, userID := range groupMemberIDList {
		groupMembers, err := GetGroupMemberInfoFromCache(groupID, userID)
		if err != nil {
			log.NewError("", utils.GetSelfFuncName(), err.Error(), groupID, userID)
			continue
		}
		groupMemberList = append(groupMemberList, groupMembers)
	}
	return groupMemberList, nil
}

func GetAllGroupMembersInfoFromCache(groupID string) ([]*db.GroupMember, error) {
	getGroupMemberInfo := func() (string, error) {
		groupMembers, err := imdb.GetGroupMemberListByGroupIDMerlin(groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupMembers)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	groupMembersStr, err := db.DB.Rc.Fetch(groupAllMemberInfoCache+groupID, time.Second*30*60, getGroupMemberInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var groupMembers []*db.GroupMember
	err = json.Unmarshal([]byte(groupMembersStr), &groupMembers)
	return groupMembers, utils.Wrap(err, "")
}

func DelAllGroupMembersInfoFromCache(groupID string) error {
	return db.DB.Rc.TagAsDeleted(groupAllMemberInfoCache + groupID)
}

func GetGroupInfoFromCache(groupID string) (*db.Group, error) {
	getGroupInfo := func() (string, error) {
		groupInfo, err := imdb.GetGroupInfoByGroupID(groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	groupInfoStr, err := db.DB.Rc.Fetch(groupInfoCache+groupID, time.Second*30*60, getGroupInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	groupInfo := &db.Group{}
	err = json.Unmarshal([]byte(groupInfoStr), groupInfo)
	return groupInfo, utils.Wrap(err, "")
}

func DelGroupInfoFromCache(groupID string) error {
	return db.DB.Rc.TagAsDeleted(groupInfoCache + groupID)
}

func GetAllFriendsInfoFromCache(userID string) ([]*db.Friend, error) {
	getAllFriendInfo := func() (string, error) {
		friendInfoList, err := imdb.GetFriendListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(friendInfoList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	allFriendInfoStr, err := db.DB.Rc.Fetch(allFriendInfoCache+userID, time.Second*30*60, getAllFriendInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var friendInfoList []*db.Friend
	err = json.Unmarshal([]byte(allFriendInfoStr), &friendInfoList)
	return friendInfoList, utils.Wrap(err, "")
}

func DelAllFriendsInfoFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(allFriendInfoCache + userID)
}

func GetAllDepartmentsFromCache() ([]db.Department, error) {
	getAllDepartments := func() (string, error) {
		departmentList, err := imdb.GetSubDepartmentList("-1")
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(departmentList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	allDepartmentsStr, err := db.DB.Rc.Fetch(allDepartmentCache, time.Second*30*60, getAllDepartments)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var allDepartments []db.Department
	err = json.Unmarshal([]byte(allDepartmentsStr), &allDepartments)
	return allDepartments, utils.Wrap(err, "")
}

func DelAllDepartmentsFromCache() error {
	return db.DB.Rc.TagAsDeleted(allDepartmentCache)
}

func GetAllDepartmentMembersFromCache() ([]db.DepartmentMember, error) {
	getAllDepartmentMembers := func() (string, error) {
		departmentMembers, err := imdb.GetDepartmentMemberList("-1")
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(departmentMembers)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	allDepartmentMembersStr, err := db.DB.Rc.Fetch(allDepartmentMemberCache, time.Second*30*60, getAllDepartmentMembers)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var allDepartmentMembers []db.DepartmentMember
	err = json.Unmarshal([]byte(allDepartmentMembersStr), &allDepartmentMembers)
	return allDepartmentMembers, utils.Wrap(err, "")
}

func DelAllDepartmentMembersFromCache() error {
	return db.DB.Rc.TagAsDeleted(allDepartmentMemberCache)
}

func GetJoinedSuperGroupListFromCache(userID string) ([]string, error) {
	getJoinedSuperGroupIDList := func() (string, error) {
		userToSuperGroup, err := db.DB.GetSuperGroupByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		if len(userToSuperGroup.GroupIDList) == 0 {
			return "", errors.New("GroupIDList == 0")
		}
		bytes, err := json.Marshal(userToSuperGroup.GroupIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	joinedSuperGroupListStr, err := db.DB.Rc.Fetch(joinedSuperGroupListCache+userID,
		time.Second*30*60, getJoinedSuperGroupIDList)
	if err != nil {
		return nil, err
	}
	var joinedSuperGroupList []string
	err = json.Unmarshal([]byte(joinedSuperGroupListStr), &joinedSuperGroupList)
	return joinedSuperGroupList, utils.Wrap(err, "")
}

func DelJoinedSuperGroupIDListFromCache(userID string) error {
	fmt.Println(string(debug.Stack()))
	log.NewDebug("", utils.GetSelfFuncName(), string(debug.Stack()))
	err := db.DB.Rc.TagAsDeleted(joinedSuperGroupListCache + userID)
	return err
}

func GetGroupMemberListHashFromCache(groupID string) (uint64, error) {
	generateHash := func() (string, error) {
		groupInfo, err := GetGroupInfoFromCache(groupID)
		if err != nil {
			return "0", utils.Wrap(err, "GetGroupInfoFromCache failed")
		}
		if groupInfo.Status == constant.GroupStatusDismissed {
			return "0", nil
		}
		groupMemberIDList, err := GetGroupMemberIDListFromCache(groupID)
		if err != nil {
			return "", utils.Wrap(err, "GetGroupMemberIDListFromCache failed")
		}
		sort.Strings(groupMemberIDList)
		var all string
		for _, v := range groupMemberIDList {
			all += v
		}
		bi := big.NewInt(0)
		bi.SetString(utils.Md5(all)[0:8], 16)
		return strconv.Itoa(int(bi.Uint64())), nil
	}
	hashCode, err := db.DB.Rc.Fetch(groupMemberListHashCache+groupID, time.Second*30*60, generateHash)
	if err != nil {
		return 0, utils.Wrap(err, "fetch failed")
	}
	hashCodeUint64, err := strconv.Atoi(hashCode)
	return uint64(hashCodeUint64), err
}

func DelGroupMemberListHashFromCache(groupID string) error {
	err := db.DB.Rc.TagAsDeleted(groupMemberListHashCache + groupID)
	return err
}

func GetGroupMemberNumFromCache(groupID string) (int64, error) {
	getGroupMemberNum := func() (string, error) {
		num, err := imdb.GetGroupMemberNumByGroupID(groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return strconv.Itoa(int(num)), nil
	}
	groupMember, err := db.DB.Rc.Fetch(groupMemberNumCache+groupID, time.Second*30*60, getGroupMemberNum)
	if err != nil {
		return 0, utils.Wrap(err, "")
	}
	num, err := strconv.Atoi(groupMember)
	return int64(num), err
}

func DelGroupMemberNumFromCache(groupID string) error {
	return db.DB.Rc.TagAsDeleted(groupMemberNumCache + groupID)
}

func GetUserConversationIDListFromCache(userID string) ([]string, error) {
	getConversationIDList := func() (string, error) {
		conversationIDList, err := imdb.GetConversationIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "getConversationIDList failed")
		}
		log.NewDebug("", utils.GetSelfFuncName(), conversationIDList)
		bytes, err := json.Marshal(conversationIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	conversationIDListStr, err := db.DB.Rc.Fetch(conversationIDListCache+userID, time.Second*30*60, getConversationIDList)
	var conversationIDList []string
	err = json.Unmarshal([]byte(conversationIDListStr), &conversationIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return conversationIDList, nil
}

func DelUserConversationIDListFromCache(userID string) error {
	return utils.Wrap(db.DB.Rc.TagAsDeleted(conversationIDListCache+userID), "DelUserConversationIDListFromCache err")
}

func GetConversationFromCache(ownerUserID, conversationID string) (*db.Conversation, error) {
	getConversation := func() (string, error) {
		conversation, err := imdb.GetConversation(ownerUserID, conversationID)
		if err != nil {
			return "", utils.Wrap(err, "get failed")
		}
		bytes, err := json.Marshal(conversation)
		if err != nil {
			return "", utils.Wrap(err, "Marshal failed")
		}
		return string(bytes), nil
	}
	conversationStr, err := db.DB.Rc.Fetch(conversationCache+ownerUserID+":"+conversationID, time.Second*30*60, getConversation)
	if err != nil {
		return nil, utils.Wrap(err, "Fetch failed")
	}
	conversation := db.Conversation{}
	err = json.Unmarshal([]byte(conversationStr), &conversation)
	if err != nil {
		return nil, utils.Wrap(err, "Unmarshal failed")
	}
	return &conversation, nil
}

func GetConversationsFromCache(ownerUserID string, conversationIDList []string) ([]db.Conversation, error) {
	var conversationList []db.Conversation
	for _, conversationID := range conversationIDList {
		conversation, err := GetConversationFromCache(ownerUserID, conversationID)
		if err != nil {
			return nil, utils.Wrap(err, "GetConversationFromCache failed")
		}
		conversationList = append(conversationList, *conversation)
	}
	return conversationList, nil
}

func GetUserAllConversationList(ownerUserID string) ([]db.Conversation, error) {
	IDList, err := GetUserConversationIDListFromCache(ownerUserID)
	if err != nil {
		return nil, err
	}
	var conversationList []db.Conversation
	log.NewDebug("", utils.GetSelfFuncName(), IDList)
	for _, conversationID := range IDList {
		conversation, err := GetConversationFromCache(ownerUserID, conversationID)
		if err != nil {
			return nil, utils.Wrap(err, "GetConversationFromCache failed")
		}
		conversationList = append(conversationList, *conversation)
	}
	return conversationList, nil
}

func DelConversationFromCache(ownerUserID, conversationID string) error {
	return utils.Wrap(db.DB.Rc.TagAsDeleted(conversationCache+ownerUserID+":"+conversationID), "DelConversationFromCache err")
}

func GetEmailUserInfo(emailID string) (*db.EmailUserSystem, error) {
	fmt.Println(emailID)
	getEmailINfo := func() (string, error) {
		emailInfo, err := imdb.GetEmailInfo(emailID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		fmt.Println(emailInfo)
		bytes, err := json.Marshal(emailInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		fmt.Println(string(bytes))
		return string(bytes), nil
	}
	fmt.Println(userEmailInfo + emailID)
	groupInfoStr, err := db.DB.Rc.Fetch(userEmailInfo+emailID, time.Second*30*60, getEmailINfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	emailInfo := &db.EmailUserSystem{}
	err = json.Unmarshal([]byte(groupInfoStr), emailInfo)
	return emailInfo, utils.Wrap(err, "")
}
func DelEmailUserInfo(emailID string) error {
	return utils.Wrap(db.DB.Rc.TagAsDeleted(userEmailInfo+emailID), "DelEmailUserInfo err")
}
func GetSpaceInfoByUser(userid string) (*db.Group, error) {
	getGroupInfo := func() (string, error) {
		groupInfo, err := imdb.GetOneGroupInfoByUserID(userid)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}

	groupInfoStr, err := db.DB.Rc.Fetch(userSpaceGroupInfo+userid, time.Minute*5, getGroupInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	groupInfo := &db.Group{}
	err = json.Unmarshal([]byte(groupInfoStr), groupInfo)
	return groupInfo, utils.Wrap(err, "")
}
func DelSpaceGroupInfoFromCache(userid string) error {
	return db.DB.Rc.TagAsDeleted(userSpaceGroupInfo + userid)
}
func GetOfficialNftFromCache() (result []*db.SystemNft, err error) {
	getSystemNftInfo := func() (string, error) {
		systemNftInfo, err := imdb.GetSystemOfficialNftInfo()
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(systemNftInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		fmt.Println(string(bytes))
		return string(bytes), nil
	}
	systemNftInfoStr, err := db.DB.Rc.Fetch(SystemOfficialNft, time.Second*30*60, getSystemNftInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	err = json.Unmarshal([]byte(systemNftInfoStr), &result)
	return result, utils.Wrap(err, "")
}
func GetOfficialNftContractFromCache() (resultStr []string, err error) {
	getSystemNftInfo := func() (string, error) {
		systemNftInfo, err := imdb.GetSystemOfficialNftInfo()
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(systemNftInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		fmt.Println(string(bytes))
		return string(bytes), nil
	}
	systemNftInfoStr, err := db.DB.Rc.Fetch(SystemOfficialNft, time.Second*30*60, getSystemNftInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var result []*db.SystemNft
	err = json.Unmarshal([]byte(systemNftInfoStr), &result)
	if err != nil {
		return nil, err
	}
	for _, value := range result {
		resultStr = append(resultStr, value.SystemNftContract)
	}
	return resultStr, nil
}

type UserBiBotFeeRate struct {
	TradeFeeRate  float64
	SniperFeeRate float64
}

func DeleteUserBiBotFeeRate(merchantUid, merchantId string) {
	db.DB.Rc.TagAsDeleted("user_id:bibot:" + merchantUid)
}
func GetUserBiBotFeeRate(merchantUid, merchantId string) (resultStr *UserBiBotFeeRate, err error) {
	BibotFeeRate := func() (string, error) {
		tradeFeeRate, SniperFeeRate := GetCurrentFee(merchantUid, merchantId, nil)
		tempData := &UserBiBotFeeRate{
			TradeFeeRate:  tradeFeeRate,
			SniperFeeRate: SniperFeeRate,
		}
		bytes, err := json.Marshal(tempData)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		fmt.Println(string(bytes))
		return string(bytes), nil
	}
	strBibotFeeRate, err := db.DB.Rc.Fetch("user_id:bibot:"+merchantUid, time.Second*30*60, BibotFeeRate)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	result := new(UserBiBotFeeRate)
	err = json.Unmarshal([]byte(strBibotFeeRate), &result)
	return result, utils.Wrap(err, "")
}

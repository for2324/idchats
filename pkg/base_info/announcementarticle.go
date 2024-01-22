package base_info

import (
	"Open_IM/pkg/common/db"
)

type AnnouncementElem struct {
	AnnouncementUrl string `json:"announcementUrl" bind:"require"` //在ipfs的公告的信息 url
	//文字、文字+视频、文字+图片、文字+图片+链接、文字+链接、图片、链接（单独一条外部链接）
	Title       string   `json:"title"`
	Text        string   `json:"text"` //文字内容
	VideoPath   []string `json:"videoPath"`
	PicturePath []string `json:"picturePath"`
	UrlPath     []string `json:"urlPath"` //是否可以对url鉴权
}

// CreateAnnouncementArticleReq 推送相关的api
type CreateAnnouncementArticleReq struct {
	OperationID      string           `json:"operatorID" bind:"require" ` //创建string
	IsGlobal         int32            `json:"isGlobal"`                   //是否全局推送
	ArticleID        *uint64          `json:"articleID"`                  //如果我要转发某个文章的情况 只需要拿到文章id
	ArticleType      *string          `json:"articleType"`                //如果我转发某个文章的情况 只需要拿到文章类型ido 或者announce
	AnnouncementMsg  AnnouncementElem `json:"announcementMsg"`            //文章内容
	AnnouncementElem AnnouncementElem `json:"announcementElem"`           //文章内容
	TxnType          string           `json:"txnType"`                    //支付方式
}
type DeleteAnnouncementArticleReq struct {
	OperationID string `json:"operatorID" bind:"require" ` //创建string
	ArticleID   string `json:"articleID" bind:"require" `
}
type GetAnnouncementArticleReq struct {
	OperationID string `json:"operatorID" bind:"require" ` //创建string
	PageIndex   int32  `json:"pageIndex"`
	PageSize    int32  `json:"pageSize"`
	ArticleID   string `json:"articleID"`
	IsGlobal    int32  `json:"isGlobal"`
}

type GetAnnouncementArticleWithIdoReq struct {
	GetAnnouncementArticleReq
	ArticleType int32 //
}
type OperatorSpaceArticleList struct {
	OperationID string `json:"operatorID" bind:"require" ` //创建string
	ID          int64  `json:"id"`
	GroupID     string `json:"groupID"`
	IsPin       int32  `json:"isPin"`
}

type GetAnnouncementArticleWithIdoResp struct {
	CommResp
	ApiSpaceArticleListData *ApiSpaceArticleList `json:"data"`
}
type ApiSpaceArticleList struct {
	Data        []interface{} `json:"data"`
	TotalCount  int64         `json:"totalCount"`
	CurrentPage int64         `json:"currentPage"`
}

type AnnouncementArticleWithGroupInfo struct {
	db.AnnouncementArticle
	GroupName string `json:"groupName"`
	FaceURL   string `json:"faceURL"`
	IsLikes   int32  `json:"isLikes"`
	IsRead    int32  `json:"isRead"`
}
type GetAnnouncementArticleResp struct {
	CommResp
	Data []*AnnouncementArticleWithGroupInfo `json:"data"`
}

// CreateAnnouncementArticleResp   返回值
type CreateAnnouncementArticleResp struct {
	CommResp
}

// GetAnnouncementArticleGroupReq 查询某个群内的所有消息 以用户token来请求，
type GetAnnouncementArticleGroupReq struct {
	OperatorID string `json:"operatorID"`
	GroupID    string `json:"groupID"`
}
type GetAnnouncementArticleGroupResp struct {
	CommResp
	//AnnouncementArticleList []*db.AnnouncementArticle `json:"data"`
}
type GetAnnouncementArticleGlobalReq struct {
	OperatorID string `json:"operatorID"`
}
type GetAnnouncementArticleGlobalResp struct {
	CommResp
	//AnnouncementArticleList []*db.AnnouncementArticle `json:"data"`
}
type DeleteFromAnnouncementArticle struct {
	OperatorID string `json:"operatorID"`
	ArticleID  string `json:"articleID"` //删除的文章的表示
}
type DeleteFromAnnouncementArticleResp struct {
	CommResp
}

// 添加 删除，修改
type AnnouncementArticleDraftReq struct {
	db.AnnouncementArticleDraft
	OperationID string `json:"operationID"`
}
type AnnouncementArticleDraftResp struct {
	CommResp
	Data []*db.AnnouncementArticleDraft `json:"data"`
}


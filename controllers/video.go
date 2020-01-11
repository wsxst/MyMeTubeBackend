package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"MeTube/models"

	"MeTube/consts"
)

func ViewVideo(c *gin.Context) {
	requestInfo := models.ViewRequest{}
	err := c.BindJSON(&requestInfo)
	if err != nil {
		fmt.Println("Parse body failed. Error:", err)
		c.JSON(http.StatusNotImplemented, gin.H{
			"status": consts.ParamErrCode,
			"msg":    "params 格式错误",
		})
	}
	models.AddVideoView(requestInfo.VideoId, requestInfo.UserId)
	err = models.AddViewCountByVideoId(requestInfo.VideoId)
	if err != nil {
		fmt.Println("Update view_count failed. Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": consts.UpdateErrCode,
			"msg":    "更新记录失败！",
		})
	}
	video := models.GetVideoByID(requestInfo.VideoId)
	if video == nil || video.Url == "" || video.Owner <= 0{
		c.JSON(http.StatusInsufficientStorage, gin.H{
			"code": consts.DataNotFoundCode,
			"msg":    "Can't find the video!",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": video,
	})
}

func GetVideoComments(c *gin.Context) {
	type UserReply struct {
		Reply     *models.MUserReply
		SendUser  *models.MUser
		RecvUser  *models.MUser
		Like      bool   `json:"like"`
		LikeCount int64  `json:"like_count"`
	}
	type VideoComment struct {
		Comment     *models.MVideoComment
		User        *models.MUser
		UserReplies []*UserReply
		Like        bool    `json:"like"`
		LikeCount   int64   `json:"like_count"`
		ReplyCount  int64     `json:"reply_count"`
	}
	requestInfo := models.GetCommentRequest{}
	err :=c.BindJSON(&requestInfo)
	if err != nil {
		fmt.Println("Parse body failed. Error:", err)
		c.JSON(http.StatusNotImplemented, gin.H{
			"status": consts.ParamErrCode,
			"msg":    "params 格式错误",
		})
	}
	var comments []*models.MVideoComment
	videoId := requestInfo.VideoId
	total := models.GetVideoByID(videoId).CommentCount
	if requestInfo.SortType == 0 {
		comments = models.GetCommentsByVideoIdSortedByHotCount(videoId, requestInfo.Start, requestInfo.Limit)
	} else if requestInfo.SortType == 1 {
		comments = models.GetCommentsByVideoIdSortedByTime(videoId, requestInfo.Start, requestInfo.Limit)
	} else {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": consts.ParamErrCode,
			"msg":    "Sort Type error!",
		})
	}
	if comments == nil {
		c.JSON(http.StatusInsufficientStorage, gin.H{
			"code": consts.DataNotFoundCode,
			"msg":    "Can't find comments of the video!",
		})
	}
	var videoComments []*VideoComment
	for _, comment := range comments {
		replies := models.GetOnePageRepliesByCommentId(comment.Id, requestInfo.ReplyStart, requestInfo.ReplyFoldLimit)
		fmt.Println(replies)
		var userReplies []*UserReply
		for _, reply := range replies {
			userReply := &UserReply{}
			userReply.Reply = reply
			userReply.SendUser = models.GetUserByID(reply.SendUserId)
			userReply.RecvUser = models.GetUserByID(reply.RecvUserId)
			userReply.LikeCount = reply.LikeCount
			userReply.Like = models.JudgeLikeReply(requestInfo.UserId, reply.Id)
			userReplies = append(userReplies, userReply)
		}
		user := models.GetUserByID(comment.UserId)
		videoComments = append(videoComments, &VideoComment{
			Comment:comment,
			User:user,
			UserReplies:userReplies,
			Like:models.JudgeLikeComment(requestInfo.UserId, comment.Id),
			LikeCount:comment.LikeCount,
			ReplyCount: models.GetCommentByCommentId(comment.Id).ReplyCount,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"comments": videoComments,
			"total": total,
		},
	})
}

func CommentAVideo(c *gin.Context)  {
	requestInfo := models.CommentRequest{}
	err := c.BindJSON(&requestInfo)
	if err != nil {
		fmt.Println("Parse body failed. Error:", err)
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": consts.ParamErrCode,
			"msg":    "params 格式错误",
		})
	}
	content := requestInfo.Content
	videoId := requestInfo.VideoId
	userId := requestInfo.UserId
	commentId, err := models.AddAComment(userId, videoId, content)
	if err != nil || commentId <=0 {
		fmt.Println("Comment failed. Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": consts.InsertErrCode,
			"msg":    "评论失败...原因是：插入记录失败.",
		})
	}
	fmt.Println(commentId)
	err = models.AddCommentCountByVideoId(videoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": consts.InsertErrCode,
			"msg":    "评论失败...原因是：更新记录失败.",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "评论成功！",
	})
}

func ReplyAUser(c *gin.Context)  {
	requestInfo := models.ReplyRequest{}
	err := c.BindJSON(&requestInfo)
	if err != nil {
		fmt.Println("Parse body failed. Error:", err)
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": consts.ParamErrCode,
			"msg":    "params 格式错误",
		})
	}
	content := requestInfo.Content
	commentId := requestInfo.CommentId
	sendUserId := requestInfo.SendUserId
	recvUserId := requestInfo.RecvUserId
	level := requestInfo.Level
	replyId, err := models.AddAReply(commentId, sendUserId, recvUserId, content, level)
	if err != nil || replyId <=0 {
		fmt.Println("Reply failed. Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": consts.InsertErrCode,
			"msg":    "回复失败...原因是：插入记录失败.",
		})
	}
	fmt.Println(replyId)
	err = models.AddReplyCountByCommentId(commentId)
	if err != nil {
		fmt.Println("Reply failed. Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": consts.InsertErrCode,
			"msg":    "回复失败...原因是：更新记录失败.",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "评论成功！",
	})
}

func ReplyLike(c *gin.Context)  {
	requestInfo := models.MReplyLike{}
	err := c.BindJSON(&requestInfo)
	if err != nil {
		fmt.Println("Parse body failed. Error:", err)
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": consts.ParamErrCode,
			"msg":    "params 格式错误",
		})
	}
	if requestInfo.Flag {
		err := models.AddLikeByReplyId(requestInfo.ReplyId, requestInfo.UserId)
		if err != nil {
			fmt.Println("Like reply failed. Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": consts.UpdateErrCode,
				"msg":    "更新记录失败！",
			})
		}
		err = models.AddLikeCountByReplyId(requestInfo.ReplyId)
		if err != nil {
			fmt.Println("Like reply failed. Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": consts.UpdateErrCode,
				"msg":    "更新记录失败！",
			})
		}
	} else {
		err := models.CancelLikeByReplyId(requestInfo.ReplyId, requestInfo.UserId)
		if err != nil {
			fmt.Println("Like reply failed. Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": consts.UpdateErrCode,
				"msg":    "更新记录失败！",
			})
		}
		err = models.MinusLikeCountByReplyId(requestInfo.ReplyId)
		if err != nil {
			fmt.Println("Like reply failed. Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": consts.UpdateErrCode,
				"msg":    "更新记录失败！",
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "点赞评论成功！",
	})
}

func CommentLike(c *gin.Context)  {
	requestInfo := models.MCommentLike{}
	err := c.BindJSON(&requestInfo)
	if err != nil {
		fmt.Println("Parse body failed. Error:", err)
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": consts.ParamErrCode,
			"msg":    "params 格式错误",
		})
	}
	if requestInfo.Flag {
		err := models.AddLikeByCommentId(requestInfo.CommentId, requestInfo.UserId)
		if err != nil {
			fmt.Println("Like comment failed. Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": consts.UpdateErrCode,
				"msg":    "更新记录失败！",
			})
		}
		err = models.AddLikeCountByCommentId(requestInfo.CommentId)
		if err != nil {
			fmt.Println("Like comment failed. Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": consts.UpdateErrCode,
				"msg":    "更新记录失败！",
			})
		}
	} else {
		err := models.CancelLikeByCommentId(requestInfo.CommentId, requestInfo.UserId)
		if err != nil {
			fmt.Println("Like comment failed. Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": consts.UpdateErrCode,
				"msg":    "更新记录失败！",
			})
		}
		err = models.MinusLikeCountByCommentId(requestInfo.CommentId)
		if err != nil {
			fmt.Println("Like comment failed. Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": consts.UpdateErrCode,
				"msg":    "更新记录失败！",
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "点赞评论成功！",
	})
}

func GetVideo(c *gin.Context)  {
	videoid := c.Param("videoid")
	videoId, err := strconv.ParseInt(videoid, 10, 64)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": consts.ParamErrCode,
			"msg":    "params 格式错误",
		})
	}
	video := models.GetVideoByID(videoId)
	if video == nil || video.Url == "" || video.Owner <= 0{
		c.JSON(http.StatusInsufficientStorage, gin.H{
			"code": consts.DataNotFoundCode,
			"msg":    "Can't find the video!",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": video,
	})
}

func RemoveComment(c *gin.Context)  {
	type RComment struct {
		Id int64 `json:"id"`
	}
	rComment := RComment{}
	err := c.BindJSON(&rComment)
	if err != nil {
		fmt.Println("Parse body failed. Error:", err)
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": consts.ParamErrCode,
			"msg":    "params 格式错误",
		})
	}
	comment := models.GetCommentByCommentId(rComment.Id)
	if comment == nil || comment.Id <= 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": consts.QueryErrCode,
			"msg": "没找到对应的评论！",
		})
	}
	err = models.MinusCommentCountByVideoId(comment.VideoId)
	if err != nil {
		fmt.Println("Update comment count failed. Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": consts.UpdateErrCode,
			"msg":    "更新记录失败！",
		})
	}
	err = models.RemoveAComment(rComment.Id)
	if err != nil {
		fmt.Println("Remove comment failed. Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": consts.DeleteErrCode,
			"msg":    "删除记录失败！",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "删除评论成功！",
	})
}

func RemoveReply(c *gin.Context)  {
	type RReply struct {
		Id int64 `json:"id"`
	}
	rReply := RReply{}
	err := c.BindJSON(&rReply)
	if err != nil {
		fmt.Println("Parse body failed. Error:", err)
		c.JSON(http.StatusNotImplemented, gin.H{
			"code": consts.ParamErrCode,
			"msg":    "params 格式错误",
		})
	}
	reply := models.GetReplyByReplyId(rReply.Id)
	if reply == nil || reply.Id <= 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": consts.QueryErrCode,
			"msg": "没找到对应的评论！",
		})
	}
	err = models.MinusReplyCountByCommentId(reply.CommentId)
	if err != nil {
		fmt.Println("Update reply count failed. Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": consts.UpdateErrCode,
			"msg":    "更新记录失败！",
		})
	}
	err = models.RemoveAReply(rReply.Id)
	if err != nil {
		fmt.Println("Remove reply failed. Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": consts.DeleteErrCode,
			"msg":    "删除记录失败！",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "删除评论成功！",
	})
}

func GetReplies(c *gin.Context)  {
	type UserReply struct {
		Reply     *models.MUserReply
		SendUser  *models.MUser
		RecvUser  *models.MUser
		Like      bool   `json:"like"`
		LikeCount int64  `json:"like_count"`
	}
	requestInfo := models.GetRepliesRequest{}
	err :=c.BindJSON(&requestInfo)
	if err != nil {
		fmt.Println("Parse body failed. Error:", err)
		c.JSON(http.StatusNotImplemented, gin.H{
			"status": consts.ParamErrCode,
			"msg":    "params 格式错误",
		})
	}
	replies := models.GetOnePageRepliesByCommentId(requestInfo.CommentId, requestInfo.Start, requestInfo.Limit)
	fmt.Println(replies)
	var userReplies []*UserReply
	for _, reply := range replies {
		userReply := &UserReply{}
		userReply.Reply = reply
		userReply.SendUser = models.GetUserByID(reply.SendUserId)
		userReply.RecvUser = models.GetUserByID(reply.RecvUserId)
		userReply.LikeCount = reply.LikeCount
		userReply.Like = models.JudgeLikeReply(requestInfo.UserId, reply.Id)
		userReplies = append(userReplies, userReply)
	}
	total := models.GetCommentByCommentId(requestInfo.CommentId).ReplyCount
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"UserReplies": userReplies,
			"total": total,
		},
	})
}

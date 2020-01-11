package models

type CommentRequest struct {
	VideoId int64 `json:"video_id"`
	UserId  int64 `json:"user_id"`
	Content string `json:"content"`
}

type ReplyRequest struct {
	CommentId  int64  `json:"comment_id"`
	SendUserId int64  `json:"send_user_id"`
	RecvUserId int64  `json:"recv_user_id"`
	Content    string `json:"content"`
	Level      int64  `json:"level"`
}

type GetCommentRequest struct {
	UserId         int64 `json:"user_id"`
	Start          int64 `json:"start"`
	Limit          int64 `json:"limit"`
	VideoId        int64 `json:"video_id"`
	SortType       int64 `json:"sort_type"`
	ReplyStart     int64 `json:"reply_start"`
	ReplyFoldLimit int64 `json:"reply_fold_limit"`
}

type ReplyLikeRequest struct {
	ReplyId   int64 `json:"reply_id"`
	UserId    int64 `json:"user_id"`
	Flag      int8  `json:"flag"`
}

type CommentLikeRequest struct {
	CommentId  int64 `json:"comment_id"`
	UserId     int64 `json:"user_id"`
	Flag       int8 `json:"flag"`
}

type ViewRequest struct {
	VideoId   int64 `json:"video_id"`
	UserId    int64 `json:"user_id"`
}

type GetRepliesRequest struct {
	UserId     int64 `json:"user_id"`
	Start      int64 `json:"start"`
	Limit      int64 `json:"limit"`
	CommentId  int64 `json:"comment_id"`
}

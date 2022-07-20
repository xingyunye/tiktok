package handlers

import (
	"net/http"
	"time"

	"douyin/v1/kitex_gen/user"
	"douyin/v1/pkg/errno"

	"github.com/gin-gonic/gin"
)

//定义回复消息中json的格式

type BaseResponse struct {
	Code    int32  `json:"status_code"`
	Message string `json:"status_msg"`
}

type Response struct {
	Code    int32       `json:"status_code"`
	Message string      `json:"status_msg"`
	Data    interface{} `json:"data"`
}

type LoginResponse struct {
	Code     int32  `json:"status_code"`
	Message  string `json:"status_msg"`
	UserId   int64  `json:"user_id"`
	UserName string `json:"name"`
	Token    string `json:"token"`
	Expire   string `json:"expire"`
}

type UserInfoResponse struct {
	Code    int32       `json:"status_code"`
	Message string      `json:"status_msg"`
	User    interface{} `json:"user"`
}

type UserListResponse struct {
	Code     int32       `json:"status_code"`
	Message  string      `json:"status_msg"`
	UserList interface{} `json:"user_list"`
}

type UserParam struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

// SendResponse pack response
func SendResponse(c *gin.Context, err error, data interface{}) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, Response{
		Code:    Err.ErrCode,
		Message: Err.ErrMsg,
		Data:    data,
	})
}

// SendLoginResponse pack response
func SendLoginResponse(c *gin.Context, err error, userId int64, userName string, token string, expire time.Time) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, LoginResponse{
		Code:     Err.ErrCode,
		Message:  Err.ErrMsg,
		UserId:   userId,
		UserName: userName,
		Token:    token,
		Expire:   expire.Format(time.RFC3339),
	})
}

// SendUserInfoResponse pack response
func SendUserInfoResponse(c *gin.Context, err error, user *user.User) {
	Err := errno.ConvertErr(err)
	users := map[string]interface{}{"user": user}
	c.JSON(http.StatusOK, UserInfoResponse{
		Code:    Err.ErrCode,
		Message: Err.ErrMsg,
		User:    users["user"],
	})
}

// SendUserListResponse
func SendUserListResponse(c *gin.Context, err error, userList []*user.User) {
	Err := errno.ConvertErr(err)
	users := map[string]interface{}{"user_list": userList}
	c.JSON(http.StatusOK, UserListResponse{
		Code:     Err.ErrCode,
		Message:  Err.ErrMsg,
		UserList: users["user_list"],
	})
}

// video ----------------------------------------------------------
type CreateVideoResponse struct {
	Code    int32  `json:"status_code"`
	Message string `json:"status_msg"`
}

type QueryByVideoList struct {
	Code    int32       `json:"status_code"`
	Message string      `json:"status_msg"`
	Data    interface{} `json:"video_list"`
}

type QueryByLastTimeResponse struct {
	Code     int32       `json:"status_code"`
	Message  string      `json:"status_msg"`
	NextTime int64       `json:"next_time"`
	Data     interface{} `json:"video_list"`
}

type PostCommentResponse struct {
	Code    int32       `json:"status_code"`
	Message string      `json:"status_msg"`
	Data    interface{} `json:"comment"`
}

type QueryCommentResponse struct {
	Code    int32       `json:"status_code"`
	Message string      `json:"status_msg"`
	Data    interface{} `json:"comment_list"`
}

func SendQueryByVideoList(c *gin.Context, err error, data interface{}) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, QueryByVideoList{
		Code:    Err.ErrCode,
		Message: Err.ErrMsg,
		Data:    data,
	})
}

func SendQueryByLastTimeResponse(c *gin.Context, err error, data interface{}, nextTime int64) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, QueryByLastTimeResponse{
		Code:     Err.ErrCode,
		Message:  Err.ErrMsg,
		NextTime: nextTime,
		Data:     data,
	})
}

func SendCreateVideoResponse(c *gin.Context, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, CreateVideoResponse{
		Code:    Err.ErrCode,
		Message: Err.ErrMsg,
	})
}

func SendPostCommentResponse(c *gin.Context, data interface{}, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, PostCommentResponse{
		Code:    Err.ErrCode,
		Message: Err.ErrMsg,
		Data:    data,
	})
}

func SendQueryCommentResponse(c *gin.Context, data interface{}, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, QueryCommentResponse{
		Code:    Err.ErrCode,
		Message: Err.ErrMsg,
		Data:    data,
	})
}

package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// JSONData 基础答复结构
type JSONData struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

// NoDataResponse 没有数据
func (j *JSONData) NoDataResponse(c *gin.Context) {
	j.Code = http.StatusNoContent
	if j.Message == "" {
		j.Message = "没有数据"
	}
	c.JSON(http.StatusOK, j)
	c.Abort()
}

// BadRequest 参数错误
func (j *JSONData) BadRequest(c *gin.Context) {
	j.Code = http.StatusBadRequest
	if j.Message == "" {
		j.Message = "参数错误"
	}
	c.JSON(http.StatusOK, j)
	c.Abort()
}

// ResponseOK 请求成功
func (j *JSONData) ResponseOK(c *gin.Context) {
	j.Code = http.StatusOK
	if j.Message == "" {
		j.Message = "操作成功"
	}
	c.JSON(http.StatusOK, j)
	c.Abort()
}

// ServerError 服务器错误
func (j *JSONData) ServerError(c *gin.Context, err error) {
	j.Code = http.StatusInternalServerError

	if err != nil {
		j.Message = err.Error()
	}

	if j.Message == "" {
		j.Message = "服务出了一点小问题"
	}
	c.JSON(http.StatusOK, j)
	c.Abort()
}

// Unauthorized 缺失认证
func (j *JSONData) Unauthorized(c *gin.Context) {
	if j.Message == "" {
		j.Message = "请先登录"
	}
	j.Code = http.StatusUnauthorized
	c.JSON(http.StatusOK, j)
	c.Abort()
}

// FailAuth 认证失败
func (j *JSONData) FailAuth(c *gin.Context) {
	j.Code = http.StatusForbidden
	c.JSON(http.StatusOK, j)
	c.Abort()
}

// NotAcceptable 无效操作
func (j *JSONData) NotAcceptable(c *gin.Context) {
	j.Code = http.StatusNotAcceptable
	if j.Message == "" {
		j.Message = "无效操作"
	}
	c.JSON(http.StatusOK, j)
	c.Abort()
}

// TimeOut 访问超时
func (j *JSONData) TimeOut(c *gin.Context) {
	j.Code = http.StatusRequestTimeout
	c.JSON(http.StatusOK, j)
	c.Abort()
}

// 自定义响应结构

// 分页结构
type pagination struct {
	Current  int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	Total    int `json:"total" form:"total"`
}

// 用户列表
type userListResponseRes struct {
	ID       string `json:"id" form:"id" xorm:"id"`
	Retired  int    `json:"retired" form:"retired" xorm:"retired"`
	Avast    string `json:"avast" form:"avast" xorm:"avast"`
	Bio      string `json:"bio" form:"bio" xorm:"bio"`
	Nickname string `json:"nickname" form:"nickname" xorm:"nickname"`
	Job      string `json:"job" form:"job" xorm:"job"`
	Order    int    `json:"order" form:"order" xorm:"order"`
	Group    string `json:"group" form:"group" xorm:"group"`
}

// 组别列表
type userGroupListResponseRes struct {
	ID   string `json:"id" form:"id" xorm:"id"`
	Name string `json:"name" form:"name" xorm:"name"`
}

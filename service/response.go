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
	j.Message = "没有数据"
	c.JSON(http.StatusOK, j)
}

// BadRequest 参数错误
func (j *JSONData) BadRequest(c *gin.Context) {
	j.Code = http.StatusBadRequest
	j.Message = "参数错误"
	c.JSON(http.StatusOK, j)
}

// ResponseOK 请求成功
func (j *JSONData) ResponseOK(c *gin.Context) {
	j.Code = http.StatusOK
	j.Message = "操作成功"
	c.JSON(http.StatusOK, j)
}

// ServerError 服务器错误
func (j *JSONData) ServerError(c *gin.Context) {
	j.Code = http.StatusInternalServerError
	j.Message = "服务出了一点小问题"
	c.JSON(http.StatusOK, j)
}

// Unauthorized 缺失认证
func (j *JSONData) Unauthorized(c *gin.Context) {
	j.Code = http.StatusUnauthorized
	j.Message = ""
	c.JSON(http.StatusOK, j)
}

// FailAuth 认证失败
func (j *JSONData) FailAuth(c *gin.Context) {
	j.Code = http.StatusForbidden
	j.Message = ""
	c.JSON(http.StatusOK, j)
}

// RepetitiveOperation 不能重复操作
func (j *JSONData) RepetitiveOperation(c *gin.Context) {
	j.Code = http.StatusNotAcceptable
	j.Message = "不能重复操作"
	c.JSON(http.StatusOK, j)
}

// TimeOut 访问超时
func (j *JSONData) TimeOut(c *gin.Context) {
	j.Code = http.StatusRequestTimeout
	j.Message = ""
	c.JSON(http.StatusOK, j)
}

// 自定义响应结构

// 分页结构
type pagination struct {
	Current  int `json:"current" form:"current"`
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

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"helm-wrapper/pkg/errcode"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

type respBody struct {
	Code  int         `json:"code"` // 0 or 1, 0 is ok, 1 is error
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}
	r.Ctx.JSON(err.StatusCode(), response)
}

func respErr(c *gin.Context, err error) {
	glog.Warningln(err)

	c.JSON(http.StatusOK, &respBody{
		Code:  1,
		Error: err.Error(),
	})
}

func respOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &respBody{
		Code: 0,
		Data: data,
	})
}

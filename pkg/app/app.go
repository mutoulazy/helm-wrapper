package app

import (
	"github.com/gin-gonic/gin"
	"helm-wrapper/pkg/errcode"
	"net/http"
)

type ResponseBody struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Details []string    `json:"details,omitempty"`
}

type Response struct {
	Ctx *gin.Context
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
	r.Ctx.JSON(http.StatusOK, &ResponseBody{
		Code: 0,
		Data: data,
	})
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := &ResponseBody{
		Code: err.Code(),
		Msg:  err.Msg(),
	}
	details := err.Details()
	if len(details) > 0 {
		response.Details = details
	}
	r.Ctx.JSON(err.StatusCode(), response)
}

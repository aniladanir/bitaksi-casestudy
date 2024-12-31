package response

import (
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v3"
)

var responsePool = &sync.Pool{
	New: func() any {
		return new(Response)
	},
}

func getResponse() *Response {
	return responsePool.Get().(*Response)
}

func putResponse(response *Response) {
	if response != nil {
		responsePool.Put(response)
	}
}

type Response struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func Success(ctx fiber.Ctx, data any) error {
	resp := getResponse()
	defer putResponse(resp)
	resp.Success = true
	resp.Code = SuccessCode
	resp.Message = SuccessMsg
	resp.Data = data
	return ctx.Status(http.StatusOK).JSON(resp)
}

func Fail(ctx fiber.Ctx, code, msg string, status int) error {
	resp := getResponse()
	defer putResponse(resp)
	resp.Success = false
	resp.Code = code
	resp.Message = msg
	return ctx.Status(status).JSON(resp)
}

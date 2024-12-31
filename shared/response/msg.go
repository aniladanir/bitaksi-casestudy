package response

const (
	// Success Codes
	SuccessCode = "BT-0000"

	// Error Codes
	ErrCodeInternal          = "BT-0001"
	ErrCodeNotFound          = "BT-0002"
	ErrCodeUnauthorized      = "BT-0003"
	ErrCodeInvalidPayload    = "BT-0004"
	ErrCodeInvalidQueryParam = "BT-0005"
	ErrCodeBadRequest        = "BT-0006"

	// Messages
	SuccessMsg             = "Success"
	ErrMsgInternal         = "Internal Error"
	ErrMsgNotFound         = "Not Found"
	ErrMsgUnauthorized     = "Failed to Authorize"
	ErrMsgInvalidPayload   = "Invalid Payload"
	ErrMgInvalidQueryParam = "Invalid Query Params"
	ErrMsgBadRequest       = "Bad Request"
)

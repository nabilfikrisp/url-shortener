package response

type Success struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Error struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type SuccessPayloadParams struct {
	Message string
	Data    any
}

func SuccessPayload(p SuccessPayloadParams) Success {
	return Success{
		Message: p.Message,
		Data:    p.Data,
	}
}

type ErrorResponseParams struct {
	Message string
	Err     string
}

func ErrorPayload(p ErrorResponseParams) Error {
	return Error{
		Message: p.Message,
		Error:   p.Err,
	}
}

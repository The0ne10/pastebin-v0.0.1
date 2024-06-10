package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func Success() Response {
	return Response{Status: StatusSuccess}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

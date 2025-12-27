package api

type HttpResponse struct {
	Status string `json:"status"` //error, ok
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func NewError(msg string) HttpResponse {
	return HttpResponse{
		Status: StatusError,
		Error:  msg,
	}
}

func NewOk() HttpResponse {
	return HttpResponse{
		Status: StatusOk,
	}
}

func (r HttpResponse) Ok() bool {
	return r.Status == StatusOk
}

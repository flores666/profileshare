package api

type HttpResponse struct {
	Status      string            `json:"status"` //error, ok
	Data        any               `json:"data,omitempty"`
	Validations []ValidationError `json:"validations,omitempty"`
	Message     string            `json:"message,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

// NewError creates [StatusError] api handlers.
// validations can be nil
func NewError(msg string, validations *ValidationErrors) HttpResponse {
	response := HttpResponse{
		Status: StatusError,
	}

	if validations != nil {
		response.Validations = validations.Validations
	}

	if msg == "" {
		response.Message = msg
	}

	return response
}

// NewOk creates [StatusOk] api handlers.
// Payload can be nil
func NewOk(payload any) HttpResponse {
	response := HttpResponse{
		Status: StatusOk,
	}

	if payload != nil {
		response.Data = payload
	}

	return response
}

// Ok checks if handlers status is [StatusOk]
func (r HttpResponse) Ok() bool {
	return r.Status == StatusOk
}

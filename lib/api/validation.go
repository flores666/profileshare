package api

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	Validations []ValidationError `json:"validations,omitempty"`
	Message     string            `json:"message,omitempty"`
}

func NewValidationErrors(message string) *ValidationErrors {
	return &ValidationErrors{
		Message: message,
	}
}

func (e *ValidationErrors) Add(field, message string) {
	if e.Validations == nil {
		e.Validations = make([]ValidationError, 0)
	}

	e.Validations = append(e.Validations, ValidationError{field, message})
}

func (e *ValidationErrors) Ok() bool {
	return len(e.Validations) == 0
}

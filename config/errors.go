package config

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ApiError struct {
	Type     string       `json:"type"`
	Title    string       `json:"title"`
	Status   uint         `json:"status"`
	Detail   string       `json:"detail"`
	Instance string       `json:"instance"`
	Errors   []FieldError `json:"errors,omitempty"`
}

func (e *ApiError) Error() string {
	return e.Detail
}
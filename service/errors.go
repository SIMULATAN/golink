package service

type ServiceError struct {
	error
	Status int
}

type ErrorResponse struct {
	Message string `json:"message"`
}

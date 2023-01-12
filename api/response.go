package api

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseOk struct {
	Message string `json:"message"`
}

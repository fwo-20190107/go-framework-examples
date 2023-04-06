package infra

type HTTPError struct {
	Title string `json:"title"`
	Body  any    `json:"body"`
}

type HandleError struct {
	HTTPError *HTTPError
	Error     error
}

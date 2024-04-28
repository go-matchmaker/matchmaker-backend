package http

type ResponseFactory interface {
	Response(isError bool, msg string, responseType int) error
}

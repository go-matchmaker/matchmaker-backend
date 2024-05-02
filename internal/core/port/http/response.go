package http

type ResponseFactory interface {
	Response(err error, msg string, responseType int) error
}

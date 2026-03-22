package error

type AppError interface {
	error
	HTTPStatus() int
	Unwrap() error
}

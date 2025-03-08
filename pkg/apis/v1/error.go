package v1

type ErrorCode int

const (
	ErrCodeMaxFileSize ErrorCode = iota + 1
)

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

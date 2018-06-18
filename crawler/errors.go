package crawler

import "fmt"

type HttpError struct {
	url string
	code int
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Http error: [code=%d] %s", e.code, e.url)
}

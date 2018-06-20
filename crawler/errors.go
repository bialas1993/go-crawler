package crawler

import "fmt"

type HttpError struct {
	parent string
	url string
	code int
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Http error: [code=%d] [from=%s] %s", e.code, e.parent, e.url)
}

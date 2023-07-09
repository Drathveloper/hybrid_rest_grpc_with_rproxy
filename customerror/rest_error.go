package customerror

import "fmt"

type RestError struct {
	Code    int
	Message string
}

func (e RestError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

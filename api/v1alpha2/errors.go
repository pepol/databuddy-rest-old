package v1alpha2

// ErrorNotFound is an error-type to signalize that the object wasn't found.
type ErrorNotFound struct{}

func (e *ErrorNotFound) Error() string {
	return "object wasn't found"
}

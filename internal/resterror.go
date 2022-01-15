package internal

type RestError struct {
	errorString string
	status      int
}

func (e RestError) Error() string {
	return e.errorString
}

func (e RestError) Status() int {
	return e.status
}

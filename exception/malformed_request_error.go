package exception

type ReqBodyMalformedError struct {
	Error string
}

func NewReqBodyMalformedError(err string) ReqBodyMalformedError {
	return ReqBodyMalformedError{Error: err}
}

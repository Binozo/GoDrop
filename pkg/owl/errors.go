package owl

type OwlError struct {
	ErrorType OwlErrorType
	Logs      []string
}

type OwlErrorType string

const Unknown OwlErrorType = "Unknown"
const UnknownInterface OwlErrorType = "UnknownInterface"
const InaccessibleInterface OwlErrorType = "InaccessibleInterface"

package auth

// Error implements error interface.
type Error string

func (err Error) Error() string {
	return string(err)
}

const (
	// ErrBadParameters is bad parameters error.
	ErrBadParameters Error = "bad parameters"
	// ErrNotFound is unsuccessful search error.
	ErrNotFound Error = "not found"
)

// BadParametersError implements error interface.
type BadParametersError struct {
	err error
}

func (bpErr BadParametersError) Error() string {
	return bpErr.err.Error()
}

package sessions

type Result struct {
	value string
	err   error
}

func MakeResult(value string, err error) Result {
	return Result{
		value: value,
		err:   err,
	}
}

func (r *Result) Err() error {
	return r.err
}

func (r *Result) String() (string, error) {
	return r.value, r.err
}

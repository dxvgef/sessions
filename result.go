package sessions

import (
	"strconv"
)

type Result struct {
	value string
	err   error
}

func NewResult(value string, err error) Result {
	return Result{
		value: value,
		err:   err,
	}
}

func (r Result) Err() error {
	return r.err
}

func (r Result) String() (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return r.value, r.err
}

func (r Result) Bytes() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	return stringToBytes(r.value), r.err
}

func (r Result) Int() (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	return strconv.Atoi(r.value)
}

func (r Result) Int8() (int8, error) {
	if r.err != nil {
		return 0, r.err
	}
	value, err := strconv.ParseInt(r.value, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(value), nil
}

func (r Result) Int16() (int16, error) {
	if r.err != nil {
		return 0, r.err
	}
	value, err := strconv.ParseInt(r.value, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(value), nil
}

func (r Result) Int32() (int32, error) {
	if r.err != nil {
		return 0, r.err
	}
	value, err := strconv.ParseInt(r.value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(value), nil
}

func (r Result) Int64() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return strconv.ParseInt(r.value, 10, 64)
}

func (r Result) Uint8() (uint8, error) {
	if r.err != nil {
		return 0, r.err
	}
	value, err := strconv.ParseUint(r.value, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(value), nil
}

func (r Result) Uint16() (uint16, error) {
	if r.err != nil {
		return 0, r.err
	}
	value, err := strconv.ParseUint(r.value, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(value), nil
}

func (r Result) Uint32() (uint32, error) {
	if r.err != nil {
		return 0, r.err
	}
	value, err := strconv.ParseUint(r.value, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(value), nil
}

func (r Result) Uint64() (uint64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return strconv.ParseUint(r.value, 10, 64)
}

func (r Result) Float32() (float32, error) {
	if r.err != nil {
		return 0, r.err
	}
	value, err := strconv.ParseFloat(r.value, 32)
	if err != nil {
		return 0, err
	}
	return float32(value), nil
}

func (r Result) Float64() (float64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return strconv.ParseFloat(r.value, 64)
}

func (r Result) Bool() (bool, error) {
	if r.err != nil {
		return false, r.err
	}
	return strconv.ParseBool(r.value)
}

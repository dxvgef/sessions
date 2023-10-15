package sessions

import (
	"strconv"
)

// Result 结果
type Result struct {
	value string
	err   error
}

// NewResult 新建一个结果
func NewResult(value string, err error) Result {
	return Result{
		value: value,
		err:   err,
	}
}

// Err 获取结果的错误信息
func (r Result) Err() error {
	return r.err
}

// String 将结果转为string类型
func (r Result) String() (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return r.value, r.err
}

// Bytes 将结果转为[]byt类型
func (r Result) Bytes() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	return stringToBytes(r.value), r.err
}

// Int 将结果转为int类型
func (r Result) Int() (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	return strconv.Atoi(r.value)
}

// Int8 将结果转为int8类型
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

// Int16 将结果转为int16类型
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

// Int32 将结果转为int32类型
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

// Int64 将结果转为int64类型
func (r Result) Int64() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return strconv.ParseInt(r.value, 10, 64)
}

// Uint8 将结果转为uint8类型
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

// Uint16 将结果转为uint16类型
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

// Uint32 将结果转为uint32类型
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

// Uint64 将结果转为uint64类型
func (r Result) Uint64() (uint64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return strconv.ParseUint(r.value, 10, 64)
}

// Float32 将结果转为float32类型
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

// Float64 将结果转为float64类型
func (r Result) Float64() (float64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return strconv.ParseFloat(r.value, 64)
}

// Bool 将结果转为bool类型
func (r Result) Bool() (bool, error) {
	if r.err != nil {
		return false, r.err
	}
	return strconv.ParseBool(r.value)
}

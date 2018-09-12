package ini

import (
	"strconv"
	"time"
)

type Key struct {
	name  string
	value string
}

func (this *Key) Name() string {
	return this.name
}

func (this *Key) Value() string {
	return this.value
}

func (this *Key) Bool() (bool, error) {
	return strconv.ParseBool(this.value)
}

func (this *Key) Float32() (float32, error) {
	f64, err := strconv.ParseFloat(this.value, 32)
	return float32(f64), err
}

func (this *Key) Float64() (float64, error) {
	return strconv.ParseFloat(this.value, 64)
}

func (this *Key) Int() (int, error) {
	return strconv.Atoi(this.value)
}

func (this *Key) Int8() (int8, error) {
	i64, err := strconv.ParseInt(this.value, 10, 8)
	return int8(i64), err
}

func (this *Key) Int16() (int16, error) {
	i64, err := strconv.ParseInt(this.value, 10, 16)
	return int16(i64), err
}

func (this *Key) Int32() (int32, error) {
	i64, err := strconv.ParseInt(this.value, 10, 32)
	return int32(i64), err
}

func (this *Key) Int64() (int64, error) {
	return strconv.ParseInt(this.value, 10, 64)
}

func (this *Key) Uint() (uint, error) {
	ui64, err := strconv.ParseUint(this.value, 10, 0)
	return uint(ui64), err
}

func (this *Key) Uint8() (uint8, error) {
	ui64, err := strconv.ParseUint(this.value, 10, 8)
	return uint8(ui64), err
}

func (this *Key) Uint16() (uint16, error) {
	ui64, err := strconv.ParseUint(this.value, 10, 16)
	return uint16(ui64), err
}

func (this *Key) Uint32() (uint32, error) {
	ui64, err := strconv.ParseUint(this.value, 10, 32)
	return uint32(ui64), err
}

func (this *Key) Uint64() (uint64, error) {
	return strconv.ParseUint(this.value, 10, 64)
}

func (this *Key) Duration() (time.Duration, error) {
	return time.ParseDuration(this.value)
}

func (this *Key) ParseTime(layout string) (time.Time, error) {
	return time.Parse(layout, this.value)
}

func (this *Key) Time() (time.Time, error) {
	return time.Parse(time.RFC3339, this.value)
}

func (this *Key) ParseTimeInLocation(layout string) (time.Time, error) {
	return time.ParseInLocation(layout, this.value, time.Local)
}

func (this *Key) TimeInLocation() (time.Time, error) {
	return time.ParseInLocation(time.RFC3339, this.value, time.Local)
}

func (this *Key) MustBool() bool {
	b, err := strconv.ParseBool(this.value)
	if err != nil {
		panic(err)
	}
	return b
}

func (this *Key) MustInt() int {
	i, err := strconv.Atoi(this.value)
	if err != nil {
		panic(err)
	}
	return i
}

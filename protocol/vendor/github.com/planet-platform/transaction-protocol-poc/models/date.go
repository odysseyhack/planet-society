package models

import (
	"fmt"
	"io"
	"time"
)

type Date struct {
	Value int64
}

func (d Date) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(time.Unix(0, d.Value).Format("02-01-2006")))
}

func (d *Date) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case string:
		t, err := time.Parse("02-01-2006", v)
		if err != nil {
			return err
		}
		d.Value = t.UnixNano()
		return nil
	default:
		return fmt.Errorf("%T is invalid", v)
	}
}

func DateFromTime(t time.Time) Date {
	return Date{Value: t.UnixNano()}
}

package utils

import (
	"time"
)

type Datetime struct {
	t time.Time
}

// ToEpochMS return epoch time in ms by normalize time to UTC first
func (d *Datetime) ToEpochMS() int64 {
	return d.t.UTC().UnixMilli()
}

// ToEpoch return epoch time in second by normalize time to UTC first
func (d *Datetime) ToEpoch() int64 {
	return d.t.Truncate(time.Second).UTC().Unix()
}

func (d *Datetime) ToString() string {
	return d.t.Format(time.RFC3339)
}

func (d Datetime) GetTime() time.Time {
	return d.t
}

func NewDatetimeNow() *Datetime {
	return &Datetime{
		t: time.Now().UTC().Truncate(time.Second),
	}
}

func NewDatetimeFromEpoch(ep int64) *Datetime {
	return &Datetime{
		t: time.Unix(ep, 0),
	}
}

func NewDatetimeFromString(t string) (*Datetime, error) {
	tn, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return nil, err
	}

	return &Datetime{
		t: tn,
	}, nil
}

func NewDatetimeFromTime(t time.Time) *Datetime {
	return &Datetime{
		t: t,
	}
}

package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tm := time.Date(2023, 2, 27, 23, 24, 0, 0, time.UTC)
	hd := humanDate(tm)
	if hd != "27 Feb 2023 at 23:24" {
		t.Errorf("got %q; want %q", hd, "27 Feb 2023 at 23:24")
	}
}
func TestHumanDate2(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2023, 2, 27, 23, 24, 0, 0, time.UTC),
			want: "27 Feb 2023 at 23:24",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2023, 2, 27, 23, 24, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "27 Feb 2023 at 22:24",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			if hd != tt.want {
				t.Errorf("got %q; want %q", hd, tt.want)
			}
		})
	}
}

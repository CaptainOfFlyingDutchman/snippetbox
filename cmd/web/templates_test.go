package main

import (
	"testing"
	"time"

	"snippetbox.manvendrask.com/internal/assert"
)

type TestCase struct {
	name string
	tm   time.Time
	want string
}

func TestHumanDate(t *testing.T) {

	tests := []TestCase{
		{
			name: "UTC",
			tm:   time.Date(2026, 6, 10, 17, 57, 0, 0, time.UTC),
			want: "10 Jun 2026 at 17:57",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2026, 6, 10, 17, 57, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "10 Jun 2026 at 16:57",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}

package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_downloadSubtitles(t *testing.T) {
	tests := []struct {
		name string
		id   string
		err  error
	}{
		{
			name: "Subtitles for video available",
			id:   "iyBAl1_grCU",
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := downloadSubtitles(tt.id)
			assert.Equal(t, tt.err, err)
		})
	}
}

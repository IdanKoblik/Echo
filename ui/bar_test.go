package ui

import (
	"testing"

	"github.com/schollz/progressbar/v3"
	"github.com/stretchr/testify/assert"
)

func TestProgressBar_Init(t *testing.T) {
	pb := ProgressBar{
		Len:         10,
		Description: "Testing",
	}

	bar := pb.Init()
	assert.NotNil(t, bar, "Progress bar should not be nil")

	assert.Equal(t, int64(10), bar.GetMax64(), "Progress bar length should match")
	assert.IsType(t, &progressbar.ProgressBar{}, bar, "Should return a pointer to progressbar.ProgressBar")

	err := bar.Add(1)
	assert.NoError(t, err, "Should be able to add to progress bar")

	err = bar.Finish()
	assert.NoError(t, err, "Finish should not error")
}

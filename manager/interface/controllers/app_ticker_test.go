package controllers

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAppTicker_New(t *testing.T) {
	appTicker := NewAppTicker()
	ticker := appTicker.New(1*time.Second, time.Second)
	for i := 0; i < 3; i++ {
		now := <-ticker.C
		t.Log(now)
		assert.WithinDuration(t, now.Truncate(time.Second), now, 10*time.Millisecond)
	}
}

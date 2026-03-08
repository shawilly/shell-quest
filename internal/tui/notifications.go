package tui

import (
	"time"

	"go.dalton.dog/bubbleup"
)

func newAlert() bubbleup.AlertModel {
	return bubbleup.NewAlertModel(52, false, 4*time.Second).
		WithUnicodePrefix().
		WithPosition(bubbleup.TopRightPosition)
}

package draftout

import (
	"fmt"
	"strconv"
)

func ParseColorString(hexCode string) (int, error) {
	rc, err := strconv.ParseInt(hexCode[1:], 16, 32)
	if err != nil {
		return -1, err
	}
	return int(rc), nil
}

func FormatDuration(ms int) string {
	totalSeconds := ms / 1000
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	if hours == 0 {
		return fmt.Sprintf("%02d:%02d", minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

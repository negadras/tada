package utils

import (
	"fmt"
	"github.com/negadras/tada/cmd/db"
	"strings"
)

// ParsePriority converts string priority to Priority enum
func ParsePriority(priorityStr string) (db.Priority, error) {
	switch strings.ToLower(strings.TrimSpace(priorityStr)) {
	case "low", "l", "1":
		return db.Low, nil
	case "medium", "m", "2":
		return db.Medium, nil
	case "high", "h", "3":
		return db.High, nil
	default:
		return db.Medium, fmt.Errorf("must be one of: low, medium, high (or l, m, h)")
	}
}

// ParseStatus converts string status to Status enum
func ParseStatus(status string) (db.Status, error) {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "open", "o", "1":
		return db.Open, nil
	case "done", "d", "2":
		return db.Done, nil
	default:
		return db.Open, fmt.Errorf("must be one of: open or done") // Default to Open status
	}
}

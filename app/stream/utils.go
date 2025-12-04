package stream

import (
	"errors"
	"strconv"
	"strings"
)

// ParseID parses a stream ID string into millisecondsTime and sequenceNumber.
// Format: <millisecondsTime>-<sequenceNumber>
func ParseID(id string) (int64, int64, error) {
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid ID format")
	}

	msTime, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, errors.New("invalid millisecondsTime")
	}

	seqNum, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, errors.New("invalid sequenceNumber")
	}

	return msTime, seqNum, nil
}

// ValidateID checks if the new ID is valid given the last ID in the stream.
// Rules:
// 1. ID must be strictly greater than lastID.
// 2. ID must be greater than 0-0.
func ValidateID(id string, lastID string) error {
	msTime, seqNum, err := ParseID(id)
	if err != nil {
		return err
	}

	// Rule: ID must be greater than 0-0
	if msTime == 0 && seqNum == 0 {
		return errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}

	// If stream is empty (lastID is empty), any ID > 0-0 is valid
	if lastID == "" {
		return nil
	}

	lastMsTime, lastSeqNum, err := ParseID(lastID)
	if err != nil {
		return err // Should not happen if stored IDs are valid
	}

	// Rule: ID must be strictly greater than lastID
	if msTime < lastMsTime {
		return errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if msTime == lastMsTime && seqNum <= lastSeqNum {
		return errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	return nil
}

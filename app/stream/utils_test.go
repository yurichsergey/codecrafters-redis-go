package stream

import (
	"testing"
)

func TestGenerateSequence(t *testing.T) {
	tests := []struct {
		name     string
		msTime   int64
		lastID   string
		expected int64
		isValid  bool
	}{
		{
			name:     "Empty Stream, msTime 0",
			msTime:   0,
			lastID:   "",
			expected: 1, // 0-1
			isValid:  true,
		},
		{
			name:     "Empty Stream, msTime > 0",
			msTime:   1000,
			lastID:   "",
			expected: 0, // 1000-0
			isValid:  true,
		},
		{
			name:     "Same Time",
			msTime:   1000,
			lastID:   "1000-0",
			expected: 1, // 1000-1
			isValid:  true,
		},
		{
			name:     "Same Time, increment",
			msTime:   1000,
			lastID:   "1000-5",
			expected: 6, // 1000-6
			isValid:  true,
		},
		{
			name:     "Different Time (Newer), msTime 0",
			msTime:   0, // Not sure if this is a valid case for "Newer" if last was > 0, but adhering to logic: msTime 0 -> 1
			lastID:   "10-0",
			expected: 1, // 0-1? Wait, if lastID is 10-0, 0-1 is smaller. But GenerateSequence just generates sequence. Validation handles if it's smaller.
			// Actually, ValidateID checks order. GenerateSequence just gives the next seq.
			// The logic says: if msTime == 0 -> return 1.
			isValid: true,
		},
		{
			name:     "Different Time (Newer), msTime > 0",
			msTime:   2000,
			lastID:   "1000-5",
			expected: 0, // 2000-0
			isValid:  true,
		},
		{
			name:     "Invalid Last ID",
			msTime:   1000,
			lastID:   "invalid",
			expected: 0,
			isValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSequence(tt.msTime, tt.lastID)
			if (err == nil) != tt.isValid {
				t.Errorf("GenerateSequence() error = %v, isValid %v", err, tt.isValid)
				return
			}
			if got != tt.expected {
				t.Errorf("GenerateSequence() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		wantMsTime int64
		wantSeqNum int64
		isValid    bool
	}{
		{"Valid ID", "1000-5", 1000, 5, true},
		{"Valid ID (0-0)", "0-0", 0, 0, true},
		{"Invalid Format (no hyphen)", "1000", 0, 0, false},
		{"Invalid Format (too many parts)", "1000-5-5", 0, 0, false},
		{"Invalid msTime", "abc-5", 0, 0, false},
		{"Invalid seqNum", "1000-def", 0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsTime, gotSeqNum, err := ParseID(tt.id)
			if (err == nil) != tt.isValid {
				t.Errorf("ParseID() error = %v, isValid %v", err, tt.isValid)
				return
			}
			if gotMsTime != tt.wantMsTime {
				t.Errorf("ParseID() gotMsTime = %v, want %v", gotMsTime, tt.wantMsTime)
			}
			if gotSeqNum != tt.wantSeqNum {
				t.Errorf("ParseID() gotSeqNum = %v, want %v", gotSeqNum, tt.wantSeqNum)
			}
		})
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		lastID  string
		isValid bool
	}{
		{"Valid > Last (Time >)", "1001-0", "1000-5", true},
		{"Valid > Last (Time ==, Seq >)", "1000-6", "1000-5", true},
		{"Valid Empty Stream (0-1)", "0-1", "", true},
		{"Valid Empty Stream (1000-0)", "1000-0", "", true},
		{"Invalid 0-0", "0-0", "", false},
		{"Invalid <= Last (Time <)", "999-5", "1000-5", false},
		{"Invalid <= Last (Time ==, Seq <)", "1000-4", "1000-5", false},
		{"Invalid <= Last (Time ==, Seq ==)", "1000-5", "1000-5", false},
		{"Invalid Format", "invalid", "", false},
		{"Invalid Last ID (should not happen)", "1000-5", "invalid", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateID(tt.id, tt.lastID); (err == nil) != tt.isValid {
				t.Errorf("ValidateID() error = %v, isValid %v", err, tt.isValid)
			}
		})
	}
}

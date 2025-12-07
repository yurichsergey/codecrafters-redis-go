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
		wantErr  bool
	}{
		{
			name:     "Empty Stream, msTime 0",
			msTime:   0,
			lastID:   "",
			expected: 1, // 0-1
			wantErr:  false,
		},
		{
			name:     "Empty Stream, msTime > 0",
			msTime:   1000,
			lastID:   "",
			expected: 0, // 1000-0
			wantErr:  false,
		},
		{
			name:     "Same Time",
			msTime:   1000,
			lastID:   "1000-0",
			expected: 1, // 1000-1
			wantErr:  false,
		},
		{
			name:     "Same Time, increment",
			msTime:   1000,
			lastID:   "1000-5",
			expected: 6, // 1000-6
			wantErr:  false,
		},
		{
			name:     "Different Time (Newer), msTime 0",
			msTime:   0, // Not sure if this is a valid case for "Newer" if last was > 0, but adhering to logic: msTime 0 -> 1
			lastID:   "10-0",
			expected: 1, // 0-1? Wait, if lastID is 10-0, 0-1 is smaller. But GenerateSequence just generates sequence. Validation handles if it's smaller.
			// Actually, ValidateID checks order. GenerateSequence just gives the next seq.
			// The logic says: if msTime == 0 -> return 1.
			wantErr: false,
		},
		{
			name:     "Different Time (Newer), msTime > 0",
			msTime:   2000,
			lastID:   "1000-5",
			expected: 0, // 2000-0
			wantErr:  false,
		},
		{
			name:     "Invalid Last ID",
			msTime:   1000,
			lastID:   "invalid",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSequence(tt.msTime, tt.lastID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSequence() error = %v, wantErr %v", err, tt.wantErr)
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
		wantErr    bool
	}{
		{"Valid ID", "1000-5", 1000, 5, false},
		{"Valid ID (0-0)", "0-0", 0, 0, false},
		{"Invalid Format (no hyphen)", "1000", 0, 0, true},
		{"Invalid Format (too many parts)", "1000-5-5", 0, 0, true},
		{"Invalid msTime", "abc-5", 0, 0, true},
		{"Invalid seqNum", "1000-def", 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsTime, gotSeqNum, err := ParseID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseID() error = %v, wantErr %v", err, tt.wantErr)
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
		wantErr bool
	}{
		{"Valid > Last (Time >)", "1001-0", "1000-5", false},
		{"Valid > Last (Time ==, Seq >)", "1000-6", "1000-5", false},
		{"Valid Empty Stream (0-1)", "0-1", "", false},
		{"Valid Empty Stream (1000-0)", "1000-0", "", false},
		{"Invalid 0-0", "0-0", "", true},
		{"Invalid <= Last (Time <)", "999-5", "1000-5", true},
		{"Invalid <= Last (Time ==, Seq <)", "1000-4", "1000-5", true},
		{"Invalid <= Last (Time ==, Seq ==)", "1000-5", "1000-5", true},
		{"Invalid Format", "invalid", "", true},
		{"Invalid Last ID (should not happen)", "1000-5", "invalid", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateID(tt.id, tt.lastID); (err != nil) != tt.wantErr {
				t.Errorf("ValidateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

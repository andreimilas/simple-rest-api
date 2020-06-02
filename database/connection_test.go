package database

import (
	"testing"
)

func TestConnect(t *testing.T) {
	// Multiple test cases
	var tests = []struct {
		name         string
		dsn          string
		isNull       bool
		errorMessage string
	}{
		{"Empty DSN",
			"",
			true,
			"Database handle should be null on empty DSN",
		},
		{"Invalid DSN",
			"mysql:",
			true,
			"Database handle should be null on invalid DSN",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dbHandle := Connect(test.dsn)
			if (dbHandle == nil) != test.isNull {
				t.Error(test.errorMessage)
			}
		})
	}
}

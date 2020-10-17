package verifychecksum

import (
	"encoding/csv"
	"io"
	"strings"
	"testing"

	cmp "github.com/google/go-cmp/cmp"
)

const csvPayload = "abc,123\nefg,456"

func TestValidateCRC32Reader(t *testing.T) {
	cases := map[string]struct {
		Payload      io.Reader
		ExpectResult [][]string
		ExpectCRC32  uint32
		ExpectErr    string
	}{
		"success": {
			Payload: strings.NewReader(csvPayload),
			ExpectResult: [][]string{
				{"abc", "123"},
				{"efg", "456"},
			},
			ExpectCRC32: 0x6405FE3A,
		},
		"failure": {
			Payload:     strings.NewReader(csvPayload[5:]),
			ExpectCRC32: 0x6405FE3A,
			ExpectErr:   "does not match",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			payload := NewValidateCRC32Reader(c.ExpectCRC32, c.Payload)

			rows, err := csv.NewReader(payload).ReadAll()
			if len(c.ExpectErr) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.ExpectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect error to contain %v, got %v", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if diff := cmp.Diff(rows, c.ExpectResult); len(diff) != 0 {
				t.Errorf("expect result match\n%s", diff)
			}
		})
	}
}

func TestValidateIntegrity(t *testing.T) {
	cases := map[string]struct {
		Payload      []byte
		ExpectResult [][]string
		ExpectCRC32  uint32
		ExpectErr    string
	}{
		"success": {
			Payload: []byte(csvPayload),
			ExpectResult: [][]string{
				{"abc", "123"},
				{"efg", "456"},
			},
			ExpectCRC32: 0x6405FE3A,
		},
		"failure": {
			Payload:     []byte(csvPayload[5:]),
			ExpectCRC32: 0x6405FE3A,
			ExpectErr:   "does not match",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := ValidateIntegrity(c.ExpectCRC32, c.Payload)
			if len(c.ExpectErr) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.ExpectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect error to contain %v, got %v", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

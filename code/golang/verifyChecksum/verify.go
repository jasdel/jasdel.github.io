package verifychecksum

import (
	"fmt"
	"hash/crc32"
	"io"
)

// ValidateIntegrity computes the CRC32 checksum of the byte slice provided,
// and compares it against the expectec checksum.
func ValidateIntegrity(expect uint32, body []byte) error {
	actual := crc32.ChecksumIEEE(body)
	if expect != actual {
		return fmt.Errorf("invalid checksum %08X does not match %08X",
			actual, expect)
	}

	return nil
}

// ValidateCRC32Reader provides a wrapper for an io.Reader that computes a
// running checksum as the wrapped reader is read. When the underlying reader
// returns EOF, the validator will compare the checksums.
type ValidateCRC32Reader struct {
	io.Reader
	expect uint32
	actual uint32
}

// NewValidateCRC32Reader initializes a new ValidateCRC32Reader wrapping the
// reader provided.
func NewValidateCRC32Reader(expect uint32, reader io.Reader) *ValidateCRC32Reader {
	return &ValidateCRC32Reader{
		expect: expect,
		Reader: reader,
	}
}

// Read reads from the wrapped reader, and updates the checksum. When the
// wrapped reader returns EOF, the expected checksum will be compared against
// the computed value. If they differ and error will be returned.
func (v *ValidateCRC32Reader) Read(p []byte) (int, error) {
	n, err := v.Reader.Read(p)

	v.actual = crc32.Update(v.actual, crc32.IEEETable, p[:n])
	if err == io.EOF && v.actual != v.expect {
		return n, fmt.Errorf("invalid checksum %08X does not match %08X",
			v.actual, v.expect)
	}

	return n, err
}

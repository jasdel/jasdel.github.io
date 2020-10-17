---
title: "Verify Payload Checksum with Streaming Reader"
date: 2020-10-17T12:00:00-07:00
draft: false
tags: [golang]
---

Sometimes you may find your self working with an API gives you the
opportunity to validate the integrity of the payload provided through
a checksum. There are a few different ways that you could go about
validating the integrity of the payload.

### Preprocess payload validation

One way to do this would be to load the full payload all into memory. Then
before consuming the payload your application would compute the checksum
of the payload and compare it against the expect value.

The library and exact method of computing the checksum may differ based on
your use case, but the resulting logic most likely will be similar to the
example code here. For these examples we'll use Go's [hash/crc32] package
with the common polynomial
[IEEE](https://golang.org/pkg/hash/crc32/#IEEE).

```go
// import "hash/crc32"

// ValidateIntegrity computes the CRC32 checksum of the byte slice provided,
// and compares it against the expectec checksum.
func ValidateIntegrity(expect uint32, body []byte) error {
	actual := crc32.ChecksumIEEE(body)
	if expect != actual {
		return fmt.Errorf("expect %08X checksum, got %08X", expect, actual)
	}

    return nil
}
```

This approach has the benefit of validating the integrity of the payload
before it is consumed. But comes with the downside that if the payload is
too large to this approach could put significant pressure on the
application's memory.

### Stream payload validation

Instead of requiring the payload to be read into memory before it can be
validated, our application could compute the checksum as the payload is
consumed upstream. This approach will allow our application to operate on
the payload as an stream, (i.e. [io.Reader]) while at the same time
compute and validate the checksum. Very little additional memory is needed
to compute the payload's checksum this way.

This is a good approach when the upstream processing in transactional and
will not operate on the payload until its been fully processed. If this is
not the case, streaming checksum validation can lead to unexpected errors,
because the integrity of the payload can not be validated until the entire
payload has been consumed.

```go
// import "hash/crc32"
// import "encoding/csv"

wrapped := NewValidateCRC32Reader(expectCRC32, payload)

rows, err := csv.NewReader(wrapped).ReadAll()
```

In the above example a `StreamValidateCRC32` function is called that takes
the expected CRC32 checksum, the payload, and returns the payload wrapped
with CRC32 checksum validation. This function wraps the payload with
a type that implements the [io.Reader] interface. This type will
iteratively compute the checksum of the payload as the downstream consumer
reads from it. When the payload has been fully consumed, the wrapper will
compare the computed checksum with the expected value. If the values are
not equal are not equal an relevant error can be returned. Otherwise the
wrapper will return `io.EOF` when the payload is valid and has been fully
read.

```go
// import "io"
// import "hash/crc32"

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
```

The above `ValidateCRC32Reader` wraps the [io.Reader] it is initialized
with, and computes the checksum of the wrapped reader as its content is
read. When the wrapped reader returns [io.EOF] the `ValidateCRC32Reader`
will return an error if the actual checksum of the content doesn't match
the expected value.

Checkout the full example [On Github](https://github.com/jasdel/jasdel.github.io/blob/b6d84b368a93fb19b608d7408c8bb2815c58f0d3/code/golang/verifyChecksum/verify.go#L21-L52)


[io.EOF]: https://golang.org/pkg/io/#EOF
[io.Reader]: https://golang.org/pkg/io/#Reader
[hash/crc32]: https://golang.org/pkg/hash/crc32/
[encoding/json]: https://golang.org/pkg/encoding/json/

---
title: "Verify Payload Checksum with Streaming Reader"
date: 2020-08-15T23:17:34-07:00
draft: true
tags: [golang]
---

Sometimes you may find your self working with an API gives you the opportunity
to validate the integrity of the payload provided through a checksum.
There are a few different ways that you could go about validating the integrity
of the payload.

### Preprocess payload validation

One way to do this would be to load the full payload all into memory. Then
before consuming the payload your application would compute the checksum of
the payload and compare it against the expect value. 

```go
// import "hash/crc32"

// ValidateIntegrity computes the CRC32 checksum of the byte slice provided,
// and compares it against the expectec checksum.
func ValidateIntegrity(expect uint32, body []byte) error {
    hash := crc32.NewIEEE()
    hash.Write(body)
    actual := hash.Sum32()

    if  expect != actual {
        return fmt.Errorf("expect %08X checksum, got %08X", expect, actual)
    }

    return nil
}
```

This approach has the benefit of validating the integrity of the payload before
it is consumed. But comes with the downside that if the payload is too large to
this approach could put significant pressure on the application's memory.

### Stream payload validation

Instead of requiring the payload to be read into memory before it can be
validated, our application could compute the checksum as the payload is
consumed upstream. This approach will allow our application to operate on the
payload as an stream, (i.e. `io.Reader`) while at the same time compute and
validate the checksum. Very little additional memory is needed to compute the
payload's checksum this way.

This is a good approach when the upstream processing in transactional and will
not operate on the payload until its been fully processed. If this is not the
case, streaming checksum validation can lead to unexpected errors, because the
integrity of the payload can not be validated until the entire payload has been
consumed. 

```go
// import "hash/crc32"
// import "encoding/json"

wrappedPayload := StreamValidateCRC32(expectCRC32, payload)

err := json.NewDecoder(wrappedPayload).Decode(&result)
```

In the above example a `StreamValidateCRC32` function is called that takes the
expected CRC32 checksum, the payload, and returns the payload wrapped with
CRC32 checksum validation.. This function wraps the payload with a type that
implements the `io.Reader` interface. This type will iteratively compute the
checksum of the payload as the downstream consumer reads from it. When the
payload has been fully consumed, the wrapper will compare the computed checksum
with the expected value. If the values are not equal are not equal an relevant
error can be returned. Otherwise the wrapper will return `io.EOF` when the
payload is valid and has been fully read.

```go
// TODO wrapper implemenration
```

* TODO describe implementation parts.

[hash/crc32]: https://golang.org/pkg/hash/crc32/
[encoding/json]: https://golang.org/pkg/encoding/json/

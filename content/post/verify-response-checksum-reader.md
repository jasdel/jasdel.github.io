---
title: "Verify Response Checksum Reader"
date: 2020-08-15T23:17:34-07:00
draft: true
tags: [golang]
---

Sometimes you'll find your self working with an API that provides a checksum of
the response body. There are a few different ways that you could validate the
integrity of the response body.

One way to do this would be to load the full response all into memory. Then use
the [hash/crc32] standard library to compute the hash of the response body.
This works, but requires the full response to be read into memory then later
be read a second time when unmarshaling it.

An alternative approach would be to compute the hash inline as the response
body is being unmarshaled. This has the benefit of not loading the full
response into memory at once. Which for some large responses, may be
problematic. Though, this approach does mean that the response cannot be
validated until after its been unmarshaled.

For this exercise we'll have the checksum in a response header, but the
approach would be similar for checksum appended after the original body as a
trailing checksum. The difference being the validation logic wouldn't know what
the expected checksum is until the last 4 bytes of the response body is read.

## Wrapping response Body io.Reader with checksum validation

We'll use [Amazon's DynamoDB][ddb] API as an example. The
DynamoDB's API may respond with a `X-Amz-Crc32` header. This header will
provide the checksum that we'll use to validate the integrity of the response
body.

Using the following HTTP response message as an example we'll use to validate
that the response body content against the CRC32 checksum in the
`X-Crc32-Checksum` response header.

```http
HTTP/1.1 200 OK
x-amzn-RequestId: <RequestId>
x-amz-crc32: 1234567
Content-Type: application/x-amz-json-1.0
Content-Length: 32
Date: <Date>

{"Item": {"A": 123, "B": "efg"}}
```

[ddb]: https://aws.amazon.com/dynamodb/

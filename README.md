# StreamingFast Errors Library
[![reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/streamingfast/derr)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This repository contains all common code for handling errors across our
various services. It is part of **[StreamingFast](https://github.com/streamingfast/streamingfast)**.


## Usage

1. Use  `derr.Wrap(err, "some  message")`  to  wrap  any calls  to  a
   sub-system that *could* have  yielded some `derr.ErrorResponse` (so
   it is passed to the user untouched).
2. Use `derr.SomeError` (see `errors.go`) wherever possible to have
   consistent errors throughout our system.  If you think the error
   you are creating can be shared, put it in
   [`errors.go`](./errors.go) and use that.
3. Craft your own custom Error builder (in your project's `errors.go`,
   see
   [eosws](https://github.com/streamingfast/dgraphql/blob/develop/errors.go))
   and use that in your code.  Craft a meaningful `derr.C` (or
   `derr.ErrorCode`) with a good name (see below), reuse where
   possible.

### gRPC APIs

We create `derr.Status` errors when we want to set a specific error
code. We call `derr.Wrap` to prefix the error message going out
(without altering the outgoing Status Code).

Returning plain `fmt.Errorf()` upstream implies returning a
`codes.Internal` gRPC status code.

Missing a `derr.Wrap()` kills the links and removes the StatusCode, so
use `Wrap` (or `Wrapf`) whenever you want to add a prefix.

### REST APIs

All error should be created by defining a properly named function that receives a `context.Context`
object as well as specialized parameters to craft a proper error message and details. This specific
error function should be re-using under the cover one the underlying pre-defined HTTP error creator,
[see list](./errors.go#36).

The most important thing is that all error should have a proper unique code defined using
`derr.C("value")` type alias which is a sugar syntax for `derr.ErrorCode("string_code_error")`.

This is of the **uttermost importance** so that we can easily grep all
our source code files to extract the used specific error codes across all our services
for documentation purposes.

The string code defined must be unique among all our services, human readable,
should clearly represent the error in few words, should be in `snake_case` format and
should end with `_error`.

For example, let's say that in your micro service, you would like to define a specialized
bad request error when a particular request is invalid due to block number being too low.

Here the piece of code that your should do:

```
func BlockNumTooLowError(ctx context.Context, blockNum uint32, thresholdBlockNum uint32) *derr.ErrorResponse {
	return HTTPBadRequestError(ctx, err, derr.C("block_num_too_low_error"), "The requested block num is too low",
        "actual_block_num", blockNum,
        "threshold_block_num", thresholdBlockNum,
    )
}
```

Each of generic HTTP error creator receives the `context.Context` object. This context is required to
extract the `traceID` from the context so that the trace ID is returned back to the user for future
analysis of the problem.

Moreover, the idiomatic way to group errors is to put them all in a file `errors.go` in the root package
of the service so they are all easily discoverable in a single location.

### JSON Format

Here the explained JSON format:

```
{
  "code": "specifc_error_code",
  "trace_id": "%s",
  "message": "Sepcific error code message, audience is the end user."
  "details": {
    "key": <value>,
    ...
  },
}
```

| Parameter | Explanation |
|-|-|
| `code` | The unique error representing this error, should be a human readable summary of the error, in snake case. |
| `trace_id` | The unique trace id to further debug that error, the trace id can also be correlated in the logs. |
| `message` | A message describing the error. The audience of the message is the end user. Should be a full sentence ending with a dot. |
| `details` | A key-value map of extra details specific to the error. Usually contains faulty parameters and extra details about the error. |

### Wrap

This package is aware of wrapped errors through the [github.com/pkg/errors](https://github.com/pkg/errors)
package.

As a convenience, the package provides shortcuts to the [Wrap](https://godoc.org/github.com/pkg/errors#Wrap)
and [Wrapf](https://godoc.org/github.com/pkg/errors#Wrapf) of the `pkgErrors`
package so you don't need to import it directly using a custom name like `pkgErrors` (to avoid
conflict with the standard `errors` package).

You can simply use `derr.Wrap` and `derr.Wrapf` to wrap your errors with newer contextual
errors.

### Write Error

The package provides a facility to write any `error` object back to the user. The
`derr.WriteError(ctx context.Context, w http.ResponseWriter, message string, err error)` will correctly
find the `derr.ErrorResponse` out of the `err` parameter received and will output to the user.

If no `derr.ErrorResponse` can be found, automatically, the error is wrapped in an `derr.UnexpectedError`
and returned like that to the user. This behavior will result in a generic error message appearing
for the user.

Note that `WriteError` is also logging the error at the same occasion, removing the burden to
log the error yourself. If the error written back to the user generates a `>= 500` error code,
the `Error` level is used. Otherwise, a `Debug` level is used to log the error.

This error logging will ultimately trickle down to our monitoring infrastructure, so if you use
`WriteError`, be sure to not log it again!


## Contributing

**Issues and PR in this repo related strictly to the derr library.**

Report any protocol-specific issues in their
[respective repositories](https://github.com/streamingfast/streamingfast#protocols)

**Please first refer to the general
[StreamingFast contribution guide](https://github.com/streamingfast/streamingfast/blob/master/CONTRIBUTING.md)**,
if you wish to contribute to this code base.


## License

[Apache 2.0](LICENSE)


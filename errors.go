// Copyright 2019 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package derr

import (
	"context"
	"net/http"
	"net/url"

	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

// Client Errors

// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
func InvalidJSONError(ctx context.Context, err error) *ErrorResponse {
	return HTTPBadRequestError(ctx, err, ErrorCode("invalid_json_error"), "The request is not a valid json.", "errors", map[string]interface{}{
		"source": err.Error(),
	})
}

// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
func MissingBodyError(ctx context.Context) *ErrorResponse {
	return HTTPBadRequestError(ctx, nil, ErrorCode("missing_body_error"), "The request body is missing.")
}

// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
func RequestValidationError(ctx context.Context, errors url.Values) *ErrorResponse {
	return HTTPBadRequestError(ctx, nil, ErrorCode("request_validation_error"), "The request is invalid.", "errors", errors)
}

// Server Errors

// ServiceUnavailableError represents a failure at the transport layer to reach a given micro-service.
// Note that while `serviceName` is required, it's not directly available to final response for now,
// will probably encrypt it into an opaque string if you ever make usage of it
//
// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
func ServiceUnavailableError(ctx context.Context, cause error, serviceName string) *ErrorResponse {
	return HTTPBadGatewayError(ctx, cause, ErrorCode("service_unavailable"), "The service your are requesting is not currently available.")
}

// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
func UnexpectedError(ctx context.Context, cause error) *ErrorResponse {
	return HTTPInternalServerError(ctx, cause, ErrorCode("unexpected_error"), "An unexpected error occurred.")
}

// Generic Request Error Classes (4XX)

var (
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPBadRequestError = newErrorClass(http.StatusBadRequest)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPUnauthorizedError = newErrorClass(http.StatusUnauthorized)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPPaymentRequiredError = newErrorClass(http.StatusPaymentRequired)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPForbiddenError = newErrorClass(http.StatusForbidden)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPNotFoundError = newErrorClass(http.StatusNotFound)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPMethodNotAllowedError = newErrorClass(http.StatusMethodNotAllowed)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPNotAcceptableError = newErrorClass(http.StatusNotAcceptable)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPProxyAuthRequiredError = newErrorClass(http.StatusProxyAuthRequired)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPRequestTimeoutError = newErrorClass(http.StatusRequestTimeout)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPConflictError = newErrorClass(http.StatusConflict)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPGoneError = newErrorClass(http.StatusGone)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPLengthRequiredError = newErrorClass(http.StatusLengthRequired)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPPreconditionFailedError = newErrorClass(http.StatusPreconditionFailed)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPRequestEntityTooLargeError = newErrorClass(http.StatusRequestEntityTooLarge)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPRequestURITooLongError = newErrorClass(http.StatusRequestURITooLong)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPUnsupportedMediaTypeError = newErrorClass(http.StatusUnsupportedMediaType)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPRequestedRangeNotSatisfiableError = newErrorClass(http.StatusRequestedRangeNotSatisfiable)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPExpectationFailedError = newErrorClass(http.StatusExpectationFailed)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPTeapotError = newErrorClass(http.StatusTeapot)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPUnprocessableEntityError = newErrorClass(http.StatusUnprocessableEntity)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPLockedError = newErrorClass(http.StatusLocked)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPFailedDependencyError = newErrorClass(http.StatusFailedDependency)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPUpgradeRequiredError = newErrorClass(http.StatusUpgradeRequired)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPPreconditionRequiredError = newErrorClass(http.StatusPreconditionRequired)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPTooManyRequestsError = newErrorClass(http.StatusTooManyRequests)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPRequestHeaderFieldsTooLargeError = newErrorClass(http.StatusRequestHeaderFieldsTooLarge)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPUnavailableForLegalReasonsError = newErrorClass(http.StatusUnavailableForLegalReasons)
)

// Generic Server Error Classes (5XX)

var (
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPInternalServerError = newErrorClass(http.StatusInternalServerError)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPNotImplementedError = newErrorClass(http.StatusNotImplemented)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPBadGatewayError = newErrorClass(http.StatusBadGateway)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPServiceUnavailableError = newErrorClass(http.StatusServiceUnavailable)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPGatewayTimeoutError = newErrorClass(http.StatusGatewayTimeout)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPHTTPVersionNotSupportedError = newErrorClass(http.StatusHTTPVersionNotSupported)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPVariantAlsoNegotiatesError = newErrorClass(http.StatusVariantAlsoNegotiates)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPInsufficientStorageError = newErrorClass(http.StatusInsufficientStorage)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPLoopDetectedError = newErrorClass(http.StatusLoopDetected)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPNotExtendedError = newErrorClass(http.StatusNotExtended)
	// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
	// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
	HTTPNetworkAuthenticationRequiredError = newErrorClass(http.StatusNetworkAuthenticationRequired)
)

// HTTPErrorFromStatus can be used to programmaticaly route the right status to one of the HTTP error class above
//
// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
func HTTPErrorFromStatus(status int, ctx context.Context, cause error, code ErrorCode, message interface{}, keyvals ...interface{}) *ErrorResponse {
	errorClass := statusToHTTPErrorClass[status]
	if errorClass == nil {
		logError(ctx, "unable to retrieved error class from status, falling back to internal server error", nil, zap.Int("status", status))
		errorClass = HTTPInternalServerError
	}

	return errorClass(ctx, cause, code, message, keyvals...)
}

var statusToHTTPErrorClass = map[int]errorClass{
	http.StatusBadRequest:                   HTTPBadRequestError,
	http.StatusUnauthorized:                 HTTPUnauthorizedError,
	http.StatusPaymentRequired:              HTTPPaymentRequiredError,
	http.StatusForbidden:                    HTTPForbiddenError,
	http.StatusNotFound:                     HTTPNotFoundError,
	http.StatusMethodNotAllowed:             HTTPMethodNotAllowedError,
	http.StatusNotAcceptable:                HTTPNotAcceptableError,
	http.StatusProxyAuthRequired:            HTTPProxyAuthRequiredError,
	http.StatusRequestTimeout:               HTTPRequestTimeoutError,
	http.StatusConflict:                     HTTPConflictError,
	http.StatusGone:                         HTTPGoneError,
	http.StatusLengthRequired:               HTTPLengthRequiredError,
	http.StatusPreconditionFailed:           HTTPPreconditionFailedError,
	http.StatusRequestEntityTooLarge:        HTTPRequestEntityTooLargeError,
	http.StatusRequestURITooLong:            HTTPRequestURITooLongError,
	http.StatusUnsupportedMediaType:         HTTPUnsupportedMediaTypeError,
	http.StatusRequestedRangeNotSatisfiable: HTTPRequestedRangeNotSatisfiableError,
	http.StatusExpectationFailed:            HTTPExpectationFailedError,
	http.StatusTeapot:                       HTTPTeapotError,
	http.StatusUnprocessableEntity:          HTTPUnprocessableEntityError,
	http.StatusLocked:                       HTTPLockedError,
	http.StatusFailedDependency:             HTTPFailedDependencyError,
	http.StatusUpgradeRequired:              HTTPUpgradeRequiredError,
	http.StatusPreconditionRequired:         HTTPPreconditionRequiredError,
	http.StatusTooManyRequests:              HTTPTooManyRequestsError,
	http.StatusRequestHeaderFieldsTooLarge:  HTTPRequestHeaderFieldsTooLargeError,
	http.StatusUnavailableForLegalReasons:   HTTPUnavailableForLegalReasonsError,

	http.StatusInternalServerError:           HTTPInternalServerError,
	http.StatusNotImplemented:                HTTPNotImplementedError,
	http.StatusBadGateway:                    HTTPBadGatewayError,
	http.StatusServiceUnavailable:            HTTPServiceUnavailableError,
	http.StatusGatewayTimeout:                HTTPGatewayTimeoutError,
	http.StatusHTTPVersionNotSupported:       HTTPHTTPVersionNotSupportedError,
	http.StatusVariantAlsoNegotiates:         HTTPVariantAlsoNegotiatesError,
	http.StatusInsufficientStorage:           HTTPInsufficientStorageError,
	http.StatusLoopDetected:                  HTTPLoopDetectedError,
	http.StatusNotExtended:                   HTTPNotExtendedError,
	http.StatusNetworkAuthenticationRequired: HTTPNetworkAuthenticationRequiredError,
}

func logError(ctx context.Context, message string, err error, fields ...zap.Field) {
	logging.Logger(ctx, zlog).Error(message, append(fields, zap.Error(err))...)
}

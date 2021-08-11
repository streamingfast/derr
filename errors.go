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

func InvalidJSONError(ctx context.Context, err error) *ErrorResponse {
	return HTTPBadRequestError(ctx, err, ErrorCode("invalid_json_error"), "The request is not a valid json.", "errors", map[string]interface{}{
		"source": err.Error(),
	})
}

func MissingBodyError(ctx context.Context) *ErrorResponse {
	return HTTPBadRequestError(ctx, nil, ErrorCode("missing_body_error"), "The request body is missing.")
}

func RequestValidationError(ctx context.Context, errors url.Values) *ErrorResponse {
	return HTTPBadRequestError(ctx, nil, ErrorCode("request_validation_error"), "The request is invalid.", "errors", errors)
}

// Server Errors

// ServiceUnavailableError represents a failure at the transport layer to reach a given micro-service.
// Note that while `serviceName` is required, it's not directly available to final response for now,
// will probably encrypt it into an opaque string if you ever make usage of it
func ServiceUnavailableError(ctx context.Context, cause error, serviceName string) *ErrorResponse {
	return HTTPBadGatewayError(ctx, cause, ErrorCode("service_unavailable"), "The service your are requesting is not currently available.")
}

func UnexpectedError(ctx context.Context, cause error) *ErrorResponse {
	return HTTPInternalServerError(ctx, cause, ErrorCode("unexpected_error"), "An unexpected error occurred.")
}

// Generic Request Error Classes (4XX)

var (
	HTTPBadRequestError                   = newErrorClass(http.StatusBadRequest)
	HTTPUnauthorizedError                 = newErrorClass(http.StatusUnauthorized)
	HTTPPaymentRequiredError              = newErrorClass(http.StatusPaymentRequired)
	HTTPForbiddenError                    = newErrorClass(http.StatusForbidden)
	HTTPNotFoundError                     = newErrorClass(http.StatusNotFound)
	HTTPMethodNotAllowedError             = newErrorClass(http.StatusMethodNotAllowed)
	HTTPNotAcceptableError                = newErrorClass(http.StatusNotAcceptable)
	HTTPProxyAuthRequiredError            = newErrorClass(http.StatusProxyAuthRequired)
	HTTPRequestTimeoutError               = newErrorClass(http.StatusRequestTimeout)
	HTTPConflictError                     = newErrorClass(http.StatusConflict)
	HTTPGoneError                         = newErrorClass(http.StatusGone)
	HTTPLengthRequiredError               = newErrorClass(http.StatusLengthRequired)
	HTTPPreconditionFailedError           = newErrorClass(http.StatusPreconditionFailed)
	HTTPRequestEntityTooLargeError        = newErrorClass(http.StatusRequestEntityTooLarge)
	HTTPRequestURITooLongError            = newErrorClass(http.StatusRequestURITooLong)
	HTTPUnsupportedMediaTypeError         = newErrorClass(http.StatusUnsupportedMediaType)
	HTTPRequestedRangeNotSatisfiableError = newErrorClass(http.StatusRequestedRangeNotSatisfiable)
	HTTPExpectationFailedError            = newErrorClass(http.StatusExpectationFailed)
	HTTPTeapotError                       = newErrorClass(http.StatusTeapot)
	HTTPUnprocessableEntityError          = newErrorClass(http.StatusUnprocessableEntity)
	HTTPLockedError                       = newErrorClass(http.StatusLocked)
	HTTPFailedDependencyError             = newErrorClass(http.StatusFailedDependency)
	HTTPUpgradeRequiredError              = newErrorClass(http.StatusUpgradeRequired)
	HTTPPreconditionRequiredError         = newErrorClass(http.StatusPreconditionRequired)
	HTTPTooManyRequestsError              = newErrorClass(http.StatusTooManyRequests)
	HTTPRequestHeaderFieldsTooLargeError  = newErrorClass(http.StatusRequestHeaderFieldsTooLarge)
	HTTPUnavailableForLegalReasonsError   = newErrorClass(http.StatusUnavailableForLegalReasons)
)

// Generic Server Error Classes (5XX)

var (
	HTTPInternalServerError                = newErrorClass(http.StatusInternalServerError)
	HTTPNotImplementedError                = newErrorClass(http.StatusNotImplemented)
	HTTPBadGatewayError                    = newErrorClass(http.StatusBadGateway)
	HTTPServiceUnavailableError            = newErrorClass(http.StatusServiceUnavailable)
	HTTPGatewayTimeoutError                = newErrorClass(http.StatusGatewayTimeout)
	HTTPHTTPVersionNotSupportedError       = newErrorClass(http.StatusHTTPVersionNotSupported)
	HTTPVariantAlsoNegotiatesError         = newErrorClass(http.StatusVariantAlsoNegotiates)
	HTTPInsufficientStorageError           = newErrorClass(http.StatusInsufficientStorage)
	HTTPLoopDetectedError                  = newErrorClass(http.StatusLoopDetected)
	HTTPNotExtendedError                   = newErrorClass(http.StatusNotExtended)
	HTTPNetworkAuthenticationRequiredError = newErrorClass(http.StatusNetworkAuthenticationRequired)
)

// HTTPErrorFromStatus can be used to programmaticaly route the right status to one of the HTTP error class above
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

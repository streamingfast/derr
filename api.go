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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Walk traverse error causes in a top to bottom fashion. Starting from the top
// `err` received, will invoke `processor(err)` with it. If the `processor`
// returns `true`, check if `err` has a cause and continue walking it like
// this recursively. If `processor` return a `non-nil` value, stop walking at
// this point. If `processor` returns an `error` stop walking from there and bubble
// up the error through `Walk` return value.
//
// Returns an `error` if `processor` returned an `error`, `nil` otherwise
func Walk(err error, processor func(err error) (bool, error)) error {
	type causer interface {
		Cause() error
	}

	shouldContinue, childErr := processor(err)
	if !shouldContinue {
		return childErr
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			return nil
		}

		err = cause.Cause()
		shouldContinue, childErr := processor(err)
		if !shouldContinue {
			return childErr
		}
	}

	return nil
}

// FindFirstMatching walks the error(s) stack (causes chain) and return the first
// error matching the `matcher` function received in argument
func FindFirstMatching(err error, matcher func(err error) bool) error {
	var matchedErr error
	Walk(err, func(candidateErr error) (bool, error) {
		if matcher(candidateErr) {
			matchedErr = candidateErr
			return false, nil
		}

		return true, nil
	})

	return matchedErr
}

// HasAny returns `true` if the `err` argument or any of its cause(s) is equal
// to `cause` argument, `false` otherwise.
func HasAny(err error, cause error) bool {
	return FindFirstMatching(err, func(candidateErr error) bool {
		return candidateErr == cause
	}) != nil
}

// ToErrorResponse turns a plain `error` interface into a proper `ErrorResponse`
// object. It does so with the following rules:
//
// - If `err` is already an `ErrorResponse`, turns it into such and returns it.
// - If `err` was wrapped, find the most cause which is an `ErrorResponse` and returns it.
// - If `err` is a status.Status (or one that was wrapped), convert it to an ErrorResponse
// - Otherwise, return an `UnexpectedError` with the cause sets to `err` received.
func ToErrorResponse(ctx context.Context, err error) *ErrorResponse {
	response := FindFirstMatching(err, isErrorResponse)
	if response != nil {
		return response.(*ErrorResponse)
	}

	response = FindFirstMatching(err, isStatusCode)
	if response != nil {
		status, _ := status.FromError(err)
		return convertStatusToErrorResponse(ctx, status)
	}

	return UnexpectedError(ctx, err)
}

func isStatusCode(err error) bool {
	if _, ok := status.FromError(err); ok {
		return true
	}

	return false
}

func isErrorResponse(err error) bool {
	if _, ok := err.(*ErrorResponse); ok {
		return true
	}

	return false
}

func convertStatusToErrorResponse(ctx context.Context, st *status.Status) *ErrorResponse {
	switch st.Code() {
	case codes.InvalidArgument:
		return HTTPBadRequestError(ctx, nil, ErrorCode("request_validation_error"), st.Message())
	case codes.Unavailable:
		return HTTPServiceUnavailableError(ctx, nil, ErrorCode("service_unavailable_error"), "Service Unavailable")
	case codes.NotFound:
		return HTTPNotFoundError(ctx, nil, ErrorCode("not_found_error"), st.Message())
	default:
		return UnexpectedError(ctx, st.Err())
	}
}

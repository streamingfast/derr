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
	"encoding/hex"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	pkgErrors "github.com/pkg/errors"

	"github.com/streamingfast/logging"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
)

func TestWriteError(t *testing.T) {
	traceID := fixedTraceID("00000000000000000000000000000001")
	testContext := func() context.Context {
		spanContext := trace.SpanContext{TraceID: traceID}
		ctx, _ := trace.StartSpanWithRemoteParent(context.Background(), "test", spanContext)

		return ctx
	}

	errInvalidJSON := func(id string) *ErrorResponse { return InvalidJSONError(testContext(), errors.New(id)) }
	errUnexpected := func(cause error) *ErrorResponse { return UnexpectedError(testContext(), cause) }

	tests := []struct {
		name               string
		err                error
		expectedStatusCode int
		expectedBody       string
	}{
		{"plain standard error", errors.New("test"), 500, `
			{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
		`},

		{"plain error response", errInvalidJSON("plain error response"), 400, `
			{"code":"invalid_json_error","trace_id":"%s","details":{"errors":{"source":"plain error response"}},"message":"The request is not a valid json."}
		`},

		{"wrapped error, no cause", pkgErrors.Wrap(nil, "test"), 500, `
			{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
		`},

		{"wrapped error, cause standard error", pkgErrors.Wrap(errors.New("wrapped"), "test"), 500, `
			{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
		`},

		{"wrapped error, cause error response", pkgErrors.Wrap(errInvalidJSON("wrapped response"), "test"), 400, `
			{"code":"invalid_json_error","trace_id":"%s","details":{"errors":{"source":"wrapped response"}},"message":"The request is not a valid json."}
		`},

		{"wrapped error, nested cause error response", pkgErrors.Wrap(pkgErrors.Wrap(errInvalidJSON("wrapped again"), "source"), "nested"), 400, `
			{"code":"invalid_json_error","trace_id":"%s","details":{"errors":{"source":"wrapped again"}},"message":"The request is not a valid json."}
		`},

		{"wrapped error, unexpected with response clause", errUnexpected(errInvalidJSON("json")), 500, `
			{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
		`},

		{"wrapped error, unexpected with wrapped response clause", errUnexpected(pkgErrors.Wrap(errInvalidJSON("json"), "nested1")), 500, `
			{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
		`},

		{"wrapped error, wrapped with unexpected with wrapped response clause", pkgErrors.Wrap(errUnexpected(pkgErrors.Wrap(errInvalidJSON("json"), "nested1")), "deep"), 500, `
		{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
	`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testContext()
			recorder := httptest.NewRecorder()
			logger := zap.NewExample()

			WriteError(logging.WithLogger(ctx, logger), recorder, "prefix", test.err)

			assert.Equal(t, test.expectedStatusCode, recorder.Code)

			expectedBody := strings.TrimSpace(test.expectedBody)
			if strings.Count(expectedBody, "%s") >= 1 {
				expectedBody = fmt.Sprintf(test.expectedBody, traceID)
			}

			assert.JSONEq(t, expectedBody, recorder.Body.String())
		})
	}
}

func fixedTraceID(hexInput string) (out trace.TraceID) {
	rawTraceID, _ := hex.DecodeString(hexInput)
	copy(out[:], rawTraceID)

	return
}

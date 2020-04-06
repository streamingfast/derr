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
	"fmt"
	"strings"
)

type ErrorCode string

// C is a sugar syntax for `derr.ErrorCode("a_string_code")` (sugared to `derr.C("a_string_code")`)
func C(code string) ErrorCode { return ErrorCode(code) }

type ErrorResponse struct {
	Code    ErrorCode              `json:"code"`
	TraceID string                 `json:"trace_id"`
	Status  int                    `json:"-"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Causer  error                  `json:"-"`
}

func (e *ErrorResponse) Cause() error { return e.Causer }

func (e *ErrorResponse) ResponseStatus() int { return e.Status }

func (e *ErrorResponse) Error() string {
	index := 0
	details := make([]string, len(e.Details))

	for k, v := range e.Details {
		details[index] = fmt.Sprintf("%s: %v", k, v)
		index++
	}

	detailsString := ""
	if len(details) > 0 {
		detailsString = fmt.Sprintf(" {%s}", strings.Join(details, ", "))
	}

	causeString := ""
	if e.Causer != nil {
		causeString = fmt.Sprintf(" (%s)", e.Causer)
	}

	return fmt.Sprintf("[%s] %d: %s%s%s", e.Code, e.Status, e.Message, causeString, detailsString)
}

type errorClass func(ctx context.Context, cause error, code ErrorCode, message interface{}, keyvals ...interface{}) *ErrorResponse

func newErrorClass(status int) errorClass {
	return func(ctx context.Context, cause error, code ErrorCode, message interface{}, keyvals ...interface{}) *ErrorResponse {
		var msg string
		switch actual := message.(type) {
		case string:
			msg = actual
		case error:
			msg = actual.Error()
		case fmt.Stringer:
			msg = actual.String()
		default:
			msg = fmt.Sprintf("%v", actual)
		}

		var details map[string]interface{}
		l := len(keyvals)
		if l > 0 {
			details = make(map[string]interface{})
		}

		for i := 0; i < l; i += 2 {
			k := keyvals[i]
			var v interface{} = "MISSING"
			if i+1 < l {
				v = keyvals[i+1]
			}

			details[fmt.Sprintf("%v", k)] = v
		}

		return &ErrorResponse{Code: code, TraceID: traceIDFromContext(ctx), Status: status, Message: msg, Details: details, Causer: cause}
	}
}

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
	"time"

	"github.com/teris-io/shortid"
	"go.opencensus.io/trace"
)

var shortIDGenerator *shortid.Shortid

func init() {
	// A new generator using the default alphabet set
	shortIDGenerator = shortid.MustNew(1, shortid.DefaultABC, uint64(time.Now().UnixNano()))
}

func traceIDFromContext(ctx context.Context) string {
	span := trace.FromContext(ctx)
	if span == nil {
		return shortIDGenerator.MustGenerate()
	}

	spanContext := span.SpanContext()
	return hex.EncodeToString(spanContext.TraceID[:])
}

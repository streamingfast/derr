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
	"encoding/json"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/streamingfast/logging"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WriteError writes the receiver error to HTTP and log it into a Zap logger at the same
// time with the right level based on the actual status code. The `WriteError` handles
// various type for the `err` parameter.
func WriteError(ctx context.Context, w http.ResponseWriter, message string, err error) {
	response := ToErrorResponse(ctx, err)
	zlogger := logging.Logger(ctx, zlog)

	if ctx.Err() != context.Canceled && response.ResponseStatus() >= 500 {
		zlogger.Error(message, zap.Error(err))
	} else {
		zlogger.Debug(message, zap.Error(err))
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(response.ResponseStatus())

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logWriteError(zlogger, "unable to serialize error response", err)
	}
}

func logWriteError(logger *zap.Logger, prefix string, err error) {
	level := zapcore.ErrorLevel
	if IsClientSideNetworkError(err) {
		level = zapcore.DebugLevel
	}

	logger.Check(level, prefix).Write(zap.Error(err))
}

// IsClientSideNetworkError returns wheter the error received is a network error caused by the client side
// that could not be possibily handled correctly on the server side anyway.
func IsClientSideNetworkError(err error) bool {
	netErr, isNetErr := err.(*net.OpError)
	if !isNetErr {
		return false
	}

	syscallErr, isSyscallErr := netErr.Err.(*os.SyscallError)
	if !isSyscallErr {
		return false
	}

	return syscallErr.Err == syscall.ECONNRESET || syscallErr.Err == syscall.EPIPE
}

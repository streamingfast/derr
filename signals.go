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
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"
)

var isShuttingDown *atomic.Bool

func init() {
	isShuttingDown = atomic.NewBool(false)
}

func IsShuttingDown() bool {
	return isShuttingDown.Load()
}

// this is a graceful delay to allow residual traffic sent by the load balancer to be processed
// without returning 500. Once the delay has passed then the service can be shutdown
func SetupSignalHandler(gracefulShutdownDelay time.Duration) <-chan os.Signal {
	outgoingSignals := make(chan os.Signal, 10)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	seen := 0

	go func() {
		for {
			s := <-signals
			switch s {
			case syscall.SIGTERM, syscall.SIGINT:
				seen++

				if seen > 3 {
					zlog.Info("Received termination signal 3 times: Forcing kill")
					zlog.Sync()
					os.Exit(1)
				}

				if !IsShuttingDown() {
					zlog.Info("Received termination signal... Ctrl+C multiple times to force kill", zap.String("signal", s.String()))
					v := true
					isShuttingDown.Store(v)
					go time.AfterFunc(gracefulShutdownDelay, func() {
						outgoingSignals <- s
					})
					break
				}
				zlog.Info("Received termination signal twice, shutting down...", zap.String("signal", s.String()))
				outgoingSignals <- s

			}
		}
	}()

	return outgoingSignals
}

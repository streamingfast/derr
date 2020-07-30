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
	"fmt"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Wrap is a shortcut for `pkgErrors.Wrap` (where `pkgErrors` is `github.com/pkg/errors`)
func Wrap(err error, message string) error {
	if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
		sts := se.GRPCStatus().Proto()
		newSts := &spb.Status{
			Code:    sts.Code,
			Message: fmt.Sprintf("%s: %s", message, sts.Message),
			Details: sts.Details,
		}
		return status.ErrorProto(newSts)
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf(message+": %w", err)
}

// Wrapf is a shortcut for `pkgErrors.Wrapf` (where `pkgErrors` is `github.com/pkg/errors`)
func Wrapf(err error, format string, args ...interface{}) error {
	if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
		sts := se.GRPCStatus().Proto()
		newSts := &spb.Status{
			Code:    sts.Code,
			Message: fmt.Sprintf("%s: %s", fmt.Sprintf(format, args...), sts.Message),
			Details: sts.Details,
		}
		return status.ErrorProto(newSts)
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf(fmt.Sprintf(format, args...)+": %w", err)
}

func WrapCode(code codes.Code, err error, message string) error {
	// ici on ajouterait le `Code` précédent dans le `err`qui est un `status.Status` en
	// [PreviousCode].
	panic("unimplemented")
}

func WrapfCode(code codes.Code, err error, format string, args ...interface{}) error {
	panic("unimplemented")
}

func Status(code codes.Code, message string) error {
	s := status.New(code, message)
	addDebugInfo(s)
	return s.Err()
}

func Statusf(code codes.Code, format string, args ...interface{}) error {
	s := status.Newf(code, format, args...)
	addDebugInfo(s)
	return s.Err()
}

func addDebugInfo(s *status.Status) {
	// Eventually, use https://godoc.org/github.com/go-stack/stack#Trace
	// to stuff in the current call stack (minus 2 levels) into a `errdetails.DebugInfo`
	// debug stack.. which will be passed
	// In the `DebugInfo.Detail`, we can put the pod name, and other contextual info, like the
	// binary program name or something.
	return
}

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
	"errors"
	"fmt"
	"strings"
	"testing"

	dedentLib "github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"
)

var errTestFake = errors.New("test")
var errTestFakeOther = errors.New("other")

var testErrNoDeep = errTestFake
var testErrOneDeep = Wrap(errTestFake, "one deep")
var testErrTwoDeep = Wrap(Wrap(errTestFake, "one deep"), "two deep")
var testErrThreeDeep = Wrap(Wrap(Wrap(errTestFake, "one deep"), "two deep"), "three deep")

func Test_Is(t *testing.T) {
	assert.True(t, Is(testErrNoDeep, errTestFake))
	assert.True(t, Is(testErrOneDeep, errTestFake))
	assert.True(t, Is(testErrTwoDeep, errTestFake))
	assert.True(t, Is(testErrThreeDeep, errTestFake))

	assert.False(t, Is(testErrNoDeep, errTestFakeOther))
	assert.False(t, Is(testErrOneDeep, errTestFakeOther))
	assert.False(t, Is(testErrTwoDeep, errTestFakeOther))
	assert.False(t, Is(testErrThreeDeep, errTestFakeOther))
}

func Test_Find(t *testing.T) {
	matcher := func(candidate error) bool {
		return strings.Contains(candidate.Error(), "two deep")
	}

	alwaysMatching := func(candidate error) bool { return true }
	neverMatching := func(candidate error) bool { return false }

	assert.Equal(t, nil, Find(testErrNoDeep, matcher))
	assert.Equal(t, nil, Find(testErrOneDeep, matcher))
	assert.Equal(t, testErrTwoDeep, Find(testErrTwoDeep, matcher))
	assert.Equal(t, testErrThreeDeep, Find(testErrThreeDeep, matcher))

	assert.Equal(t, testErrNoDeep, Find(testErrNoDeep, alwaysMatching))
	assert.Equal(t, testErrOneDeep, Find(testErrOneDeep, alwaysMatching))
	assert.Equal(t, testErrTwoDeep, Find(testErrTwoDeep, alwaysMatching))
	assert.Equal(t, testErrThreeDeep, Find(testErrThreeDeep, alwaysMatching))

	assert.Equal(t, nil, Find(testErrNoDeep, neverMatching))
	assert.Equal(t, nil, Find(testErrOneDeep, neverMatching))
	assert.Equal(t, nil, Find(testErrTwoDeep, neverMatching))
	assert.Equal(t, nil, Find(testErrThreeDeep, neverMatching))
}

func TestDebugErrorChain(t *testing.T) {
	tests := []struct {
		name string
		on   error
		want string
	}{
		{"nil error", nil, dedent(`<nil>`)},
		{"single error", errors.New("end"), dedent(`
			*errors.errorString | end
		`)},
		{"multi errors", fmt.Errorf("root: %w", fmt.Errorf("middle: %w", errors.New("end"))), dedent(`
			*fmt.wrapError | root: middle: end
			*fmt.wrapError | middle: end
			*errors.errorString | end
		`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, DebugErrorChain(tt.on))
		})
	}
}

func dedent(input string) string {
	return strings.TrimSpace(dedentLib.Dedent(input))
}

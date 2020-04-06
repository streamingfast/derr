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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errTestFake = errors.New("test")
var errTestFakeOther = errors.New("other")

var testErrNoDeep = errTestFake
var testErrOneDeep = Wrap(errTestFake, "one deep")
var testErrTwoDeep = Wrap(Wrap(errTestFake, "one deep"), "two deep")
var testErrThreeDeep = Wrap(Wrap(Wrap(errTestFake, "one deep"), "two deep"), "three deep")

func Test_HasAny(t *testing.T) {
	assert.True(t, HasAny(testErrNoDeep, errTestFake))
	assert.True(t, HasAny(testErrOneDeep, errTestFake))
	assert.True(t, HasAny(testErrTwoDeep, errTestFake))
	assert.True(t, HasAny(testErrThreeDeep, errTestFake))

	assert.False(t, HasAny(testErrNoDeep, errTestFakeOther))
	assert.False(t, HasAny(testErrOneDeep, errTestFakeOther))
	assert.False(t, HasAny(testErrTwoDeep, errTestFakeOther))
	assert.False(t, HasAny(testErrThreeDeep, errTestFakeOther))
}

func Test_FindFirstMatching(t *testing.T) {
	matcher := func(candidate error) bool {
		return strings.Contains(candidate.Error(), "two deep")
	}

	alwaysMatching := func(candidate error) bool { return true }
	neverMatching := func(candidate error) bool { return false }

	assert.Equal(t, nil, FindFirstMatching(testErrNoDeep, matcher))
	assert.Equal(t, nil, FindFirstMatching(testErrOneDeep, matcher))
	assert.Equal(t, testErrTwoDeep, FindFirstMatching(testErrTwoDeep, matcher))
	assert.Equal(t, testErrThreeDeep, FindFirstMatching(testErrThreeDeep, matcher))

	assert.Equal(t, testErrNoDeep, FindFirstMatching(testErrNoDeep, alwaysMatching))
	assert.Equal(t, testErrOneDeep, FindFirstMatching(testErrOneDeep, alwaysMatching))
	assert.Equal(t, testErrTwoDeep, FindFirstMatching(testErrTwoDeep, alwaysMatching))
	assert.Equal(t, testErrThreeDeep, FindFirstMatching(testErrThreeDeep, alwaysMatching))

	assert.Equal(t, nil, FindFirstMatching(testErrNoDeep, neverMatching))
	assert.Equal(t, nil, FindFirstMatching(testErrOneDeep, neverMatching))
	assert.Equal(t, nil, FindFirstMatching(testErrTwoDeep, neverMatching))
	assert.Equal(t, nil, FindFirstMatching(testErrThreeDeep, neverMatching))
}

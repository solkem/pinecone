// Copyright 2021 The Matrix.org Foundation C.I.C.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"bytes"
	"testing"
)

func TestSwitchPorts(t *testing.T) {
	expected := []byte{0, 5, 1, 2, 3, 159, 32}
	input := SwitchPorts{1, 2, 3, 4000}
	b, err := input.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, expected) {
		t.Fatalf("MarshalBinary produced %v, expected %v", b, expected)
	}
	var output SwitchPorts
	if _, err := output.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if !input.EqualTo(output) {
		t.Fatalf("Expected %v, got %v", input, output)
	}
}

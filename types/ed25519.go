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
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

type PublicKey [ed25519.PublicKeySize]byte
type PrivateKey [ed25519.PrivateKeySize]byte
type Signature [ed25519.SignatureSize]byte

func (a PublicKey) EqualTo(b PublicKey) bool {
	return bytes.Equal(a[:], b[:])
}

func (a PublicKey) CompareTo(b PublicKey) int {
	return bytes.Compare(a[:], b[:])
}

func (a PublicKey) String() string {
	return fmt.Sprintf("%v", hex.EncodeToString(a[:]))
}

func (a PublicKey) Network() string {
	return "ed25519"
}

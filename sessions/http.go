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

package sessions

import (
	"net/http"
	"time"
)

type HTTP struct {
	httpServer    *http.Server
	httpMux       *http.ServeMux
	httpTransport *http.Transport
	httpClient    *http.Client
}

func (q *Sessions) HTTP() *HTTP {
	h := &HTTP{
		httpServer: &http.Server{
			IdleTimeout: time.Second * 30,
		},
		httpMux: &http.ServeMux{},
		httpTransport: &http.Transport{
			Dial:           q.Dial,
			DialTLS:        q.DialTLS,
			DialContext:    q.DialContext,
			DialTLSContext: q.DialTLSContext,
		},
	}

	h.httpServer.Handler = h.httpMux
	h.httpClient = &http.Client{
		Transport: h.httpTransport,
	}

	go h.httpServer.Serve(q) // nolint:errcheck
	return h
}

func (h *HTTP) Mux() *http.ServeMux {
	return h.httpMux
}

func (h *HTTP) Client() *http.Client {
	return h.httpClient
}

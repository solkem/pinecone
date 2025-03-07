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
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"time"

	"github.com/matrix-org/pinecone/router"
	"github.com/neilalexander/utp"
)

type Sessions struct {
	r            *router.Router
	log          *log.Logger        // logger
	context      context.Context    // router context
	cancel       context.CancelFunc // shut down the router
	streams      chan net.Conn      // accepted connections
	tlsCert      *tls.Certificate   //
	tlsServerCfg *tls.Config        //
	utpSocket    *utp.Socket        //
}

func NewSessions(log *log.Logger, r *router.Router) *Sessions {
	ctx, cancel := context.WithCancel(context.Background())
	s, err := utp.NewSocketFromPacketConnNoClose(r)
	if err != nil {
		panic(err)
	}
	q := &Sessions{
		r:         r,
		log:       log,
		context:   ctx,
		cancel:    cancel,
		streams:   make(chan net.Conn, 16),
		utpSocket: s,
	}

	q.tlsCert = q.generateTLSCertificate()
	q.tlsServerCfg = &tls.Config{
		Certificates: []tls.Certificate{*q.tlsCert},
		ClientAuth:   tls.RequireAnyClientCert,
	}

	go q.listener()
	return q
}

func (q *Sessions) Sessions() []ed25519.PublicKey {
	var sessions []ed25519.PublicKey
	//for _, s := range q.sessions {
	//sessions = append(sessions, s.RemoteAddr())
	/*
		if certs := s.ConnectionState().PeerCertificates; len(certs) > 0 {
			sessions = append(sessions, certs[0].PublicKey.(ed25519.PublicKey))
		}
	*/
	//}
	return sessions
}

func (q *Sessions) generateTLSCertificate() *tls.Certificate {
	private, public := q.r.PrivateKey(), q.r.PublicKey()
	id := hex.EncodeToString(public[:])

	template := x509.Certificate{
		Subject: pkix.Name{
			CommonName: id,
		},
		SerialNumber: big.NewInt(1),
		NotAfter:     time.Now().Add(time.Hour * 24 * 365),
		DNSNames:     []string{id},
	}

	certDER, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		ed25519.PublicKey(public[:]),
		ed25519.PrivateKey(private[:]),
	)
	if err != nil {
		panic(err)
	}
	privateKey, err := x509.MarshalPKCS8PrivateKey(ed25519.PrivateKey(private[:]))
	if err != nil {
		panic(err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKey})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}

	return &tlsCert
}

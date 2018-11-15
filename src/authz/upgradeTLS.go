/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
package authz

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"net/http"
)

// a listener that can listen to encrypted and unencrypted connections
type dualListener struct {
	net.Listener
	TLSConfig *tls.Config
}

type wrappedConnection struct {
	io.Reader
	net.Conn
}

func ListenAndUpgradeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	netListener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer netListener.Close()

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	listener := &dualListener{netListener, &tls.Config{
		Certificates: []tls.Certificate{cert},
	}}

	// redirect from http to https
	handleRedirect := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			r.URL.Host = r.Host
			r.URL.Scheme = "https"
			http.Redirect(w, r, r.URL.String(), http.StatusFound)
		} else {
			handler.ServeHTTP(w, r)
		}
	})

	err = http.Serve(listener, handleRedirect)
	if err != nil {
		return err
	}

	return nil
}

// override dualListener.Accept to handle tls and non-tls connections
func (l *dualListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)
	wrapper := &wrappedConnection{reader, conn}

	firstByte, err := reader.Peek(1)
	if err == nil && firstByte[0] == 0x16 {
		// the first byte of a TLS handshake is 0x16
		return tls.Server(wrapper, l.TLSConfig), nil
	}

	return wrapper, nil
}

func (c *wrappedConnection) Read(b []byte) (n int, err error) {
	return c.Reader.Read(b)
}

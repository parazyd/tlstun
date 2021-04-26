// Copyright (c) 2019-2021 Ivan J. <parazyd@dyne.org>
//
// This file is part of tlstun
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

var (
	cacert  = flag.String("cacert", "ca.pem", "Path for CA certificate file")
	cert    = flag.String("cert", "server.pem", "Path for Certificate file")
	key     = flag.String("key", "server-key.pem", "Path for Key file")
	listen  = flag.String("listen", "127.0.0.1:7443", "Listen address")
	forward = flag.String("forward", "127.0.0.1:72", "Forward address")
	fwtls   = flag.Bool("forwardtls", false, "Forward using TLS")
	client  = flag.Bool("verifyclient", false, "Do client verification")
	verbose = flag.Bool("verbose", false, "Verbose mode")
	notls   = flag.Bool("notls", false, "Disable TLS and tunnel plain TCP")
	tlsver  = flag.Int("tlsver", 13, "TLS version to use (11, 12, 13)")
)

func tlsConfig(cert, key string) (*tls.Config, error) {
	creds, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	tlscfg := &tls.Config{Certificates: []tls.Certificate{creds}}

	if *client {
		certpool, _ := x509.SystemCertPool()
		if certpool == nil {
			certpool = x509.NewCertPool()
		}
		pem, err := ioutil.ReadFile(*cacert)
		if err != nil {
			return nil, err
		}
		if !certpool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("can't parse client certificate authority")
		}
		tlscfg.ClientCAs = certpool
		tlscfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	switch *tlsver {
	case 11:
		tlscfg.MinVersion = tls.VersionTLS11
	case 12:
		tlscfg.MinVersion = tls.VersionTLS12
	case 13:
		tlscfg.MinVersion = tls.VersionTLS13
	default:
		log.Fatal("Unsupported TLS version:", *tlsver)
	}

	return tlscfg, nil
}

func tunnel(conn net.Conn, tlsCfg *tls.Config) {
	var client net.Conn
	var err error

	if *fwtls {
		client, err = tls.Dial("tcp", *forward, tlsCfg)
	} else {
		client, err = net.Dial("tcp", *forward)
	}

	if err != nil {
		log.Fatal(err)
	}

	if *verbose {
		log.Printf("Connected to localhost for %s\n", conn.RemoteAddr())
	}

	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		if *verbose {
			defer log.Printf("Closed connection from %s\n", conn.RemoteAddr())
		}
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

func server(tlsCfg *tls.Config) (net.Listener, error) {
	t, err := net.Listen("tcp", *listen)
	if err != nil {
		return nil, err
	}

	if *notls {
		return t, nil
	}

	return tls.NewListener(t, tlsCfg), nil
}

func main() {
	flag.Parse()

	var tlsCfg *tls.Config
	var err error

	if *notls {
		tlsCfg = nil
	} else {
		tlsCfg, err = tlsConfig(*cert, *key)
		if err != nil {
			log.Fatal(err)
		}
	}

	tcpsock, err := server(tlsCfg)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := tcpsock.Accept()
		if err != nil {
			log.Fatal(err)
		}
		if *verbose {
			log.Printf("Accepted connection from %s\n", conn.RemoteAddr())
		}
		go tunnel(conn, tlsCfg)
	}
}

package main

/*
 * Copyright (c) 2019 Dyne.org Foundation
 * tlstun is written and maintained by Ivan Jelincic <parazyd@dyne.org>
 *
 * This file is part of tlstun
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
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
	client  = flag.Bool("verifyclient", false, "Do client verification")
	verbose = flag.Bool("verbose", false, "Verbose mode")
	tlsver  = flag.Int("tlsver", 13, "TLS version to use (11, 12, 13)")
)

func tlsConfig(cert, key string) (*tls.Config, error) {
	creds, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	tlscfg := &tls.Config{Certificates: []tls.Certificate{creds}}

	if *client {
		certpool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(*cacert)
		if err != nil {
			return nil, err
		}
		if !certpool.AppendCertsFromPEM(pem) {
			return nil, errors.New("can not parse client certificate authority")
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

func tunnel(conn net.Conn) {
	client, err := net.Dial("tcp", *forward)
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

func server() (net.Listener, error) {
	t, err := net.Listen("tcp", *listen)
	if err != nil {
		return nil, err
	}

	cfg, err := tlsConfig(*cert, *key)
	if err != nil {
		return nil, err
	}

	return tls.NewListener(t, cfg), nil
}

func main() {
	flag.Parse()

	tcpsock, err := server()
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
		go tunnel(conn)
	}
}

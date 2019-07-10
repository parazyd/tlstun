tlstun
======

tlstun is a simple Go program that will add TLS support for your
programs that do not have it.

It simply proxies from one TLS-listening host:port to another plaintext
host:port.


Installation
------------

```
$ go get github.com/parazyd/tlstun
```

Make sure you generate or acquire a TLS certificate keypair to use with
tlstun.


Usage
-----

```
Usage of ./tlstun:
  -cacert string
        Path for CA certificate file (default "ca.pem")
  -cert string
        Path for Certificate file (default "server.pem")
  -forward string
        Forward address (default "127.0.0.1:72")
  -key string
        Path for Key file (default "server-key.pem")
  -listen string
        Listen address (default "127.0.0.1:7443")
  -tlsver int
        TLS version to use (11, 12, 13) (default 13)
  -verbose
        Verbose mode
  -verifyclient
        Do client verification
```

tlstun supports two different ways of multiplexing, one being normal TLS
proxying, and the other being TLS proxying with client certificate
authentication.


### Without client verification

Start tlstun with `-cert` and `-key`, and it will simply provide a TLS
forward to its destination with the given TLS certificate.


### With client verification

With client verification, start tlstun with `-cacert`, `-cert`, `-key`,
and `-verifyclient` and it will do client certificate verification. This
means it will only allow access from clients providing a certificate
signed by the CA certificate that is being loaded/used with tlstun on
startup with `-cacert`.

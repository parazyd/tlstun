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
  -c string
        Path for Certificate file (default "server.pem")
  -f string
        Forward address (default "127.0.0.1:72")
  -k string
        Path for Key file (default "server-key.pem")
  -l string
        Listen address (default "127.0.0.1:7443")
  -v    Verbose mode
```

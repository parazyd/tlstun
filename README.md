tlstun
======

tlstun is a simple Go program that will add TLS support for your
programs that do not have it.

It simply proxies from one TLS-listening host:port to another plaintext
host:port.


Usage
-----

```
Usage of ./tlstun:
  -cert string
        Path for Certificate file (default "server.pem")
  -forward string
        Forward address (default "127.0.0.1:72")
  -key string
        Path for Key file (default "server-key.pem")
  -listen string
        Listen address (default "127.0.0.1:7443")
  -verbose
        Verbose mode

```

#!/bin/sh
keyfile="server-key.pem"
certfile="server.pem"

openssl genrsa -out "$keyfile" 4096
openssl req -new -x509 -key "$keyfile" -out "$certfile" -days 365

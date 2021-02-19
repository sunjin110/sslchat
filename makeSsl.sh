#!/bin/sh

# golangはSSL/TLSの機能は自前で作っているらしい
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

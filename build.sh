#!/bin/bash

go build -v -ldflags="-w -s" -o tcpproxy ./cmd/tcpproxy/
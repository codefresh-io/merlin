#!/bin/bash
set -e
OUTFILE=/usr/local/bin/merlindev
go build -o $OUTFILE main.go

chmod +x $OUTFILE
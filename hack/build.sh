#!/bin/bash
set -e
OUTFILE=/usr/local/bin/merlin
go build -o $OUTFILE main.go

chmod +x $OUTFILE
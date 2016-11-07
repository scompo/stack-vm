#!/bin/bash
REPORT_DIR=".reports"
go test -coverprofile $REPORT_DIR/coverage.out
go tool cover -html=$REPORT_DIR/coverage.out

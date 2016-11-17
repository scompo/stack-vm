#!/bin/bash

# coverage-reports.sh

# this file uses gocov and gocov-html to generate an html report and opens it
# in firefox.

gocov test ./... | gocov-html > .reports/cover.html
firefox ./.reports/cover.html

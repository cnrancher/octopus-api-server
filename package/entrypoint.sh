#!/bin/bash
set -e

exec tini -- edge-api-server --http-listen-port=8080 --https-listen-port=8443 "${@}"

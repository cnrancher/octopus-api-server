#!/bin/bash
set -e

exec tini -- octopus-api-server --http-listen-port=8080 --https-listen-port=8443 "${@}"

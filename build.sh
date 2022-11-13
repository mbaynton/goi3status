#!/bin/bash

set -euo pipefail

. $HOME/.config/barista/build_secrets.sh

go build -ldflags="-X 'main.GithubClientId=$GITHUB_CLIENT_ID' -X 'main.GithubClientSecret=$GITHUB_CLIENT_SECRET' -X 'main.OwmApiKey=$OWM_API_KEY' -X 'main.GoogleClientId=$GOOGLE_CLIENT_ID' -X 'main.GoogleClientSecret=$GOOGLE_CLIENT_SECRET'" \
 "$1"

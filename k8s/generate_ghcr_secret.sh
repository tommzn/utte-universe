#!/bin/bash

kubectl create secret ghcr-credentials  -n utte \
  --docker-server=ghcr.io \
  --docker-username=USERNAME \
  --docker-password=YOUR_PAT \
  --docker-email=unused@example.com


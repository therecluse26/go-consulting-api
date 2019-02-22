#!/bin/sh

mkdir -p /var/www/go_consulting_api

docker pull quay.io/letsencrypt/letsencrypt

docker run -it --rm --name letsencrypt \
  -v "/etc/letsencrypt:/etc/letsencrypt" \
  -v "/var/lib/letsencrypt:/var/lib/letsencrypt" \
  --volumes-from {{ proxy_docker_container }} \
  quay.io/letsencrypt/letsencrypt \
  certonly \
  --webroot \
  --webroot-path /var/www/go_consulting_api \
  --agree-tos \
  --renew-by-default \
  -d http://go_consulting_api.eastus.azurecontainer.io \
  -m {{ user_email }}

docker kill --signal=HUP {{ proxy_docker_container }}
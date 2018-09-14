#!/bin/sh

mkdir -p /var/www/fortisureapi

docker pull quay.io/letsencrypt/letsencrypt

docker run -it --rm --name letsencrypt \
  -v "/etc/letsencrypt:/etc/letsencrypt" \
  -v "/var/lib/letsencrypt:/var/lib/letsencrypt" \
  --volumes-from {{ proxy_docker_container }} \
  quay.io/letsencrypt/letsencrypt \
  certonly \
  --webroot \
  --webroot-path /var/www/fortisureapi \
  --agree-tos \
  --renew-by-default \
  -d http://fortisureapi.eastus.azurecontainer.io \
  -m brad.magyar@fortisureit.com

docker kill --signal=HUP {{ proxy_docker_container }}
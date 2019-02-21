#!/bin/bash

# iron/go is the alpine image with only ca-certificates added
FROM iron/go

WORKDIR /app

# Now just add the binary
COPY dist/go_consulting_api /app/

COPY nginx/go_consulting_api-nginx.conf /etc/nginx/conf.d/go_consulting_api-nginx.conf

CMD /bin/sh -c "apk update && apk add nginx certbot memcached"

CMD /bin/sh -c "wget https://github.com/jwilder/docker-gen/releases/download/0.7.3/docker-gen-alpine-linux-amd64-0.7.3.tar.gz nginx/docker-gen-alpine-linux-amd64-0.7.3.tar.gz && tar xvzf nginx/docker-gen-alpine-linux-amd64-0.7.3.tar.gz"

CMD /bin/sh -c "wget https://dl.eff.org/certbot-auto && chmod a+x certbot-auto && ./certbot-auto --nginx"

CMD /bin/sh -c "export CFGJSON=$(echo $gocfg64 | base64 -d)"

CMD /bin/sh -c "memcached -u root -d -p 11211"

# CMD /bin/sh -c "echo $CFGJSON"

ENTRYPOINT ["./go_consulting_api"]

EXPOSE 80/tcp
EXPOSE 80/udp
EXPOSE 443/tcp
EXPOSE 443/udp
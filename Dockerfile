#!/bin/bash

# iron/go is the alpine image with only ca-certificates added
FROM iron/go

WORKDIR /app

# Now just add the binary
ADD dist/fortisure-api /app/

ENTRYPOINT ["./fortisure-api"]

CMD echo "testing 123"

CMD echo $config64

EXPOSE 80/tcp
EXPOSE 80/udp
EXPOSE 443/tcp
EXPOSE 443/udp
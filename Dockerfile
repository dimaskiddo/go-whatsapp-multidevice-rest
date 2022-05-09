# Builder Image
# ---------------------------------------------------
FROM dimaskiddo/debian-buster:go-1.17 AS go-builder

WORKDIR /usr/src/app

COPY . ./

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -o main cmd/main/main.go


# Final Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:base
MAINTAINER Dimas Restu Hidayanto <dimas.restu@student.upi.edu>

ARG SERVICE_NAME="go-whatsapp-multidevice-rest"

ENV PATH $PATH:/usr/app/${SERVICE_NAME}

WORKDIR /usr/app/${SERVICE_NAME}

RUN mkdir -p dbs \
    && chmod 755 dbs
COPY --from=go-builder /usr/src/app/.env.default ./.env
COPY --from=go-builder /usr/src/app/main ./main

EXPOSE 3000

VOLUME ["/usr/app/${SERVICE_NAME}/dbs"]
CMD ["main"]

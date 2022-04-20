# Builder Image
# ---------------------------------------------------
FROM dimaskiddo/alpine:go-1.17 AS go-builder

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

COPY --from=go-builder /usr/src/app/.env.production ./.env
COPY --from=go-builder /usr/src/app/main ./main

EXPOSE 3000

CMD ["main"]

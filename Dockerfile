FROM golang:1.16-alpine AS local-dev
RUN apk update && apk add --no-cache g++ gcc make bash ca-certificates
ENV TZ=Asia/Tokyo
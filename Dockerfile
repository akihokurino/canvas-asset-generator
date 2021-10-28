FROM golang:1.15-alpine AS local-dev
RUN apk update && apk add --no-cache g++ gcc make bash ca-certificates
ENV TZ=Asia/Tokyo
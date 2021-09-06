FROM golang:1.16-alpine AS build-env

ENV PACKAGES make

RUN apk add --no-cache $PACKAGES

WORKDIR /go/src/github.com/jtieri/HabbGo

# Copy the files from host
COPY . .

RUN make install

FROM alpine:edge

RUN apk add --no-cache ca-certificates

ENV HABBGO /habbgo

RUN addgroup habbgo && adduser -S -G habbgo habbgo -h "$HABBGO"

USER habbgo

# Define working directory
WORKDIR $HABBGO

# Copy binary from BUILD
COPY --from=build-env /go/bin/habbgo /usr/bin/habbgo

ENTRYPOINT ["/usr/bin/habbgo"]
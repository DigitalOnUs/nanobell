# Build code
FROM golang:latest AS COMPILATION
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o nanobell .

# running stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /bellator
COPY --from=COMPILATION /build/nanobell .
CMD ["/bellator/nanobell"]
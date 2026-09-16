FROM golang:1.27.1 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
COPY cmd ./cmd
COPY internal ./internal
RUN go test ./...
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/tunnexo .

FROM scratch

COPY LICENSE /LICENSE

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /out/tunnexo /tunnexo
USER 65532:65532
WORKDIR /app
ENTRYPOINT ["/tunnexo"]
CMD ["--help"]

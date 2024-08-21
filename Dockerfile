FROM docker.io/golang:1.22 as builder

COPY . /opt
WORKDIR /opt

RUN env CGO_ENABLED=0 go build -o /opt/bin/near_exporter ./near_exporter/cmd/near_exporter

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /opt/bin/near_exporter /

CMD ["/near_exporter"]

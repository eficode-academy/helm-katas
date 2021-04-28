FROM golang:1.16 AS builder
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -a -o sentence-regex .

FROM scratch
COPY --from=builder /build/sentence-regex /sentence-regex
ENTRYPOINT ["./sentence-regex"]

FROM golang:1.17.7-alpine as builder
RUN apk add git
WORKDIR /wled-build

RUN git clone https://github.com/widapro/pixel-led-sender.git .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /wled-build/pixel-led-sender

FROM alpine
COPY --from=builder /wled-build/pixel-led-sender /pixel-led-sender

CMD ["/pixel-led-sender"]
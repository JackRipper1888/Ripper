FROM golang:latest AS builder

WORKDIR /Users/jackripper/project/tools
COPY go.mod .
COPY go.sum .
RUN go mod tidy
COPY . .

RUN apt-get update && apt-get install -y xz-utils  && rm -rf /var/lib/apt/lists/*
ADD https://github.com/upx/upx/releases/download/v3.95/upx-3.95-amd64_linux.tar.xz /usr/local
RUN xz -d -c /usr/local/upx-3.95-amd64_linux.tar.xz | tar -xOf - upx-3.95-amd64_linux/upx > /bin/upx && chmod a+x /bin/upx

WORKDIR /Users/jackripper/project/Ripper
COPY go.mod .
COPY go.sum .
RUN go mod tidy
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix  cgo  -o Ripper && chmod a+rx Ripper
RUN strip --strip-unneeded Ripper
RUN upx --ultra-brute Ripper

FROM alpine:latest
EXPOSE 5003:5003/udp
WORKDIR /root
COPY --from=builder /Users/jackripper/project/Ripper/conf ./conf
COPY --from=builder /Users/jackripper/project/Ripper/Ripper .
RUN mkdir log/ 

CMD ["./Ripper"]
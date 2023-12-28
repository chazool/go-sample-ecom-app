
FROM golang:1.20.1-bullseye 

RUN apt-get clean
RUN apt-get update && apt-get install -y ca-certificates openssl git tzdata
ARG cert_location=/usr/local/share/ca-certificates

RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.ctr
RUN openssl s_client -showcerts -connect gitlab.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/gitlab.ctr
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/proxy.golang.ctr
RUN openssl s_client -showcerts -connect gopkg.in:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/gopkg.ctr
RUN openssl s_client -showcerts -connect storage.googleapis.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/storage.googleapis.ctr
RUN openssl s_client -showcerts -connect sum.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/sum.golang.ctr

RUN update-ca-certificates

WORKDIR /app

#COPY ./project/* .
COPY . .

RUN ls

RUN go mod tidy

RUN go build  -a cmd/main.go

EXPOSE 8181

CMD ["./main"]

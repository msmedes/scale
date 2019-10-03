FROM golang:1.11.5 as builder
MAINTAINER msmedes
ENV GO111MODULE=on
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ARG PKG=/go/src/github.com/msmedes/scale
RUN mkdir -p ${PKG}
WORKDIR ${PKG}
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/scale ./cmd/scale

FROM scratch
COPY --from=builder /go/bin/scale /go/bin/
EXPOSE 3000
CMD ["/go/bin/scale"]

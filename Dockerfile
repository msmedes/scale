FROM golang:1.11.5
MAINTAINER msmedes
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ARG PKG=/go/src/github.com/msmedes/scale
RUN mkdir -p ${PKG}
WORKDIR ${PKG}
COPY . ${PKG}
RUN go build -o /go/bin/scale ./cmd/scale
EXPOSE 3000
CMD ["/go/bin/scale"]

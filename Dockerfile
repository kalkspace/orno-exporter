FROM golang:1.13-alpine AS build

COPY . /src
WORKDIR /src

RUN go mod download
ENV CGO_ENABLED=0
RUN go build -o orno-exporter

FROM alpine
COPY --from=build /src/orno-exporter /bin/orno-exporter

ENTRYPOINT [ "/bin/orno-exporter" ]

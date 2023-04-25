FROM golang:1 AS build-stage

ARG VERSION

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /bookprint -ldflags="-X 'main.Version=$VERSION'" -trimpath ./cmd/...

FROM gcr.io/distroless/base-debian11 AS release-stage

WORKDIR /mnt/local
VOLUME /mnt/local

COPY --from=build-stage /bookprint /bookprint

ENTRYPOINT ["/bookprint"]
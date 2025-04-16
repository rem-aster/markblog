FROM golang:1.24.2 AS build
WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /app app
COPY web/dist web/dist
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
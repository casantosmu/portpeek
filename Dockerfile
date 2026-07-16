FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /portpeek .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /portpeek /portpeek
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/portpeek"]

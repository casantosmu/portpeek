FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /portpeek ./cmd/portpeek

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /portpeek /portpeek
EXPOSE 8080
ENTRYPOINT ["/portpeek"]

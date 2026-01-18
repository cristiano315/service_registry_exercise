# STAGE 1: Common builder
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

# STAGE 2: Registry
FROM builder AS registry_builder
RUN go build -o /registry_bin ./registry/main.go

# STAGE 3: Weather Service and Counter Service
FROM builder AS weather_builder
RUN go build -o /weather_bin ./weather-service/main.go

FROM builder AS counter_builder
RUN go build -o /counter_bin ./counter-service/main.go

# STAGE 4: Client
FROM builder AS client_builder
RUN go build -o /client_bin ./client/main.go

# --- FINAL IMAGES ---

FROM alpine:latest AS registry
COPY --from=registry_builder /registry_bin /app/
EXPOSE 1234
CMD ["/app/registry_bin"]

FROM alpine:latest AS weather
COPY --from=weather_builder /weather_bin /app/
CMD ["/app/weather_bin"]

FROM alpine:latest AS counter
COPY --from=counter_builder /counter_bin /app/
CMD ["/app/counter_bin"]

FROM alpine:latest AS client
COPY --from=client_builder /client_bin /app/
CMD ["/app/client_bin"]
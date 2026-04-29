# Stage 1: Build frontend
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.23-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/internal/web/dist/ internal/web/dist/
RUN CGO_ENABLED=0 go build -o /aegis ./cmd/apigateway

# Stage 3: Final image
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=backend /aegis /usr/local/bin/aegis
COPY config.yaml /etc/aegis/config.yaml
EXPOSE 8080
ENTRYPOINT ["aegis"]

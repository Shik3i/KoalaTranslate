# Stage 1: Build the SvelteKit frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.26.9-alpine AS backend-builder
WORKDIR /app/server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
# Copy built frontend assets to the build folder for embedding
COPY --from=frontend-builder /app/web/build ./cmd/server/build
# Compile static CGO-free binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server_bin ./cmd/server/main.go

# Stage 3: Runner
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /app/server_bin ./koalatranslate

EXPOSE 8080
ENV PORT=8080
ENV DATABASE_PATH=/data/koalatranslate.db
ENV SESSION_SECRET=default_session_secret_change_me_in_prod

CMD ["./koalatranslate"]

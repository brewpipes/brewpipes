# Stage 1: Build frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /app/service/www
# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate
# Copy package files and install dependencies
COPY service/www/package.json service/www/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
# Copy source and build
COPY service/www/ ./
RUN pnpm build

# Stage 2: Build Go binary
FROM golang:1.25.5-alpine AS go-builder
WORKDIR /app
# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download
# Copy source code
COPY . .
# Copy built frontend assets
COPY --from=frontend-builder /app/service/www/dist ./service/www/dist
# Build the binary
RUN CGO_ENABLED=0 go build -o brewpipes ./cmd/monolith

# Stage 3: Runtime
FROM alpine:latest
WORKDIR /app
# Copy the binary (migrations are embedded, no external files needed)
COPY --from=go-builder /app/brewpipes .
EXPOSE 8080
ENTRYPOINT ["./brewpipes"]

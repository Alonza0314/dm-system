# -------------------------------
# Stage 1: Build Go binary
# -------------------------------
FROM golang:1.26.2 AS builder
WORKDIR /src/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
RUN go build -o /out/dm .

# -------------------------------
# Stage 2: Build the frontend
# -------------------------------
FROM node:20 AS consolebuilder
WORKDIR /src/frontend

COPY frontend/ .
RUN yarn install && yarn build

# -------------------------------
# Stage 3: Build the final image
# -------------------------------
FROM debian:bookworm-slim
WORKDIR /dm

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/dm ./dm
COPY --from=consolebuilder /src/frontend/dist ./frontend

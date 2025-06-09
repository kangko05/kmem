# frontend
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/yarn.lock ./
RUN yarn install
COPY frontend/ ./
RUN yarn build

# backend
FROM golang:1.24-alpine AS backend
WORKDIR /app
RUN apk add --no-cache gcc musl-dev 
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .


# final
FROM alpine:latest
RUN apk --no-cache add ffmpeg ca-certificates
WORKDIR /app
COPY config.yml.example ./config.yml

RUN echo "# Environment variables for development" > .env

COPY --from=backend /app/main .
COPY --from=frontend /app/frontend/dist ./frontend/dist
EXPOSE 8000

CMD ["./main"]

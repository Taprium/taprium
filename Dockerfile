FROM node:22 AS fe-builder
COPY frontend/ /app
WORKDIR /app
ENV PUBLIC_PB_ADDR="/"
COPY frontend/. .
RUN npm ci
RUN npm run build

FROM golang:alpine AS be-builder
WORKDIR /src
COPY backend/main.go backend/go.mod backend/go.sum  ./
COPY backend/migrations ./migrations
COPY backend/hooks ./hooks
RUN go build -o /bin/pocketbase

FROM alpine:latest

WORKDIR /pb
COPY --from=be-builder /bin/pocketbase ./
COPY --from=fe-builder /app/build ./pb_public

EXPOSE 8080
VOLUME [ "/pb/pb_data" ]

CMD ["/pb/pocketbase", "serve", "--http=0.0.0.0:8080"]
FROM golang:alpine AS builder
COPY . /app
WORKDIR /app/server
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /naive-blockchain .

# Build the React application
FROM node:alpine AS node_builder
COPY --from=builder /app/client ./
RUN npm install
RUN npm run build

# Final stage build, this will be the container
# that we will deploy to production
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
COPY --from=builder /naive-blockchain ./
COPY --from=node_builder /build ./web
RUN chmod +x ./naive-blockchain
EXPOSE 8080
CMD ./naive-blockchain

# Build the React application
FROM node:alpine AS node_builder
RUN mkdir -p /usr/local/app/client
WORKDIR /usr/local/app/client
COPY . .
RUN npm install

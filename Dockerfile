# Golang builder.
FROM golang:1.15-alpine AS golang-builder
WORKDIR /go/src/armeria
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o build/armeria cmd/armeria/main.go

# Nodejs builder.
FROM node:16-alpine AS node-builder
WORKDIR /go/src/armeria
COPY . .
RUN yarn install
RUN yarn build

# Armeria container.
FROM scratch
COPY --from=golang-builder /go/src/armeria/build/armeria /go/bin/armeria
COPY --from=node-builder /go/src/armeria/dist /opt/armeria/client
EXPOSE 8081
ENTRYPOINT ["/go/bin/armeria"]
CMD ["/go/bin/armeria"]
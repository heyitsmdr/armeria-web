FROM golang:1.15-alpine

# Build from source.
COPY . $GOPATH/src/armeria
RUN \
	cd $GOPATH/src/armeria && \
	go build -o $GOPATH/bin/armeria cmd/armeria/main.go

# Remove the source code from the image.
RUN rm -rf $GOPATH/src/armeria

# Entrypoint.
CMD ["$GOPATH/bin/armeria"]
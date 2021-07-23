FROM SCRATCH

# Move server to container.
COPY ./build/armeria /go/bin/armeria

# Move client to container.
COPY ./dist /opt/armeria/client

# Expose port 8081.
EXPOSE 8081

# Entrypoint.
CMD ["/go/bin/armeria"]
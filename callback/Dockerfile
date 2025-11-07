# Run this image from the parent directory using "-f"

# Using alpine for size optimisation
FROM golang:1.25.3-alpine

# Install necessary packages
RUN apk add --no-cache make

# Copy the source code and compile
WORKDIR /src
COPY . .
RUN make oauthd

# Setup configuration files
WORKDIR /config
RUN cp /src/config/oauth_example.json ./oauth.json

# Copy the app binary and entrypoint, then create necessary folders
WORKDIR /app
RUN cp /src/build/oauthd . &&\
    cp /src/callback/docker/entrypoint.sh . &&\
    chmod +x entrypoint.sh

# Forward ports
EXPOSE 8000/tcp

# Set volumes
VOLUME ["/config"]

# Set binary and parameters
ENTRYPOINT ["/app/entrypoint.sh", "/app/oauthd"]
CMD ["--config", "/config/oauth.json"]
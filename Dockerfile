# Compile stage
FROM golang:1.21-bookworm  as build-env

# Creates an app directory to hold your app’s source code
WORKDIR /app

# Copies everything from your root directory into /app
COPY . .

# Installs Go dependencies
RUN go mod download

# Builds your app with optional configuration
RUN go build -o /server

# Final stage
FROM debian:bookworm

# Installs g++ that contains GLIBC_2.34 used by the server binary
RUN apt update
RUN apt install g++ --quiet --yes

# Copies binary from compile stage to final stage
WORKDIR /
COPY --from=build-env /server /

# Tells Docker which network port your container listens on
EXPOSE 8080

# Specifies the executable command that runs when the container starts
CMD ["/server"]
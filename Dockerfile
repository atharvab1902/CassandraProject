# Stage 1: Compile the binary in a containerized Golang environment
#
FROM golang:1.21 as build

# Copy the source files from the host
COPY . /go/src

WORKDIR /go/src/reader
RUN go build -o reader

WORKDIR /go/src/writer
RUN go build -o writer

# Stage 2: Build the image for client
#
FROM ubuntu:22.04 as image

# Copy the binary from the build container
WORKDIR /client
COPY --from=build /go/src/reader/reader .
COPY --from=build /go/src/writer/writer .

ENV PATH="/client:$PATH"

# Let it run forever and use exec to run reader and writer
CMD tail -f /dev/null

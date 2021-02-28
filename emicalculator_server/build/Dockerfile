FROM golang:1.13-alpine

# Update packages and install git.
RUN apk update && apk add git

# Move to project root.
WORKDIR /go/src/github.com/gkjoyes/emi-calculator

# Install dependencies.
RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go

# Copy source code to the container.
COPY emicalculator_server/server.go ./
COPY emicalculatorpb emicalculatorpb

# Run the application.
RUN go install .
ENTRYPOINT ["/go/bin/emi-calculator"]
EXPOSE 5300
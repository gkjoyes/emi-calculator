FROM golang:1.13-alpine

# Update packages and install git.
RUN apk update && apk add git

# Move to project root.
WORKDIR /go/src/github.com/gkjoyes/emi-calculator

# Install dependencies.
RUN go get -u google.golang.org/grpc

# Copy source code to the container.
COPY emicalculator_client/client.go ./
COPY emicalculatorpb emicalculatorpb

# Run the application.
RUN go install .
CMD ["emi-calculator"]
EXPOSE 5200
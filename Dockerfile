# Build stage
FROM golang:alpine AS build-env
ADD . $GOPATH/src/github/platformwg/GoRunner
RUN apk update && apk add git
RUN cd $GOPATH/src/github/platformwg/GoRunner/ && go get && go build -o /GoRunner

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /GoRunner /app/
ENTRYPOINT ["./GoRunner"]

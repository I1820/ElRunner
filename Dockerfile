# Build stage
FROM golang:alpine AS build-env
ADD . $GOPATH/src/github/aiotrc/GoRunner
RUN apk update && apk add git
RUN cd $GOPATH/src/github/aiotrc/GoRunner/ && go get && go build -o /GoRunner

# Final stage
FROM python:3
WORKDIR /app
COPY --from=build-env /GoRunner /app/
ENTRYPOINT ["./GoRunner"]

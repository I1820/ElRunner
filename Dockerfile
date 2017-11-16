# Build stage
FROM golang:alpine AS build-env
ADD . $GOPATH/src/github.com/aiotrc/GoRunner
RUN apk update && apk add git
RUN cd $GOPATH/src/github.com/aiotrc/GoRunner/ && go get -v && go build -v -o /GoRunner

# Final stage
FROM python:3-alpine
WORKDIR /app
COPY --from=build-env /GoRunner /app/
ADD runtime.py /app/runtime.py
RUN cd /app/runtime.py && python3 setup.py install
ENTRYPOINT ["./GoRunner"]

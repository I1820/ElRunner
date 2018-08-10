# Build stage
FROM golang:alpine AS build-env
COPY . $GOPATH/src/github.com/I1820/ElRunner
RUN apk update && apk add git
WORKDIR $GOPATH/src/github.com/I1820/ElRunner
RUN go get -v && go build -v -o /ElRunner

# Final stage
FROM alpine:latest

# Metadata
LABEL org.i1820.build-date=$BUILD_DATE
LABEL maintainer="Parham Alvani <parham.alvani@gmail.com>"

EXPOSE 8080/tcp
WORKDIR /app
COPY --from=build-env /ElRunner /app/
COPY runtime.py /app/runtime.py
# Install python stuffs
RUN apk update && apk add ca-certificates && update-ca-certificates
RUN apk update && apk add --no-cache python3 && \
            python3 -m ensurepip && \
            rm -r /usr/lib/python*/ensurepip && \
            pip3 install --upgrade pip setuptools && \
            if [ ! -e /usr/bin/pip ]; then ln -s pip3 /usr/bin/pip ; fi && \
            if [[ ! -e /usr/bin/python ]]; then ln -sf /usr/bin/python3 /usr/bin/python; fi && \
            rm -r /root/.cache
RUN apk update && apk add build-base python3-dev
# Install runtime.py
RUN cd /app/runtime.py && python3 setup.py install
# Remove python stuffs
RUN apk del build-base python3-dev && \
            rm -rf /var/cache/apk/*
ENTRYPOINT ["./ElRunner"]

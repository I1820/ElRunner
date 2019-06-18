# Build stage
FROM golang:1.12-alpine AS build-env
WORKDIR /app
RUN apk --no-cache add git ca-certificates git gcc g++ libc-dev
COPY . .
RUN go build -v -o /ElRunner

# Final stage
FROM alpine:latest

# Metadata
ARG BUILD_DATE
ARG BUILD_COMMIT
ARG BUILD_COMMIT_MSG
LABEL maintainer="Parham Alvani <parham.alvani@gmail.com>"
LABEL org.i1820.build-date=$BUILD_DATE
LABEL org.i1820.build-commit-sha=$BUILD_COMMIT
LABEL org.i1820.build-commit-msg=$BUILD_COMMIT_MSG

EXPOSE 8080/tcp
WORKDIR /app
COPY --from=build-env /ElRunner /app/
COPY runtime.py /app/runtime.py
# Install python stuffs
RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN apk add --no-cache python3 && \
            python3 -m ensurepip && \
            rm -r /usr/lib/python*/ensurepip && \
            pip3 install --upgrade pip setuptools && \
            if [ ! -e /usr/bin/pip ]; then ln -s pip3 /usr/bin/pip ; fi && \
            if [ ! -e /usr/bin/python ]; then ln -sf /usr/bin/python3 /usr/bin/python; fi && \
            rm -r /root/.cache
RUN apk add --no-cache build-base python3-dev
# Install runtime.py
WORKDIR /app/runtime.py
RUN python3 setup.py install
# Remove python stuffs
RUN apk del build-base python3-dev
ENTRYPOINT ["/app/ElRunner"]

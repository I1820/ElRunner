# Build stage
FROM golang:alpine AS build-env
COPY . $GOPATH/src/github.com/I1820/ElRunner
RUN apk --no-cache add git
WORKDIR $GOPATH/src/github.com/I1820/ElRunner
RUN go get -v && go build -v -o /ElRunner

# Final stage
FROM python:3.7-alpine3.8

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
RUN apk add --no-cache build-base python3-dev
# Install runtime.py
WORKDIR /app/runtime.py
RUN python3 setup.py install
# Remove python stuffs
RUN apk del build-base python3-dev
ENTRYPOINT ["/app/ElRunner"]

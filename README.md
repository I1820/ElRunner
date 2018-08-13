# ElRunner
[![Travis branch](https://img.shields.io/travis/com/I1820/ElRunner/master.svg?style=flat-square)](https://travis-ci.com/I1820/ElRunner)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/ElRunner?style=flat-square)](https://goreportcard.com/report/github.com/I1820/ElRunner)
[![Maintainability](https://api.codeclimate.com/v1/badges/581d84707830162ab4a8/maintainability)](https://codeclimate.com/github/I1820/ElRunner/maintainability)
[![Docker Pulls](https://img.shields.io/docker/pulls/i1820/elrunner.svg?style=flat-square)]()

## Introduction

ElRunner runs Python3 code when specific events come, but it ensure you, you will have one instance of your code in running state.
its provides runtime python library [runtime.py] for you application in order to have required packages and functions.

## Environment Variables

These variables are avaiable in container:

- Project ID: `NAME`
- User Email: `USER`
- Redis Host: `REDIS_HOST`
- Mongo URL: `MONGO_URL`

## Decode/Encode

ElRunner can decode/encode your data with your given codec in python.

```python
from codec import Codec
import cbor

class ISRC(Codec):
    def decode(self, data):
        return cbor.loads(data)
    def encode(self, data):
        pass
```

## Scenario

ElRunner can run your given scenario on data comming events or periodically.

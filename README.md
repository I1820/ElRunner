# ElRunner

[![Travis branch](https://img.shields.io/travis/aiotrc/GoRunner/master.svg?style=flat-square)](https://travis-ci.org/aiotrc/GoRunner)
[![Docker Pulls](https://img.shields.io/docker/pulls/aiotrc/gorunner.svg?style=flat-square)]()

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

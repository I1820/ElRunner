# GoRunner
[![Travis branch](https://img.shields.io/travis/aiotrc/GoRunner/master.svg?style=flat-square)](https://travis-ci.org/aiotrc/GoRunner)
[![Docker Pulls](https://img.shields.io/docker/pulls/aiotrc/gorunner.svg?style=flat-square)]()

## Introduction
GoRunner runs Python3 code when specific events come, but it ensure you, you will have one instance of your code in running state.
its provides runtime python library [runtime.py] for you application in order to have required packages and ...

## Decode/Encode
GoRunner can decode/encode your data with your given codec in python.

```python
class ISRC(Codec, requirements=["cbor"]):
    def decode(self, data):
        return self.cbor.loads(data)
    def encode(self, data):
        pass
```

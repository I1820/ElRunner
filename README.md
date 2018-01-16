# GoRunner
[![Travis branch](https://img.shields.io/travis/aiotrc/GoRunner/master.svg?style=flat-square)](https://travis-ci.org/aiotrc/GoRunner)
[![Docker Pulls](https://img.shields.io/docker/pulls/aiotrc/gorunner.svg?style=flat-square)]()
[![Codacy grade](https://img.shields.io/codacy/grade/b0b53df0a7264498a760232425be52e4.svg?style=flat-square)](https://www.codacy.com/app/1995parham/GoRunner?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=aiotrc/GoRunner&amp;utm_campaign=Badge_Grade)

## Introduction
GoRunner runs Python3 code when specific events come, but it ensure you, you will have one instance of your code in running state.
its provides runtime python library [runtime.py] for you application in order to have required packages and ...

## Environment Variables
- Project ID
- User ID

## Decode/Encode
GoRunner can decode/encode your data with your given codec in python.

```python
from codec import Codec
import cbor

class ISRC(Codec):
    def decode(self, data):
        return cbor.loads(data)
    def encode(self, data):
        pass
```

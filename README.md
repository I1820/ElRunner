# ElRunner
[![Travis branch](https://img.shields.io/travis/com/I1820/ElRunner/master.svg?style=flat-square)](https://travis-ci.com/I1820/ElRunner)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/ElRunner?style=flat-square)](https://goreportcard.com/report/github.com/I1820/ElRunner)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/eada226f7b04403380cb7dc8dd517e5b)](https://www.codacy.com/app/i1820/ElRunner?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=I1820/ElRunner&amp;utm_campaign=Badge_Grade)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)
[![Docker Pulls](https://img.shields.io/docker/pulls/i1820/elrunner.svg?style=flat-square)]()

## Introduction

ElRunner runs Python3 code when specific events come, but it ensure you, you will have one instance of your code in running state.
its provides runtime python library [runtime.py] for you application in order to have required packages and functions.

## Environment Variables

These variables are avaiable in container:

- Project ID: `PROJECT`
- Project Owner Email Address: `OWNER`
- Redis Host: `REDIS_HOST`
- Mongo URL: `DB_URL`
- Broker URL: `BROKER_URL`

## Runtime.py
Python runtime environment for ElRunner scripts. This environment is written in python 3.7. In order to install
it do following commands in ubuntu 18.04:

```sh
sudo apt install python3.7 python3.7-venv python3.7-doc
cd runtime.py
python3.7 -mvenv .
rm bin/python3
ln -s /usr/bin/python3.7 bin/python3
pip3.7 install -U setuptools
python3.7 setup.py install
```

Please note that following packages have C dependencies so they require python3.7-dev for installation.

- cbor
- aiohttp

## Logging
Each I1820Core classes have a Logger that sends logs to MongoDB using [log4mongo](https://pypi.org/project/log4mongo/).

## Decode/Encode

ElRunner can decode/encode your data with your given codec in python.

```python
from codec import Codec
import cbor

class Fanco(Codec):
    def decode(self, data):
        return cbor.loads(data)
    def encode(self, data):
        return cbor.stores(data)
```

## Scenario

ElRunner can run your given scenario on data comming events or periodically.

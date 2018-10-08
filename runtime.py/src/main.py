#!/usr/bin/env python3
# In The Name Of God
# ========================================
# [] File Name : main.py
#
# [] Creation Date : 13-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import click
import base64
import runpy
import json
import traceback
import sys
import contextlib

from codec import Codec
from scenario import Scenario


@click.command()
@click.argument('target', type=click.Path())
@click.option('--job', type=click.Choice(['decode', 'encode', 'rule']))
@click.option('--id', type=str, default="")
def run(target, job, id):
    '''
    run given target in isolated environment. job specifies target type.
    id specifies a thing we are working on :)
    '''
    g = runpy.run_path(target, run_name='ucodec')
    for value in g.values():
        if isinstance(value, type) and issubclass(value, Codec) and \
                value.__module__ == 'ucodec':
            codec = value

    g = runpy.run_path(target, run_name='uscenario')
    for value in g.values():
        if isinstance(value, type) and issubclass(value, Scenario) and \
                value.__module__ == 'uscenario':
            scenario = value

    if (job == 'decode' or job == 'encode') and 'codec' not in locals():
        raise Exception("Invalid codec class")

    if job == 'scenario' and 'scenario' not in locals():
        raise Exception("Invalid scenario class")

    if job == 'decode':
        s = input()
        with contextlib.redirect_stdout(sys.stderr):
            d = codec().decode(base64.b64decode(s))
        if codec.thing_location != '' and codec.thing_location in d:
            d['_location'] = d[codec.thing_location]
        print(json.dumps(d))

    if job == 'encode':
        s = input()
        with contextlib.redirect_stdout(sys.stderr):
            e = codec().encode(json.loads(s))
        print(base64.b64encode(e).decode('ascii'))

    if job == 'rule':
        s = input()
        scenario(id).run(json.loads(s))


def main():
    try:
        run()
    except Exception:
        traceback.print_exc()
        sys.exit(1)

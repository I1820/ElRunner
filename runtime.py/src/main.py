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

from codec import Codec
from analyze import run_pylint_json


@click.command()
@click.argument('target', type=click.Path())
@click.option('--job', type=click.Choice(['decode', 'encode', 'rule', 'lint']))
def run(target, job):
    '''
    run given target in provided environment
    '''
    if job == 'lint':
        r = run_pylint_json(target)
        print(r)
        return

    g = runpy.run_path(target, run_name='ucodec')
    for value in g.values():
        if isinstance(value, type) and issubclass(value, Codec) and \
                value.__module__ == 'ucodec':
            codec = value

    if job == 'decode':
        s = input()
        d = codec().decode(base64.b64decode(s))
        print(d)
    if job == 'encode':
        s = input()
        e = codec().encode(s)
        print(base64.b64encode(e).decode('ascii'))


def main():
    run()

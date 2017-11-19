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
from codec import Codec


@click.command()
@click.argument('target', type=click.Path())
@click.option('--job', type=click.Choice(['decode', 'encode', 'rule']))
def run(target, job):
    '''
    run given target in provided environment
    '''
    try:
        exec(compile(open(target, 'r').read(), target, 'exec'))
    except Exception as e:
        print('Target Error: ', e)
        return
    if job == 'decode':
        s = input()
        d = Codec.get()().decode(base64.b64decode(s))
        print(d)
    if job == 'encode':
        s = input()
        e = Codec.get()().encode(s)
        print(base64.b64encode(e))


def main():
    run()

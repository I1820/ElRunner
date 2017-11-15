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
@click.option('--target', help='Target python source code', type=click.Path())
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
        s = base64.b64decode(s)
        d = Codec.get().decode(s)
        print(d)


if __name__ == '__main__':
    run()

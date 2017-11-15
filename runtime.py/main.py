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
from codec import Codec


@click.command()
@click.option('--target', help='Target python source code', type=click.Path())
def run(target):
    """ run given target in provided environment """
    try:
        exec(compile(open(target, "rb").read(), target, 'exec'))
    except Exception as e:
        print("Target Error: ", e)


if __name__ == '__main__':
    run()

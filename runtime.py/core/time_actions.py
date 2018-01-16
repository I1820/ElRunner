import time
from threading import Timer


def sleep(seconds):
    time.sleep(seconds)


def schedule(delay_seconds, action_function, args=()):
    Timer(delay_seconds, action_function, args).start()

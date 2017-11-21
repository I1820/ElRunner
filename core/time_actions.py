import time
from threading import Timer


def sleep(seconds):
    time.sleep(seconds)


def schedule(delay_seconds, action_function):
    Timer(delay_seconds, action_function, ()).start()

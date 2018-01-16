import datetime

from core import time_actions


def test_sleep():
    time_actions.sleep(seconds=1)
    time_actions.sleep(seconds=0.1)


def action_function():
    print("action function:")
    print(datetime.datetime.now())


def test_schedule():
    time_actions.schedule(delay_seconds=1, action_function=action_function)
    time_actions.schedule(delay_seconds=0.5, action_function=action_function)

import time

import datetime

from core import time_actions


def sleep_test():
    print("test_sleep:")
    print(time.time())
    time_actions.sleep(seconds=1)
    print(time.time())
    time_actions.sleep(seconds=0.1)
    print(time.time())


def action_function():
    print("action_function:")
    print(datetime.datetime.now())


def schedule_test():
    print("schedule_test:")
    print("Before Schedule:" + str(datetime.datetime.now()))
    time_actions.schedule(delay_seconds=1, action_function=action_function)
    time_actions.schedule(delay_seconds=0.5, action_function=action_function)
    print("After Schedule:" + str(datetime.datetime.now()))


sleep_test()
schedule_test()

import datetime
import pytest
import scenario


class TestScenario(scenario.Scenario):
    def run():
        pass


def test_sleep(ts):
    ts.sleep(seconds=1)
    ts.sleep(seconds=0.1)


def action_function():
    print("action function:")
    print(datetime.datetime.now())


@pytest.fixture(scope="session")
def ts():
    s = TestScenario("")
    return s


def test_schedule(ts):
    ts.schedule(delay_seconds=1, action_function=action_function)
    ts.schedule(delay_seconds=0.5, action_function=action_function)

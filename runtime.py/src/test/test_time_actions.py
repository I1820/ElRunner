import datetime
import pytest
import scenario


class TestScenario(scenario.Scenario):
    def run():
        pass


def test_sleep(scenario):
    scenario.sleep(seconds=1)
    scenario.sleep(seconds=0.1)


def action_function():
    print("action function:")
    print(datetime.datetime.now())


@pytest.fixture(scope="session")
def scenario():
    s = TestScenario()
    return s


def test_schedule(scenario):
    scenario.schedule(delay_seconds=1, action_function=action_function)
    scenario.schedule(delay_seconds=0.5, action_function=action_function)

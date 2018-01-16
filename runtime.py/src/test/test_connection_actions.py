import threading
import time
import pytest
import scenario

from core.rpc_server import start_server

SERVER_DATA_RESPONSE = 'Test Message'
SERVER_ACK_RESPONSE = 'OK'


def wait_for_data():
    return SERVER_DATA_RESPONSE


def send_to_down_link(message):
    print(message)
    return SERVER_ACK_RESPONSE


class TestScenario(scenario.Scenario, requirements=[]):
    def run():
        pass


@pytest.fixture(scope="session")
def rpc():
    thread = threading.Thread(target=start_server,
                              daemon=True,
                              args=(wait_for_data, send_to_down_link))
    thread.start()
    # wait for server to load
    time.sleep(2)

    return thread


@pytest.fixture(scope="session")
def scenario():
    s = TestScenario()
    return s


@pytest.mark.asyncio
async def test_wait_for_data(rpc, scenario):
    response = await scenario.wait_for_data(timeout=30)
    print(response)
    assert response['result'] == 'Test Message'


@pytest.mark.asyncio
async def test_send_to_down_link(rpc, scenario):
    response = await scenario.send_to_down_link(
        message=SERVER_DATA_RESPONSE, timeout=30)
    print(response)
    assert response['result'] == 'OK'

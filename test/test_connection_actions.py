import threading
import time
import pytest

from core import connection_actions
from core.rpc_server import start_server

SERVER_DATA_RESPONSE = 'Test Message'
SERVER_ACK_RESPONSE = 'OK'


def wait_for_data():
    return SERVER_DATA_RESPONSE


def send_to_down_link(message):
    print(message)
    return SERVER_ACK_RESPONSE


@pytest.fixture(scope="session")
def rpc():
    thread = threading.Thread(target=start_server,
                              daemon=True,
                              args=(wait_for_data, send_to_down_link))
    thread.start()
    # wait for server to load
    time.sleep(2)

    return thread


@pytest.mark.asyncio
async def test_wait_for_data(rpc):
    response = await connection_actions.wait_for_data(timeout=30)
    print(response)
    assert response['result'] == 'Test Message'


@pytest.mark.asyncio
async def test_send_to_down_link(rpc):
    response = await connection_actions.send_to_down_link(
        message=SERVER_DATA_RESPONSE, timeout=30)
    print(response)
    assert response['result'] == 'OK'

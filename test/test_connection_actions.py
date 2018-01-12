import threading
import pytest
import time

from core import connection_actions
from core.rpc_server import start_server

server_data_response = 'Test Message'
server_ack_response = 'OK'


def wait_for_data():
    return server_data_response


def send_to_down_link(message):
    return server_ack_response


@pytest.fixture(scope="session")
def rpc():
    t = threading.Thread(target=start_server,
                         daemon=True,
                         args=(wait_for_data, send_to_down_link))
    t.start()
    # wait for server to load
    time.sleep(2)

    return t


@pytest.mark.asyncio
async def test_wait_for_data(rpc):
    response = await connection_actions.wait_for_data(timeout=30)
    print(response)
    assert response['result'] == 'Test Message'


@pytest.mark.asyncio
async def test_send_to_down_link(rpc):
    response = await connection_actions.send_to_down_link(
            message=server_data_response, timeout=30)
    print(response)
    assert response['result'] == 'OK'

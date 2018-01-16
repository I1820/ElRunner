import threading
import time
import pytest
import scenario

from jsonrpc import JSONRPCResponseManager, dispatcher
from werkzeug.serving import run_simple
from werkzeug.wrappers import Request, Response


SERVER_DATA_RESPONSE = 'Test Message'
SERVER_ACK_RESPONSE = 'OK'


def wait_for_data():
    return SERVER_DATA_RESPONSE


def send_to_down_link(message):
    print(message)
    return SERVER_ACK_RESPONSE


def start_server(wait_for_data, send_to_down_link):
    @Request.application
    def application(request):
        # Dispatcher is dictionary {<method_name>: callable}
        dispatcher["Endpoint.WaitForData"] = wait_for_data
        dispatcher["Endpoint.SendToDownLink"] = send_to_down_link

        response = JSONRPCResponseManager.handle(request.data, dispatcher)
        return Response(response.json, mimetype='application/json')

    run_simple(scenario.RPC_SERVER, scenario.PRC_PORT, application)


class TestScenario(scenario.Scenario):
    def run():
        pass


@pytest.fixture(scope="session")
def rpc():
    thread = threading.Thread(target=start_server,
                              daemon=True,
                              args=(wait_for_data, send_to_down_link))
    return thread


@pytest.fixture(scope="session")
def scenario():
    s = TestScenario()
    return s


def test_rpc_server(rpc):
    rpc.start()
    # wait for server to load
    time.sleep(2)


@pytest.mark.asyncio
async def test_wait_for_data(scenario):
    response = await scenario.wait_for_data(timeout=30)
    print(response)
    assert response['result'] == 'Test Message'


@pytest.mark.asyncio
async def test_send_to_down_link(scenario):
    response = await scenario.send_to_down_link(
        message=SERVER_DATA_RESPONSE, timeout=30)
    print(response)
    assert response['result'] == 'OK'

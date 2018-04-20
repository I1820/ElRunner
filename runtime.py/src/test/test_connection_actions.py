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

    run_simple(scenario.RPC_SERVER, scenario.RPC_PORT, application)


class TestScenario(scenario.Scenario):
    def run(self, data=None):
        pass


@pytest.fixture(scope="session")
def rpc():
    run_rpc_server(wait_for_data, send_to_down_link)


def run_rpc_server(wait_for_data, send_to_down_link):
    thread = threading.Thread(target=start_server,
                              daemon=True,
                              args=(wait_for_data, send_to_down_link))
    thread.start()
    # wait for server to load
    time.sleep(2)
    print("RPC is ready")

    return thread


@pytest.fixture(scope="session")
def ts():
    s = TestScenario("")
    return s


@pytest.mark.asyncio
@pytest.mark.usefixtures("rpc")
async def test_wait_for_data(ts):
    response = await ts.wait_for_data(timeout=30)
    print(response)
    if not response == 'Test Message':
        raise AssertionError()


@pytest.mark.asyncio
@pytest.mark.usefixtures("rpc")
async def test_send_to_down_link(ts):
    response = await ts.send_to_down_link(
            message=SERVER_DATA_RESPONSE, timeout=30)
    print(response)
    if not response == 'OK':
        raise AssertionError()

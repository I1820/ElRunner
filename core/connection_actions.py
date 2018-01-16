import aiohttp

from core.connection_config import server_name, server_port


URL = 'http://{}:{}/'.format(server_name, server_port)
HEADERS = {'content-type': 'application/json'}
PAYLOAD = {'jsonrpc': '2.0'}


async def wait_for_data(timeout):
    request_payload = PAYLOAD.copy()
    request_payload['method'] = 'Endpoint.WaitForData'
    request_payload['params'] = []
    request_payload['id'] = 0

    async with aiohttp.ClientSession() as session:
        session.headers = HEADERS
        response = await session.post(URL, json=request_payload,
                                      timeout=timeout)
        return await response.json()


async def send_to_down_link(message, timeout):
    request_payload = PAYLOAD.copy()
    request_payload['method'] = 'Endpoint.SendToDownLink'
    request_payload['params'] = [message]
    request_payload['id'] = 1

    async with aiohttp.ClientSession() as session:
        session.headers = HEADERS
        response = await session.post(URL, json=request_payload,
                                      timeout=timeout)
        return await response.json()

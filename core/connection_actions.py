import aiohttp

from core.connection_config import server_name, server_port


url = 'http://{}:{}/'.format(server_name, server_port)
headers = {'content-type': 'application/json'}
payload = {'jsonrpc': '2.0'}


async def wait_for_data(timeout):
    request_payload = payload.copy()
    request_payload['method'] = 'Endpoint.WaitForData'
    request_payload['params'] = []
    request_payload['id'] = 0

    async with aiohttp.ClientSession() as session:
        session.headers = headers
        response = await session.post(url, json=request_payload,
                                      timeout=timeout)
        return response.json()


async def send_to_down_link(message, timeout):
    request_payload = payload.copy()
    request_payload['method'] = 'Endpoint.SendToDownLink'
    request_payload['params'] = [message]
    request_payload['id'] = 1

    async with aiohttp.ClientSession() as session:
        session.headers = headers
        response = await session.post(url, json=request_payload,
                                      timeout=timeout)
        return response.json()

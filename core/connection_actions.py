import json
import sys
from requests_threads import AsyncSession

from core.connection_config import server_name, server_port, debug


url = 'http://{}:{}/'.format(server_name, server_port)
session = AsyncSession()
session.headers = {'content-type': 'application/json'}
payload = dict(method='method', params=[], jsonrpc='2.0', id=0)


async def wait_for_data(timeout_seconds):
    request_payload = payload.copy()
    request_payload['method'] = 'wait_for_data'
    request_payload['params'] = []
    request_payload['id'] = 0
    response = await session.post(url, data=json.dumps(request_payload), timeout=timeout_seconds).json()
    if response:
        return response
    else:
        if debug:
            print(sys.stderr, 'No Response Received!')
        return None


async def send_to_down_link(message, timeout_seconds):
    request_payload = payload.copy()
    request_payload['method'] = 'send_to_down_link'
    request_payload['params'] = [message]
    request_payload['id'] = 1
    response = await session.post(url, data=json.dumps(request_payload), timeout=timeout_seconds).json()
    if response:
        return response
    else:
        if debug:
            print(sys.stderr, 'No Response Received!')
        return None

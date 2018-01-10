import json
import sys

import requests

from core.connection_config import server_name, server_port, debug


url = 'http://{}:{}/'.format(server_name, server_port)
headers = {'content-type': 'application/json'}
payload = dict(method='method', params=[], jsonrpc='2.0', id=0)


def wait_for_data(timeout_seconds):
    request_payload = payload.copy()
    request_payload['method'] = 'wait_for_data'
    request_payload['params'] = []
    request_payload['id'] = 0
    response = requests.post(url, data=json.dumps(request_payload), headers=headers, timeout=timeout_seconds).json()
    if response:
        return response
    else:
        if debug:
            print(sys.stderr, 'No Response Received!')
        return None


def send_to_down_link(message, timeout_seconds):
    request_payload = payload.copy()
    request_payload['method'] = 'send_to_down_link'
    request_payload['params'] = [message]
    request_payload['id'] = 1
    response = requests.post(url, data=json.dumps(request_payload), headers=headers, timeout=timeout_seconds).json()
    if response:
        return response
    else:
        if debug:
            print(sys.stderr, 'No Response Received!')
        return None

# In The Name Of God
# ========================================
# [] File Name : codec.py
#
# [] Creation Date : 15-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import abc
import aiohttp
import time
import threading


RPC_SERVER = '127.0.0.1'
RPC_PORT = 1373
URL = 'http://{}:{}/'.format(RPC_SERVER, RPC_PORT)
HEADERS = {'content-type': 'application/json'}
PAYLOAD = {'jsonrpc': '2.0'}


class Scenario(metaclass=abc.ABCMeta):
    def sleep(self, seconds):
        time.sleep(seconds)

    def schedule(self, delay_seconds, action_function, args=()):
        threading.Timer(delay_seconds, action_function, args).start()

    async def wait_for_data(self, timeout):
        request_payload = PAYLOAD.copy()
        request_payload['method'] = 'Endpoint.WaitForData'
        request_payload['params'] = []
        request_payload['id'] = 0

        async with aiohttp.ClientSession() as session:
            session.headers = HEADERS
            response = await session.post(URL, json=request_payload,
                                          timeout=timeout)
            return await response.json()

    async def send_to_down_link(self, message, timeout):
        request_payload = PAYLOAD.copy()
        request_payload['method'] = 'Endpoint.SendToDownLink'
        request_payload['params'] = [message]
        request_payload['id'] = 1

        async with aiohttp.ClientSession() as session:
            session.headers = HEADERS
            response = await session.post(URL, json=request_payload,
                                          timeout=timeout)
            return await response.json()

    @abc.abstractmethod
    def run(self, data=None):
        pass

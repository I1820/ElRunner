# - اگر داده سنسور شماره x بر روی شی a آمد، یک ایمیل ارسل شده و این رخداد را اطلاع دهد

import json
import socket
import _thread
import asyncio

from scenario import Scenario
from test.test_connection_actions import rpc, run_rpc_server

thing_id = 'a'
sensor_id = 'x'
server_data_response = '{"thing_id":"a","sensor_id":"x","data":"100"}'
server_ack_response = 'ACK'


def wait_for_data():
    return server_data_response


def send_to_down_link(message):
    print('server got message: ' + message)
    return server_ack_response


class Scenario1(Scenario):

    def __init__(self):
        self.received = False

    def action(self, data):
        print("action:")
        data_parsed_json = json.loads(data)

        if data_parsed_json["thing_id"] != thing_id or \
                        data_parsed_json["sensor_id"] != sensor_id:
            print("not expected thing and sensor! expected[" + thing_id + ":" + sensor_id + "] got[" +
                  data_parsed_json["thing_id"] + ":" + data_parsed_json["sensor_id"] + "]")
            return
        self.received = True
        sender = 'ceitiotlabtest@gmail.com'
        receivers = ['ceitiotlabtest@gmail.com']

        message = 'From: From Person <ceitiotlabtest@gmail.com>\n' \
                  'To: To Person <ceitiotlabtest@gmail.com>\n' \
                  'Subject: Rule Engine Notification\n\n' \
                  'Data:' + data + '\n' \
                                   'Sent by Rule Engine. Scenario:1.'
        self.send_email(host='smtp.gmail.com', port=587, username="ceitiotlabtest", password="ceit is the best",
                   sender=sender,
                   receivers=receivers, message=message)

    async def wait_for_data_wrapper(self, future, timeout):
        future.set_result(await self.wait_for_data(timeout))

    def run(self, data=None):
        while True:
            try:
                print("wait for data...")
                loop = asyncio.get_event_loop()
                future = asyncio.Future()
                loop.run_until_complete(self.wait_for_data_wrapper(future, timeout=30))
                loop.close()
                response = future.result()
                if response:
                    print('Response:' + str(response))
                    self.action(response['result'])
                else:
                    print('No Response!')
            except socket.timeout:
                print("wait for data: Timeout!")
            if self.received:
                break

run_rpc_server(wait_for_data, send_to_down_link)

scenario_1 = Scenario1()
scenario_1.run()

# coding=utf-8
# - اگر داده سنسور شماره x بر روی شی a آمد، یک ایمیل ارسل شده و این رخداد را اطلاع دهد
import json
import socket
import thread

import jsonrpclib

from core import connection_actions
from core.notification_actions import send_email
from core.rpc_server import start_server

thing_id = 'a'
sensor_id = 'x'
server_data_response = '{"thing_id":"a","sensor_id":"x","data":"100"}'
server_ack_response = 'ACK'
received = False


def wait_for_data():
    return server_data_response


def send_to_down_link(message):
    print('server got message: ' + message)
    return server_ack_response


def action(data):
    print("action:")
    data_parsed_json = json.loads(data)
    if data_parsed_json["thing_id"] != thing_id or data_parsed_json["sensor_id"] != sensor_id:
        print("not expected thing and sensor! expected[" + thing_id + ":" + sensor_id + "] got[" +
              data_parsed_json["thing_id"] + ":" + data_parsed_json["sensor_id"] + "]")
        return
    global received
    received = True
    sender = 'ceitiotlabtest@gmail.com'
    receivers = ['ceitiotlabtest@gmail.com']

    message = 'From: From Person <ceitiotlabtest@gmail.com>\n' \
              'To: To Person <ceitiotlabtest@gmail.com>\n' \
              'Subject: Rule Engine Notification\n\n' \
              'Data:' + data + '\n' \
                               'Sent by Rule Engine. Scenario:1.'
    send_email(host='smtp.gmail.com', port=587, username="ceitiotlabtest", password="ceit is the best", sender=sender,
               receivers=receivers, message=message)


thread.start_new(start_server, (wait_for_data, send_to_down_link))

while True:
    try:
        print("wait for data...")
        response = connection_actions.wait_for_data(timeout_seconds=30)
        print('Request:' + jsonrpclib.history.request)
        if response:
            print('Response:' + jsonrpclib.history.response)
            action(response)
        else:
            print('No Response!')
    except socket.timeout:
        print("wait for data: Timeout!")
    if received:
        break

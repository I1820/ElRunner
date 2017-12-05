# coding=utf-8
# - اگر داده سنسور شماره x بر روی شی a آمد، یک ایمیل ارسل شده و این رخداد را اطلاع دهد
import socket
import json

from core.connection_actions import wait_for_data, send_to_down_link
from core.notification_actions import send_email

thing_id = "a"
sensor_id = "x"
packet_message = """{"thing_id":"a","sensor_id":"x","data":"100"}"""
ack_message = """ACK"""
received = False


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


while True:
    try:
        print("wait for data...")
        wait_for_data(data_received_function=action, read_bytes=len(packet_message.encode('utf-8')),
                      ack_message=ack_message, timeout_seconds=30)
    except socket.timeout:
        print("wait for data: Timeout!")
    if received:
        break

# try:
#     print("send to down link...")
#     send_to_down_link(message=packet_message, expected_ack_message=ack_message, ack_timeout_seconds=30)
# except socket.timeout:
#     print("send to down link: Timeout!")

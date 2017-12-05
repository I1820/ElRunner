# coding=utf-8
# - اگر اول داده سنسور x از شی a آمد، بعد داده سنسور y از شی b آمد میانگین آنها را در پایگاه داده ذخیره کند.
import socket

from core.connection_actions import wait_for_data

packet_message = """{"thing_id":"a","sensor_id":"x","data":"..."}"""
ack_message = """ACK"""


def action(data):
    print("action:")



try:
    print("wait for data...")
    wait_for_data(data_received_function=action, read_bytes=len(packet_message.encode('utf-8')),
                  ack_message=ack_message, timeout_seconds=30)
except socket.timeout:
    print("wait for data: Timeout!")
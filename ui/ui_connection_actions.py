import socket

from core.connection_actions import wait_for_data, send_to_down_link

test_message = "Test Message"


def data_received_function(data):
    print("data_received_function:")
    print("Data:")
    print(data)


def wait_for_data_test():
    print("wait_for_data_test:")
    wait_for_data(data_received_function=data_received_function, read_bytes=len(test_message.encode('utf-8')),
                  ack_message="ACK", timeout_seconds=1)


def send_to_down_link_test():
    print("send_to_down_link_test:")
    send_to_down_link(message=test_message, expected_ack_message="ACK", ack_timeout_seconds=1)


try:
    wait_for_data_test()
except socket.timeout:
    print("wait_for_data_test Timeout!")

try:
    send_to_down_link_test()
except socket.timeout:
    print("send_to_down_link_test Timeout!")

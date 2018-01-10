import socket
import thread

from core import connection_actions
from core.rpc_server import start_server

server_data_response = 'Test Message'
server_ack_response = 'OK'


def wait_for_data():
    return server_data_response


def send_to_down_link(message):
    print('server got message: ' + message)
    return server_ack_response


def response_received_function(data):
    print("response_received_function:")
    print("Response:")
    print(data)


def wait_for_data_test():
    print("wait_for_data_test:")
    response = connection_actions.wait_for_data(timeout_seconds=30)
    callback(response)


# def wait_for_data_async_test():
#     print("wait_for_data_async_test:")
#     if __name__ == '__main__':
#         connection_actions.get_data_async(timeout_seconds=30, callback=callback)


def send_to_down_link_test():
    print("send_to_down_link_test:")
    response = connection_actions.send_to_down_link(message=server_data_response, timeout_seconds=30)
    callback(response)


def callback(response):
    print('callback:')
    if response:
        response_received_function(response)
    else:
        print('No Response!')


thread.start_new(start_server, (wait_for_data, send_to_down_link))

try:
    wait_for_data_test()
except socket.timeout:
    print("wait_for_data_test Timeout!")

try:
    send_to_down_link_test()
except socket.timeout:
    print("send_to_down_link_test Timeout!")

# wait_for_data_async_test()

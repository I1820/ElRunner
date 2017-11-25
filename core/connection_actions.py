import socket
import sys

from core.connection_config import receive_data_server_name, receive_data_server_port, send_data_server_name, \
    send_data_server_port, debug


def wait_for_data(data_received_function, read_bytes, ack_message, timeout_seconds):
    # Create a TCP/IP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.settimeout(timeout_seconds)
    # Bind the socket to the address given on the command line
    server_address = (receive_data_server_name, receive_data_server_port)
    if debug:
        print >> sys.stderr, 'starting up on %s port %s' % server_address
    sock.bind(server_address)
    sock.listen(1)

    if debug:
        print >> sys.stderr, 'waiting for a connection'
    connection, client_address = sock.accept()
    try:
        if debug:
            print >> sys.stderr, 'client connected:', client_address

        data = connection.recv(read_bytes)
        if debug:
            print >> sys.stderr, 'received "%s"' % data
        if data:
            data_received_function(data)
            connection.sendall(ack_message)
        else:
            if debug:
                print >> sys.stderr, 'No Data Received!'
    finally:
        connection.close()


def send_to_down_link(message, expected_ack_message, ack_timeout_seconds):
    # Create a TCP/IP socket
    sock = socket.create_connection((send_data_server_name, send_data_server_port))
    sock.settimeout(ack_timeout_seconds)
    try:
        # Send data
        if debug:
            print >> sys.stderr, 'sending "%s"' % message
        sock.sendall(message)

        # Receive Acknowledge
        data = sock.recv(len(expected_ack_message.encode('utf-8')))
        if data != expected_ack_message:
            print(sys.stderr, "Real and Expected Ack message are not equal. Real:[" + data + "] \
                Expected:[" + expected_ack_message + "]")
        if debug:
            print >> sys.stderr, 'received "%s"' % data

    finally:
        if debug:
            print >> sys.stderr, 'closing socket'
        sock.close()

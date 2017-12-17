import httplib

import jsonrpclib
import sys
from jsonrpclib.jsonrpc import SafeTransport

from core.connection_config import server_name, server_port, debug


class TimeoutTransport(SafeTransport):
    timeout = 10.0

    def set_timeout(self, timeout):
        self.timeout = timeout

    def make_connection(self, host):
        h = httplib.HTTPConnection(host, timeout=self.timeout)
        return h


transport = TimeoutTransport()
server = jsonrpclib.Server('http://' + server_name + ':' + str(server_port),
                           transport=transport)


def wait_for_data(timeout_seconds):
    transport.set_timeout(timeout_seconds)
    response = server.wait_for_data()
    if response:
        return response
    else:
        if debug:
            print >> sys.stderr, 'No Response Received!'
        return None


def send_to_down_link(message, timeout_seconds):
    transport.set_timeout(timeout_seconds)
    response = server.send_to_down_link(message)
    if response:
        return response
    else:
        if debug:
            print >> sys.stderr, 'No Response Received!'
        return None

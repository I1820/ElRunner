from jsonrpclib.SimpleJSONRPCServer import SimpleJSONRPCServer

from core.connection_config import server_name, server_port


def start_server(wait_for_data, send_to_down_link):
    server = SimpleJSONRPCServer((server_name, server_port))
    server.register_function(wait_for_data)
    server.register_function(send_to_down_link)
    server.serve_forever()

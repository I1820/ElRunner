from jsonrpc import JSONRPCResponseManager, dispatcher
from werkzeug.serving import run_simple
from werkzeug.wrappers import Request, Response

from core.connection_config import server_name, server_port


def start_server(wait_for_data, send_to_down_link):
    @Request.application
    def application(request):
        # Dispatcher is dictionary {<method_name>: callable}
        dispatcher["Endpoint.WaitForData"] = wait_for_data
        dispatcher["Endpoint.SendToDownLink"] = send_to_down_link

        response = JSONRPCResponseManager.handle(request.data, dispatcher)
        return Response(response.json, mimetype='application/json')

    run_simple(server_name, server_port, application)

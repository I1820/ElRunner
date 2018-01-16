# coding=utf-8
# - اگر اول داده سنسور x از شی a آمد، بعد داده سنسور y از شی b آمد میانگین آنها را در پایگاه داده ذخیره کند.
import datetime
import json
import socket
import _thread

from core import connection_actions
from core.db_crud import create_one, read_one, update_one
from core.rpc_server import start_server

things_series = [{'thing_id': 'a', 'sensor_id': 'x'}, {'thing_id': 'b', 'sensor_id': 'y'}]
server_data_responses = ['{"thing_id":"a","sensor_id":"x","data":"200"}',
                         '{"thing_id":"b","sensor_id":"y","data":"300"}']
server_ack_response = 'ACK'
db_sensor_document = dict(thing_id="b", sensor_id="y", data="50", user="Mike", tags=["iot", "temperature"],
                          date=datetime.datetime.utcnow())
received = False
state = 0
final_state = 2

db_sensor_partial_doc = db_sensor_document.copy()
for e in ['data', 'user', 'tags', 'date']:
    db_sensor_partial_doc.pop(e)


def wait_for_data():
    return server_data_responses[state]


def send_to_down_link(message):
    print('server got message: ' + message)
    return server_ack_response


def init_db():
    print("create one:")
    document_id = create_one(db_sensor_document)
    print("Created with ID:", document_id)


def read_from_db():
    print("read one:")
    document = read_one(db_sensor_partial_doc)
    if document is not None:
        return document
    else:
        print("Nothing Found!")


def store_result_in_db(new_data):
    print("update one:")
    result = update_one(db_sensor_partial_doc, {"$set": {"data": new_data}})
    print(result.raw_result)


def action(data):
    print("action:")
    global state
    data_parsed_json = json.loads(data)
    thing_id = things_series[state]["thing_id"]
    sensor_id = things_series[state]["sensor_id"]
    if data_parsed_json["thing_id"] != thing_id or data_parsed_json["sensor_id"] != sensor_id:
        print("not expected thing and sensor! expected[" + thing_id + ":" + sensor_id + "] got[" +
              data_parsed_json["thing_id"] + ":" + data_parsed_json["sensor_id"] + "]")
        return
    print("Received:" + data)
    things_series[state]["data"] = data_parsed_json["data"]
    state += 1
    if state == final_state:
        doc = read_from_db()
        print("last data in DB:" + str(doc))
        new_data_value = 0
        for thing in things_series:
            new_data_value += int(thing["data"])
        new_data_value /= len(things_series)
        store_result_in_db(new_data_value)
        new_doc = read_from_db()
        print("New Doc in DB:" + str(new_doc))
        global received
        received = True


_thread.start_new(start_server, (wait_for_data, send_to_down_link))

init_db()

while True:
    try:
        print("wait for data...")
        response = connection_actions.wait_for_data(timeout_seconds=30)
        if response:
            print('Response:' + str(response))
            action(response['result'])
        else:
            print('No Response!')
    except socket.timeout:
        print("wait for data: Timeout!")
    if received:
        break

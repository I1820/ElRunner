# coding=utf-8
# - اگر داده سنسور شماره x بر روی شی a آمد،
#  به اندازه w ثانیه منتظر بماند،
#  آخرین مقدار سنسور y از شی b را خوانده و با مقدار سنسور x جمع کرده و در پایگاه داده ذخیره کند
import datetime
import json
import socket

import _thread

from core import connection_actions
from core.db_crud import read_one, create_one, update_one
from core.rpc_server import start_server
from core.time_actions import sleep

thing_id = 'a'
sensor_id = 'x'
server_data_response = '{"thing_id":"a","sensor_id":"x","data":"100"}'
server_ack_response = 'ACK'
w = 3
received = False

db_sensor_document = dict(thing_id="b", sensor_id="y", data="50", user="Mike", tags=["iot", "temperature"],
                          date=datetime.datetime.utcnow())

db_sensor_partial_doc = db_sensor_document.copy()
for e in ['data', 'user', 'tags', 'date']:
    db_sensor_partial_doc.pop(e)


def wait_for_data():
    return server_data_response


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
    print(document)
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
    data_parsed_json = json.loads(data)
    if data_parsed_json["thing_id"] != thing_id or data_parsed_json["sensor_id"] != sensor_id:
        print("not expected thing and sensor! expected[" + thing_id + ":" + sensor_id + "] got[" +
              data_parsed_json["thing_id"] + ":" + data_parsed_json["sensor_id"] + "]")
        return
    global received
    received = True
    print("before wait: " + str(datetime.datetime.utcnow()))
    sleep(w)
    print("after wait: " + str(datetime.datetime.utcnow()))
    data_parsed_json = json.loads(data)
    doc = read_from_db()
    new_data_value = int(doc["data"]) + int(data_parsed_json["data"])
    store_result_in_db(new_data_value)


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

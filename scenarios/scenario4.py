# coding=utf-8
# - اگر اول داده سنسور x از شی a آمد، بعد داده سنسور y از شی b آمد میانگین آنها را در پایگاه داده ذخیره کند.
import json
import socket

import datetime

from core.connection_actions import wait_for_data, send_to_down_link
from core.db_crud import create_one, read_one, update_one

things_series = [{'thing_id': 'a', 'sensor_id': 'x'}, {'thing_id': 'b', 'sensor_id': 'y'}]
packet_message = '{"thing_id":"b","sensor_id":"y","data":"200"}'
ack_message = 'ACK'
db_sensor_document = dict(thing_id="b", sensor_id="y", data="50", user="Mike", tags=["iot", "temperature"],
                          date=datetime.datetime.utcnow())
received = False
state = 0
final_state = 2

db_sensor_partial_doc = db_sensor_document.copy()
for e in ['data', 'user', 'tags', 'date']:
    db_sensor_partial_doc.pop(e)


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


# init_db()

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

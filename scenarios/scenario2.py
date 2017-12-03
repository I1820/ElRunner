# coding=utf-8
# - اگر داده سنسور شماره x بر روی شی a آمد،
#  به اندازه w ثانیه منتظر بماند،
#  آخرین مقدار سنسور y از شی b را خوانده و با مقدار سنسور x جمع کرده و در پایگاه داده ذخیره کند
import json
import socket

import datetime

from core.connection_actions import wait_for_data, send_to_down_link
from core.db_crud import read_one, create_one, update_one
from core.time_actions import sleep

packet_message = """{"thing_id":"a","sensor_id":"x","data":"100"}"""
ack_message = """ACK"""
w = 3

db_sensor_document = {"thing_id": "b",
                      "sensor_id": "y",
                      "data": "50",
                      "User": "Mike",
                      "tags": ["iot", "temperature"],
                      "date": datetime.datetime.utcnow()}


def init_db():
    print("create one:")
    document_id = create_one(db_sensor_document)
    print("Created with ID:", document_id)


def read_from_db():
    print("read one:")
    document = read_one({"thing_id": "b", "sensor_id": "y"})
    if document is not None:
        return document
    else:
        print("Nothing Found!")


def store_result_in_db(new_data):
    print("update one:")
    result = update_one({"thing_id": "b", "sensor_id": "y"}, {"$set": {"data": new_data}})
    print(result.raw_result)


def action(json_data):
    print("action:")
    print("before wait: " + str(datetime.datetime.utcnow()))
    sleep(w)
    print("after wait: " + str(datetime.datetime.utcnow()))
    data_parsed_json = json.loads(json_data)
    doc = read_from_db()
    new_data_value = int(doc["data"]) + int(data_parsed_json["data"])
    store_result_in_db(new_data_value)


# init_db()

try:
    print("wait for data...")
    wait_for_data(data_received_function=action, read_bytes=len(packet_message.encode('utf-8')),
                  ack_message=ack_message, timeout_seconds=30)
except socket.timeout:
    print("wait for data: Timeout!")

# try:
#     print("send to down link...")
#     send_to_down_link(message=packet_message, expected_ack_message=ack_message, ack_timeout_seconds=30)
# except socket.timeout:
#     print("send to down link: Timeout!")

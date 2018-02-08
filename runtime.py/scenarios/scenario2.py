# coding=utf-8
# - اگر داده سنسور شماره x بر روی شی a آمد،
#  به اندازه w ثانیه منتظر بماند،
#  آخرین مقدار سنسور y از شی b را خوانده و با مقدار سنسور x جمع کرده و در پایگاه داده ذخیره کند
import datetime
import json
import socket

import _thread
from time import sleep

import asyncio

from scenario import Scenario
from test.test_connection_actions import run_rpc_server

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


class Scenario2(Scenario):
    def init_db(self):
        print("create one:")
        document_id = self.create_one(db_sensor_document)
        print("Created with ID:", document_id)

    def read_from_db(self):
        print("read one:")
        document = self.read_one(db_sensor_partial_doc)
        print(document)
        if document is not None:
            return document
        else:
            print("Nothing Found!")

    def store_result_in_db(self, new_data):
        print("update one:")
        result = self.update_one(db_sensor_partial_doc, {"$set": {"data": new_data}})
        print(result.raw_result)

    def action(self, data):
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
        doc = self.read_from_db()
        new_data_value = int(doc["data"]) + int(data_parsed_json["data"])
        self.store_result_in_db(new_data_value)

    def run(self, data=None):
        self.init_db()

        while True:
            try:
                print("wait for data...")
                loop = asyncio.get_event_loop()
                future = asyncio.Future()
                loop.run_until_complete(self.wait_for_data_wrapper(future, timeout=30))
                loop.close()
                response = future.result()
                if response:
                    print('Response:' + str(response))
                    self.action(response['result'])
                else:
                    print('No Response!')
            except socket.timeout:
                print("wait for data: Timeout!")
            if received:
                break

run_rpc_server(wait_for_data, send_to_down_link)
scenario_2 = Scenario2()
scenario_2.run()

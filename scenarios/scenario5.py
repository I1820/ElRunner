# coding=utf-8
# - هر w ثانیه یک بار،
#  میانگین سنسور x از شی a را با میانگین سنسور y از شی a مقایسه کند و اگر بزرگتر بود دستور z را ارسال کند
import datetime
import socket
import thread

from core import connection_actions
from core.db_crud import create_one, read_many
from core.rpc_server import start_server
from core.time_actions import sleep

w = 2
num_recent_data = 5
things = [{'thing_id': 'a', 'sensor_id': 'x'}, {'thing_id': 'a', 'sensor_id': 'y'}]
down_link_message = '{"thing_id":"b","sensor_id":"y","data":"200"}'
server_ack_response = 'ACK'
command_sent = False


def wait_for_data():
    return 'Data from Server'


def send_to_down_link(message):
    print('server got message: ' + message)
    return server_ack_response


def db_create_one(document):
    print("create one:")
    document_id = create_one(document)
    print("Created with ID:", document_id)


def read_from_db(partial_doc):
    print("read many:")
    documents = read_many(partial_doc).sort([("date", -1)]).limit(num_recent_data)
    print("\n".join(str(doc) for doc in documents))
    documents.rewind()
    if documents is not None:
        return documents
    else:
        print("Nothing Found!")


def init_db():
    docs = []
    for i in range(1, num_recent_data + 1):
        docs.append(dict(thing_id="a", sensor_id="x", data=str(200 * i), user="Mike", tags=["iot", "temperature"],
                         date=datetime.datetime.utcnow()))
        docs.append(dict(thing_id="a", sensor_id="y", data=str(100 * i), user="Mike", tags=["iot", "temperature"],
                         date=datetime.datetime.utcnow()))
    for doc in docs:
        db_create_one(doc)


def get_average(db_cursor):
    sum_value = 0
    num = 0
    for doc in db_cursor:
        data = int(doc['data'])
        sum_value += data
        num += 1
    avg = sum_value / float(num)
    return avg


def print_command(data):
    print(data)


thread.start_new(start_server, (wait_for_data, send_to_down_link))

init_db()

while True:
    things_cursors = list()
    for j in range(len(things)):
        things_cursors.append(read_from_db(things[j]))
    averages = list()
    for cursor in things_cursors:
        averages.append(get_average(cursor))
    if averages[0] > averages[1]:
        try:
            print("send to down link...")
            response = connection_actions.send_to_down_link(message=down_link_message, timeout_seconds=30)
            if response:
                print('Response:' + str(response))
                command_sent = True
            else:
                print('No Response!')
        except socket.timeout:
            print("send to down link: Timeout!")
    if command_sent:
        print("End Scenario")
        break
    print("before wait: " + str(datetime.datetime.utcnow()))
    sleep(w)
    print("after wait: " + str(datetime.datetime.utcnow()))

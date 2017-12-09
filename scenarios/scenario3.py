# coding=utf-8
# اگر میانگین ۵ داده اخیر سنسور شماره x از شی a بزرگ از w بود یک ایمیل بفرستد
import datetime

from core.db_crud import read_many, create_one
from core.notification_actions import send_email

w = 6


def init_db():
    print("create one:")
    for i in range(10):
        document_id = create_one({"thing_id": "a",
                                  "sensor_id": "x",
                                  "data": i + 1,
                                  "User": "Mike",
                                  "tags": ["iot", "temperature"],
                                  "date": datetime.datetime.utcnow()})
        print("Created with ID:", document_id)


def read_data_from_db():
    print("read_data_from_db:")
    return read_many({"thing_id": "a", "sensor_id": "x"}).sort({"date": -1}).limit(5)


def action(data):
    print("action:")
    sender = 'ceitiotlabtest@gmail.com'
    receivers = ['ceitiotlabtest@gmail.com']

    message = 'From: From Person <ceitiotlabtest@gmail.com>\n' \
              'To: To Person <ceitiotlabtest@gmail.com>\n' \
              'Subject: Rule Engine Notification\n\n' \
              'Data:' + data + '\n' \
              'Sent by Rule Engine. Scenario:1.'
    send_email(host='smtp.gmail.com', port=587, username="ceitiotlabtest", password="ceit is the best", sender=sender,
               receivers=receivers, message=message)


# init_db()

docs = read_data_from_db()
print sum(docs["data"])/len(docs)

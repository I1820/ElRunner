# coding=utf-8
# - هر w ثانیه یک بار،
#  میانگین سنسور x از شی a را با میانگین سنسور y از شی a مقایسه کند و اگر بزرگتر بود دستور z را ارسال کند
import datetime

from core.db_crud import create_one, read_many


def db_create_one(document):
    print("create one:")
    document_id = create_one(document)
    print("Created with ID:", document_id)


def read_from_db(partial_doc):
    print("read one:")
    documents = read_many(partial_doc)
    if documents is not None:
        return documents
    else:
        print("Nothing Found!")


def init_db():
    docs = []
    for i in range(1, 6):
        docs.append(dict(thing_id="a", sensor_id="x", data=str(100 * i), user="Mike", tags=["iot", "temperature"],
                         date=datetime.datetime.utcnow()))
        docs.append(dict(thing_id="a", sensor_id="y", data=str(200 * i), user="Mike", tags=["iot", "temperature"],
                         date=datetime.datetime.utcnow()))
    for doc in docs:
        db_create_one(doc)

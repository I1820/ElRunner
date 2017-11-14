from pymongo import MongoClient

from core.db_config import db_ip, db_port, db_name, db_collection

client = MongoClient(db_ip, db_port)
db = client[db_name]
collection = db[db_collection]


def create_one(document):
    document_id = collection.insert_one(document).inserted_id
    return document_id


def create_many(documents_list):
    documents_id_list = collection.insert_many(documents_list).inserted_ids
    return documents_id_list


def read_one(partial_document):
    return collection.find_one(partial_document)


def read_many(partial_document):
    return collection.find(partial_document)


def update_one(partial_document, update_instructions):
    return collection.update_one(partial_document, update_instructions)


def update_many(partial_document, update_instructions):
    return collection.update_many(partial_document, update_instructions)


def delete_one(partial_document):
    return collection.delete_one(partial_document)


def delete_many(partial_document):
    return collection.delete_many(partial_document)

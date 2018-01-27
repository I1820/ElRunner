from pymongo import MongoClient

from core.db_config import db_ip, db_port, db_name, db_collection

client = MongoClient(db_ip, db_port)
db = client[db_name]
collection = db[db_collection]


def create_one(document):
    """
    Create a single document in database.
    :param document: Documents as a python dictionary
    :return: Created document id in database
    """
    document_id = collection.insert_one(document).inserted_id
    return document_id


def create_many(documents_list):
    """
    Create a list of documents in database.
    :param documents_list: Documents as a list of python dictionaries.
    :return: A list containing created documents ids.
    """
    documents_id_list = collection.insert_many(documents_list).inserted_ids
    return documents_id_list


def read_one(partial_document):
    """
    Read a single document from database.
    :param partial_document: A dictionary specifying the query to be performed.
    :return: A single document, or ``None`` if no matching document is found.
    """
    return collection.find_one(partial_document)


def read_many(partial_document):
    """
    Read a list of documents from database.
    :param partial_document: A dictionary specifying the query to be performed.
    :return: A Cursor object containing requested documents.
    """
    return collection.find(partial_document)


def update_one(partial_document, update_instructions):
    """
    Update a single document from database.
    :param partial_document: A query as a dictionary that matches the document to update.
    :param update_instructions: The modifications to apply.
    :return: An UpdateResult object containing update results.
    """
    return collection.update_one(partial_document, update_instructions)


def update_many(partial_document, update_instructions):
    """
    Update one or more documents that match the partial_document.
    :param partial_document: A query that matches the documents to update.
    :param update_instructions: The modifications to apply.
    :return: An UpdateResult object containing update results.
    """
    return collection.update_many(partial_document, update_instructions)


def delete_one(partial_document):
    """
    Delete a single document matching the partial_document.
    :param partial_document: A query that matches the document to delete.
    :return: A DeleteResult object containing delete results.
    """
    return collection.delete_one(partial_document)


def delete_many(partial_document):
    """
    Delete one or more documents matching the partial_document.
    :param partial_document: A query that matches the documents to delete.
    :return: A DeleteResult object containing delete results.
    """
    return collection.delete_many(partial_document)

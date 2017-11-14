import datetime

from core.db_crud import create_one, read_one, create_many, read_many, update_one

test_document = {"ID": "123456789",
                 "User": "Mike",
                 "Device ID": "IAU250",
                 "Message": "device message test 1",
                 "tags": ["iot", "temperature"],
                 "date": datetime.datetime.utcnow()}

test_document_2 = {"ID": "1234567890",
                   "User": "Mike",
                   "Device ID": "EOA150",
                   "Message": "device message test 2",
                   "tags": ["iot", "temperature"],
                   "date": datetime.datetime.utcnow()}

test_document_3 = {"ID": "12345678900",
                   "User": "John",
                   "Device ID": "EOA150",
                   "Message": "device message test 3",
                   "tags": ["iot", "temperature"],
                   "date": datetime.datetime.utcnow()}


def create_one_test():
    print("create_one_test:")
    document_id = create_one(test_document)
    print("Created with ID:", document_id)


def create_many_test():
    print("create_many_test:")
    documents_id_list = create_many([test_document_2, test_document_3])
    print("Created with IDs:", documents_id_list)


def read_one_test():
    print("read_one_test:")
    document = read_one({"ID": "123456789"})
    if document is not None:
        print("Found:", document)
    else:
        print("Nothing Found!")


def read_many_test():
    print("read_many_test:")
    documents = read_many({"User": "Mike"})
    print("Documents:")
    for document in documents:
        print(document)


def update_one_test():
    print("update_one_test:")
    result = update_one({"ID": "123456789"}, {"$set": {"Message": "UPDATED device message test 1"}})
    print(result.raw_result)


create_one_test()
create_many_test()
read_one_test()
read_many_test()
update_one_test()

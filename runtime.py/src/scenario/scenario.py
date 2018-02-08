# In The Name Of God
# ========================================
# [] File Name : codec.py
#
# [] Creation Date : 15-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import abc
import aiohttp
import time
import threading
import smtplib
from urllib.error import HTTPError
from urllib.request import urlopen
from pymongo import MongoClient

import async_timeout

RPC_SERVER = '127.0.0.1'
RPC_PORT = 1373
URL = 'http://{}:{}/'.format(RPC_SERVER, RPC_PORT)
HEADERS = {'content-type': 'application/json'}
PAYLOAD = {'jsonrpc': '2.0'}

DB_IP = "localhost"
DB_PORT = 27017
DB_NAME = "datadb"
DB_COLLECTION = "datacollection"

client = MongoClient(DB_IP, DB_PORT)
db = client[DB_NAME]
collection = db[DB_COLLECTION]


class Scenario(metaclass=abc.ABCMeta):
    def sleep(self, seconds):
        time.sleep(seconds)

    def schedule(self, delay_seconds, action_function, args=()):
        threading.Timer(delay_seconds, action_function, args).start()

    async def wait_for_data(self, timeout):
        """
        Requests data from server.
        :param timeout: Timeout time until response (data) in seconds.
        """
        request_payload = PAYLOAD.copy()
        request_payload['method'] = 'Endpoint.WaitForData'
        request_payload['params'] = []
        request_payload['id'] = 0

        async with aiohttp.ClientSession() as session:
            session.headers = HEADERS
            with async_timeout.timeout(timeout):
                async with session.post(URL, json=request_payload,
                                        timeout=timeout) as response:
                    return await response.json()

    async def send_to_down_link(self, message, timeout):
        """
        Send data containing commands to server.
        :param message: Message data to be sent to server.
        :param timeout: Timeout time until response (acknowledge) in seconds.
        :return: Response that shows whether message is received or not.
        """
        request_payload = PAYLOAD.copy()
        request_payload['method'] = 'Endpoint.SendToDownLink'
        request_payload['params'] = [message]
        request_payload['id'] = 1

        async with aiohttp.ClientSession() as session:
            session.headers = HEADERS
            response = await session.post(URL, json=request_payload,
                                          timeout=timeout)
            return await response.json()

    def send_email(self, host, port, username, password, sender, receivers, message):
        """
        Send email using given host to some receivers.
        :param host: Email server host name or ip.
        :param port: Email server port number.
        :param username: Email account username
        :param password: Email account password.
        :param sender: Sender email address
        :param receivers: Receivers' email list
        :param message: Email body message containing from address, to address, subject and body.
        :return: True if email is sent or False otherwise.
        """
        successful = False
        try:
            smtp_obj = smtplib.SMTP(host=host, port=port)
            smtp_obj.starttls()
            smtp_obj.login(user=username, password=password)
            smtp_obj.sendmail(sender, receivers, message)
            smtp_obj.quit()
            print("Successfully sent email")
            successful = True
        except smtplib.SMTPException as e:
            print(e)

        return successful

    def send_sms(self, username, password, from_number, to_number, message):
        """
        Send SMS message to a cell phone using SMS panel.
        :param username: SMS panel username.
        :param password: SMS panel password.
        :param from_number: Sender panel number.
        :param to_number: Receiver phone number.
        :param message: Message body to be sent.
        :return: True if SMS is sent or False otherwise.
        """
        successful = False
        # Use payam-resan.com sms panel url format to send SMS
        url = "http://www.payam-resan.com/APISend.aspx?Username={0}&Password={1}&From={2}&To={3}&Text={4}" \
            .format(username, password, from_number, to_number, message).replace(" ", "%20")
        try:
            # Send GET request to send SMS with SMS Panel
            data = urlopen(url=url).read()
            if str(data) == '0':
                print("Error Sending SMS: check url and credentials")
            else:
                print("Successfully sent SMS")
                successful = True
        except HTTPError as e:
            print(e)

        return successful

    def create_one(self, document):
        """
        Create a single document in database.
        :param document: Documents as a python dictionary
        :return: Created document id in database
        """
        document_id = collection.insert_one(document).inserted_id
        return document_id

    def create_many(self, documents_list):
        """
        Create a list of documents in database.
        :param documents_list: Documents as a list of python dictionaries.
        :return: A list containing created documents ids.
        """
        documents_id_list = collection.insert_many(documents_list).inserted_ids
        return documents_id_list

    def read_one(self, partial_document):
        """
        Read a single document from database.
        :param partial_document: A dictionary specifying the query to be performed.
        :return: A single document, or ``None`` if no matching document is found.
        """
        return collection.find_one(partial_document)

    def read_many(self, partial_document):
        """
        Read a list of documents from database.
        :param partial_document: A dictionary specifying the query to be performed.
        :return: A Cursor object containing requested documents.
        """
        return collection.find(partial_document)

    def update_one(self, partial_document, update_instructions):
        """
        Update a single document from database.
        :param partial_document: A query as a dictionary that matches the document to update.
        :param update_instructions: The modifications to apply.
        :return: An UpdateResult object containing update results.
        """
        return collection.update_one(partial_document, update_instructions)

    def update_many(self, partial_document, update_instructions):
        """
        Update one or more documents that match the partial_document.
        :param partial_document: A query that matches the documents to update.
        :param update_instructions: The modifications to apply.
        :return: An UpdateResult object containing update results.
        """
        return collection.update_many(partial_document, update_instructions)

    def delete_one(self, partial_document):
        """
        Delete a single document matching the partial_document.
        :param partial_document: A query that matches the document to delete.
        :return: A DeleteResult object containing delete results.
        """
        return collection.delete_one(partial_document)

    def delete_many(self, partial_document):
        """
        Delete one or more documents matching the partial_document.
        :param partial_document: A query that matches the documents to delete.
        :return: A DeleteResult object containing delete results.
        """
        return collection.delete_many(partial_document)

    @abc.abstractmethod
    def run(self, data=None):
        pass

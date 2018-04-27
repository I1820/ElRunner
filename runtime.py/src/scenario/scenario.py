# In The Name Of God
# ========================================
# [] File Name : codec.py
#
# [] Creation Date : 15-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import abc
import time
import threading
import smtplib
import os
import redis
import aiohttp
import pymongo

# RPC requirements
RPC_SERVER = '127.0.0.1'
RPC_PORT = 1373
URL = 'http://{}:{}/'.format(RPC_SERVER, RPC_PORT)
HEADERS = {'Content-Type': 'application/json', 'Accept': 'application/json'}
PAYLOAD = {'jsonrpc': '2.0'}


class Scenario(metaclass=abc.ABCMeta):
    counter = 0

    def __init__(self, id):
        self.id = id
        try:
            self.redis = redis.Redis(host=os.environ['REDIS_HOST'])
        except KeyError:
            self.redis = None

        try:
            self.data_db = pymongo.MongoClient(os.environ['MONGO_URL'])
        except Exception:
            self.data_db = None

    @staticmethod
    def sleep(seconds):
        time.sleep(seconds)

    @staticmethod
    def schedule(delay_seconds, action_function, args=()):
        threading.Timer(delay_seconds, action_function, args).start()

    def find_data(self, thingid):
        return self.data_db.find_many({
            'thingid': thingid,
            'project': os.environ['NAME'],
        })

    async def wait_for_data(self, timeout):
        """
        Requests data from server.
        :param timeout: Timeout time until response (data) in seconds.
        """
        request_payload = PAYLOAD.copy()
        request_payload['method'] = 'Endpoint.WaitForData'
        request_payload['id'] = self.counter
        self.counter += 1

        async with aiohttp.ClientSession(headers=HEADERS) as session:
            response = await session.post(URL, json=request_payload,
                                          timeout=timeout)
            json = await response.json()
            return json['result']

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
        request_payload['id'] = self.counter
        self.counter += 1

        async with aiohttp.ClientSession(headers=HEADERS) as session:
            response = await session.post(URL, json=request_payload,
                                          timeout=timeout)
            json = await response.json()
            return json['result']

    @staticmethod
    def send_email(host, port, username, password, sender,
                   receivers, message):
        """
        Send email using given host to some receivers.
        :param host: Email server host name or ip.
        :param port: Email server port number.
        :param username: Email account username
        :param password: Email account password.
        :param sender: Sender email address
        :param receivers: Receivers' email list
        :param message: Email body message containing from address,
        to address, subject and body.
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
        except smtplib.SMTPException as exception:
            print(exception)

        return successful

    @abc.abstractmethod
    def run(self, data=None):
        pass

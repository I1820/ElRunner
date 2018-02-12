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
import redis
import os

from urllib.error import HTTPError
from urllib.request import urlopen


# RPC requirements
RPC_SERVER = '127.0.0.1'
RPC_PORT = 1373
URL = 'http://{}:{}/'.format(RPC_SERVER, RPC_PORT)
HEADERS = {'Content-Type': 'application/json', 'Accept': 'application/json'}
PAYLOAD = {'jsonrpc': '2.0'}


class Scenario(metaclass=abc.ABCMeta):
    counter = 0

    def __init__(self):
        try:
            self.redis = redis.Redis(host=os.environ['REDIS_HOST'])
        except KeyError:
            self.redis = None

    @staticmethod
    def sleep(seconds):
        time.sleep(seconds)

    @staticmethod
    def schedule(delay_seconds, action_function, args=()):
        threading.Timer(delay_seconds, action_function, args).start()

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

    def send_email(self, host, port, username, password, sender,
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
        url = "http://www.payam-resan.com/APISend.aspx?Username={0}&Password={1}'\
            '&From={2}&To={3}&Text={4}" \
            .format(username, password, from_number, to_number, message)\
            .replace(" ", "%20")
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

    @abc.abstractmethod
    def run(self, data=None):
        pass

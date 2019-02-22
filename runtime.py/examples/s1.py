'''
Scenario 1:
This scenario sends email that contains the given data.
Message contains current time and date.
'''
from scenario import Scenario
from datetime import datetime


class S1(Scenario):
    async def run(self, data=None):
        sender = ''
        receivers = ['parham.alvani@gmail.com']

        message = 'From: From I1820 <>\n' \
                  'To: To Parham Alvani <parham.alvani@gmail.com>\n' \
                  'Subject: Rule Engine Notification \
            [Thing: ' + self.id + ']\n\n' \
                  'Data:' + str(data) + '\n' \
                  'Date-Time:' + str(datetime.now()) + '\n' \
            'Sent by Rule Engine. Scenario-1.'
        self.send_email(host='', port=587,
                        username='',
                        password='',
                        sender=sender,
                        receivers=receivers, message=message)

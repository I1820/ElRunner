'''
Scenario 6:
This scenario relay on two asset, time and count.
it sends email when count is multiple of 10 * 60 * 1000.
Avidnet uses this script for its load tests. In these tests avidnet compares
datetime with at in each mail.
'''

from scenario import Scenario
from datetime import datetime


class S6(Scenario):
    async def run(self, data=None):
        if data['asset'] != 'count':
            return
        count = data['raw']
        if int(count) % (10 * 60 * 100) != 0:
            return

        sender = ''
        receivers = ['parham.alvani@gmail.com']

        message = 'From: From I1820 <>\n' \
                  'To: To Parham Alvani <parham.alvani@gmail.com>\n' \
                  'Subject: Load Generator Notification' \
            '[Thing: ' + self.id + ']\n\n' \
                  'Count: ' + count + '\n' \
                  'Date-Time: ' + str(datetime.now()) + '\n' \
                  'At: ' + data['at'] + '\n' \
            'Sent by Rule Engine. Scenario-6.'
        self.send_email(host='', port=587,
                        username='',
                        password='',
                        sender=sender,
                        receivers=receivers, message=message)

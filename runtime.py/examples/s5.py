import asyncio
import time

from scenario import Scenario
from kavenegar import KavenegarAPI


class S5(Scenario):
    def run(self, data=None):
        # Wait for data
        loop = asyncio.get_event_loop()
        t = asyncio.ensure_future(self.wait_for_data(timeout=30))
        loop.run_until_complete(t)
        loop.close()
        response = t.result()
        if response:
            sender = 'ceitiotlabtest@gmail.com'
            receivers = ['parham.alvani@gmail.com']

            message = 'From: From Person <ceitiotlabtest@gmail.com>\n' \
                'To: To Person <ceitiotlabtest@gmail.com>\n' \
                'Subject: Rule Engine Notification\n\n' \
                'Data:' + str(response) + '\n' \
                'Sent by Rule Engine. Scenario:1.'
            self.send_email(host='smtp.gmail.com', port=587,
                            username="ceitiotlabtest",
                            password="ceit is the best",
                            sender=sender,
                            receivers=receivers, message=message)

        # Wait 10 milliseconds
        time.sleep(10)
        api = KavenegarAPI('Your APIKey')
        params = {
            # optional
            'sender': '',
            # multiple mobile number, split by comma
            'receptor': '09390909540',
            'message': 'helloooo',
        }
        api.sms_send(params)

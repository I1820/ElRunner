from scenario import Scenario


class S1(Scenario):
    async def run(self, data=None):
        sender = 'platform.avidnetco@gmail.com'
        receivers = ['parham.alvani@gmail.com']

        message = 'From: From Avidnet <platform.avidnetco@gmail.com>\n' \
                  'To: To Parham Alvani <parham.alvani@gmail.com>\n' \
                  'Subject: Rule Engine Notification [Thing: ' + self.id + ']\n\n' \
                  'Data:' + str(data) + '\n' \
                                        'Sent by Rule Engine. Scenario-1.'
        self.send_email(host='smtp.gmail.com', port=587,
                        username="platform.avidnetco@gmail.com",
                        password="fancopass(1397)",
                        sender=sender,
                        receivers=receivers, message=message)

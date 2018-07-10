from scenario import Scenario


class S1(Scenario):
    def run(self, data=None):
        sender = 'ceitiotlabtest@gmail.com'
        receivers = ['parham.alvani@gmail.com']

        message = 'From: From AIoTRC <ceitiotlabtest@gmail.com>\n' \
                  'To: To Parham Alvani <parham.alvani@gmail.com>\n' \
                  'Subject: Rule Engine Notification [Thing: ' + self.id + ']\n\n' \
                  'Data:' + str(data) + '\n' \
                                        'Sent by Rule Engine. Scenario:1.'
        self.send_email(host='smtp.gmail.com', port=587,
                        username="ceitiotlabtest",
                        password="ceit is the best",
                        sender=sender,
                        receivers=receivers, message=message)

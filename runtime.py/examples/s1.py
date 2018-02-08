from scenario import Scenario


class S1(Scenario):
    def run(self, data):
        sender = 'ceitiotlabtest@gmail.com'
        receivers = ['parham.alvani@gmail.com']

        message = 'From: From Person <ceitiotlabtest@gmail.com>\n' \
                  'To: To Person <ceitiotlabtest@gmail.com>\n' \
                  'Subject: Rule Engine Notification\n\n' \
                  'Data:' + str(data) + '\n' \
                                        'Sent by Rule Engine. Scenario:1.'
        self.send_email(host='smtp.gmail.com', port=587,
                        username="ceitiotlabtest",
                        password="ceit is the best",
                        sender=sender,
                        receivers=receivers, message=message)

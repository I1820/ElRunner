from core.notification_actions import send_email, send_sms


def send_email_test():
    print("send_email_test")
    sender = 'test@gmail.com'
    receivers = ['test@gmail.com']

    message = """From: From Person <test@gmail.com>
    To: To Person <test@gmail.com>
    Subject: SMTP e-mail test

    This is a test e-mail message.
    """
    send_email(host='smtp.gmail.com', port=587, username="username", password="password", sender=sender,
               receivers=receivers, message=message)


def send_sms_test():
    print("send_sms_test")
    send_sms(username="user", password="pass", from_number="-1", to_number="09121111111",
             message="test message!")


send_email_test()
send_sms_test()

import smtplib


def send_email(host, port, username, password, sender, receivers, message):
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


data = '{"thing_id":"a","sensor_id":"x","data":"100"}'
sender = 'ceitiotlabtest@gmail.com'
receivers = ['ceitiotlabtest@gmail.com']

message = 'From: From Person <ceitiotlabtest@gmail.com>\n' \
          'To: To Person <ceitiotlabtest@gmail.com>\n' \
          'Subject: Rule Engine Notification\n\n' \
          'Data:' + data + '\n' \
                           'Sent by Rule Engine. Scenario:1.'
send_email(host='smtp.gmail.com', port=587, username="ceitiotlabtest", password="ceit is the best", sender=sender,
           receivers=receivers, message=message)
import smtplib
from urllib.error import HTTPError
from urllib.request import urlopen


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


def send_sms(username, password, from_number, to_number, message):
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
    url = "http://www.payam-resan.com/APISend.aspx?Username={0}&Password={1}&From={2}&To={3}&Text={4}" \
        .format(username, password, from_number, to_number, message).replace(" ", "%20")
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

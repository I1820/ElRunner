import smtplib
from urllib.error import HTTPError
from urllib.request import urlopen


def send_email(host, port, username, password, sender, receivers, message):
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

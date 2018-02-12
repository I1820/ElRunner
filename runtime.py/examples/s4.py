from scenario import Scenario
from kavenegar import KavenegarAPI


class S4(Scenario):
    def run(self, data=None):
        v = self.redis.get("Her")
        if v is None:
            self.redis.set("Her", 0)
        else:
            self.redis.set("Her", int(v) + 1)

        api = KavenegarAPI('Your APIKey')
        params = {
            # optional
            'sender': '',
            # multiple mobile number, split by comma
            'receptor': '09390909540',
            'message': str(v),
        }
        api.sms_send(params)

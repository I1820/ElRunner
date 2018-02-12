from scenario import Scenario


class S3(Scenario):
    def run(self, data=None):
        v = self.redis.get("Her")
        if v is None:
            self.redis.set("Her", 0)
        else:
            self.redis.set("Her", int(v) + 1)

        f = open('/tmp/redis', 'w+')
        f.write(str(v))
        f.write('\n')

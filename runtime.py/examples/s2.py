import asyncio

from scenario import Scenario


class S2(Scenario):
    def run(self, data=None):
        f = open('/tmp/rpc', 'w+')
        f.write(str(data))
        f.write('\n')

        loop = asyncio.get_event_loop()
        t = asyncio.ensure_future(self.wait_for_data(timeout=30))
        loop.run_until_complete(t)
        loop.close()
        response = t.result()
        if response:
            f.write('Response:' + str(response))
        else:
            f.write('No Response!')
        f.close()

# In The Name Of God
# ========================================
# [] File Name : codec.py
#
# [] Creation Date : 15-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import abc


class Codec(metaclass=abc.ABCMeta):
    thing_location = ''

    @staticmethod
    def create_location(lat, lng):
        return {
            'type': 'Point',
            'coordinates': [lat, lng]
        }

    @abc.abstractmethod
    def decode(self, data):
        pass

    @abc.abstractmethod
    def encode(self, data):
        pass

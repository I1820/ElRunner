# In The Name Of God
# ========================================
# [] File Name : codec.py
#
# [] Creation Date : 15-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import abc


class CodecBase(abc.ABCMeta):
    def __init__(self, name, bases, namespace, **kwargs):
        abc.ABCMeta.__init__(self, name, bases, namespace)

    def __new__(cls, name, bases, namespace):
        instance = abc.ABCMeta.__new__(
            cls, name, bases, namespace)

        print(instance)
        globals()['codecInstance'] = instance

        return instance


class Codec(metaclass=CodecBase):
    @abc.abstractmethod
    def decode(self, data):
        pass

    @abc.abstractmethod
    def encode(self, data):
        pass

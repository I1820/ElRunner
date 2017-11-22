# In The Name Of God
# ========================================
# [] File Name : codec.py
#
# [] Creation Date : 15-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import subprocess
import abc
import importlib


class CodecBase(abc.ABCMeta):
    sub_class = None

    def __init__(self, name, bases, namespace, **kwargs):
        abc.ABCMeta.__init__(self, name, bases, namespace)

    def __new__(cls, name, bases, namespace, requirements):
        instance = abc.ABCMeta.__new__(
            cls, name, bases, namespace)

        if len(requirements) != 0:
            for requirement in requirements:
                setattr(instance, requirement,
                        importlib.import_module(requirement))

        cls.sub_class = instance
        return instance


class Codec(metaclass=CodecBase, requirements=[]):
    @staticmethod
    def get():
        return CodecBase.sub_class

    @abc.abstractmethod
    def decode(self, data):
        pass

    @abc.abstractmethod
    def encode(self, data):
        pass

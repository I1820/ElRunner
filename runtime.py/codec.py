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


class Codec(metaclass=abc.ABCMeta):
    sub_class = None

    @classmethod
    def __init_subclass__(cls, requirements, **kwargs):
        super().__init_subclass__(**kwargs)
        subprocess.run(['pip3', '-qqq', 'install'] + requirements)
        for requirement in requirements:
            setattr(cls, requirement, importlib.import_module(requirement))
        Codec.sub_class = cls

    @classmethod
    def get(cls):
        return cls.sub_class

    @abc.abstractmethod
    def decode(self, data):
        pass

    @abc.abstractmethod
    def encode(self, data):
        pass

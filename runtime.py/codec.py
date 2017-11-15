# In The Name Of God
# ========================================
# [] File Name : codec.py
#
# [] Creation Date : 15-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import pip


class Codec:
    @classmethod
    def __init_subclass__(cls, requirements, **kwargs):
        super().__init_subclass__(**kwargs)
        pip.main(['install', *requirements])
        print(cls)

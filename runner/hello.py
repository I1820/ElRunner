# In The Name Of God
# ========================================
# [] File Name : hello.py
#
# [] Creation Date : 05-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import time
import base64

s = input()
s = base64.b64decode(s).decode('ascii')
time.sleep(1)
print("hello from python", s)

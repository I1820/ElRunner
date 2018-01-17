# In The Name Of God
# ========================================
# [] File Name : hello.py
#
# [] Creation Date : 05-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
'''
Hello module just for saying hello from python to go
'''
import time
import base64

# Input
s = input()
s = base64.b64decode(s).decode('ascii')

# Thinking ...
time.sleep(1)

# Output
print("hello from python", s)

from enum import Enum
import sys

class Level(Enum):
    Debug=0
    Info=1
    Warn=2
    Error=3
    Fatal=4
    Panic=5
    NoLevel=6

def Log(lvl:Level, msg:str):
    bytes_lvl = lvl.value.to_bytes(1, 'little', signed=False)
    sys.stderr.buffer.write(bytes_lvl)
    msg_encoded = str.encode(msg)
    bytes_size = len(msg_encoded).to_bytes(4, 'little', signed=False)
    sys.stderr.buffer.write(bytes_size)
    sys.stderr.buffer.write(msg_encoded)

def Result(code:int, msg:str):
    bytes_val = code.to_bytes(4, 'little', signed=True)
    sys.stdout.buffer.write(bytes_val)
    sys.stdout.buffer.write(str.encode(msg))
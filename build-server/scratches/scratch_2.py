#!/usr/bin/python3


def mult(a, b):
    def inner(b,base):
        if (b==0): return base
        else: return inner(b-1,base+a)
    return inner(b,0)

print(mult(8,2))


for x in range(10):
    for y in range(x):
        print(y)

for x in range(1, 10, 3):
    print(x)

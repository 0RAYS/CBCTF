from Crypto.Util.number import bytes_to_long, getPrime
from secret import flag

m = bytes_to_long(flag)
p = getPrime(2048)
q = getPrime(2048)
n1 = p*q
e1 = 0x3
c1 = pow(m,e1,n1)
print(n1)
print(c1)

'''
n1 = {given_n1}
c1 = {given_c1}
'''

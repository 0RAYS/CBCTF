from Crypto.Util.number import bytes_to_long, getPrime
from secret import flag
m = bytes_to_long(flag)
p = getPrime(128)
q = getPrime(128)
n = p * q
e = 65537
c = pow(m,e,n)
print(n,c)
# {given_n}
# {given_c}

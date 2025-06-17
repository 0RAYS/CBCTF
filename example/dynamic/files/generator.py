import io
import zipfile


def generate_attachment(flags: bytes) -> bytes:
    n1, c1 = generate_given(flags)
    with open("template.py", "r") as f:
        template = f.read()
    attachment = template.format(given_n1=n1, given_c1=c1)
    byte = io.BytesIO()
    with zipfile.ZipFile(byte, "w", compression=zipfile.ZIP_DEFLATED) as f:
        f.writestr("attachment.py", attachment)
    byte.seek(0)
    return byte.read()


def generate_given(flag: bytes):
    from Crypto.Util.number import bytes_to_long, getPrime

    m = bytes_to_long(flag)
    p = getPrime(2048)
    q = getPrime(2048)
    n1 = p*q
    e1 = 0x3
    c1 = pow(m,e1,n1)
    return n1, c1

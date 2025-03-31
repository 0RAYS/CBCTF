import sys
import os
import zipfile
import base64


def generate_given(flag: str):
    from Crypto.Util.number import bytes_to_long, getPrime

    m = bytes_to_long(flag.encode())
    p = getPrime(64)
    q = getPrime(64)
    n = p * q
    e = 65537
    c = pow(m, e, n)
    return n, c


def generate_attachment(uuid: str, **kwargs):
    with open("template.py", "r") as f:
        template = f.read()
    attachment = template.format(**kwargs)
    if not os.path.exists("attachments"):
        os.mkdir("attachments")
    with zipfile.ZipFile(f"attachments/{uuid}.zip", "w", compression=zipfile.ZIP_DEFLATED) as f:
        f.writestr("attachment.py", attachment)


if __name__ == "__main__":
    uuid, flag = sys.argv[1], base64.b64decode(sys.argv[2]).split(",")[0]
    n, c = generate_given(flag)
    generate_attachment(uuid, given_n=n, given_c=c)

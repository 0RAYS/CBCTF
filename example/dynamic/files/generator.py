import sys
import os
import zipfile
import base64


def generate_given(flag: bytes):
    from Crypto.Util.number import bytes_to_long, getPrime

    m = bytes_to_long(flag)
    p = getPrime(2048)
    q = getPrime(2048)
    n1 = p*q
    e1 = 0x3
    c1 = pow(m,e1,n1)
    return n1, c1


def generate_attachment(team_id: str, **kwargs):
    with open("template.py", "r") as f:
        template = f.read()
    attachment = template.format(**kwargs)
    if not os.path.exists("mnt/attachments"):
        os.mkdir("mnt/attachments")
    with zipfile.ZipFile(f"mnt/attachments/{team_id}.zip", "w", compression=zipfile.ZIP_DEFLATED) as f:
        f.writestr("attachment.py", attachment)


if __name__ == "__main__":
    team_id, flag = sys.argv[1], base64.b64decode(sys.argv[2])
    n1, c1 = generate_given(flag)
    generate_attachment(team_id, given_n1=n1, given_c1=c1)

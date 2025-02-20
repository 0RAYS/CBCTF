FROM docker.0rays.club/test/python:3.10-slim

COPY ./files/* /root/
COPY ./requirements.txt /root/requirements.txt

RUN apt-get update &&  \
    apt-get install -y unzip zip &&  \
    pip install -r /root/requirements.txt &&  \
    chmod +x /root/run.sh

WORKDIR /root

CMD ["ls -al /root"]

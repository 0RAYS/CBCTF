FROM python:3.10-slim

COPY ./files/* /root/

RUN apt-get update &&  \
    apt-get install -y unzip zip &&  \
    pip install -r /root/requirements.txt &&  \
    chmod +x /root/run.sh

WORKDIR /root

CMD ["ls -al /root"]

FROM python:3.10-slim

COPY ./files/* /root/

RUN apt-get update &&  \
    apt-get install -y unzip &&  \
    pip install -r /root/requirements.txt

EXPOSE 8000

WORKDIR /root

CMD ["python", "app.py"]

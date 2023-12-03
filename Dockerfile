FROM selenium/standalone-chrome:4.15.0
USER 0

RUN apt update && apt install -y python3 python3-pip

WORKDIR /work
COPY ./requirements.txt /work
RUN ["pip3", "install", "-r", "./requirements.txt"]

COPY . /work

CMD ["python3", "main.py"]

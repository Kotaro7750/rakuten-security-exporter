FROM mcr.microsoft.com/playwright/python:v1.52.0
USER 0

WORKDIR /debug_info
WORKDIR /csv

WORKDIR /work
COPY ./requirements.txt /work
RUN ["pip3", "install", "-r", "./requirements.txt"]
RUN ["playwright", "install", "chromium"]

COPY . /work

CMD ["python3", "main.py"]

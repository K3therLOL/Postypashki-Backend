FROM python:3.12-slim

WORKDIR /app
COPY . .
RUN python -m pip install pika --upgrade
CMD ["python", "consumer.py"]

#!/usr/bin/env python3

import pika, os, sys
import json

from image.image import Storage, download, process

storage = Storage(
    os.getenv("ENDPOINT", ""),
    os.getenv("ACCESS_KEY", ""),
    os.getenv("SECRET_KEY", ""),
    os.getenv("BUCKET", ""),
    os.getenv("REGION", ""),
    os.getenv("DB_CONN_STRING", "")
)

def test_image():
    print("test_image started")

    img_url = "https://i.ytimg.com/vi/fbByrLJ3ehw/maxresdefault.jpg"
    img_src = download(img_url)
    print(img_src)
    new_img = process(img_src, "sharpen")
    print(new_img)
    #storage.save(new_img)


def execute_image_pipeline(body: str):
    img_attrs = json.loads(body)
    img_src = download(img_attrs["image_url"])

    params = img_attrs.get("parameters")
    if params is None:
        params = {}
    new_img = process(img_src, img_attrs["filter"], **params)
    storage.save(new_img, img_attrs["task_id"])


def main():
    print("Connecting to rabbitmq...")
    connection = pika.BlockingConnection(pika.ConnectionParameters("rabbitmq"))
    print("Connected")
    channel = connection.channel()
    queue_name = os.getenv("QUEUE_NAME")
    channel.queue_declare(queue=queue_name, durable=True, arguments={"x-queue-type": "quorum"})

    def callback(ch, method, properties, body):
        print(f" [x] Received {body}")
        execute_image_pipeline(body)

    channel.basic_consume(queue=queue_name, on_message_callback=callback, auto_ack=True)
    channel.start_consuming()


if __name__ == "__main__":
    try:
        test_image()
        main()
    except KeyboardInterrupt:
        print("Interrupted")
        try:
            sys.exit(0)
        except SystemExit:
            os._exit(0)

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


def execute_image_pipeline(body: str) -> dict[str, str]:
    img_attrs = json.loads(body)
    img_src = download(img_attrs["image_url"])

    params = img_attrs.get("parameters")
    if params is None:
        params = {}
    
    print("image process start")
    new_img = process(img_src, img_attrs["filter"], **params)
    print("image processed")
    result = storage.save(new_img, img_attrs["task_id"])
    return result


def reply_back(ch, method, properties, result):
    ch.basic_ack(delivery_tag=method.delivery_tag)
    ch.basic_publish(
        exchange="",
        routing_key=properties.reply_to,
        properties=pika.BasicProperties(
            correlation_id=properties.correlation_id
        ),
        body=json.dumps(result)
    )
    print("Replied back")


def main():
    print("Connecting to rabbitmq...")
    connection = pika.BlockingConnection(pika.ConnectionParameters("rabbitmq"))
    print("Connected")
    channel = connection.channel()
    queue_name = os.getenv("QUEUE_NAME")
    reply_queue_name = os.getenv("REPLY_QUEUE_NAME")
    channel.queue_declare(queue=queue_name, durable=True, arguments={"x-queue-type": "quorum"})
    channel.queue_declare(queue=reply_queue_name, durable=False)

    def callback(ch, method, properties, body):
        print(f" [x] Received {body}")
        result = execute_image_pipeline(body)
        reply_back(ch, method, properties, result)

    channel.basic_consume(queue=queue_name, on_message_callback=callback)
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

import tls_client
import psycopg
import boto3
import uuid

from urllib.parse import urlparse
from PIL import Image
from io import BytesIO
from datetime import datetime, timedelta, timezone

class ImageStorage:
    def __init__(self, endpoint_url: str, access_key: str, secret_key: str, bucket: str, db_url: str):
        self.bucket = bucket
        self.s3 = boto3.client(
            "s3",
            endpoint_url=endpoint_url,
            access_key=access_key,
            secret_key=secret_key
        )

        self.db = psycopg.connect(db_url)


    def save(self, img: Image.Image):
        key = f"images/{uuid.uuid4()}"

        buffer = BytesIO()
        img.save(buffer)
        buffer.seek(0)
        self.s3.upload_fileobj(
            buffer,
            self.bucket,
            key
        )

        url = self.s3.generate_presigned_url(
            "get_object",
            Params={
                "Bucket": self.bucket,
                "Key": key,
            },
        )

        expires_at = datetime.now(timezone.utc) + timedelta(minutes=60)
        with self.db.cursor() as cur:
            cur.execute(
                """
                INSERT INTO images (url, expires_at)
                VALUES (%s, %s)
                """,
                (url, expires_at)
            )

        self.db.commit()


def download(url: str) -> Image.Image:
    session = tls_client.Session(
        client_identifier="chrome_120",
        random_tls_extension_order=True
    )

    parsed_url = urlparse(url)
    session.headers.update({
        "Host": f"cdn.{parsed_url}",
        "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:149.0) Gecko/20100101 Firefox/149.0",
        "Accept": "image/avif,image/webp,image/png,image/svg+xml,image/*;q=0.8,*/*;q=0.5",
        "Accept-Language": "en-US,en;q=0.9",
        "Accept-Encoding": "gzip, deflate, br, zstd",
        "Referer": f"https://{parsed_url}",
    })

    resp = session.get(url, timeout=10)
    resp.raise_for_status()
    return Image.open(BytesIO(resp.content))


def process(img: Image.Image, filter: str, **args: dict[str, str]) -> Image.Image:
    # TO DO
    return img


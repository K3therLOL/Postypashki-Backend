import tls_client
import psycopg
import boto3
import uuid
import requests

from urllib.parse import urlparse
from PIL import Image, ImageFilter, ImageOps
from io import BytesIO
from datetime import datetime, timedelta, timezone


class Storage:
    # access to s3 and postgres
    def __init__(self, endpoint_url: str, access_key: str, secret_key: str, bucket: str, region: str, db_dsn: str):
        self.bucket = bucket
        self.s3 = boto3.client(
            "s3",
            endpoint_url=endpoint_url,
            aws_access_key_id=access_key,
            aws_secret_access_key=secret_key,
            region_name=region
        )

        self.db = psycopg.connect(db_dsn)


    # uploading to s3
    def save(self, img: Image.Image, task_id: str) -> dict[str, str]:
        key = f"images/{uuid.uuid4()}"

        buffer = BytesIO()
        fmt = img.format or "PNG"
        img.save(buffer, format=fmt)
        buffer.seek(0)
        self.s3.upload_fileobj(
            buffer,
            self.bucket,
            key,
            ExtraArgs={
                "ContentType": f"image/{fmt}",
                "ContentDisposition": "inline",
            },
        )

        saved_img_url = self.s3.generate_presigned_url(
            "get_object",
            Params={
                "Bucket": self.bucket,
                "Key": key,
            },
            ExpiresIn=3600,
        )

        print(f"saved img url {saved_img_url}")
        expires_at = datetime.now(timezone.utc) + timedelta(minutes=60)
        with self.db.cursor() as cur:
            cur.execute(
                """
                INSERT INTO images (url, task_id, expires_at)
                VALUES (%s, %s, %s)
                """,
                (saved_img_url, task_id, expires_at)
            )

        self.db.commit()
        return {
            "status": "ready",
            "task_id": task_id
        }


# downloading using tls_client to prevent server blocks
def download(url: str) -> Image.Image:
    session = tls_client.Session(
        client_identifier="chrome_120",
        random_tls_extension_order=True
    )

    parsed_url = urlparse(url)
    headers = {
        "Host": parsed_url.netloc,
        "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:149.0) Gecko/20100101 Firefox/149.0",
        "Accept": "image/avif,image/webp,image/png,image/svg+xml,image/*;q=0.8,*/*;q=0.5",
        "Accept-Language": "en-US,en;q=0.9",
        "Accept-Encoding": "gzip, deflate, br, zstd",
        "Referer": f"{parsed_url.scheme}://{parsed_url.netloc}/",
    }

    tls_resp = session.get(url, headers=headers, timeout_seconds=10)
    cookies = {c.name: c.value for c in tls_resp.cookies}
    resp = requests.get(url, headers=headers, cookies=cookies, timeout=10)
    resp.raise_for_status()

    content_type = resp.headers.get("content-type", "")
    if not content_type.startswith("image/"):
        raise ValueError(f"Expected image, instead got {content_type}")

    return Image.open(BytesIO(resp.content))


# image processing function
def process(img: Image.Image, filter: str, **args: int) -> Image.Image:
    new_img = img.copy()
    match filter:
        case sharp if sharp.lower().startswith("sharp"):
            new_img = new_img.filter(ImageFilter.UnsharpMask(**args))
        case neg if neg.lower().startswith("neg"):
            new_img = ImageOps.invert(new_img.convert("RGB"))
        case blur if blur.lower() == "blur":
            new_img = new_img.filter(ImageFilter.GaussianBlur(**args))
        case proj if proj.lower() == "projection":
            new_img = new_img.transpose(Image.FLIP_LEFT_RIGHT)
        case gray if gray.lower() == "grayscale":
            new_img = new_img.convert("L")
        case _:
            pass

    img.show()
    return new_img

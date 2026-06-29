import pytest
import time
import requests


@pytest.fixture
def base_url():
    return "https://localhost:8080"


@pytest.fixture
def register_base_credentials():
    return { "username": "abc", "password": "abc" }


@pytest.fixture
def registered_user(base_url):
    creds = { "username": "def", "password": "def" }

    requests.post(f"{base_url}/register", json=creds, verify=False)
    return creds


class TestAuthentication:
    def test_register_user(self, base_url, register_base_credentials):
        resp = requests.post(f"{base_url}/register", json=register_base_credentials, verify=False)

        assert resp.status_code == 201
        body = resp.json()
        assert "token" in body


    def test_login_user(self, base_url, registered_user):
        resp = requests.post(f"{base_url}/login", json=registered_user, verify=False)

        assert resp.status_code == 200
        body = resp.json()
        assert "token" in body


    def test_repeated_register_user(self, base_url, registered_user):
        resp = requests.post(f"{base_url}/register", json=registered_user, verify=False)

        assert resp.status_code == 401
        body = resp.json()
        assert "error" in body


    def test_no_register_user(self, base_url):
        data = { "username": "no_register", "password": "no_register" }
        resp = requests.post(f"{base_url}/login", json=data, verify=False)

        assert resp.status_code == 401
        body = resp.json()
        assert "error" in body


@pytest.fixture
def auth_headers(base_url, register_base_credentials):
    resp = requests.post(f"{base_url}/register", json=register_base_credentials, verify=False)

    assert resp.status_code == 201
    body = resp.json()
    assert "token" in body

    resp = requests.post(f"{base_url}/login", json=register_base_credentials, verify=False)

    assert resp.status_code == 200
    body = resp.json()
    assert "token" in body

    headers = { "Authorization": f"Bearer {body["token"]}" }
    return headers


@pytest.fixture
def img_url():
    return "https://upload.wikimedia.org/wikipedia/commons/thumb/3/36/Grizzly_Denali_edit.jpg/960px-Grizzly_Denali_edit.jpg"


class TestImageProcess:
    def test_image_task_start(self, base_url, auth_headers, img_url):
        body = { 
            "image_url": img_url,
            "filter": "sharp",
        }

        resp = requests.post(f"task_url", headers=auth_headers, json=body, verify=False)

        task_body = resp.json()
        assert "task_id" in task_body


@pytest.fixture
def task_url(base_url):
    return f"{base_url}/task"


@pytest.fixture
def status_url(base_url):
    return f"{base_url}/status"


@pytest.fixture
def result_url(base_url):
    return f"{base_url}/result"


@pytest.fixture
def image_auth_headers(base_url, register_base_credentials):
    resp = requests.post(f"{base_url}/register", json=register_base_credentials, verify=False)

    if resp.status_code == 401:
        login_resp = requests.post(f"{base_url}/login", json=register_base_credentials, verify=False)
        assert login_resp.status_code == 200
        login_body = login_resp.json()
        assert "token" in login_body
        return { "Authorization": f"Bearer {login_body["token"]}" }


    assert resp.status_code == 201
    register_body = resp.json()
    assert "token" in register_body
    headers = { "Authorization": f"Bearer {register_body["token"]}" }
    return headers


class TestImageFilters:
    def _run_filter_test(self, task_url, status_url, result_url, image_auth_headers, img_url, filter_name):
        body = { 
            "image_url": img_url,
            "filter": filter_name,
        }

        resp = requests.post(f"{task_url}", headers=image_auth_headers, json=body, verify=False)

        task_body = resp.json()
        assert "task_id" in task_body

        task_id = task_body["task_id"]
        resp = requests.get(f"{status_url}/{task_id}", headers=image_auth_headers, verify=False)
        while resp.json()["status"] != "ready":
            time.sleep(1)
            resp = requests.get(f"{status_url}/{task_id}", headers=image_auth_headers, verify=False)

        assert resp.status_code == 200

        resp = requests.get(f"{result_url}/{task_id}", headers=image_auth_headers, verify=False)
        result_body = resp.json()
        assert "result" in result_body, result_body["result"]


    @pytest.mark.parametrize("filter_name", ["sharp", "neg", "blur", "projection", "grayscale"])
    def test_image_filter(self, task_url, status_url, result_url, image_auth_headers, img_url, filter_name):
        self._run_filter_test(task_url, status_url, result_url, image_auth_headers, img_url, filter_name)


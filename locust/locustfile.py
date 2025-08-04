from locust import HttpUser, task, between
import logging
import random
import string
import time

class GettingOwnershipRequests(HttpUser):
    # Add wait time between requests to simulate realistic user behavior
    # wait_time = between(1, 3)  # Wait 1-3 seconds between requests
    token = None

    def on_start(self):
        try:
            response = self.client.post("/api/login", json={
                "username": "aleksa",
                "password": "123456"
            }, timeout=60)

            if response.status_code == 200:
                self.token = response.json().get("token")
            else:
                self.token = None
        except Exception as e:
            self.token = None

    def get_auth_headers(self):
        if self.token:
            return {"Authorization": f"Bearer {self.token}"}
        return {}

    @task(2)
    def get_ownership_requests(self):
        if self.token:
            params = {
                "page": 1,
                "pageSize": 5,
                "sortBy": "created_at",
                "sortOrder": "",
                "search": '{"city":"Užice"}'
            }

            try:
                with self.client.get("/api/ownership/requests/8",
                                   params=params,
                                   headers=self.get_auth_headers(),
                                   timeout=60,
                                   catch_response=True) as response:
                    if response.status_code == 200:
                        response.success()
                    elif response.status_code == 401:
                        response.failure("Unauthorized - token expired")
                        self.token = None
                    else:
                        response.failure(f"Unexpected status code: {response.status_code}")
            except Exception as e:
                logging.error(f"Error in get_ownership_requests: {e}")

    @task(1)
    def search_ownership_requests(self):
        if self.token:
            params = {
                "page": 1,
                "pageSize": 5,
                "sortBy": "created_at",
                "sortOrder": "",
                "search": '{"city":"Užice","floor":2}'
            }

            try:
                with self.client.get("/api/ownership/requests/8",
                                   params=params,
                                   headers=self.get_auth_headers(),
                                   timeout=60,
                                   catch_response=True) as response:
                    if response.status_code == 200:
                        response.success()
                    elif response.status_code == 401:
                        response.failure("Unauthorized - token expired")
                        self.token = None
                    else:
                        response.failure(f"Unexpected status code: {response.status_code}")
            except Exception as e:
                logging.error(f"Error in search_ownership_requests: {e}")

class UserRegistration(HttpUser):
    wait_time = between(1, 2)

    def generate_unique_user_data(self):
        timestamp = str(int(time.time()))
        random_suffix = ''.join(random.choices(string.ascii_lowercase + string.digits, k=4))  # Reduced from 6 to 4
        unique_id = f"{timestamp[-6:]}_{random_suffix}"  # Use only last 6 digits of timestamp

        first_names = ["Ana", "Marko", "Jovana", "Stefan", "Milica", "Nikola", "Tijana", "Petar", "Jelena", "Milos"]
        last_names = ["Petrovic", "Jovanovic", "Nikolic", "Stojanovic", "Milic", "Djordjevic", "Stankovic", "Radovic", "Milosevic", "Savic"]

        first_name = random.choice(first_names)
        last_name = random.choice(last_names)

        return {
            "username": f"LT_{unique_id}",  # Shortened: LT_123456_abc1 (about 13-14 chars)
            "email": f"loadtest_{unique_id}@testmail.com",
            "password": "TestPass123!",
            "first_name": first_name,
            "last_name": last_name,
            "role": "Regular"
        }

    @task(1)
    def register_user(self):
        user_data = self.generate_unique_user_data()

        try:
            with self.client.post("/api/register",
                                json=user_data,
                                timeout=60,
                                catch_response=True) as response:
                if response.status_code == 200:
                    response.success()
                elif response.status_code == 400:
                    response.failure(f"Registration validation failed: {response.text}")
                elif response.status_code == 409:
                    response.failure(f"User already exists: {user_data['username']}")
                else:
                    response.failure(f"Unexpected status code: {response.status_code}")
        except Exception as e:
            logging.error(f"Error in register_user: {e}")

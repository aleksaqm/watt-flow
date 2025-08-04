from locust import HttpUser, task, between
import logging
import random
import string
import time
import json

class GettingOwnershipRequests(HttpUser):
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

    @task
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


class UserRegistration(HttpUser):
    wait_time = between(1, 2)
    
    def generate_unique_user_data(self):
        timestamp = str(int(time.time()))
        random_suffix = ''.join(random.choices(string.ascii_lowercase + string.digits, k=4))
        unique_id = f"{timestamp[-6:]}_{random_suffix}"
        
        first_names = ["Ana", "Marko", "Jovana", "Stefan", "Milica", "Nikola", "Tijana", "Petar", "Jelena", "Milos"]
        last_names = ["Petrovic", "Jovanovic", "Nikolic", "Stojanovic", "Milic", "Djordjevic", "Stankovic", "Radovic", "Milosevic", "Savic"]
        
        first_name = random.choice(first_names)
        last_name = random.choice(last_names)
        
        return {
            "username": f"LT_{unique_id}",
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


class RegisteringNewAdmin(HttpUser):
    token = None

    def on_start(self):
        try:
            response = self.client.post("/api/login", json={
                "username": "admin",
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

    @task
    def create_new_admin(self):
        if self.token:
            timestamp = str(int(time.time()))
            random_suffix = ''.join(random.choices(string.ascii_lowercase + string.digits, k=4))
            unique_id = f"{timestamp[-6:]}_{random_suffix}"
            data = {
                "username": f"AT_{unique_id}",
                "password": "123456",
                "role": "Admin",
                "email": f"admin_test_{unique_id}@e.com",
            }
            
            try:
                with self.client.post("/api/user/admin",
                                   json=data,
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

class CreateOwnershipRequest(HttpUser):
    wait_time = between(1, 3)
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
    
    @task
    def get_all_unowned_househodls(self):
        if self.token:
            search_data = {
                "city": "Aranđelovac",
                "withoutowner": True
            }
            params = {
                "page": 1,
                "pageSize": 10,
                "sortBy": "city",
                "sortOrder": ""
            }

            try:
                with self.client.post("/api/household/query",
                                    json=search_data,
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


    @task
    def create_ownership_request(self):
        if self.token:
            ownership_data = {
                "owner_id": 8,
                "household_id": 6,
                "images": ["iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="],
                "documents": ["JVBERi0xLjQKJdPr6eEKMSAwIG9iago8PAovVHlwZSAvQ2F0YWxvZwovUGFnZXMgMiAwIFIKPj4KZW5kb2JqCjIgMCBvYmoKPDwKL1R5cGUgL1BhZ2VzCi9LaWRzIFszIDAgUl0KL0NvdW50IDEKPJ4KZW5kb2JqCjMgMCBvYmoKPDwKL1R5cGUgL1BhZ2UKL1BhcmVudCAyIDAgUgovTWVkaWFCb3ggWzAgMCA2MTIgNzkyXQo+PgplbmRvYmoKeHJlZgowIDQKMDAwMDAwMDAwMCA2NTUzNSBmIAowMDAwMDAwMDA5IDAwMDAwIG4gCjAwMDAwMDAwNTggMDAwMDAgbiAKMDAwMDAwMDExNSAwMDAwMCBuIAp0cmFpbGVyCjw8Ci9TaXplIDQKL1Jvb3QgMSAwIFIKPj4Kc3RhcnR4cmVmCjE3NQolJUVPRg=="]
            }
            
            try:
                with self.client.post("/api/household/owner",
                                    json=ownership_data,
                                    headers=self.get_auth_headers(),
                                    timeout=60,
                                    catch_response=True) as response:
                    if response.status_code == 201:
                        response.success()
                    elif response.status_code == 400:
                        response.failure(f"Ownership request validation failed: {response.text}")
                    elif response.status_code == 401:
                        response.failure("Unauthorized - token expired")
                        self.token = None
                    else:
                        response.failure(f"Unexpected status code: {response.status_code}")
            except Exception as e:
                logging.error(f"Error in create_ownership_request: {e}")

class DeclineOwnershipRequest(HttpUser):
    wait_time = between(2, 5)
    token = None
    admin_token = None
    created_request_id = None

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

    def get_auth_headers(self, use_admin=False):
        if use_admin and self.admin_token:
            return {"Authorization": f"Bearer {self.admin_token}"}
        elif not use_admin and self.token:
            return {"Authorization": f"Bearer {self.token}"}
        return {}

    def login_as_admin(self):
        try:
            response = self.client.post("/api/login", json={
                "username": "admin",
                "password": "123456"
            }, timeout=60)

            if response.status_code == 200:
                self.admin_token = response.json().get("token")
                return True
            else:
                self.admin_token = None
                return False
        except Exception as e:
            self.admin_token = None
            return False

    @task
    def complete_ownership_request_flow(self):
        if self.token:
            ownership_data = {
                "owner_id": 8,
                "household_id": 16,
                "images": ["iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="],
                "documents": ["JVBERi0xLjQKJdPr6eEKMSAwIG9iago8PAovVHlwZSAvQ2F0YWxvZwovUGFnZXMgMiAwIFIKPj4KZW5kb2JqCjIgMCBvYmoKPDwKL1R5cGUgL1BhZ2VzCi9LaWRzIFszIDAgUl0KL0NvdW50IDEKPJ4KZW5kb2JqCjMgMCBvYmoKPDwKL1R5cGUgL1BhZ2UKL1BhcmVudCAyIDAgUgovTWVkaWFCb3ggWzAgMCA2MTIgNzkyXQo+PgplbmRvYmoKeHJlZgowIDQKMDAwMDAwMDAwMCA2NTUzNSBmIAowMDAwMDAwMDA5IDAwMDAwIG4gCjAwMDAwMDAwNTggMDAwMDAgbiAKMDAwMDAwMDExNSAwMDAwMCBuIAp0cmFpbGVyCjw8Ci9TaXplIDQKL1Jvb3QgMSAwIFIKPj4Kc3RhcnR4cmVmCjE3NQolJUVPRg=="]
            }
            
            try:
                with self.client.post("/api/household/owner",
                                    json=ownership_data,
                                    headers=self.get_auth_headers(use_admin=False),
                                    timeout=60,
                                    catch_response=True,
                                    name="Create Ownership Request") as response:
                    if response.status_code == 201:
                        response_data = response.json()
                        self.created_request_id = response_data.get("data").get("id")
                        response.success()
                    else:
                        response.failure(f"Failed to create ownership request: {response.status_code}")
                        return
            except Exception as e:
                logging.error(f"Error creating ownership request: {e}")
                return

            if not self.login_as_admin():
                logging.error("Failed to login as admin")
                return

            if self.created_request_id and self.admin_token:
                try:
                    reason = {"message": "Ne moze"}
                    with self.client.put(f"/api/ownership/decline/{self.created_request_id}",
                                         headers=self.get_auth_headers(use_admin=True),
                                         json=reason,
                                         timeout=60,
                                         catch_response=True,
                                         name="Decline Ownership Request") as response:
                        if response.status_code == 200:
                            response.success()
                        elif response.status_code == 400:
                            response.failure(f"Failed to approve request: {response.text}")
                        elif response.status_code == 401:
                            response.failure("Admin unauthorized")
                        else:
                            response.failure(f"Unexpected status code: {response.status_code}")
                except Exception as e:
                    logging.error(f"Error approving ownership request: {e}")
            else:
                logging.error("No request ID to approve or no admin token")
    
class GetAllOwneredHouseholds(HttpUser):
    wait_time = between(1, 3)
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

    @task
    def get_owned_households(self):
        if self.token:
            params = {
                "page": 1,
                "pageSize": 5,
                "sortBy": "city",
                "sortOrder": "",
            }
            search_data = {
                "ownerid": "8",
            }
            
            try:
                with self.client.post("/api/household/query",
                                   json=search_data,
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


class GetHouseholdDetails(HttpUser):
    wait_time = between(1, 3)
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

    @task(3)
    def get_household_basic_info(self):
        if self.token:
            household_id = 2
            
            try:
                with self.client.get(f"/api/household/my/{household_id}",
                                   headers=self.get_auth_headers(),
                                   timeout=60,
                                   catch_response=True,) as response:
                    if response.status_code == 200:
                        response.success()
                    elif response.status_code == 401:
                        response.failure("Unauthorized - token expired")
                        self.token = None
                    elif response.status_code == 404:
                        response.failure("Household not found")
                    else:
                        response.failure(f"Unexpected status code: {response.status_code}")
            except Exception as e:
                logging.error(f"Error in get_household_basic_info: {e}")

    @task(2)
    def get_household_consumption_12months(self):
        if self.token:
            household_id = 2
            current_year = 2025
            current_month = 8
            
            try:
                with self.client.get(f"/api/household/{household_id}/consumption/12months",
                                   params={
                                       "endYear": current_year,
                                       "endMonth": current_month
                                   },
                                   headers=self.get_auth_headers(),
                                   timeout=60,
                                   catch_response=True,) as response:
                    if response.status_code == 200:
                        response.success()
                    elif response.status_code == 401:
                        response.failure("Unauthorized - token expired")
                        self.token = None
                    elif response.status_code == 404:
                        response.failure("Household consumption not found")
                    else:
                        response.failure(f"Unexpected status code: {response.status_code}")
            except Exception as e:
                logging.error(f"Error in get_household_consumption_12months: {e}")

    


class GetDeviceConsumption(HttpUser):
    wait_time = between(1, 3)
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
    
    @task
    def get_device_consumption_query(self):
        if self.token:
            device_data = random.choice([{
                "DeviceId": "be7ffb42-c3b0-475b-bdc5-cb467d0f4111",
                "GroupPeriod": "1d",
                "Precision": "m",
                "Realtime": False,
                "TimePeriod": "90d"
            },
            {
                "DeviceId": "be7ffb42-c3b0-475b-bdc5-cb467d0f4111",
                "GroupPeriod": "7d",
                "Precision": "m",
                "Realtime": False,
                "TimePeriod": "365d"
            },
            {
                "DeviceId": "be7ffb42-c3b0-475b-bdc5-cb467d0f4111",
                "GroupPeriod": "1h",
                "Precision": "m",
                "Realtime": False,
                "TimePeriod": "12h"
            }])
            
            try:
                with self.client.post("/api/device-consumption/query-consumption",
                                    json=device_data,
                                    headers=self.get_auth_headers(),
                                    timeout=60,
                                    catch_response=True) as response:
                    if response.status_code == 200:
                        response.success()
                    elif response.status_code == 401:
                        response.failure("Unauthorized - token expired")
                        self.token = None
                    elif response.status_code == 400:
                        response.failure("Invalid query parameters")
                    else:
                        response.failure(f"Unexpected status code: {response.status_code}")
            except Exception as e:
                logging.error(f"Error in get_device_consumption_query: {e}")


class GetAllUsersMeetings(HttpUser):
    wait_time = between(1, 3)
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

    @task
    def get_users_meetings(self):
        if self.token:
            params = {
                "page": 1,
                "pageSize": 5,
                "sortBy": "start_time",
                "sortOrder": "",
                "search": '{}'
            }
            
            try:
                with self.client.get("/api/user/meetings/8",
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
    

class CreateMeeting(HttpUser):
    wait_time = between(1, 3)
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
    
    @task
    def get_available_timeslots(self):
        if self.token:
            future_days = random.randint(1, 14)
            future_date = time.strftime('%Y-%m-%d', time.localtime(time.time() + future_days * 24 * 3600))
            
            params = {
                "date": future_date,
                "clerk_id": 9,
            }
            
            try:
                with self.client.get("/api/timeslot",
                                   params=params,
                                   headers=self.get_auth_headers(),
                                   timeout=60,
                                   catch_response=True,
                                   name="/api/timeslot") as response:
                    if response.status_code == 200:
                        response.success()
                    elif response.status_code == 404:
                        if response.text == "timeslot not found":
                            response.success()
                    elif response.status_code == 401:
                        response.failure("Unauthorized - token expired")
                        self.token = None
                    else:
                        response.failure(f"Unexpected status code: {response.status_code}")
            except Exception as e:
                logging.error(f"Error in get_available_timeslots: {e}")

    @task
    def create_meeting(self):
        if self.token:
            future_days = random.randint(1, 14)
            meeting_date = time.strftime('%Y-%m-%d', time.localtime(time.time() + future_days * 24 * 3600))
            
            time_slots = [
                "09:00:00", "09:30:00", "10:00:00", "10:30:00", "11:00:00", "11:30:00",
                "14:00:00", "14:30:00", "15:00:00", "15:30:00"
            ]
            
            slot_number = random.randint(0, 14)
            occupied_ids = [slot_number]
            
            slot = {
                "Date": f"{meeting_date}T00:00:00Z",
                "ClerkId": 9,
                "Occupied": occupied_ids
            }
            
            start_time_js = f"{meeting_date}T{random.choice(time_slots)}.000Z"
            
            meeting = {
                "user_id": 8,
                "duration": 30,
                "clerk_id": 9,
                "start_time": start_time_js,
                "time_slot_id": slot_number
            }
            
            meeting_data = {
                "timeslot": slot,
                "meeting": meeting
            }
            
            try:
                with self.client.post("/api/meeting",
                                    json=meeting_data,
                                    headers=self.get_auth_headers(),
                                    timeout=60,
                                    catch_response=True,) as response:
                    if response.status_code == 201:
                        response.success()
                    elif response.status_code == 401:
                        response.failure("Unauthorized - token expired")
                        self.token = None
                    else:
                        response.failure(f"Unexpected status code: {response.status_code}")
            except Exception as e:
                logging.error(f"Error in create_meeting: {e}")

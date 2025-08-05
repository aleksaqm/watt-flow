import random
import logging
import string
from locust import HttpUser, task, between, constant
import uuid
import json
import time

logging.basicConfig(level=logging.INFO)

class PropertySearchUser(HttpUser):
    wait_time = between(1, 3)
    token = None
    user_id = None
    
    user_credentials = {"username": "vladimir", "password": "password"}

    def on_start(self):
        try:
            response = self.client.post("/api/login", json={
                "username": "vladimir",
                "password": "password"
            })
            response.raise_for_status()
            self.token = response.json().get("token")
            logging.info(f"User logged in successfully. Token acquired.")
        except Exception as e:
            self.token = None
            logging.error(f"Login failed: {e}")

    def get_auth_headers(self):
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}

    @task(10)
    def list_my_properties(self):
        if not self.token: return

        search_params = {
            "ownerId": self.user_id
        }
        
        params = {
            "page": 1,
            "pageSize": 5,
            "sortBy": "created_at",
            "sortOrder": "desc",
            "search": json.dumps(search_params)
        }
        
        self.client.get(
            "/api/property/query",
            params=params,
            headers=self.get_auth_headers(),
            name="/api/property/query"
        )

    @task(5)
    def filter_by_city(self):
        if not self.token: return
        
        search_params = {
            "ownerId": self.user_id,
            "city": random.choice(["Novi Sad", "Beograd", "Užice"])
        }
        
        params = {
            "page": 1,
            "pageSize": 5,
            "sortBy": "street",
            "sortOrder": "asc",
            "search": json.dumps(search_params)
        }
        
        self.client.get(
            "/api/property/query",
            params=params,
            headers=self.get_auth_headers(),
            name="/api/property/query"
        )
        
    @task(2)
    def complex_search(self):
        if not self.token: return

        search_params = {
            "ownerId": self.user_id,
            "city": "Novi Sad",
            "street": random.choice(["Bulevar", "Gogoljeva", "Test"]),
            "withoutOwner": False
        }
        
        params = {
            "page": 1,
            "pageSize": 10,
            "sortBy": "number",
            "sortOrder": "asc",
            "search": json.dumps(search_params)
        }
        
        self.client.get(
            "/api/property/query",
            params=params,
            headers=self.get_auth_headers(),
            name="/api/property/query"
        )

class PropertyRegistrationUser(HttpUser):
    wait_time = between(3, 5)
    token = None
    user_id = 14

    def on_start(self):
        try:
            response = self.client.post("/api/login", json={
                "username": "vladimir",
                "password": "password"
            })
            response.raise_for_status()
            self.token = response.json().get("token")
            logging.info(f"User logged in successfully. Token acquired.")
        except Exception as e:
            self.token = None
            logging.error(f"Login failed: {e}")

    def get_auth_headers(self):
        if not self.token:
            return {}
        return {"Authorization": f"Bearer {self.token}", "Content-Type": "application/json"}

    @task
    def create_property_request(self):
        if not self.token:
            logging.warning("No token available, skipping task.")
            return

        unique_suffix = str(uuid.uuid4())[:8]
        
        dummy_base64_image = "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"

        payload = {
            "floors": random.randint(1, 10),
            "status": 0,
            "owner_id": self.user_id,
            "images": [dummy_base64_image],
            "documents": [dummy_base64_image],
            
            "household": [
                {
                    "floor": 1,
                    "suite": f"1a-{unique_suffix}",
                    "status": 0,
                    "sq_footage": random.uniform(30.0, 90.0),
                    "cadastral_number": f"CN-{unique_suffix}-1",
                    "device_status": {
                        "device_id": str(uuid.uuid4()),
                        "is_active": False
                    }
                },
                {
                    "floor": 2,
                    "suite": f"2b-{unique_suffix}",
                    "status": 0,
                    "sq_footage": random.uniform(50.0, 120.0),
                    "cadastral_number": f"CN-{unique_suffix}-2",
                    "device_status": {
                        "device_id": str(uuid.uuid4()),
                        "is_active": False
                    }
                }
            ],
            
            "address": {
                "city": random.choice(["Novi Sad", "Beograd", "Niš"]),
                "street": f"Test Ulica {unique_suffix}",
                "number": str(random.randint(1, 200)),
                "latitude": 45.25 + random.uniform(-0.1, 0.1),
                "longitude": 19.83 + random.uniform(-0.1, 0.1)
            }
        }

        try:
            with self.client.post(
                "/api/property",
                json=payload,
                headers=self.get_auth_headers(),
                name="Create Property",
                catch_response=True
            ) as response:
                if response.status_code == 201:
                    response.success()
                else:
                    response.failure(
                        f"Failed to create property. Status: {response.status_code}, Response: {response.text}"
                    )
        except Exception as e:
            logging.error(f"Exception during property creation: {e}")
            if 'response' in locals():
                response.failure(str(e))


class AdminPropertyAcceptor(HttpUser):
    wait_time = between(5, 10)
    admin_token = None
    
    admin_credentials = {"username": "admin", "password": "password"}
    user_credentials = {"username": "vladimir", "password": "password"}
    user_id = 14

    def on_start(self):
        try:
            response = self.client.post("/api/login", json=self.admin_credentials)
            response.raise_for_status()
            self.admin_token = response.json().get("token")
            logging.info("Admin logged in successfully.")
        except Exception as e:
            self.admin_token = None
            logging.error(f"Admin login failed: {e}")

    def get_auth_headers(self, token):
        return {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}

    def _create_property_as_user(self):
        try:
            response = self.client.post("/api/login", json=self.user_credentials)
            response.raise_for_status()
            user_token = response.json().get("token")
            if not user_token:
                raise Exception("User token not found in login response.")
        except Exception as e:
            logging.error(f"Failed to log in as regular user for property creation: {e}")
            return None

        unique_suffix = str(uuid.uuid4())[:8]
        payload = {
            "floors": 1,
            "owner_id": self.user_id,
            "household": [{"cadastral_number": f"KN-E2E-{unique_suffix}-1", "device_status": {"device_id": str(uuid.uuid4())}}],
            "address": {"city": "Beograd", "street": f"E2E Test Ulica {unique_suffix}", "number": "1"},
            "images": ["R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"],
            "documents": ["R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"],
            "status": 0
        }

        try:
            with self.client.post("/api/property", json=payload, headers=self.get_auth_headers(user_token), name="Create Property as User", catch_response=True) as response:
                if response.status_code == 201:
                    created_property_id = response.json().get("data", {}).get("id")
                    if created_property_id:
                        logging.info(f"Property {created_property_id} created successfully by user.")
                        return created_property_id
                
                response.failure(f"Property creation failed with status {response.status_code}")
                return None
        except Exception as e:
            logging.error(f"Exception during property creation as user: {e}")
            return None

    @task
    def create_and_accept_property_flow(self):
        if not self.admin_token:
            return

        property_id = self._create_property_as_user()

        if not property_id:
            logging.error("Stopping flow because property creation failed.")
            return
        
        with self.client.put(
            f"/api/property/{property_id}/accept",
            headers=self.get_auth_headers(self.admin_token),
            name="Accept Property as Admin",
            catch_response=True
        ) as response:
            if response.status_code == 200:
                response.success()
                logging.info(f"Property {property_id} accepted successfully by admin.")
            else:
                response.failure(f"Failed to accept property {property_id}. Status: {response.status_code}")


class AdminPropertyDecliner(HttpUser):
    wait_time = between(5, 10)
    admin_token = None
    
    admin_credentials = {"username": "admin", "password": "password"}
    user_credentials = {"username": "vladimir", "password": "password"}
    user_id = 14

    def on_start(self):
        try:
            response = self.client.post("/api/login", json=self.admin_credentials)
            response.raise_for_status()
            self.admin_token = response.json().get("token")
            logging.info("Admin logged in successfully.")
        except Exception as e:
            self.admin_token = None
            logging.error(f"Admin login failed: {e}")

    def get_auth_headers(self, token):
        return {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}

    def _create_property_as_user(self):
        try:
            response = self.client.post("/api/login", json=self.user_credentials)
            response.raise_for_status()
            user_token = response.json().get("token")
            if not user_token:
                raise Exception("User token not found in login response.")
        except Exception as e:
            logging.error(f"Failed to log in as regular user for property creation: {e}")
            return None

        unique_suffix = str(uuid.uuid4())[:8]
        payload = {
            "floors": 1,
            "owner_id": self.user_id,
            "household": [{"cadastral_number": f"KN-E2E-{unique_suffix}-1", "device_status": {"device_id": str(uuid.uuid4())}}],
            "address": {"city": "Beograd", "street": f"E2E Test Ulica {unique_suffix}", "number": "1"},
            "images": ["R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"],
            "documents": ["R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"],
            "status": 0
        }

        try:
            with self.client.post("/api/property", json=payload, headers=self.get_auth_headers(user_token), name="Create Property as User", catch_response=True) as response:
                if response.status_code == 201:
                    created_property_id = response.json().get("data", {}).get("id")
                    if created_property_id:
                        logging.info(f"Property {created_property_id} created successfully by user.")
                        return created_property_id
                
                response.failure(f"Property creation failed with status {response.status_code}")
                return None
        except Exception as e:
            logging.error(f"Exception during property creation as user: {e}")
            return None

    @task
    def create_and_decline_property_flow(self):
        if not self.admin_token:
            return

        property_id = self._create_property_as_user()

        decline_payload = {
            "message": "Automated test decline."
        }

        if not property_id:
            logging.error("Stopping flow because property creation failed.")
            return
        
        with self.client.put(
            f"/api/property/{property_id}/decline",
            headers=self.get_auth_headers(self.admin_token),
            json=decline_payload,
            name="Decline Property as Admin",
            catch_response=True
        ) as response:
            if response.status_code == 200:
                response.success()
                logging.info(f"Property {property_id} declined successfully by admin.")
            else:
                response.failure(f"Failed to accept property {property_id}. Status: {response.status_code}")


class HouseholdAccessGranter(HttpUser):
    wait_time = between(1, 4)
    token = None
    
    owner_credentials = {"username": "vladimir", "password": "password"}
    owner_id = None
    owned_household_ids = []

    grantable_user_ids = [] 
    
    def on_start(self):
        try:
            with self.client.post(
                "/api/login", 
                json=self.owner_credentials, 
                name="/api/login",
                catch_response=True
            ) as response:
                response.raise_for_status()
                login_data = response.json()
                self.token = login_data.get("token")
                self.owner_id = '3'
                if not self.token or not self.owner_id:
                    raise Exception("Token or User ID not found in login response.")
            
            query_payload = {
                "ownerid": str(self.owner_id)
            }
            query_params= {
                "page": 1,
                "pageSize": 500
            }

            with self.client.post(
                "/api/household/query",
                json=query_payload,
                params=query_params,
                headers=self.get_auth_headers(),
                name="/api/household/query"
            ) as res:
                res.raise_for_status()
                response_data = res.json()
                households = response_data.get("households", [])
                self.owned_household_ids = [h['id'] for h in households]
                
                if not self.owned_household_ids:
                    logging.warning(f"Owner {self.owner_id} has no households.")

        except Exception as e:
            self.token = None
            logging.error(f"Setup failed for owner: {e}")

    def get_auth_headers(self):
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}

    def generate_random_string(length=3):
        characters = string.ascii_lowercase + string.digits
        return ''.join(random.choice(characters) for _ in range(length))

    @task(5)
    def search_for_users(self):
        if not self.token:
            return

        
        search_payload = {
            "Role": "Regular",
            "Username": "an"
        }
        
        params = {
            "sortBy": "username"
        }

        try:
            with self.client.post(
                "/api/user/query",
                json=search_payload,
                params=params,
                headers=self.get_auth_headers(),
                name="/api/user/query (Search and Prepare)",
                catch_response=True
            ) as response:
                response.raise_for_status()
                
                response_data = response.json()
                found_users = response_data.get("users", [])
                self.grantable_user_ids = [user['id'] for user in found_users if 'id' in user]
                if not self.grantable_user_ids:
                    response.success() 
                    return

        except Exception as e:
            logging.error(f"Failed during user search step: {e}")
            if 'response' in locals():
                response.failure(str(e))
            return

    @task(2)
    def grant_access_to_household(self):
        if not self.token or not self.owned_household_ids or not self.grantable_user_ids:
            return

        household_id = random.choice(self.owned_household_ids)
        user_id_to_grant = random.choice(self.grantable_user_ids)
        
        payload = {"userId": user_id_to_grant}

        with self.client.post(
            f"/api/household/{household_id}/access",
            json=payload,
            headers=self.get_auth_headers(),
            name="/api/households/[id]/access",
            catch_response=True
        ) as response:
            if response.status_code in [200, 201]:
                response.success()
                logging.info(
                    f"User {user_id_to_grant} granted access to household {household_id}."
                )
            elif response.status_code in [400, 404]:
                response.success()
                logging.info(
                    f"User {user_id_to_grant} likely already has access to household {household_id}."
                )
            else:
                response.failure(f"Failed to grant access. Status: {response.status_code}")



class HouseholdAccessRevoker(HttpUser):
    wait_time = between(2, 5)
    token = None
    
    owner_credentials = {"username": "vladimir", "password": "password"}
    owner_id = '3'
    owned_household_ids = []

    @task
    def grant_and_revoke_access(self):
        if not self.token or not self.owned_household_ids:
            return

        # 1. Search for grantable users
        search_payload = {
            "Role": "Regular",
            "Username": "an"
        }
        params = {"sortBy": "username"}

        try:
            with self.client.post(
                "/api/user/query",
                json=search_payload,
                params=params,
                headers=self.get_auth_headers(),
                name="/api/user/query (Search)",
                catch_response=True
            ) as response:
                response.raise_for_status()
                response_data = response.json()
                found_users = response_data.get("users", [])
                grantable_user_ids = [user['id'] for user in found_users if 'id' in user]
                
                if not grantable_user_ids:
                    logging.info("No users found to grant access.")
                    response.success()
                    return

        except Exception as e:
            logging.error(f"User search failed: {e}")
            if 'response' in locals():
                response.failure(str(e))
            return

        user_id = random.choice(grantable_user_ids)
        household_id = random.choice(self.owned_household_ids)

        grant_payload = {"userId": user_id}
        with self.client.post(
            f"/api/household/{household_id}/access",
            json=grant_payload,
            headers=self.get_auth_headers(),
            name="/api/household/[id]/access (Grant)",
            catch_response=True
        ) as grant_response:
            if grant_response.status_code in [200, 201, 400, 409]:
                grant_response.success()
                logging.info(f"Access granted to user {user_id} for household {household_id}")
            else:
                grant_response.failure(f"Grant failed. Status: {grant_response.status_code}")
                return

        time.sleep(2)

        with self.client.delete(
            f"/api/household/{household_id}/access/revoke/{user_id}",
            headers=self.get_auth_headers(),
            name="/api/household/[id]/access/revoke/[userId]",
            catch_response=True
        ) as revoke_response:
            if revoke_response.status_code in [200, 204]:
                revoke_response.success()
                logging.info(f"Access revoked from user {user_id} for household {household_id}")
            else:
                revoke_response.failure(f"Revoke failed. Status: {revoke_response.status_code}")

    def on_start(self):
        try:
            with self.client.post("/api/login", json=self.owner_credentials, name="/api/login") as response:
                response.raise_for_status()
                self.token = response.json().get("token")
                if not self.token:
                    raise Exception("Token not found")
            
            with self.client.post(
                "/api/household/query",
                json={"ownerid": self.owner_id},
                params={"pageSize": 100},
                headers=self.get_auth_headers(),
                name="/api/household/query"
            ) as res:
                res.raise_for_status()
                households = res.json().get("households", [])
                self.owned_household_ids = [h['id'] for h in households]
                if not self.owned_household_ids:
                    logging.warning(f"Owner {self.owner_id} has no households.")

        except Exception as e:
            logging.error(f"Login/setup failed: {e}")

    def get_auth_headers(self):
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}
    

class ConsumptionQueryUser(HttpUser):
    wait_time = between(0.1, 1.0)
    connection_timeout = 30.0 
    network_timeout = 30.0
    token = None
    
    user_credentials = {"username": "admin", "password": "password"}
    
    cities_to_test = ["Novi Sad", "Beograd", "Niš", "Subotica", "Kragujevac"]

    def on_start(self):
        self.client.verify = False
        try:
            with self.client.post("/api/login", json=self.user_credentials, name="/api/login") as response:
                response.raise_for_status()
                self.token = response.json().get("token")
                if not self.token: raise Exception("Token not found")
        except Exception as e:
            logging.error(f"Login failed: {e}")
            self.environment.runner.quit()

    def get_auth_headers(self):
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}

    def _send_consumption_query(self, time_period, city, name):
        group_period_map = {
            "1h": "5m",
            "12h": "1h",
            "1y": "1mo"
        }
        precision_map = {
            "1h": "m",
            "12h": "h",
            "1y": "d"
        }

        query_payload = {
            "TimePeriod": time_period,
            "GroupPeriod": group_period_map.get(time_period, "1h"),
            "Precision": precision_map.get(time_period, "m"),
            "City": city,
            "Realtime": False
        }

        with self.client.post(
            "/api/device-status/query-consumption",
            json=query_payload,
            headers=self.get_auth_headers(),
            name=name
        ) as response:
            if response.status_code != 200:
                logging.warning(f"Query '{name}' failed with status {response.status_code}")


    @task(10)
    def query_last_hour(self):
        if not self.token: return
        city = random.choice(self.cities_to_test)
        self._send_consumption_query("1h", city, "/api/consumption?period=1h")

    @task(5)
    def query_last_12_hours(self):
        if not self.token: return
        city = random.choice(self.cities_to_test)
        self._send_consumption_query("12h", city, "/api/consumption?period=12h")

    @task(1)
    def query_last_year(self):
        if not self.token: return
        city = random.choice(self.cities_to_test)
        self._send_consumption_query("1y", city, "/api/consumption?period=1y")



class BillDetailsViewer(HttpUser):
    wait_time = between(1, 2)
    token = None
    
    user_credentials = {"username": "vladimir", "password": "password"}
    my_bill_ids = []
    
    def on_start(self):
        try:
            with self.client.post("/api/login", json=self.user_credentials, name="/api/login") as response:
                response.raise_for_status()
                self.token = response.json().get("token")
                if not self.token: raise Exception("Token not found")
            
            search_params = {}
            params = {
                "page": 1,
                "pageSize": 500,
                "search": json.dumps(search_params)
            }
            with self.client.get(
                "/api/bills/search",
                params=params,
                headers=self.get_auth_headers(),
                name="/api/bills/search"
            ) as response:
                response.raise_for_status()
                response_data = response.json()
                
                bills = response_data.get("bills")
                if isinstance(bills, list):
                    self.my_bill_ids = [bill['payment_reference'] for bill in bills]
                
                if not self.my_bill_ids:
                    logging.warning("Logged-in user has no bills to view. Test may not run tasks.")
                else:
                    logging.info(f"Fetched {len(self.my_bill_ids)} bill IDs for user to test with.")

        except Exception as e:
            logging.error(f"Setup failed for BillDetailsViewer: {e}")

    def get_auth_headers(self):
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}

    @task
    def get_bill_details(self):
        if not self.token or not self.my_bill_ids:
            return

        bill_id = random.choice(self.my_bill_ids)
        
        params = {
            "id": bill_id
        }
        
        self.client.get(
            "/api/bills",
            params=params,
            headers=self.get_auth_headers(),
            name="/api/bills?id=[billId]"
        )

class BillSearchUser(HttpUser):
    wait_time = between(1, 3)
    token = None
    
    user_credentials = {"username": "vladimir", "password": "password"}

    def on_start(self):
        try:
            with self.client.post("/api/login", json=self.user_credentials, name="/api/login") as response:
                response.raise_for_status()
                self.token = response.json().get("token")
        except Exception as e:
            logging.error(f"Login failed: {e}")
            self.environment.runner.quit()

    def get_auth_headers(self):
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}

    @task(10)
    def list_first_page_of_bills(self):
        if not self.token: return

        search_params = {}
        
        params = {
            "page": 1,
            "pageSize": 10,
            "sortBy": "issue_date",
            "sortOrder": "desc",
            "search": json.dumps(search_params)
        }
        
        self.client.get(
            "/api/bills/search",
            params=params,
            headers=self.get_auth_headers(),
            name="/api/bills/search"
        )

    @task(8)
    def filter_unpaid_bills(self):
        if not self.token: return
        
        search_params = {
            "status": "Delivered"
        }
        
        params = {
            "page": random.randint(1, 3),
            "pageSize": 10,
            "sortBy": "issue_date",
            "sortOrder": "asc",
            "search": json.dumps(search_params)
        }
        
        self.client.get(
            "/api/bills/search",
            params=params,
            headers=self.get_auth_headers(),
            name="/api/bills/search"
        )

    @task(3)
    def complex_search_and_sort(self):
        if not self.token: return

        search_params = {
            "status": random.choice(["Paid", "Delivered"]),
            "minPrice": random.randint(1000, 5000),
            "maxPrice": random.randint(6000, 10000),
            "billingDate": f"2024-{random.randint(1, 12):02d}"
        }
        
        params = {
            "page": 1,
            "pageSize": 20,
            "sortBy": "price",
            "sortOrder": "desc",
            "search": json.dumps(search_params)
        }
        
        self.client.get(
            "/api/bills/search",
            params=params,
            headers=self.get_auth_headers(),
            name="/api/bills/search"
        )



class BillPayerUser(HttpUser):
    wait_time = between(2, 6)
    token = None
    
    user_credentials = {"username": "vladimir", "password": "password"}
    unpaid_bill_ids = []
    
    def on_start(self):
        try:
            with self.client.post("/api/login", json=self.user_credentials, name="/api/login") as response:
                response.raise_for_status()
                self.token = response.json().get("token")
                if not self.token: raise Exception("Token not found")
            
            self.fetch_unpaid_bills()

        except Exception as e:
            logging.error(f"Setup failed for BillPayerUser: {e}")

    def get_auth_headers(self):
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}

    def fetch_unpaid_bills(self):
        if not self.token: return

        search_params = {"status": "Delivered"}
        params = {
            "page": 1,
            "pageSize": 100,
            "search": json.dumps(search_params)
        }
        
        try:
            with self.client.get(
                "/api/bills/search",
                params=params,
                headers=self.get_auth_headers(),
                name="/api/bills/search"
            ) as response:
                response.raise_for_status()
                response_data = response.json()
                
                bills = response_data.get("bills")

                if isinstance(bills, list):
                    self.unpaid_bill_ids = [bill['payment_reference'] for bill in bills]
                else:
                    self.unpaid_bill_ids = []

                logging.info(f"Fetched {len(self.unpaid_bill_ids)} unpaid bills.")
        except Exception as e:
            logging.error(f"Failed to fetch unpaid bills: {e}")


    @task
    def pay_bill(self):
        if not self.token: return

        if not self.unpaid_bill_ids:
            self.fetch_unpaid_bills()
            if not self.unpaid_bill_ids:
                logging.info("No unpaid bills to pay. Skipping task.")
                return

        bill_id_to_pay = self.unpaid_bill_ids.pop(
            random.randrange(len(self.unpaid_bill_ids))
        )
        
        with self.client.put(
            f"/api/bills/pay/{bill_id_to_pay}",
            headers=self.get_auth_headers(),
            name="/api/bills/pay/[billId]",
            catch_response=True
        ) as response:
            if response.status_code == 200:
                response.success()
                logging.info(f"Successfully paid bill {bill_id_to_pay}")
            
            elif response.status_code in [400, 404, 409]:
                response.success()
                logging.warning(
                    f"Attempted to pay bill {bill_id_to_pay}, but it was likely "
                    f"already processed. Status: {response.status_code}"
                )
            
            else:
                response.failure(
                    f"Unexpected failure when paying bill {bill_id_to_pay}. "
                    f"Status: {response.status_code}"
                )
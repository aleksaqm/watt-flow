import json
import logging
import random
import string
import time
import uuid

import pika

from locust import HttpUser, User, between, task
from locust.env import Environment

logging.basicConfig(level=logging.INFO)

ADMIN_CREDENTIALS = {
    "username": "admin",
    "password": "admin123",
}

CLERK_CREDENTIALS = {
    "username": "test_clerk",
    "password": "2301002153154",
    "id": 999999999,
}

USER_CREDENTIALS = {
    "username": "test_user",
    "password": "2301002153154",
    "id": 999999998,
}


class AdminSimulatorAvailability(HttpUser):
    wait_time = between(1, 5)
    token = None
    admin_credentials = ADMIN_CREDENTIALS

    device_ids = [
        "be3ffb42-c3b0-475b-bdc5-cb467d0f4111",
        "be4ffb42-c3b0-475b-bdc5-cb467d0f4111",
        "be781b42-c3b0-475b-bdc5-cb467d0f4fa1",
        "be7ffb42-c3b0-475b-bdc5-cb467d0f4111",
        "be2ffb42-c3b0-475b-bdc5-cb467d0f4111",
        "ff781b42-c3b0-475b-bdc5-cb467d0f4222",
        "be1ffb42-c3b0-475b-bdc5-cb467d0f4111",
    ]

    def on_start(self):
        try:
            # 1. Login
            with self.client.post(
                "/api/login", json=self.admin_credentials, name="/api/login"
            ) as response:
                response.raise_for_status()
                self.token = response.json().get("token")
                if not self.token:
                    raise Exception("Admin token not found in login response.")
                logging.info("Admin logged in successfully for availability check.")

            # # 2. Dobavi listu uredjaja
            # with self.client.post(
            #     "/api/household/query",
            #     json={"Street": "Sim"},
            #     params={"page": 1, "pageSize": 50},
            #     headers=self.get_auth_headers(),
            #     name="/api/household/query (init)",
            # ) as res:
            #     res.raise_for_status()
            #     households = res.json().get("households", [])
            #     self.device_ids = [
            #         h["device_address"] for h in households if "device_address" in h
            #     ]
            #     if not self.device_ids:
            #         logging.warning("No device IDs found for availability testing.")
            #         self.environment.runner.quit()

        except Exception as e:
            logging.error(f"Setup failed for AdminSimulatorAvailability: {e}")
            self.environment.runner.quit()

    def get_auth_headers(self):
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}

    def _query_availability(self, period, group, precision, name):
        if not self.device_ids:
            return

        payload = {
            "TimePeriod": period,
            "GroupPeriod": group,
            "Precision": precision,
            "DeviceId": random.choice(self.device_ids),
            "Realtime": False,
        }

        with self.client.post(
            "/api/device-status/query-status",
            json=payload,
            headers=self.get_auth_headers(),
            name=name,
            catch_response=True,
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(
                    f"Query '{name}' failed with status {response.status_code}: {response.text}"
                )

    @task(10)
    def query_short_term_availability(self):
        """Testira upite za kraće periode (poslednja 3 sata)."""
        self._query_availability("3h", "1h", "m", "/api/device-status?period=short")

    @task(5)
    def query_medium_term_availability(self):
        """Testira upite za srednje periode (poslednja 24 sata)."""
        self._query_availability("24h", "3h", "h", "/api/device-status?period=medium")

    @task(1)
    def query_long_term_availability(self):
        """Testira upite za duge periode (poslednjih 30 dana) - ovo je stres test."""
        self._query_availability("30d", "1d", "d", "/api/device-status?period=long")


class AdminClerkManagement(HttpUser):
    wait_time = between(1, 5)
    token = None
    admin_credentials = {"username": "admin", "password": "admin123"}
    created_clerk_ids = []
    suspended_clerk_ids = []

    def on_start(self):
        try:
            with self.client.post(
                "/api/login", json=self.admin_credentials, name="/api/login"
            ) as response:
                response.raise_for_status()
                self.token = response.json().get("token")
                if not self.token:
                    raise Exception("Admin token not found")
                logging.info("Admin logged in for clerk management.")
        except Exception as e:
            logging.error(f"Login failed for AdminClerkManagement: {e}")

    def get_auth_headers(self):
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}

    def generate_unique_clerk_data(self):
        """Generiše jedinstvene podatke za novog službenika."""
        unique_id = str(uuid.uuid4())[:8]
        return {
            "username": f"clerk_lt_{unique_id}",
            "email": f"clerk_lt_{unique_id}@test.com",
            "first_name": "LoadTest",
            "last_name": "Clerk",
            "jmbg": "".join(
                random.choices(string.digits, k=13)
            ),  # JMBG je predefinisana lozinka
        }

    @task(5)
    def search_clerks_task(self):
        """Scenario: Pretraga postojećih službenika."""
        search_payload = {
            "Role": "Clerk",
            "Username": "clerk",
        }
        params = {"page": 1, "pageSize": 20, "sortBy": "username"}

        with self.client.post(
            "/api/user/query",
            json=search_payload,
            params=params,
            headers=self.get_auth_headers(),
            name="/api/user/query (Clerk)",
            catch_response=True,
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(
                    f"Clerk search failed with status {response.status_code}"
                )

    @task(2)
    def create_and_suspend_clerk_flow(self):
        """E2E Scenario: Kreiraj novog službenika, a zatim ga suspenduj."""
        # 1. Kreiraj službenika
        clerk_data = self.generate_unique_clerk_data()
        clerk_id = None

        with self.client.post(
            "/api/user/clerk/new",
            json=clerk_data,
            headers=self.get_auth_headers(),
            name="/api/user/clerk/new",
            catch_response=True,
        ) as response:
            if response.status_code in [200, 201]:
                response.success()
                try:
                    clerk_id = response.json().get("data", {}).get("id")
                    if clerk_id:
                        logging.info(f"Successfully created clerk with ID: {clerk_id}")
                        self.created_clerk_ids.append(clerk_id)
                except (json.JSONDecodeError, AttributeError):
                    logging.warning("Could not parse clerk ID from creation response.")
                    response.failure(
                        "Clerk creation response is not valid JSON or lacks ID"
                    )
                    return
            else:
                response.failure(
                    f"Clerk creation failed with status {response.status_code}"
                )
                return

        # 2. Suspenduj službenika (samo ako je uspešno kreiran)
        if clerk_id:
            with self.client.get(
                f"/api/user/suspend-clerk/{clerk_id}",
                headers=self.get_auth_headers(),
                name="/api/user/suspend-clerk/[id]",
                catch_response=True,
            ) as suspend_response:
                if suspend_response.status_code == 200:
                    suspend_response.success()
                    logging.info(f"Successfully suspended clerk {clerk_id}")
                else:
                    suspend_response.failure(
                        f"Failed to suspend clerk {clerk_id}, status: {suspend_response.status_code}"
                    )

        @task(1)
        def unsuspend_clerk_flow(self):
            """Scenario: Ponovo aktiviraj suspendovanog službenika."""
            if not self.suspended_clerk_ids:
                logging.warning("No suspended clerks to unsuspend.")
                return

            clerk_id = random.choice(self.suspended_clerk_ids)

            with self.client.get(
                f"/api/user/unsuspend/{clerk_id}",
                headers=self.get_auth_headers(),
                name="/api/user/unsuspend-clerk/[id]",
                catch_response=True,
            ) as unsuspend_response:
                if unsuspend_response.status_code == 200:
                    unsuspend_response.success()
                    logging.info(f"Successfully unsuspended clerk {clerk_id}")
                else:
                    unsuspend_response.failure(
                        f"Failed to unsuspend clerk {clerk_id}, status: {unsuspend_response.status_code}"
                    )


class SchedulingUser(HttpUser):
    wait_time = between(1, 5)

    clerk_token = None
    user_token = None

    clerk_id = 2
    available_slots = []

    def on_start(self):
        try:
            # Login kao službenik
            with self.client.post(
                "/api/login", json=CLERK_CREDENTIALS, name="/api/login (Clerk)"
            ) as response:
                response.raise_for_status()
                data = response.json()
                self.clerk_token = data.get("token")
                if not self.clerk_token or not self.clerk_id:
                    raise Exception("Clerk token/ID not found")

            logging.info("Clerk logged in for scheduling tests.")
        except Exception as e:
            logging.error(f"Login failed for SchedulingUser: {e}")

    def get_auth_headers(self, user_type="clerk"):
        token = self.clerk_token if user_type == "clerk" else self.user_token
        return {"Authorization": f"Bearer {token}"} if token else {}

    @task(3)
    def clerk_searches_user(self):
        """Scenario: Službenik pretrazuje kornisnike za sastanak."""
        if not self.clerk_token:
            return
        names = ["marko, aleksandar, dunja, marija"]

        payload = {
            "Role": "Regular",
            "Status": "Active",
            "Username": random.choice(names),
        }

        with self.client.post(
            "/api/user/query",
            json=payload,
            params={"page": 1, "pageSize": 20, "sortBy": "username"},
            headers=self.get_auth_headers("clerk"),
            name="/api/user/query (Clerk for user)",
            catch_response=True,
        ) as response:
            if response.status_code in [200]:
                response.success()
            else:
                response.failure(
                    f"Regular user search failed with status {response.status_code}"
                )

    @task(3)
    def clerk_schedules_private_meeting(self):
        """Scenario: Službenik zakazuje privatni sastanak."""
        if not self.clerk_token:
            return

        future_days = random.randint(1, 14)
        meeting_date = time.strftime(
            "%Y-%m-%d", time.localtime(time.time() + future_days * 24 * 3600)
        )

        time_slots = [
            "09:00:00",
            "09:30:00",
            "10:00:00",
            "10:30:00",
            "11:00:00",
            "11:30:00",
            "14:00:00",
            "14:30:00",
            "15:00:00",
            "15:30:00",
        ]

        slot_number = random.randint(0, 14)
        occupied_ids = [slot_number]

        slot = {
            "Date": f"{meeting_date}T00:00:00Z",
            "ClerkId": USER_CREDENTIALS["id"],
            "Occupied": occupied_ids,
        }

        start_time_js = f"{meeting_date}T{random.choice(time_slots)}.000Z"

        meeting = {
            "user_id": USER_CREDENTIALS["id"],
            "duration": 30,
            "clerk_id": CLERK_CREDENTIALS["id"],
            "start_time": start_time_js,
            "time_slot_id": slot_number,
        }

        payload = {"timeslot": slot, "meeting": meeting}

        with self.client.post(
            "/api/meeting/",
            json=payload,
            headers=self.get_auth_headers("clerk"),
            name="/api/schedule/ (clerk)",
            catch_response=True,
        ) as response:
            # Službenik može pokušati da zakaže već zauzet termin, što je OK
            if response.status_code in [200, 201, 409]:
                response.success()
            else:
                response.failure(
                    f"Private meeting scheduling failed with status {response.status_code}"
                )


class AdminPricelistManagement(HttpUser):
    wait_time = between(1, 5)
    token = None

    def on_start(self):
        try:
            with self.client.post(
                "/api/login", json=ADMIN_CREDENTIALS, name="/api/login"
            ) as response:
                response.raise_for_status()
                self.token = response.json().get("token")
                if not self.token:
                    raise Exception("Admin token not found")
                logging.info("Admin logged in for pricelist management.")
        except Exception as e:
            logging.error(f"Login failed for AdminPricelistManagement: {e}")

    def get_auth_headers(self):
        return {"Authorization": f"Bearer {self.token}"} if self.token else {}

    def generate_pricelist_data(self):
        """Generiše podatke za novi cjenovnik za budući mesec."""
        year = random.randint(2025, 2030)
        month = random.randint(1, 12)
        return {
            "month": month,
            "year": year,
            "red": round(random.uniform(15.0, 25.0), 2),
            "blue": round(random.uniform(8.0, 14.0), 2),
            "green": round(random.uniform(4.0, 7.0), 2),
            "tax": 20.0,
            "bill_power": round(
                random.uniform(50.0, 100.0), 2
            ),
        }

    @task
    def create_new_pricelist(self):
        """Scenario: Administrator kreira novi cenovnik."""
        payload = self.generate_pricelist_data()

        with self.client.post(
            "/api/pricelist",
            json=payload,
            headers=self.get_auth_headers(),
            name="/api/pricelist [POST]",
            catch_response=True,
        ) as response:
            # Očekujemo da će ponekad pokušati da kreira cenovnik za mesec koji već ima definisan cenovnik
            if response.status_code in [201, 409, 500]:
                response.success()
            else:
                response.failure(
                    f"Pricelist creation failed with status {response.status_code}"
                )


class AmqpClient:
    """
    Custom AMQP klijent za Locust koji se integriše sa RabbitMQ.
    Meri vreme slanja poruke i prijavljuje uspehe/neuspehe.
    """

    def __init__(self, environment: Environment, amqp_url: str, exchange_name: str):
        self.env = environment
        try:
            params = pika.URLParameters(amqp_url)
            self.connection = pika.BlockingConnection(params)
            self.channel = self.connection.channel()
            self.exchange_name = exchange_name
            self.channel.exchange_declare(
                exchange=self.exchange_name, exchange_type="topic", durable=True
            )
        except Exception as e:
            # Ako konekcija ne uspe, zaustavi ceo test
            self.env.runner.quit()
            raise e

    def publish(self, routing_key: str, body: dict, name: str):
        """Šalje poruku na exchange sa datim routing key-em."""
        start_time = time.time()
        try:
            # Konvertuj telo poruke u JSON string
            message_body = json.dumps(body)

            self.channel.basic_publish(
                exchange=self.exchange_name,
                routing_key=routing_key,
                body=message_body,
                properties=pika.BasicProperties(
                    delivery_mode=2,  # Poruka je perzistentna
                ),
            )
        except Exception as e:
            # Ako slanje ne uspe
            total_time = int((time.time() - start_time) * 1000)
            self.env.events.request.fire(
                request_type="AMQP",
                name=name,
                response_time=total_time,
                response_length=0,
                exception=e,
                context={},
            )
        else:
            # Ako je slanje uspešno
            total_time = int((time.time() - start_time) * 1000)
            self.env.events.request.fire(
                request_type="AMQP",
                name=name,
                response_time=total_time,
                response_length=len(message_body),
                exception=None,
                context={},
            )

    def close(self):
        """Zatvara konekciju."""
        if self.connection and self.connection.is_open:
            self.connection.close()


class AmqpUser(User):
    """
    Abstraktna klasa za korisnike koji komuniciraju preko AMQP.
    """

    abstract = True

    def __init__(self, environment: Environment):
        super().__init__(environment)
        amqp_url = "amqp://guest:guest@localhost:5672/"
        exchange_name = "watt-flow"
        self.client = AmqpClient(environment, amqp_url, exchange_name)

    def on_stop(self):
        """Poziva se kada se korisnik zaustavi."""
        self.client.close()


class SimulatorUser(AmqpUser):
    """
    Simulira jedan uređaj (pametno brojilo) koji šalje podatke na AMQP broker.
    """

    wait_time = between(1, 2)

    def on_start(self):
        """Inicijalizuje jedinstvene podatke za svaki simulator."""
        self.device_id = f"device_{uuid.uuid4()}"
        self.city = random.choice(["Beograd", "Novi Sad", "Niš", "Kragujevac"])
        self.current_time = time.time()
        logging.info(f"Simulator {self.device_id} for city {self.city} started.")

    @task(1)
    def send_measurement(self):
        """
        Simulira slanje merenja potrošnje.
        Ovo je glavni test opterećenja za upis u InfluxDB.
        """
        # Generiši podatke za merenje
        measurement_data = {
            "DeviceId": self.device_id,
            "Value": round(
                random.uniform(0.1, 5.0), 4
            ),
            "Timestamp": time.strftime(
                "%Y-%m-%dT%H:%M:%SZ", time.gmtime(self.current_time)
            ),
            "Address": {
                "city": self.city,
                "street": "test",
                "number": str(random.randint(1, 100)),
            },
        }

        routing_key = f"measurement.{self.city}"

        self.client.publish(
            routing_key=routing_key, body=measurement_data, name="send_measurement"
        )

        self.current_time += 3600

    @task(10)
    def send_heartbeat(self):
        """
        Simulira slanje heartbeat poruke.
        """
        heartbeat_data = {
            "DeviceId": self.device_id,
            "Timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        }

        routing_key = f"heartbeat.{self.city}"

        self.client.publish(
            routing_key=routing_key, body=heartbeat_data, name="send_heartbeat"
        )

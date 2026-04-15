# Practice04: 
## Deployment (Docker & Docker Compose)

The project is fully containerized and ready to be launched locally with a single command.

### Basic Docker Principles in this Project
1. **Multi-stage Dockerfile:**
   * **Stage 1 (Builder):** Uses a heavy `golang:1.25-alpine` image to download dependencies (`go mod download`) and compile the binary file.
   * **Stage 2 (Release):** Uses a completely clean and minimalist `alpine:latest` image (~5 MB). Only the compiled binary is copied from the first stage. This keeps the final image highly secure and lightweight.
2. **.dockerignore:**
   Excludes local artifacts (like `.git` folders, tests, and local `.env` files) from the build context. This protects secrets and speeds up the build process.
3. **Docker Compose & Orchestration:**
   Instead of packing the database and the application into a single container (which is a Docker anti-pattern), the project uses `docker-compose.yml` to manage three independent services:
   * **db**: PostgreSQL database.
   * **migrator**: A temporary container that applies SQL migrations to the database and then gracefully exits (Exit 0).
   * **api**: The Go backend server.
   
   The `depends_on` mechanism with `condition: service_completed_successfully` ensures strict execution order: the API will never start until the migrations are successfully applied.

---

### How to Run the Project

1. Ensure Docker Desktop is running on your machine.
2. Run the following command in the root directory:
   ```bash
   docker compose up -d
   ```
3. Shut down :
   ```bash
   docker compose down
   ```

### API:

#### 1. Health Check


```bash
curl -X GET http://localhost:8080/health
```


#### 2. Create a Meetup
Creates a new meetup in the **Draft** status. 

```bash
curl -X POST http://localhost:8080/meetups \
     -H "Content-Type: application/json" \
     -d '{"title": "Golang Architecture Workshop", "capacity": 100}'
```

#### 3. Change Meetup Status

```bash
curl -X PATCH http://localhost:8080/meetups/<MEETUP_ID>/status \
     -H "Content-Type: application/json" \
     -d '{"status": "Published"}'
```

#### 4. Get Meetup by ID

```bash
curl -X GET http://localhost:8080/meetups/<MEETUP_ID>
```



--- 

## Pracitce 5

```bash
docker-compose down -v --remove-orphans
```

```bash
COMMIT_HASH=$(git rev-parse --short HEAD) docker-compose build --no-cache
```

```bash
docker-compose up
```

```bash
curl -i -X GET http://localhost:8080/health
```


```bash
curl -i -X GET http://localhost:8080/users/health/live
```
```bash
curl -i -X GET http://localhost:8080/users/health/ready
```


```bash
curl -i -X GET http://localhost:8080/meetups/health/live
```
```bash
curl -i -X GET http://localhost:8080/meetups/health/ready
```


```bash
curl -i -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{
           "first_name": "Anna",
           "last_name": "Nechytailenko",
           "email": "anna.test@example.com"
         }'
```

```bash
curl -i -X GET http://localhost:8080/users/<USER_ID>
```

```bash
curl -i -X POST http://localhost:8080/meetups \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Golang Microservices Workshop",
           "capacity": 50,
           "owner_user_id": "<USER_ID>"
         }'
```



```bash
curl -i -X GET http://localhost:8080/meetups/<MEETUP_ID>
```

```bash
curl -i -X PATCH http://localhost:8080/meetups/<MEETUP_ID>/status \
     -H "Content-Type: application/json" \
     -d '{
           "status": "Published"
         }'
```


--- 

## Pracitce 6

```bash
docker-compose down -v --remove-orphans
```

```bash
COMMIT_HASH=$(git rev-parse --short HEAD) docker-compose build --no-cache
```

```bash
docker-compose up
```

```bash
curl -i -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{
           "first_name": "Anna",
           "last_name": "Nechytailenko",
           "email": "anna.test@example.com"
         }'
```

```bash
curl -i -X POST http://localhost:8080/meetups \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Golang Microservices Workshop",
           "capacity": 50,
           "owner_user_id": "<USER_ID>"
         }'
```

After meetup created user will have notificatoin that was delivered through http client: 

```bash
curl -i -X GET http://localhost:8080/notifications/<USER_ID>
```



## Practice 7

Linux:
```bash
chmod +x installation/setup.sh
./installation/setup.sh
```

Windows:
```bash
.\installation\setup.ps1
```



```bash
curl -i -X GET http://localhost:8080/health/live
```
```bash
curl -i -X GET http://localhost:8080/health/ready
```

```bash
curl -i -X GET http://localhost:8080/users/health/live
```
```bash
curl -i -X GET http://localhost:8080/users/health/ready
```

```bash
curl -i -X GET http://localhost:8080/meetups/health/live
```
```bash
curl -i -X GET http://localhost:8080/meetups/health/ready
```

```bash
curl -i -X GET http://localhost:8080/notifications/health/live
```
```bash
curl -i -X GET http://localhost:8080/notifications/health/ready
``

```bash
curl -i -X GET http://localhost:8080/workflows/health/live
```

```bash
curl -i -X GET http://localhost:8080/workflows/health/ready
``



```bash
curl -i -X POST http://localhost:8080/meetups \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Go Kubernetes Workshop",
           "capacity": 100,
           "owner_user_id": "<USER_ID>"
         }'
```

```bash
curl -X GET http://localhost:8080/meetups/<MEETUP_ID>
```

```bash
curl -X PATCH http://localhost:8080/meetups/<MEETUP_ID>/status \
     -H "Content-Type: application/json" \
     -d '{"status": "<STATUS>"}'
```

```bash
curl -i -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{
           "first_name": "John",
           "last_name": "Doe",
           "email": "john.doe@example.com"
         }'
```

```bash
curl -X GET http://localhost:8080/users/<USER_ID>
```


```bash
curl -X GET http://localhost:8080/notifications/<USER_ID>
```

```bash
curl -X POST http://localhost:8080/workflows/join-meetup \
     -H "Content-Type: application/json" \
     -d '{"userId": "<USER_ID>", "meetupId": "<MEETUP_ID"}'
```


```bash
curl -X GET http://localhost:8080/workflows/<USER_ID>
```

**To see compensation:**
- 1 stage compensation : not existed meetupId
- 2 stage compensation : not existed userId

### Saga Orchestration Service Demonstration Report

This report documents the successful demonstration of the **Saga Orchestration Service** for the `join-meetup` workflow. The demo validates the system's ability to handle distributed transactions, trigger compensations when a stage fails, and complete the full flow under correct conditions.

---

#### 1. System Readiness & Health Check
Before initiating the orchestration tests, a full cluster health check was performed via the Gateway. All microservices and their respective databases were confirmed as **UP**.

* **Gateway:** `HTTP 200` (Database skipped)
* **Users Service:** `HTTP 200` (Database OK)
* **Meetups Service:** `HTTP 200` (Database OK)
* **Notifications Service:** `HTTP 200` (Database OK)
* **Workflow Service:** `HTTP 200` (Database OK)

---

### 2. Orchestration Test Scenarios

#### **Scenario A: Stage 1 Compensation (Meetup Validation Failure)**
In this stage, the orchestrator attempts to reserve a seat, but the Meetup service cannot find the requested ID.
* **Action:** POST to `/workflows/join-meetup` with a non-existent `meetupId`.
* **Result:** The Saga state transitioned to `Failed`.
* **Compensation Logic:** The system successfully captured the error: `"meetup with id '01688...65' not found"`. Since the first step failed, the transaction was safely aborted before affecting other services.
* **Correlation ID:** `039b51dc-e51f-43e3-9593-92a21df9fca5`

#### **Scenario B: Stage 2 Compensation (User Validation Failure)**
In this stage, the Meetup is valid, but the User service fails to validate the participant.
* **Action:** POST to `/workflows/join-meetup` with an invalid `userId`.
* **Result:** The Saga state transitioned to `Failed`.
* **Compensation Logic:** The orchestrator received a failure from the User stage: `"User 3bfa4...28 does not exist"`.
* **Automatic Rollback:** The system successfully triggered compensation to ensure no "ghost" reservations remained in the Meetup service despite the user error.

#### **Scenario C: Successful Orchestration (Happy Path)**
The final test confirms that the orchestrator can coordinate all services when data is valid.
* **Action:** POST to `/workflows/join-meetup` with valid `userId` and `meetupId`.
* **Result:** The Saga state successfully transitioned from `Initializing` to **`Completed`**.
* **Timeline:** The entire distributed transaction (Validation -> Reservation -> Notification) was finalized within ~3 seconds.
* **Verification:** Final GET request confirmed the state: `{"state":"Completed"}`.

---

#### 3. Summary of Results

| Scenario | Target State | Actual State | Error Captured |
| :--- | :--- | :--- | :--- |
| **Invalid Meetup** | Failed | **Failed** | Meetup not found |
| **Invalid User** | Failed | **Failed** | User does not exist |
| **Valid Data** | Completed | **Completed** | None |


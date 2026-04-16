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
curl -i -X GET http://localhost:8080/notifications/<USER_ID>
```



## Practice 7


**Linux:**
```bash
chmod +x installation/setup.sh
./installation/setup.sh
```


**Windows:**
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
           "owner_user_id": "b34afd33-84d0-4ed1-9ed5-cb00b6287619"
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
     -d '{"userId": "<USER_ID>", "meetupId": "<MEETUP_ID>"}'
```


```bash
curl -X GET http://localhost:8080/workflows/<WORKFLOW_ID>
```

**To see compensation:**
- 1 stage compensation : not existed meetupId
- 2 stage compensation : not existed userId


## Practice 8

To run minikube : 

**Linux:**
```bash
chmod +x installation/setup.sh
./installation/setup.sh
```
**Windows:**
```bash
.\installation\setup.ps1
```

The Version 2.0 update includes the following images, which have been hardened with correlation propagation and centralized logging handling:

* **Gateway Service:** `anna13nechytailenko/saga-gateway:v2.0`
* **Users Service:** `anna13nechytailenko/saga-users:v2.0`
* **Notifications Service:** `anna13nechytailenko/saga-notifications:v2.0`
* **Workflow Service:** `anna13nechytailenko/saga-workflow:v2.0`

### Kubernetes Deployment & Resilience Demonstration

> **Note:** The complete, raw terminal output for this demonstration can be found in the folder [demo/practice_8](./demo/practice_8/terminal_demo.txt).

### Objective
To demonstrate the production-grade capabilities of Kubernetes for the Meetups microservice, specifically focusing on **elastic scaling**, **zero-downtime deployments (Rolling Updates)**, **system observability (Centralized Logging with Correlation IDs)**, and **instant rollbacks**.

---

#### Step 1: Demonstrating Elasticity (Scaling Up)
**What we did:** We checked the initial state of the `meetups-service` (which had 3 running pods) and dynamically scaled it up to 5 replicas using the imperative `kubectl scale` command.
* **Command:** `kubectl scale deploy meetups-service --replicas=5 -n saga-system`
* **Purpose:** To demonstrate how quickly Kubernetes can provision new resources and adapt to sudden spikes in network traffic without requiring code changes or downtime.

#### Step 2: Establishing the Baseline (Version 1.1 - The Problem)
**What we did:** We sent successful HTTP POST requests to create a User and a Meetup, generating traffic. We then inspected the logs of the `meetups-service`.
* **Command:** `kubectl logs -l app=meetups-service -n saga-system --tail=20`
* **Purpose:** To highlight the observability flaw in Version 1.1. The logs showed application events, but lacked any tracking context. In a microservice architecture, without a Correlation ID, it is nearly impossible to trace a single user's request across multiple services.

#### Step 3: The Zero-Downtime Rolling Update (Deploying Version 2.0)
**What we did:** We triggered a live update of the Meetups service, pointing the deployment to the newly built `v2.0` Docker image. We then monitored the rollout status.
* **Commands:** * `kubectl set image deployment/meetups-service app=anna13nechytailenko/saga-meetups:v2.0 -n saga-system`
  * `kubectl rollout status deploy/meetups-service -n saga-system`
* **Purpose:** To demonstrate the `RollingUpdate` strategy. Kubernetes gracefully replaced the old `v1.1` pods with the new `v2.0` pods one by one, ensuring that at least some instances were always online to serve user traffic. There was zero downtime during the upgrade.

#### Step 4: Verifying the Fix (Centralized Observability)
**What we did:** After the rollout finished, we sent new HTTP requests to the system (including an intentional 404 error) and checked the logs again.
* **Command:** `kubectl logs -l app=meetups-service -n saga-system --tail=20`
* **Purpose:** To prove the architectural improvement of Version 2.0. The logs now prominently displayed `[CorrID: <uuid>]` at the start of every line. We successfully achieved distributed tracing, allowing us to connect Gateway requests, Meetups logic, and RabbitMQ events using a single ID.

#### Step 5: The Instant Rollback (Disaster Recovery)
**What we did:** We simulated a scenario where the new version contained a critical bug and needed to be immediately reverted. We executed an undo command and verified the image version.
* **Commands:**
  * `kubectl rollout undo deployment/meetups-service -n saga-system`
  * `kubectl describe deploy meetups-service -n saga-system | grep Image:`
* **Purpose:** To showcase Kubernetes' built-in safety net. Instead of manually rebuilding and redeploying the old code, Kubernetes kept a history of the deployment and instantly rolled the infrastructure back to the stable `v1.1` state safely and automatically.


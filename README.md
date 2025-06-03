# 🚀 ESP32 Backend API

A secure and extensible backend service built in Go, designed to support a distributed sensor network based on ESP32 devices. It handles user authentication, sensor/KPI data storage, decryption, and serves as the bridge between mobile/web apps and your IoT infrastructure.

---

## ✨ Features

✅ JWT-based login/logout  
🔐 Role-based access control (`admin`, `user`)  
📦 Sensor & KPI data retrieval from MongoDB  
🛡️ Real-time and historical metrics (decrypted)  
🔁 Token TTL for offline mobile interactions  
📄 OpenAPI (Swagger) documentation  
⚙️ Auto-admin seeding on startup  
🧪 Unit-tested endpoints  
🧵 WebSocket live stream via MQTT subscription  

---

## 🧰 Tech Stack

- **Go 1.20+**
- **MongoDB**
- **Gin Web Framework**
- **JWT (Authentication)**
- **Swagger (API Docs)**
- **Docker-compatible**

---

## ⚙️ Environment Setup

### 1. Clone Project

```bash
git clone https://github.com/rednexx46/esp32-backend-api.git
cd esp32-backend-api
go mod tidy
````

### 2. Configure `.env`

```dotenv
# MongoDB Configuration
MONGO_HOST=mongodb
MONGO_PORT=27017
MONGO_USER=mongo_admin
MONGO_PASS=mongo_password
MONGO_DATABASE=iot_mesh
MONGO_SENSORS_COLLECTION=sensor_data
MONGO_USERS_COLLECTION=users
MONGO_DEVICES_COLLECTION=devices
MONGO_TOKENS_COLLECTION=tokens
MONGO_KPIS_COLLECTION=kpis

# MQTT Broker Configuration
MQTT_BROKER=mosquitto
MQTT_PORT=1883
MQTT_TOPIC=mesh/data/
MQTT_USERNAME=backend_api
MQTT_PASSWORD=backend_password
MQTT_TOPIC_SENSORS_DATA=mesh/data/

# Encryption Service
ENCRYPTION=true
ENCRYPT_API_URL=http://cipher-api:8080/

# Authentication & Security
JWT_SECRET=supersecretkey
ADMIN_USERNAME=portal-admin
ADMIN_PASSWORD=changeme123
TOKEN_TTL_MINUTES=30

# Server
PORT=8080
```

> ✅ On first run, the backend **automatically seeds** the database with the admin user from `.env`.

---

## 🔑 Authentication

### POST `/api/login`

**Body:**

```json
{
  "username": "portal-admin",
  "password": "changeme123"
}
```

**Returns:**

```json
{
  "token": "<JWT_TOKEN>"
}
```

Use it in `Authorization: Bearer <token>` for protected routes.

---

## 🔐 Protected Endpoints (Admin Only)

| Endpoint                          | Description                   |
| --------------------------------- | ----------------------------- |
| `GET /api/data/:device_id`        | All sensor data from a device |
| `GET /api/data`                   | All sensor data               |
| `GET /api/devices`                | List of active device IDs     |
| `GET /api/kpis`                   | All KPI entries               |
| `GET /api/kpis/device/:device_id` | KPI data per device           |
| `GET /api/profile`                | Authenticated user info       |

All responses are decrypted if `ENCRYPTION=true`.

---

## 📡 Real-Time via WebSocket

### WebSocket Endpoint

```http
GET /ws/live-data
```

**Header:**

```
Authorization: Bearer <JWT_TOKEN>
```

Streams all MQTT-sourced sensor data live to connected clients.

---

## 📄 Swagger Documentation

### Generate Docs

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

Then access:

> 🌐 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 🧪 Running Tests

```bash
go test ./tests/...
```

Tested via `testify` and seeded user credentials.

---

## 🧩 Project Structure

```
esp32-backend-api/
├── internal/
│   ├── db/          # Mongo connection, seed logic
│   ├── handlers/    # HTTP handlers
│   ├── middleware/  # JWT / Role guards
│   ├── models/      # Structs (User, Requests, Claims)
│   ├── mqtt/        # MQTT listener for live data
│   ├── utils/       # Hashing, TTL helpers
│   └── ws/          # WebSocket manager
├── tests/           # Unit tests
├── main.go          # Entry point
└── go.mod
```

---

## 🔒 Security Practices

* Hashed passwords using `bcrypt`
* JWT with `exp` and server-side validation
* Never exposes encrypted payloads to client
* `.env` secrets (not committed to repo)

---

## 🧠 Future Improvements

* ✅ Token refresh endpoint
* 🔄 BLE handshake for config upload
* 📊 InfluxDB integration for metrics
* 🧠 Role-based dashboards

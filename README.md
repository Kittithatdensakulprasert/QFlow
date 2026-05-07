# 🚀 QFlow — Queue Management API

ระบบจัดการคิวออนไลน์ (Queue Management System) พัฒนาด้วย Go + Gin Framework รองรับการจองคิว, จัดการคิว, และระบบแจ้งเตือนแบบครบวงจร พร้อมโครงสร้างแบบ Clean Architecture และรองรับ CI/CD

---

## 📌 Overview

QFlow เป็น RESTful API สำหรับระบบจัดการคิว ที่รองรับหลาย role:

* Guest
* User
* Provider
* Admin

รองรับ workflow:

* สมัครสมาชิก + OTP Authentication
* จองคิว (Queue Booking)
* จัดการคิว (Queue Management)
* ระบบแจ้งเตือน (Notification)

---

## 🧱 Architecture

โครงสร้างถูกปรับเป็น **Clean Architecture**

```
internal/
├── domain/        # Entity + Interface
├── service/       # Business Logic
├── repository/    # Database Layer (GORM)
├── handler/       # HTTP Layer (Gin)
├── middleware/    # JWT + Rate Limit Middleware
└── router/        # Route Registration
```

### ✅ Benefits

* แยกความรับผิดชอบชัดเจน
* Test ได้ง่าย (mock dependency)
* Maintain ง่าย
* Scale ได้ในอนาคต

---

## ⚙️ Tech Stack

* Language: Go 1.25
* Framework: Gin
* Database: PostgreSQL + GORM
* Authentication: JWT + OTP
* Container: Docker + Docker Compose
* API Docs: Swagger
* API Testing: Postman + Newman
* CI/CD: GitHub Actions
* Registry: DockerHub
* Deployment: Render

---

## 🔗 System Flow (CI/CD + Deployment)

```
Developer → Push Code → GitHub
                ↓
        GitHub Actions (CI)
        ├── Unit Tests + Coverage (≥80%)
        ├── Security Scan (gosec + Trivy)
        └── Integration Tests (Postman/Newman)
                ↓
        Build & Push Image → DockerHub
                ↓
        Deploy → Render
                ↓
        Test via Swagger / Postman
```

---

## 🗄️ Database Indexes

| Table         | Column(s)                    | Type            | เหตุผล                          |
| ------------- | ---------------------------- | --------------- | -------------------------------- |
| users         | phone                        | UNIQUE          | Login lookup                     |
| categories    | name                         | UNIQUE          | ป้องกันชื่อซ้ำ                  |
| otps          | phone                        | INDEX           | OTP lookup ทุก request           |
| otps          | expires_at                   | INDEX           | Cleanup expired OTPs             |
| providers     | category_id                  | INDEX           | FK filter                        |
| zones         | provider_id                  | INDEX           | FK, query บ่อยมาก               |
| queues        | (zone_id, queue_number)      | UNIQUE COMPOSITE | Queue number unique ต่อ zone    |
| queues        | zone_id                      | INDEX           | FK, query บ่อยมาก               |
| queues        | user_id                      | INDEX           | FK, query บ่อยมาก               |
| queues        | status                       | INDEX           | Filter by status ทุก query       |
| notifications | user_id                      | INDEX           | FK, query บ่อยมาก               |

---

## 📦 API Documentation

### 🔹 Swagger

ใช้ดู API และทดสอบแบบ interactive

```
http://localhost:3000/swagger/index.html
```

สามารถ:

* ดู endpoints ทั้งหมด
* กดทดลองยิง API
* ใส่ JWT Token

---

### 🔹 Postman

ไฟล์ collection:

```
postman/qflow_api.json
```

วิธีใช้:

1. เปิด Postman
2. Import file
3. ตั้ง Base URL:

```
http://localhost:3000
```

---

## 🐳 Docker & DockerHub

### 🔹 Run ด้วย Docker

```
docker-compose up --build
```

จะรัน:

* Backend (Go API)
* PostgreSQL Database

---

### 🔹 DockerHub

Docker image จะถูก build และ push อัตโนมัติจาก CI/CD เมื่อ merge เข้า `main`

```
docker pull kittithat/qflow:latest
```

---

## 🔄 CI/CD (GitHub Actions)

Pipeline ทำงานอัตโนมัติเมื่อ push code ไปที่ `main`, `develop`, หรือ `feature/*`

### Jobs:

1. **Test and Quality** — Unit tests (`./internal/service/...`, `./internal/swagger/...`), coverage ≥ 80%, formatting, go vet
2. **Security Scan** — gosec (static analysis) + Trivy (dependency vulnerability scan)
3. **Integration Test** — รัน Postman collection ผ่าน Newman กับ PostgreSQL จริง
4. **Docker Build + Push** — Build multi-arch image (amd64/arm64) แล้ว push DockerHub *(เฉพาะ push ไป `main`)*
5. **Deploy to Render** — Trigger deploy hook อัตโนมัติ *(เฉพาะ push ไป `main`)*

---

## 🛠 Installation

### 1. Clone Project

```
git clone <repo-url>
cd QFlow
```

---

### 2. Setup Environment

สร้างไฟล์ `.env` (ดูตัวอย่างได้จาก `.env.example`)

```
PORT=3000
DATABASE_URL=postgres://user:password@localhost:5432/qflow
JWT_SECRET=your-super-secret-key

# Bootstrap users — สร้าง Admin และ Provider ตอน startup ครั้งแรก
BOOTSTRAP_ADMIN_PHONE=0800000001
BOOTSTRAP_ADMIN_NAME=Admin

BOOTSTRAP_PROVIDER_PHONE=0800000002
BOOTSTRAP_PROVIDER_NAME=Provider
```

> **Bootstrap users**: เมื่อ app เริ่มทำงาน จะสร้าง user ที่มี role `admin` และ `provider` จาก env vars เหล่านี้โดยอัตโนมัติ ใช้สำหรับเข้าระบบครั้งแรก

---

### 3. Run ด้วย Go

```
go mod tidy
go run main.go
```

---

### 4. Run ด้วย Docker (แนะนำ)

```
docker-compose up --build
```

---

## 🧪 Testing

### Run Unit Tests

```
go test ./internal/service/... ./internal/swagger/...
```

ครอบคลุม:

* Service Layer
* Swagger Layer

Coverage threshold: **≥ 80%**

---

## 🔐 Authentication Flow

1. POST /api/auth/request-otp
2. POST /api/auth/verify-otp → ได้ JWT
3. ใช้ JWT ใน Header:

```
Authorization: Bearer <token>
```

> **Rate Limiting**: OTP request มี rate limit เพื่อป้องกัน abuse

---

## 📡 API Modules

### Auth (5 endpoints)

* Request OTP
* Verify OTP
* Register
* Get Profile
* Update Profile

---

### Category (5 endpoints)

* CRUD Category

---

### Provider & Zone (5 endpoints)

* Create provider
* Manage zone

---

### Queue Booking (4 endpoints)

* Book queue
* Check status
* History
* Cancel

---

### Queue Management (4 endpoints)

* Call queue
* Skip queue
* Complete queue

---

### Notification (4 endpoints)

* Get notifications
* Mark as read
* Delete
* Send

---

## 👥 Git Workflow

```
main        # production
develop     # integration
feature/*   # feature branches
```

### Flow:

1. checkout จาก develop
2. สร้าง feature branch
3. พัฒนา + commit
4. เปิด Pull Request → develop
5. review + merge

---

## 💡 Highlights

* Clean Architecture
* Modular Design (6 Modules)
* Unit Test Coverage ≥ 80%
* Security Scan (gosec + Trivy)
* Integration Tests (Postman/Newman)
* Dockerized Application (multi-arch)
* Swagger API Docs
* Postman Collection
* CI/CD Pipeline (GitHub Actions)
* Deploy to Render

---

## 📌 Notes

* ต้องตั้ง `JWT_SECRET` เป็นค่าที่ปลอดภัย (อย่างน้อย 32 characters)
* ต้องมี Docker ติดตั้งก่อนใช้งาน docker-compose
* ใช้ Swagger หรือ Postman สำหรับทดสอบ API
* Bootstrap users จะถูกสร้างอัตโนมัติตอน startup — ต้องกำหนด env vars ให้ครบ

---

## 👨‍💻 Contributors

| Name                      | Module           |
| ------------------------- | ---------------- |
| กิตติธัช เด่นสกุลประเสริฐ | Auth             |
| พิรญาณ์ เอนอ่อน           | Category         |
| ณัฏฐ์ ศรีสุวรรณกุล        | Provider & Zone  |
| ธนกฤต พิบูลย์สวัสดิ์      | Queue Booking    |
| พชร พรพงศ์                | Queue Management |
| กิตติภณ คำนวล             | Notification     |

---

## 🏁 Conclusion

QFlow เป็นระบบ Queue Management ที่ออกแบบให้:

* ใช้งานได้จริง
* โครงสร้างดี (Clean Architecture)
* รองรับการ deploy จริง (Docker + CI/CD + Render)
* มีเครื่องมือครบ (Swagger + Postman)

🚀 พร้อมต่อยอดสู่ production system

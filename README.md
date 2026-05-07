# 🚀 QFlow — Queue Management API (Refactored Edition)

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

## 🧱 Architecture (Refactored)

โครงสร้างถูกปรับเป็น **Clean Architecture**

```
internal/
├── domain/        # Entity + Interface
├── service/       # Business Logic
├── repository/    # Database Layer (GORM)
├── handler/       # HTTP Layer (Gin)
├── middleware/    # JWT Middleware
└── router/        # Route Registration
```

### ✅ Benefits

* แยกความรับผิดชอบชัดเจน
* Test ได้ง่าย (mock dependency)
* Maintain ง่าย
* Scale ได้ในอนาคต

---

## ⚙️ Tech Stack

* Language: Go
* Framework: Gin
* Database: PostgreSQL + GORM
* Authentication: JWT + OTP
* Container: Docker + Docker Compose
* API Docs: Swagger
* API Testing: Postman
* CI/CD: GitHub Actions
* Deployment: DockerHub

---

## 🔗 System Flow (CI/CD + Deployment)

```
Developer → Push Code → GitHub
                ↓
        GitHub Actions (CI)
        - Run Unit Tests
        - Build Docker Image
                ↓
        Push Image → DockerHub
                ↓
        Deploy (Docker Compose / Server)
                ↓
        Test via Swagger / Postman
```

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
postman/QFlow.postman_collection.json
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

Docker image จะถูก build และ push อัตโนมัติจาก CI/CD

```
docker pull <your-dockerhub-username>/qflow:latest
```

---

## 🔄 CI/CD (GitHub Actions)

Pipeline ทำงานอัตโนมัติเมื่อ push code:

### Steps:

1. Run Unit Tests
2. Build Docker Image
3. Push Image ไป DockerHub

### Example Workflow

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [develop, main]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'

      - name: Install dependencies
        run: go mod tidy

      - name: Run tests
        run: go test ./...

      - name: Build Docker image
        run: docker build -t qflow .

      - name: Login to DockerHub
        run: echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin

      - name: Push image
        run: |
          docker tag qflow ${{ secrets.DOCKER_USERNAME }}/qflow:latest
          docker push ${{ secrets.DOCKER_USERNAME }}/qflow:latest
```

---

## 🛠 Installation

### 1. Clone Project

```
git clone <repo-url>
cd QFlow
```

---

### 2. Setup Environment

สร้างไฟล์ `.env`

```
PORT=3000
DATABASE_URL=postgres://user:password@db:5432/qflow
JWT_SECRET=your-super-secret-key

BOOTSTRAP_ADMIN_PHONE=
BOOTSTRAP_ADMIN_NAME=Admin

BOOTSTRAP_PROVIDER_PHONE=
BOOTSTRAP_PROVIDER_NAME=Provider
```

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

### Run Unit Test

```
go test ./...
```

มี test ครอบคลุม:

* Service Layer
* Handler Layer
* Repository Layer

---

## 🔐 Authentication Flow

1. POST /api/auth/request-otp
2. POST /api/auth/verify-otp → ได้ JWT
3. ใช้ JWT ใน Header:

```
Authorization: Bearer <token>
```

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

* Clean Architecture (Refactored)
* Modular Design (6 Modules)
* Unit Test Coverage
* Dockerized Application
* Swagger API Docs
* Postman Collection
* CI/CD Pipeline (GitHub Actions)
* DockerHub Deployment Ready

---

## 📌 Notes

* ต้องตั้ง `JWT_SECRET` เป็นค่าที่ปลอดภัย
* ต้องมี Docker ติดตั้งก่อนใช้งาน docker-compose
* ใช้ Swagger หรือ Postman สำหรับทดสอบ API

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
* รองรับการ deploy จริง (Docker + CI/CD)
* มีเครื่องมือครบ (Swagger + Postman)

🚀 พร้อมต่อยอดสู่ production system

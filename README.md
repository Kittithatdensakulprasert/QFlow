# QFlow — Queue Management API

ระบบจัดการคิวออนไลน์ พัฒนาด้วย Go + Gin Framework สำหรับวิชา CS367 Web Service Development Concepts

**6 Modules — 27 Endpoints — 4 User Roles**

---

## Tech Stack

- **Language:** Go
- **Framework:** Gin
- **Database:** PostgreSQL + GORM
- **Auth:** JWT + OTP
- **Container:** Docker

---

## โครงสร้างโปรเจกต์

```
QFlow/
├── main.go
├── config/
│   └── config.go               ← โหลด environment variables
├── internal/
│   ├── domain/                 ← entities และ interfaces
│   │   ├── auth.go
│   │   ├── category.go
│   │   ├── provider.go
│   │   ├── queue.go
│   │   └── notification.go
│   ├── handler/                ← HTTP handlers (Gin)
│   │   ├── auth_handler.go
│   │   ├── category_handler.go
│   │   ├── provider_handler.go
│   │   ├── queue_handler.go
│   │   └── notification_handler.go
│   ├── service/                ← business logic
│   ├── repository/             ← database access layer
│   ├── middleware/
│   │   └── auth.go             ← JWT middleware
│   └── router/
│       └── router.go           ← ลงทะเบียน routes ทั้งหมด
└── db/
    └── migrations/             ← SQL migration files
```

---

## สมาชิกกลุ่มและการแบ่งงาน

### Auth Module

| ชื่อ | งานที่รับผิดชอบ | Branch |
|------|----------------|--------|
| กิตติธัช เด่นสกุลประเสริฐ | `POST /api/auth/request-otp`, `POST /api/auth/verify-otp`, `POST /api/auth/register`, `GET /api/auth/me`, `PUT /api/auth/me`, **Unit Test: Auth** | `feature/auth` |

### Category Module

| ชื่อ | งานที่รับผิดชอบ | Branch |
|------|----------------|--------|
| พิรญาณ์ เอนอ่อน | `GET /api/categories`, `GET /api/categories/:id`, `POST /api/categories`, `PUT /api/categories/:id`, `DELETE /api/categories/:id`, **Unit Test: Category** | `feature/category` |

### Provider & Zone Module

| ชื่อ | งานที่รับผิดชอบ | Branch |
|------|----------------|--------|
| ณัฏฐ์ ศรีสุวรรณกุล | `POST /api/providers`, `POST /api/providers/:id/zones`, `GET /api/providers`, `GET /api/providers/:id/zones`, `PATCH /api/zones/:id/toggle`, **Unit Test: Provider & Zone** | `feature/provider-zone` |

### Queue Booking Module

| ชื่อ | งานที่รับผิดชอบ | Branch |
|------|----------------|--------|
| ธนกฤต พิบูลย์สวัสดิ์ | `POST /api/queues/book`, `GET /api/queues/:queueNumber`, `GET /api/queues/history`, `PATCH /api/queues/:id/cancel`, **Docker (Dockerfile + docker-compose)**, **Unit Test: Queue Booking** | `feature/queue-booking` |

### Queue Management Module

| ชื่อ | งานที่รับผิดชอบ | Branch |
|------|----------------|--------|
| พชร พรพงศ์ | `GET /api/manage/queues/:zoneId`, `PATCH /api/manage/queues/:id/call`, `PATCH /api/manage/queues/:id/complete`, `PATCH /api/manage/queues/:id/skip`, **JWT Authentication Middleware**, **Unit Test: Queue Management** | `feature/queue-management` |

### Notification Module

| ชื่อ | งานที่รับผิดชอบ | Branch |
|------|----------------|--------|
| กิตติภณ คำนวล | `GET /api/notifications`, `PATCH /api/notifications/:id/read`, `DELETE /api/notifications/:id`, `POST /api/notifications/send`, **Database Schema + Migration**, **Unit Test: Notification** | `feature/notification` |

---

## API Endpoints

### Auth Module (5 endpoints)
| Method | Endpoint | คำอธิบาย | Role |
|--------|----------|----------|------|
| `POST` | `/api/auth/request-otp` | ขอ OTP | Guest |
| `POST` | `/api/auth/verify-otp` | ยืนยัน OTP → JWT | Guest |
| `POST` | `/api/auth/register` | ลงทะเบียน | Guest |
| `GET` | `/api/auth/me` | ดู Profile | User |
| `PUT` | `/api/auth/me` | แก้ไข Profile | User |

### Category Module (5 endpoints)
| Method | Endpoint | คำอธิบาย | Role |
|--------|----------|----------|------|
| `GET` | `/api/categories` | ดูประเภทร้านทั้งหมด | Guest |
| `GET` | `/api/categories/:id` | รายละเอียดของประเภท | Guest |
| `POST` | `/api/categories` | สร้างประเภท | Admin |
| `PUT` | `/api/categories/:id` | แก้ไขประเภท | Admin |
| `DELETE` | `/api/categories/:id` | ลบประเภท | Admin |

### Provider & Zone Module (5 endpoints)
| Method | Endpoint | คำอธิบาย | Role |
|--------|----------|----------|------|
| `POST` | `/api/providers` | สร้างผู้ให้บริการ | Admin |
| `GET` | `/api/providers` | ดูผู้ให้บริการทั้งหมด | Guest |
| `POST` | `/api/providers/:id/zones` | เพิ่มโซนใหม่ | Provider |
| `GET` | `/api/providers/:id/zones` | ดูโซน + จำนวนคิว | Guest |
| `PATCH` | `/api/zones/:id/toggle` | เปิด/ปิดโซน | Provider |

### Queue Booking Module (4 endpoints)
| Method | Endpoint | คำอธิบาย | Role |
|--------|----------|----------|------|
| `POST` | `/api/queues/book` | จองคิว → ได้รับเลขคิว | User |
| `GET` | `/api/queues/:queueNumber` | ดูสถานะคิว | User |
| `GET` | `/api/queues/history` | ประวัติการจองทั้งหมด | User |
| `PATCH` | `/api/queues/:id/cancel` | ยกเลิกคิว | User |

### Queue Management Module (4 endpoints)
> หมายเหตุ: endpoint ชุดนี้อยู่ในงานของ `feature/queue-management` และยังไม่ถูกรวมใน branch `feature/queue-booking`

### Notification Module (4 endpoints)
| Method | Endpoint | คำอธิบาย | Role |
|--------|----------|----------|------|
| `GET` | `/api/notifications` | ดูแจ้งเตือนทั้งหมด | User |
| `PATCH` | `/api/notifications/:id/read` | ทำเครื่องหมายว่าอ่านแล้ว | User |
| `DELETE` | `/api/notifications/:id` | ลบแจ้งเตือน | User |
| `POST` | `/api/notifications/send` | สร้างการแจ้งเตือนให้ผู้ใช้ที่ล็อกอินอยู่ | User |

---

## วิธีการติดตั้งและรัน

### Environment Variables

สร้างไฟล์ `.env` ที่ root:

```env
PORT=3000
DATABASE_URL=postgres://user:password@localhost:5432/qflow
JWT_SECRET=
BOOTSTRAP_ADMIN_PHONE=
BOOTSTRAP_ADMIN_NAME=Bootstrap Admin
BOOTSTRAP_PROVIDER_PHONE=
BOOTSTRAP_PROVIDER_NAME=Bootstrap Provider
```

`JWT_SECRET` ต้องตั้งเป็นค่าสุ่มจริงที่ยาวพอ และต้องไม่เป็นค่า default เช่น `secret` หรือ `your-secret-key-here` เพราะระบบจะไม่ start ถ้าใช้ค่าที่ไม่ปลอดภัย

ถ้าต้องทดสอบ endpoint ที่ต้องใช้ role `admin` หรือ `provider` ให้กำหนด `BOOTSTRAP_ADMIN_PHONE` หรือ `BOOTSTRAP_PROVIDER_PHONE` ก่อน start app จากนั้นเรียก `/api/auth/request-otp` และ `/api/auth/verify-otp` ด้วยเบอร์นั้นเพื่อรับ JWT ที่มี role ตามที่ bootstrap ไว้

### รันด้วย Go

```bash
go run main.go
```

Server จะรันที่ `http://localhost:3000`

### รันด้วย Docker

```bash
JWT_SECRET=replace-with-a-real-random-secret docker-compose up --build
```

หรือกำหนด `JWT_SECRET` ในไฟล์ `.env` ก่อน แล้วรัน:

```bash
docker-compose up --build
```

---

## Git Workflow

```
main        ← final version (merge จาก develop เมื่อเสร็จสิ้น)
develop     ← รวมงานจากทุก feature branch
feature/*   ← พัฒนาแต่ละ feature
```

1. แต่ละคน checkout จาก `develop` ไปยัง `feature/<feature-name>`
2. พัฒนาและ commit อย่างสม่ำเสมอ
3. เปิด Pull Request เข้า `develop`
4. รอ review จากเพื่อนในกลุ่มก่อน merge

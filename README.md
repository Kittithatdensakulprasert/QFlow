# QFlow — Queue Management API

ระบบจัดการคิวออนไลน์ พัฒนาด้วย Go (net/http) สำหรับวิชา CS367 Web Service Development Concepts

**6 Modules — 27 Endpoints — 4 User Roles**

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
| กิตติภณ คำนวล | `GET /api/notifications`, `PATCH /api/notifications/:id/read`, `DELETE /api/notifications/:id`, `POST /api/notifications/send`, **Database (Schema + Migration)**, **Unit Test: Notification** | `feature/notification` |

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
| `POST` | `/api/providers/:id/zones` | เพิ่มโซนใหม่ | Provider |
| `GET` | `/api/providers` | ดูผู้ให้บริการทั้งหมด | Guest |
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
| Method | Endpoint | คำอธิบาย | Role |
|--------|----------|----------|------|
| `GET` | `/api/manage/queues/:zoneId` | ดูรายการคิวทั้งหมดในโซน | Provider |
| `PATCH` | `/api/manage/queues/:id/call` | เรียกคิว + แจ้งเตือน | Provider |
| `PATCH` | `/api/manage/queues/:id/complete` | ปิดคิว (เสร็จสิ้น) | Provider |
| `PATCH` | `/api/manage/queues/:id/skip` | ข้ามคิว | Provider |

### Notification Module (4 endpoints)
| Method | Endpoint | คำอธิบาย | Role |
|--------|----------|----------|------|
| `GET` | `/api/notifications` | ดูแจ้งเตือนทั้งหมด | User |
| `PATCH` | `/api/notifications/:id/read` | ทำเครื่องหมายว่าอ่านแล้ว | User |
| `DELETE` | `/api/notifications/:id` | ลบแจ้งเตือน | User |
| `POST` | `/api/notifications/send` | ส่งการแจ้งเตือน (ระบบ) | System |

---

## Tech Stack

- **Language:** Go
- **Database:** (TBD)
- **Auth:** JWT + OTP
- **Container:** Docker

---

## วิธีการติดตั้งและรัน

### รันด้วย Go
```bash
go run main.go
```
Server จะรันที่ `http://localhost:3000`

### รันด้วย Docker
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

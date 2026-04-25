# QFlow — Queue Management API

ระบบจัดการคิวออนไลน์ พัฒนาด้วย Go (net/http) สำหรับวิชา CS367 Web Service Development Concepts

---

## สมาชิกกลุ่มและการแบ่งงาน

### Queue Booking Module

| ชื่อ | งานที่รับผิดชอบ | Branch |
|------|----------------|--------|
| กิตติธัช เด่นสกุลประเสริฐ | `POST /api/queues/book` จองคิว, `GET /api/queues/:queueNumber` ดูสถานะคิว, **Unit Test: Queue Booking** | `feature/queue-booking` |
| พิรญาณ์ เอนอ่อน | `GET /api/queues/history` ประวัติการจอง, `PATCH /api/queues/:id/cancel` ยกเลิกคิว, **Unit Test: Queue History & Cancel** | `feature/queue-history-cancel` |

### Queue Management Module

| ชื่อ | งานที่รับผิดชอบ | Branch |
|------|----------------|--------|
| ณัฏฐ์ ศรีสุวรรณกุล | `GET /api/manage/queues/:zoneId` ดูรายการคิวในโซน, `PATCH /api/manage/queues/:id/call` เรียกคิว + แจ้งเตือน, **Unit Test: Queue Call** | `feature/manage-queue-call` |
| ธนกฤต พิบูลย์สวัสดิ์ | `PATCH /api/manage/queues/:id/complete` ปิดคิว, `PATCH /api/manage/queues/:id/skip` ข้ามคิว, **Docker (Dockerfile + docker-compose)** | `feature/manage-queue-complete-skip` |

### Notification Module

| ชื่อ | งานที่รับผิดชอบ | Branch |
|------|----------------|--------|
| พชร พรพงศ์ | `GET /api/notifications` ดูแจ้งเตือนทั้งหมด, `PATCH /api/notifications/:id/read` อ่านการแจ้งเตือน, **JWT Authentication Middleware** | `feature/notification-read` |
| กิตติภณ คำนวล | `DELETE /api/notifications/:id` ลบแจ้งเตือน, `POST /api/notifications/send` ส่งการแจ้งเตือน, **Database (เชื่อมต่อ DB + Schema)** | `feature/notification-delete-send` |

---

## API Endpoints

### Queue Booking
| Method | Endpoint | คำอธิบาย |
|--------|----------|----------|
| `POST` | `/api/queues/book` | จองคิว → ได้รับเลขคิว |
| `GET` | `/api/queues/:queueNumber` | ดูสถานะคิว |
| `GET` | `/api/queues/history` | ประวัติการจองทั้งหมด |
| `PATCH` | `/api/queues/:id/cancel` | ยกเลิกคิว |

### Queue Management
| Method | Endpoint | คำอธิบาย |
|--------|----------|----------|
| `GET` | `/api/manage/queues/:zoneId` | ดูรายการคิวทั้งหมดในโซน |
| `PATCH` | `/api/manage/queues/:id/call` | เรียกคิว + แจ้งเตือน |
| `PATCH` | `/api/manage/queues/:id/complete` | ปิดคิว (เสร็จสิ้น) |
| `PATCH` | `/api/manage/queues/:id/skip` | ข้ามคิว |

### Notifications
| Method | Endpoint | คำอธิบาย |
|--------|----------|----------|
| `GET` | `/api/notifications` | ดูแจ้งเตือนทั้งหมด |
| `PATCH` | `/api/notifications/:id/read` | ทำเครื่องหมายว่าอ่านแล้ว |
| `DELETE` | `/api/notifications/:id` | ลบแจ้งเตือน |
| `POST` | `/api/notifications/send` | ส่งการแจ้งเตือน (ระบบ) |

---

## Tech Stack

- **Language:** Go
- **Database:** (TBD)
- **Auth:** JWT
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

# QFlow — GitHub Issues

> รายการ Issue ทั้งหมด 6 ข้อ สำหรับการพัฒนาระบบ QFlow

---

## Issue #1 — [กิตติธัช] Auth Module — OTP, JWT, Register, Profile

**สถานะ:** OPEN
**Branch:** `feature/auth`
**ผู้รับผิดชอบ:** กิตติธัช เด่นสกุลประเสริฐ

### งานที่ต้องทำ

- [ ] `POST /api/auth/request-otp` — ขอ OTP
- [ ] `POST /api/auth/verify-otp` — ยืนยัน OTP → ได้รับ JWT
- [ ] `POST /api/auth/register` — ลงทะเบียนผู้ใช้ใหม่
- [ ] `GET /api/auth/me` — ดู Profile ของตัวเอง
- [ ] `PUT /api/auth/me` — แก้ไข Profile
- [ ] Unit Test: Auth Module (coverage ≥ 80%)

### ไฟล์ที่เกี่ยวข้อง

```
internal/handler/auth_handler.go
internal/service/auth_service.go
internal/repository/auth_repository.go
internal/domain/auth.go
internal/middleware/auth.go
```

---

## Issue #2 — [พิรญาณ์] Category Module — CRUD ประเภทร้านอาหาร

**สถานะ:** OPEN
**Branch:** `feature/category`
**ผู้รับผิดชอบ:** พิรญาณ์ เอนอ่อน

### งานที่ต้องทำ

- [ ] `GET /api/categories` — ดูประเภทร้านทั้งหมด (Guest)
- [ ] `GET /api/categories/:id` — รายละเอียดของประเภท (Guest)
- [ ] `POST /api/categories` — สร้างประเภทใหม่ (Admin)
- [ ] `PUT /api/categories/:id` — แก้ไขประเภท (Admin)
- [ ] `DELETE /api/categories/:id` — ลบประเภท (Admin)
- [ ] Unit Test: Category Module (coverage ≥ 80%)

### ไฟล์ที่เกี่ยวข้อง

```
internal/handler/category_handler.go
internal/service/category_service.go
internal/repository/category_repository.go
internal/domain/category.go
```

---

## Issue #3 — [ณัฏฐ์] Provider & Zone Module — ผู้ให้บริการและโซน

**สถานะ:** OPEN
**Branch:** `feature/provider-zone`
**ผู้รับผิดชอบ:** ณัฏฐ์ ศรีสุวรรณกุล

### งานที่ต้องทำ

- [ ] `POST /api/providers` — สร้างผู้ให้บริการ (Admin)
- [ ] `GET /api/providers` — ดูผู้ให้บริการทั้งหมด (Guest)
- [ ] `POST /api/providers/:id/zones` — เพิ่มโซนใหม่ (Provider)
- [ ] `GET /api/providers/:id/zones` — ดูโซน + จำนวนคิว (Guest)
- [ ] `PATCH /api/zones/:id/toggle` — เปิด/ปิดโซน (Provider)
- [ ] Unit Test: Provider & Zone Module (coverage ≥ 80%)

### ไฟล์ที่เกี่ยวข้อง

```
internal/handler/provider_handler.go
internal/service/provider_service.go
internal/repository/provider_repository.go
internal/domain/provider.go
```

---

## Issue #4 — [ธนกฤต] Queue Booking Module และ Docker

**สถานะ:** OPEN
**Branch:** `feature/queue-booking`
**ผู้รับผิดชอบ:** ธนกฤต พิบูลย์สวัสดิ์

### งานที่ต้องทำ

- [ ] `POST /api/queues/book` — จองคิว → ได้รับเลขคิว (User)
- [ ] `GET /api/queues/:queueNumber` — ดูสถานะคิว (User)
- [ ] `GET /api/queues/history` — ประวัติการจองทั้งหมด (User)
- [ ] `PATCH /api/queues/:id/cancel` — ยกเลิกคิว (User)
- [ ] `Dockerfile` — build และ run ระบบ
- [ ] `docker-compose.yml` — รัน app + database
- [ ] Unit Test: Queue Booking Module (coverage ≥ 80%)

### ไฟล์ที่เกี่ยวข้อง

```
internal/handler/queue_handler.go   (BookQueue, GetHistory, GetQueue, CancelQueue)
internal/service/queue_service.go
internal/repository/queue_repository.go
internal/domain/queue.go
Dockerfile
docker-compose.yml
```

---

## Issue #5 — [พชร] Queue Management Module และ JWT Middleware

**สถานะ:** OPEN
**Branch:** `feature/queue-management`
**ผู้รับผิดชอบ:** พชร พรพงศ์

### งานที่ต้องทำ

- [ ] `GET /api/manage/queues/:zoneId` — ดูรายการคิวทั้งหมดในโซน (Provider)
- [ ] `PATCH /api/manage/queues/:id/call` — เรียกคิว + แจ้งเตือน (Provider)
- [ ] `PATCH /api/manage/queues/:id/complete` — ปิดคิว (เสร็จสิ้น) (Provider)
- [ ] `PATCH /api/manage/queues/:id/skip` — ข้ามคิว (Provider)
- [ ] JWT Authentication Middleware — ตรวจสอบ Token และ Role
- [ ] Unit Test: Queue Management Module (coverage ≥ 80%)

### ไฟล์ที่เกี่ยวข้อง

```
internal/handler/queue_handler.go   (GetQueuesByZone, CallQueue, CompleteQueue, SkipQueue)
internal/service/queue_service.go
internal/repository/queue_repository.go
internal/middleware/auth.go
```

---

## Issue #6 — [กิตติภณ] Notification Module และ Database Schema

**สถานะ:** OPEN
**Branch:** `feature/notification`
**ผู้รับผิดชอบ:** กิตติภณ คำนวล

### งานที่ต้องทำ

- [ ] `GET /api/notifications` — ดูแจ้งเตือนทั้งหมด (User)
- [ ] `POST /api/notifications/send` — ส่งการแจ้งเตือน (System)
- [ ] `PATCH /api/notifications/:id/read` — ทำเครื่องหมายว่าอ่านแล้ว (User)
- [ ] `DELETE /api/notifications/:id` — ลบแจ้งเตือน (User)
- [ ] Database Schema — GORM + PostgreSQL (ครอบคลุมทุก Module)
- [ ] Migration Files — `db/migrations/`
- [ ] Unit Test: Notification Module (coverage ≥ 80%)

### ไฟล์ที่เกี่ยวข้อง

```
internal/handler/notification_handler.go
internal/service/notification_service.go
internal/repository/notification_repository.go
internal/domain/notification.go
db/migrations/
```

---

## สรุป

| Issue | Module | ผู้รับผิดชอบ | Branch | Endpoints |
|-------|--------|-------------|--------|-----------|
| #1 | Auth | กิตติธัช | `feature/auth` | 5 |
| #2 | Category | พิรญาณ์ | `feature/category` | 5 |
| #3 | Provider & Zone | ณัฏฐ์ | `feature/provider-zone` | 5 |
| #4 | Queue Booking + Docker | ธนกฤต | `feature/queue-booking` | 4 |
| #5 | Queue Management + JWT | พชร | `feature/queue-management` | 4 |
| #6 | Notification + Database | กิตติภณ | `feature/notification` | 4 |

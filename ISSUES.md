# QFlow — GitHub Issues

---

## Issue #1 — [กิตติธัช] Auth Module — OTP, JWT, Register, Profile

**สถานะ:** OPEN
**Branch:** `feature/auth`

### ผู้รับผิดชอบ
กิตติธัช เด่นสกุลประเสริฐ

### งานที่ต้องทำ
- [ ] `POST /api/auth/request-otp` — ขอ OTP
- [ ] `POST /api/auth/verify-otp` — ยืนยัน OTP → JWT
- [ ] `POST /api/auth/register` — ลงทะเบียน
- [ ] `GET /api/auth/me` — ดู Profile
- [ ] `PUT /api/auth/me` — แก้ไข Profile
- [ ] Unit Test: Auth Module (coverage ≥ 80%)

---

## Issue #2 — [พิรญาณ์] Category Module — CRUD ประเภทร้านอาหาร

**สถานะ:** OPEN
**Branch:** `feature/category`

### ผู้รับผิดชอบ
พิรญาณ์ เอนอ่อน

### งานที่ต้องทำ
- [ ] `GET /api/categories` — ดูประเภทร้านทั้งหมด
- [ ] `GET /api/categories/:id` — รายละเอียดประเภท
- [ ] `POST /api/categories` — สร้างประเภท (Admin)
- [ ] `PUT /api/categories/:id` — แก้ไขประเภท (Admin)
- [ ] `DELETE /api/categories/:id` — ลบประเภท (Admin)
- [ ] Unit Test: Category Module (coverage ≥ 80%)

---

## Issue #3 — [ณัฏฐ์] Provider & Zone Module — ผู้ให้บริการและโซน

**สถานะ:** OPEN
**Branch:** `feature/provider-zone`

### ผู้รับผิดชอบ
ณัฏฐ์ ศรีสุวรรณกุล

### งานที่ต้องทำ
- [ ] `POST /api/providers` — สร้างผู้ให้บริการ (Admin)
- [ ] `POST /api/providers/:id/zones` — เพิ่มโซนใหม่ (Provider)
- [ ] `GET /api/providers` — ดูผู้ให้บริการทั้งหมด
- [ ] `GET /api/providers/:id/zones` — ดูโซน + จำนวนคิว
- [ ] `PATCH /api/zones/:id/toggle` — เปิด/ปิดโซน (Provider)
- [ ] Unit Test: Provider & Zone (coverage ≥ 80%)

---

## Issue #4 — [ธนกฤต] Queue Booking Module และ Docker

**สถานะ:** OPEN
**Branch:** `feature/queue-booking`

### ผู้รับผิดชอบ
ธนกฤต พิบูลย์สวัสดิ์

### งานที่ต้องทำ
- [ ] `POST /api/queues/book` — จองคิว → ได้รับเลขคิว
- [ ] `GET /api/queues/:queueNumber` — ดูสถานะคิว
- [ ] `GET /api/queues/history` — ประวัติการจอง
- [ ] `PATCH /api/queues/:id/cancel` — ยกเลิกคิว
- [ ] `Dockerfile` — build และ run ระบบ
- [ ] `docker-compose.yml` — รัน app + database
- [ ] Unit Test: Queue Booking (coverage ≥ 80%)

---

## Issue #5 — [พชร] Queue Management Module และ JWT Middleware

**สถานะ:** OPEN
**Branch:** `feature/queue-management`

### ผู้รับผิดชอบ
พชร พรพงศ์

### งานที่ต้องทำ
- [ ] `GET /api/manage/queues/:zoneId` — ดูรายการคิวในโซน (Provider)
- [ ] `PATCH /api/manage/queues/:id/call` — เรียกคิว + แจ้งเตือน (Provider)
- [ ] `PATCH /api/manage/queues/:id/complete` — ปิดคิว (Provider)
- [ ] `PATCH /api/manage/queues/:id/skip` — ข้ามคิว (Provider)
- [ ] JWT Authentication Middleware — ป้องกัน route ตาม Role
- [ ] Unit Test: Queue Management (coverage ≥ 80%)

---

## Issue #6 — [กิตติภณ] Notification Module และ Database Schema

**สถานะ:** OPEN
**Branch:** `feature/notification`

### ผู้รับผิดชอบ
กิตติภณ คำนวล

### งานที่ต้องทำ
- [x] `DELETE /api/notifications/:id` — ลบแจ้งเตือน ✅
- [x] `POST /api/notifications/send` — ส่งการแจ้งเตือน ✅
- [x] `GET /api/notifications` — ดูแจ้งเตือนทั้งหมด ✅
- [x] `PATCH /api/notifications/:id/read` — ทำเครื่องหมายว่าอ่านแล้ว ✅
- [x] Database Schema — GORM + PostgreSQL ✅ (PR #11)
- [x] Unit Test: Notification (coverage 96-100%) ✅ (PR #10)

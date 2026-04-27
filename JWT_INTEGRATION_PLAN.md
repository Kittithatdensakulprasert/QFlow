# JWT Middleware Integration Plan

## 🎯 Current Status
✅ **Authentication handlers are ready** - Support context + fallback
⏳ **JWT middleware is in separate issue** - Not yet implemented

## 📋 Integration Steps

### Phase 1: JWT Middleware Implementation (Separate Issue)
**คนทำ JWT middleware ต้อง:**
```go
// JWT middleware ต้องใส่ user_id ใน context
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Validate JWT token
        // 2. Extract user ID from token
        // 3. Set user ID in context
        c.Set("user_id", userID) // 🔑 IMPORTANT: Use "user_id" key
        
        c.Next()
    }
}
```

### Phase 2: Remove Fallback (After JWT Ready)
**เมื่อ JWT middleware เสร็จแล้ว:**
1. **ลบ fallback logic** จาก `GetProfile()` และ `UpdateProfile()`
2. **เหลือเฉพาะ context** จาก JWT middleware
3. **ทำให้ strict** ตาม authentication จริง

## 🔄 Handler Changes After JWT Ready

### Current Handler (with fallback):
```go
func (h *AuthHandler) GetProfile(c *gin.Context) {
    // Try to get user ID from context (JWT middleware) first
    userIDInterface, exists := c.Get("user_id")
    var userIDStr string
    
    if exists {
        // From context (JWT middleware)
        // ... type conversion logic
    } else {
        // Fallback: try header first, then query parameter
        userIDStr = c.GetHeader("X-User-ID")
        if userIDStr == "" {
            userIDStr = c.Query("user_id")
        }
    }
    // ... rest of logic
}
```

### Future Handler (strict - no fallback):
```go
func (h *AuthHandler) GetProfile(c *gin.Context) {
    // Only get from context (JWT middleware)
    userIDInterface, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }
    
    // Convert to uint (JWT middleware should provide consistent type)
    userID, err := convertToUint(userIDInterface)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
        return
    }
    
    // ... rest of logic
}
```

## 🎯 Key Requirements for JWT Middleware

### 1. Context Key
- **ต้องใช้ key**: `"user_id"`
- **Type consistency**: แนะนำให้ใช้ `uint` หรือ `int`

### 2. Error Handling
- JWT middleware ต้อง handle token validation
- Return 401 ถ้า token invalid/missing
- Set context ถ้า token valid

### 3. Route Protection
- Apply JWT middleware กับ routes ที่ต้องการ auth:
```go
// In router setup
api.Use(JWTMiddleware())
api.GET("/auth/me", authHandler.GetProfile)
api.PUT("/auth/me", authHandler.UpdateProfile)
```

## 📝 Testing Plan

### Before JWT (Current State):
```bash
# Test with fallback
curl -X GET http://localhost:3000/api/auth/me -H "X-User-ID: 1"
curl -X GET "http://localhost:3000/api/auth/me?user_id=1"
```

### After JWT (Future State):
```bash
# Test with JWT token only
curl -X GET http://localhost:3000/api/auth/me -H "Authorization: Bearer <jwt_token>"
```

## 🚀 Migration Timeline

1. **Now**: ✅ Handler ready with fallback
2. **JWT Issue**: 🔨 Implement JWT middleware
3. **Integration**: 🔗 Test JWT + handler together
4. **Cleanup**: 🧹 Remove fallback logic
5. **Production**: 🎯 Strict authentication only

## 📞 Communication

**คนทำ JWT middleware ต้องรู้:**
- Handler รองรับ `user_id` จาก context แล้ว
- ใช้ key `"user_id"` ใน `c.Set()`
- พร้อมทดสอบ integration เมื่อ middleware เสร็จ

**คนทำ handler ต้องรู้:**
- รอ JWT middleware เสร็จก่อน
- จะลบ fallback logic เมื่อพร้อม
- ทำให้ strict ตาม production requirement

---

**Status**: 🟡 Ready for JWT middleware integration
**Next Step**: 🔨 JWT middleware implementation (separate issue)

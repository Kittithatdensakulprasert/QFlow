package schema

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Phone     string    `gorm:"uniqueIndex;not null;size:20" json:"phone"`
	FirstName string    `gorm:"not null;size:100" json:"firstName"`
	LastName  string    `gorm:"not null;size:100" json:"lastName"`
	Email     string    `gorm:"size:255" json:"email"`
	Role      string    `gorm:"not null;default:'user';size:20" json:"role"` // user, provider, admin
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null;size:100" json:"name"` // ร้านอาหารจีน, ชาบู, ปิ้งย่าง
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Providers []Provider `gorm:"foreignKey:CategoryID" json:"providers,omitempty"`
}

type Provider struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"not null;size:200" json:"name"`
	CategoryID uint      `gorm:"not null" json:"categoryId"`
	OwnerID    uint      `gorm:"not null" json:"ownerId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Owner    User     `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Zones    []Zone   `gorm:"foreignKey:ProviderID" json:"zones,omitempty"`
}

type Zone struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderID uint      `gorm:"not null" json:"providerId"`
	Name       string    `gorm:"not null;size:100" json:"name"`
	IsOpen     bool      `gorm:"not null;default:true" json:"isOpen"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	Provider Provider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	Queues   []Queue  `gorm:"foreignKey:ZoneID" json:"queues,omitempty"`
}

type Queue struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ZoneID      uint      `gorm:"not null" json:"zoneId"`
	UserID      uint      `gorm:"not null" json:"userId"`
	QueueNumber int       `gorm:"not null" json:"queueNumber"`
	Status      string    `gorm:"not null;default:'waiting';size:20" json:"status"` // waiting, called, completed, cancelled, skipped
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	Zone Zone `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Notification struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"userId"`
	Message   string    `gorm:"not null" json:"message"`
	IsRead    bool      `gorm:"not null;default:false" json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

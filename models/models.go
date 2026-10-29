package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RolePartner UserRole = "partner"
	RoleGuest   UserRole = "guest"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:120;not null" json:"name"`
	Email        string    `gorm:"size:150;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         UserRole  `gorm:"type:enum('admin','partner','guest');default:'guest';not null" json:"role"`
	Phone        string    `gorm:"size:30" json:"phone"`
	Country      string    `gorm:"size:30;default:'Singapore'" json:"country"`
	Avatar       string    `gorm:"size:255" json:"avatar"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type HotelStatus string

const (
	HotelStatusPending   HotelStatus = "pending"
	HotelStatusActive    HotelStatus = "active"
	HotelStatusSuspended HotelStatus = "suspended"
)

type Hotel struct {
	ID                uint        `gorm:"primaryKey" json:"id"`
	Name              string      `gorm:"size:150;not null;index" json:"name"`
	Tagline           string      `gorm:"size:255" json:"tagline"`
	Type              string      `gorm:"size:50;index" json:"type"` // Hotel, Resort, Boutique Hotel, Villa
	City              string      `gorm:"size:100;not null;index" json:"city"`
	Area              string      `gorm:"size:150" json:"area"`
	Country           string      `gorm:"size:30;not null;index" json:"country"` // Singapore, Indonesia, Malaysia
	Address           string      `gorm:"size:255;not null" json:"address"`
	DistanceText      string      `gorm:"size:150" json:"distance_text"`
	Badge             string      `gorm:"size:50" json:"badge"` // Preferred Partner, Heritage Luxury, Eco Luxury, Top Choice
	Latitude          float64     `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude         float64     `gorm:"type:decimal(10,7)" json:"longitude"`
	StarRating        int         `gorm:"default:5" json:"star_rating"`
	Rating            float64     `gorm:"type:decimal(3,1);default:4.8" json:"rating"`
	ReviewCount       int         `gorm:"default:0" json:"review_count"`
	PriceStart        float64     `gorm:"type:decimal(15,2);not null" json:"price_start"`
	OriginalPrice     float64     `gorm:"type:decimal(15,2)" json:"original_price"`
	DiscountPercent   int         `gorm:"default:0" json:"discount_percent"`
	FlashSale         bool        `gorm:"default:false" json:"flash_sale"`
	FreeBreakfast     bool        `gorm:"default:true" json:"free_breakfast"`
	FreeCancellation  bool        `gorm:"default:true" json:"free_cancellation"`
	PayAtHotel        bool        `gorm:"default:false" json:"pay_at_hotel"`
	HasPool           bool        `gorm:"default:true" json:"has_pool"`
	HasWifi           bool        `gorm:"default:true" json:"has_wifi"`
	HasSpa            bool        `gorm:"default:true" json:"has_spa"`
	HasRestaurant     bool        `gorm:"default:true" json:"has_restaurant"`
	HasGym            bool        `gorm:"default:true" json:"has_gym"`
	ReviewCleanliness float64     `gorm:"type:decimal(3,1);default:9.6" json:"review_cleanliness"`
	ReviewComfort     float64     `gorm:"type:decimal(3,1);default:9.5" json:"review_comfort"`
	ReviewService     float64     `gorm:"type:decimal(3,1);default:9.6" json:"review_service"`
	ReviewLocation    float64     `gorm:"type:decimal(3,1);default:9.8" json:"review_location"`
	Currency          string      `gorm:"size:10;default:'SGD'" json:"currency"`
	Image             string      `gorm:"type:text;not null" json:"image"`
	Gallery           string      `gorm:"type:text" json:"gallery"`       // JSON Array of image URLs
	Facilities        string      `gorm:"type:text" json:"facilities"`    // JSON Array of strings
	Description       string      `gorm:"type:text" json:"description"`
	PartnerID         uint        `gorm:"index;not null" json:"partner_id"`
	Partner           *User       `gorm:"foreignKey:PartnerID" json:"partner,omitempty"`
	Status            HotelStatus `gorm:"type:enum('pending','active','suspended');default:'active';not null" json:"status"`
	Rooms             []Room      `gorm:"foreignKey:HotelID" json:"rooms,omitempty"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

type Room struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	HotelID          uint      `gorm:"index;not null" json:"hotel_id"`
	Name             string    `gorm:"size:120;not null" json:"name"`
	Type             string    `gorm:"size:60;not null" json:"type"`
	Price            float64   `gorm:"type:decimal(15,2);not null" json:"price"`
	OriginalPrice    float64   `gorm:"type:decimal(15,2)" json:"original_price"`
	MaxGuests        int       `gorm:"default:2;not null" json:"max_guests"`
	BedType          string    `gorm:"size:80;default:'1 King Bed'" json:"bed_type"`
	SizeSqm          int       `gorm:"default:35" json:"size_sqm"`
	RemainingRooms   int       `gorm:"default:5;not null" json:"remaining_rooms"`
	TotalRooms       int       `gorm:"default:10;not null" json:"total_rooms"`
	FreeBreakfast    bool      `gorm:"default:true" json:"free_breakfast"`
	FreeCancellation bool      `gorm:"default:true" json:"free_cancellation"`
	Amenities        string    `gorm:"type:text" json:"amenities"` // JSON array
	Image            string    `gorm:"type:text" json:"image"`
	Status           string    `gorm:"size:30;default:'available'" json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type BookingStatus string

const (
	BookingConfirmed  BookingStatus = "confirmed"
	BookingCheckedIn  BookingStatus = "checked_in"
	BookingCompleted  BookingStatus = "completed"
	BookingCancelled  BookingStatus = "cancelled"
)

type Booking struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	BookingCode     string        `gorm:"size:40;uniqueIndex;not null" json:"booking_code"`
	HotelID         uint          `gorm:"index;not null" json:"hotel_id"`
	Hotel           *Hotel        `gorm:"foreignKey:HotelID" json:"hotel,omitempty"`
	RoomID          uint          `gorm:"index;not null" json:"room_id"`
	Room            *Room         `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	UserID          *uint         `gorm:"index" json:"user_id,omitempty"`
	GuestName       string        `gorm:"size:120;not null" json:"guest_name"`
	GuestEmail      string        `gorm:"size:150;not null;index" json:"guest_email"`
	GuestPhone      string        `gorm:"size:30;not null" json:"guest_phone"`
	CheckIn         time.Time     `gorm:"type:date;not null" json:"check_in"`
	CheckOut        time.Time     `gorm:"type:date;not null" json:"check_out"`
	Nights          int           `gorm:"default:1;not null" json:"nights"`
	GuestsCount     int           `gorm:"default:1;not null" json:"guests_count"`
	TotalAmount     float64       `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	Currency        string        `gorm:"size:10;default:'SGD'" json:"currency"`
	Status          BookingStatus `gorm:"type:enum('confirmed','checked_in','completed','cancelled');default:'confirmed';not null" json:"status"`
	PaymentMethod   string        `gorm:"size:50;default:'credit_card'" json:"payment_method"`
	PaymentStatus   string        `gorm:"size:30;default:'paid'" json:"payment_status"`
	SpecialRequests string        `gorm:"type:text" json:"special_requests"`
	QrCodeData      string        `gorm:"size:255" json:"qr_code_data"`
	CheckedInAt     *time.Time    `json:"checked_in_at"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type SettlementStatus string

const (
	SettlementPending     SettlementStatus = "pending"
	SettlementApproved    SettlementStatus = "approved"
	SettlementTransferred SettlementStatus = "transferred"
	SettlementRejected    SettlementStatus = "rejected"
)

type Settlement struct {
	ID                 uint             `gorm:"primaryKey" json:"id"`
	SettlementCode     string           `gorm:"size:50;uniqueIndex;not null" json:"settlement_code"`
	PartnerID          uint             `gorm:"index;not null" json:"partner_id"`
	HotelID            uint             `gorm:"index;not null" json:"hotel_id"`
	Hotel              *Hotel           `gorm:"foreignKey:HotelID" json:"hotel,omitempty"`
	BookingID          uint             `gorm:"index;not null" json:"booking_id"`
	Booking            *Booking         `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	GrossAmount        float64          `gorm:"type:decimal(15,2);not null" json:"gross_amount"`
	PlatformCommission float64          `gorm:"type:decimal(15,2);not null" json:"platform_commission"`
	NetPayout          float64          `gorm:"type:decimal(15,2);not null" json:"net_payout"`
	Currency           string           `gorm:"size:10;default:'SGD'" json:"currency"`
	Status             SettlementStatus `gorm:"type:enum('pending','approved','transferred','rejected');default:'pending';not null" json:"status"`
	PayoutDate         *time.Time       `json:"payout_date"`
	BankName           string           `gorm:"size:80" json:"bank_name"`
	BankAccount        string           `gorm:"size:80" json:"bank_account"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}

type PartnerApplication struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	BusinessName  string    `gorm:"size:150;not null" json:"business_name"`
	ContactPerson string    `gorm:"size:120;not null" json:"contact_person"`
	Email         string    `gorm:"size:150;not null" json:"email"`
	Phone         string    `gorm:"size:40;not null" json:"phone"`
	Country       string    `gorm:"size:30;not null" json:"country"`
	City          string    `gorm:"size:100;not null" json:"city"`
	PropertyCount int       `gorm:"default:1" json:"property_count"`
	Status        string    `gorm:"type:enum('pending','approved','rejected');default:'pending';not null" json:"status"`
	Notes         string    `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Hotel{},
		&Room{},
		&Booking{},
		&Settlement{},
		&PartnerApplication{},
	)
}

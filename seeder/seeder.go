package seeder

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"hotels-api/config"
	"hotels-api/models"

	"golang.org/x/crypto/bcrypt"
)

//go:embed hotels_seed.json
var hotelsSeedJSON []byte

type SeedRoom struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Size           string   `json:"size"`
	BedType        string   `json:"bedType"`
	MaxGuests      int      `json:"maxGuests"`
	FreeBreakfast  bool     `json:"freeBreakfast"`
	FreeCancellation bool   `json:"freeCancellation"`
	OriginalPrice  float64  `json:"originalPrice"`
	Price          float64  `json:"price"`
	Image          string   `json:"image"`
	Amenities      []string `json:"amenities"`
	RemainingRooms int      `json:"remainingRooms"`
}

type SeedHotel struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Type             string     `json:"type"`
	StarRating       int        `json:"starRating"`
	Country          string     `json:"country"`
	City             string     `json:"city"`
	Area             string     `json:"area"`
	DistanceText     string     `json:"distanceText"`
	Badge            string     `json:"badge"`
	Rating           float64    `json:"rating"`
	ReviewsCount     int        `json:"reviewsCount"`
	OriginalPrice    float64    `json:"originalPrice"`
	Price            float64    `json:"price"`
	DiscountPercent  int        `json:"discountPercent"`
	FlashSale        bool       `json:"flashSale"`
	FreeBreakfast    bool       `json:"freeBreakfast"`
	FreeCancellation bool       `json:"freeCancellation"`
	PayAtHotel       bool       `json:"payAtHotel"`
	HasPool          bool       `json:"hasPool"`
	HasWifi          bool       `json:"hasWifi"`
	HasSpa           bool       `json:"hasSpa"`
	HasRestaurant    bool       `json:"hasRestaurant"`
	HasGym           bool       `json:"hasGym"`
	MainImage        string     `json:"mainImage"`
	Thumbnails       []string   `json:"thumbnails"`
	Description      string     `json:"description"`
	ReviewBreakdown  struct {
		Cleanliness float64 `json:"cleanliness"`
		Comfort     float64 `json:"comfort"`
		Service     float64 `json:"service"`
		Location    float64 `json:"location"`
	} `json:"reviewBreakdown"`
	Coordinates struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"coordinates"`
	Rooms []SeedRoom `json:"rooms"`
}

func Seed() {
	log.Println("[INFO] Seeding / refreshing database with 17 luxury properties and 100% unique authentic photography...")

	// Truncate/Clean old tables to avoid conflicts
	config.DB.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	config.DB.Exec("TRUNCATE TABLE settlements;")
	config.DB.Exec("TRUNCATE TABLE bookings;")
	config.DB.Exec("TRUNCATE TABLE rooms;")
	config.DB.Exec("TRUNCATE TABLE hotels;")
	config.DB.Exec("TRUNCATE TABLE users;")
	config.DB.Exec("TRUNCATE TABLE partner_applications;")
	config.DB.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	// 1. Users
	hashedAdmin, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	hashedPartner, _ := bcrypt.GenerateFromPassword([]byte("partner123"), bcrypt.DefaultCost)
	hashedGuest, _ := bcrypt.GenerateFromPassword([]byte("guest123"), bcrypt.DefaultCost)

	admin := models.User{
		Name:         "Superadmin Director",
		Email:        "admin@hotels.com",
		PasswordHash: string(hashedAdmin),
		Role:         models.RoleAdmin,
		Phone:        "+65 8920 1192",
		Country:      "Singapore",
		Avatar:       "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80",
	}

	partnerSG := models.User{
		Name:         "Marina Bay Hospitality Ltd",
		Email:        "partner.sg@marinabay.com",
		PasswordHash: string(hashedPartner),
		Role:         models.RolePartner,
		Phone:        "+65 6688 8868",
		Country:      "Singapore",
		Avatar:       "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&auto=format&fit=crop&q=80",
	}

	partnerID := models.User{
		Name:         "Nusantara Resorts & Spas",
		Email:        "partner.id@ayana.com",
		PasswordHash: string(hashedPartner),
		Role:         models.RolePartner,
		Phone:        "+62 361 702222",
		Country:      "Indonesia",
		Avatar:       "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150&auto=format&fit=crop&q=80",
	}

	partnerMY := models.User{
		Name:         "Kuala Lumpur Prestige Hotels",
		Email:        "partner.my@kualalumpur.com",
		PasswordHash: string(hashedPartner),
		Role:         models.RolePartner,
		Phone:        "+60 3 2142 8000",
		Country:      "Malaysia",
		Avatar:       "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&auto=format&fit=crop&q=80",
	}

	guest := models.User{
		Name:         "Michael Tan",
		Email:        "guest@example.com",
		PasswordHash: string(hashedGuest),
		Role:         models.RoleGuest,
		Phone:        "+65 9123 4567",
		Country:      "Singapore",
		Avatar:       "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150&auto=format&fit=crop&q=80",
	}

	config.DB.Create(&admin)
	config.DB.Create(&partnerSG)
	config.DB.Create(&partnerID)
	config.DB.Create(&partnerMY)
	config.DB.Create(&guest)

	// Helper for JSON strings
	toJSON := func(v interface{}) string {
		b, _ := json.Marshal(v)
		return string(b)
	}

	var rawHotels []SeedHotel
	if err := json.Unmarshal(hotelsSeedJSON, &rawHotels); err != nil {
		log.Fatalf("[FATAL] Failed to parse hotelsSeedJSON: %v", err)
	}

	log.Printf("[INFO] Seeding %d real luxury properties from JSON...\n", len(rawHotels))

	var seededHotels []models.Hotel

	for _, sh := range rawHotels {
		var pID uint = partnerID.ID
		currency := "IDR"
		if sh.Country == "Singapore" {
			pID = partnerSG.ID
			currency = "SGD"
		} else if sh.Country == "Malaysia" {
			pID = partnerMY.ID
			currency = "MYR"
		}

		// Collect facilities
		var facs []string
		if sh.HasPool {
			facs = append(facs, "Swimming Pool")
		}
		if sh.HasWifi {
			facs = append(facs, "Free High-Speed Wi-Fi")
		}
		if sh.HasSpa {
			facs = append(facs, "Wellness Spa")
		}
		if sh.HasRestaurant {
			facs = append(facs, "Fine Dining Restaurant")
		}
		if sh.HasGym {
			facs = append(facs, "Fitness Center")
		}

		h := models.Hotel{
			Name:              sh.Name,
			Tagline:           sh.DistanceText,
			Type:              sh.Type,
			City:              sh.City,
			Area:              sh.Area,
			Country:           sh.Country,
			Address:           sh.Area,
			DistanceText:      sh.DistanceText,
			Badge:             sh.Badge,
			Latitude:          sh.Coordinates.Lat,
			Longitude:         sh.Coordinates.Lng,
			StarRating:        sh.StarRating,
			Rating:            sh.Rating,
			ReviewCount:       sh.ReviewsCount,
			PriceStart:        sh.Price,
			OriginalPrice:     sh.OriginalPrice,
			DiscountPercent:   sh.DiscountPercent,
			FlashSale:         sh.FlashSale,
			FreeBreakfast:     sh.FreeBreakfast,
			FreeCancellation:  sh.FreeCancellation,
			PayAtHotel:        sh.PayAtHotel,
			HasPool:           sh.HasPool,
			HasWifi:           sh.HasWifi,
			HasSpa:            sh.HasSpa,
			HasRestaurant:     sh.HasRestaurant,
			HasGym:            sh.HasGym,
			ReviewCleanliness: sh.ReviewBreakdown.Cleanliness,
			ReviewComfort:     sh.ReviewBreakdown.Comfort,
			ReviewService:     sh.ReviewBreakdown.Service,
			ReviewLocation:    sh.ReviewBreakdown.Location,
			Currency:          currency,
			Image:             sh.MainImage,
			Gallery:           toJSON(sh.Thumbnails),
			Facilities:        toJSON(facs),
			Description:       sh.Description,
			PartnerID:         pID,
			Status:            models.HotelStatusActive,
		}

		config.DB.Create(&h)
		seededHotels = append(seededHotels, h)

		for _, sr := range sh.Rooms {
			rem := sr.RemainingRooms
			if rem <= 0 {
				rem = 5
			}
			r := models.Room{
				HotelID:          h.ID,
				Name:             sr.Name,
				Type:             h.Type,
				Price:            sr.Price,
				OriginalPrice:    sr.OriginalPrice,
				MaxGuests:        sr.MaxGuests,
				BedType:          sr.BedType,
				SizeSqm:          38,
				RemainingRooms:   rem,
				TotalRooms:       rem + 5,
				FreeBreakfast:    sr.FreeBreakfast,
				FreeCancellation: sr.FreeCancellation,
				Amenities:        toJSON(sr.Amenities),
				Image:            sr.Image,
				Status:           "available",
			}
			config.DB.Create(&r)
		}
	}

	// 3. Seed Realistic Bookings & Settlements
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)
	inThreeDays := today.AddDate(0, 0, 3)

	type seedBookingItem struct {
		hotelIndex  int
		guestName   string
		guestEmail  string
		guestPhone  string
		checkIn     time.Time
		checkOut    time.Time
		nights      int
		status      models.BookingStatus
		code        string
		settleCode  string
		settleState models.SettlementStatus
	}

	mockBookings := []seedBookingItem{
		{
			hotelIndex:  0, // Marina Bay Sands
			guestName:   "Alexander Wright",
			guestEmail:  "a.wright@techglobal.com",
			guestPhone:  "+65 9234 5678",
			checkIn:     today,
			checkOut:    tomorrow,
			nights:      1,
			status:      models.BookingConfirmed,
			code:        "HTL-SG-8812",
			settleCode:  "STL-SG-5011",
			settleState: models.SettlementPending,
		},
		{
			hotelIndex:  0, // Marina Bay Sands
			guestName:   "Sophie Laurent",
			guestEmail:  "sophie.l@luxurygroup.fr",
			guestPhone:  "+33 612 3456",
			checkIn:     yesterday,
			checkOut:    tomorrow,
			nights:      2,
			status:      models.BookingCheckedIn,
			code:        "HTL-SG-9104",
			settleCode:  "STL-SG-5012",
			settleState: models.SettlementApproved,
		},
		{
			hotelIndex:  11, // The Apurva Kempinski Bali
			guestName:   "Budi Santoso",
			guestEmail:  "budi.santoso@jakarta.id",
			guestPhone:  "+62 812 3456 7890",
			checkIn:     today,
			checkOut:    inThreeDays,
			nights:      3,
			status:      models.BookingConfirmed,
			code:        "HTL-ID-7721",
			settleCode:  "STL-ID-6019",
			settleState: models.SettlementPending,
		},
		{
			hotelIndex:  5, // Mandarin Oriental KL
			guestName:   "Ahmad Farhan",
			guestEmail:  "farhan.ahmad@mycorp.com",
			guestPhone:  "+60 12 345 6789",
			checkIn:     tomorrow,
			checkOut:    inThreeDays,
			nights:      2,
			status:      models.BookingConfirmed,
			code:        "HTL-MY-4402",
			settleCode:  "STL-MY-7023",
			settleState: models.SettlementPending,
		},
	}

	for _, mb := range mockBookings {
		if mb.hotelIndex < len(seededHotels) {
			h := seededHotels[mb.hotelIndex]
			var room models.Room
			config.DB.Where("hotel_id = ?", h.ID).First(&room)

			totalAmount := room.Price * float64(mb.nights)
			commission := totalAmount * 0.15
			netPayout := totalAmount - commission

			var checkedInTime *time.Time
			if mb.status == models.BookingCheckedIn {
				t := yesterday.Add(2 * time.Hour)
				checkedInTime = &t
			}

			b := models.Booking{
				BookingCode:     mb.code,
				HotelID:         h.ID,
				RoomID:          room.ID,
				UserID:          &guest.ID,
				GuestName:       mb.guestName,
				GuestEmail:      mb.guestEmail,
				GuestPhone:      mb.guestPhone,
				CheckIn:         mb.checkIn,
				CheckOut:        mb.checkOut,
				Nights:          mb.nights,
				GuestsCount:     2,
				TotalAmount:     totalAmount,
				Currency:        h.Currency,
				Status:          mb.status,
				PaymentMethod:   "credit_card",
				PaymentStatus:   "paid",
				SpecialRequests: "Quiet high-floor room preferred",
				QrCodeData:      fmt.Sprintf("HTL_CHECKIN:%s:%s", mb.code, mb.guestName),
				CheckedInAt:     checkedInTime,
			}
			config.DB.Create(&b)

			s := models.Settlement{
				SettlementCode:     mb.settleCode,
				PartnerID:          h.PartnerID,
				HotelID:            h.ID,
				BookingID:          b.ID,
				GrossAmount:        totalAmount,
				PlatformCommission: commission,
				NetPayout:          netPayout,
				Currency:           h.Currency,
				Status:             mb.settleState,
				BankName:           "Standard Chartered Bank",
				BankAccount:        "992-012-38419",
			}
			config.DB.Create(&s)
		}
	}

	// 4. Partner Applications
	mockApps := []models.PartnerApplication{
		{
			BusinessName:  "Desa Potato Head Seminyak",
			ContactPerson: "Ronald Tan",
			Email:         "partners@potatohead.co",
			Phone:         "+62 361 4737979",
			Country:       "Indonesia",
			City:          "Bali",
			PropertyCount: 2,
			Status:        "pending",
			Notes:         "Sustainable creative village and luxury suites in Seminyak beach.",
		},
		{
			BusinessName:  "The Majestic Hotel Malacca",
			ContactPerson: "Grace Lim",
			Email:         "grace@majestic.com.my",
			Phone:         "+60 6 289 8000",
			Country:       "Malaysia",
			City:          "Malacca",
			PropertyCount: 1,
			Status:        "pending",
			Notes:         "Historic 1920s Straits settlement mansion by the Melaka River.",
		},
	}
	for _, app := range mockApps {
		config.DB.Create(&app)
	}

	log.Println("[INFO] Successfully seeded all 17 luxury properties, rooms, bookings, and applications!")
}

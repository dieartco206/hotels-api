package seeder

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"hotels-api/config"
	"hotels-api/models"

	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	var count int64
	config.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("[INFO] Database already seeded. Skipping initial seed.")
		return
	}

	log.Println("[INFO] Seeding initial database with luxury hotels, partners, and bookings...")

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
		Country:      "SG",
		Avatar:       "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80",
	}

	partnerSG := models.User{
		Name:         "Marina Bay Hospitality Ltd",
		Email:        "partner.sg@marinabay.com",
		PasswordHash: string(hashedPartner),
		Role:         models.RolePartner,
		Phone:        "+65 6688 8868",
		Country:      "SG",
		Avatar:       "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&auto=format&fit=crop&q=80",
	}

	partnerID := models.User{
		Name:         "Nusantara Resorts & Spas",
		Email:        "partner.id@ayana.com",
		PasswordHash: string(hashedPartner),
		Role:         models.RolePartner,
		Phone:        "+62 361 702222",
		Country:      "ID",
		Avatar:       "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150&auto=format&fit=crop&q=80",
	}

	partnerMY := models.User{
		Name:         "Kuala Lumpur Prestige Hotels",
		Email:        "partner.my@kualalumpur.com",
		PasswordHash: string(hashedPartner),
		Role:         models.RolePartner,
		Phone:        "+60 3 2142 8000",
		Country:      "MY",
		Avatar:       "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&auto=format&fit=crop&q=80",
	}

	guest := models.User{
		Name:         "Michael Tan",
		Email:        "guest@example.com",
		PasswordHash: string(hashedGuest),
		Role:         models.RoleGuest,
		Phone:        "+65 9123 4567",
		Country:      "SG",
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

	// 2. Hotels Data
	hotelsData := []struct {
		Hotel models.Hotel
		Rooms []models.Room
	}{
		// --- SINGAPORE ---
		{
			Hotel: models.Hotel{
				Name:        "Marina Bay Sands",
				Tagline:     "Iconic rooftop infinity pool & skyline vistas",
				Type:        "luxury",
				City:        "Singapore",
				Area:        "Marina Bay",
				Country:     "SG",
				Address:     "10 Bayfront Avenue, Marina Bay, Singapore 018956",
				StarRating:  5,
				Rating:      4.9,
				ReviewCount: 2480,
				PriceStart:  650,
				Currency:    "SGD",
				Image:       "https://images.unsplash.com/photo-1525625293386-3f8f99389edd?w=1000&auto=format&fit=crop&q=80",
				Gallery:     toJSON([]string{"https://images.unsplash.com/photo-1525625293386-3f8f99389edd?w=800", "https://images.unsplash.com/photo-1582719508461-905c673771fd?w=800"}),
				Facilities:  toJSON([]string{"Infinity Pool", "SkyPark Observation Deck", "Banyan Tree Spa", "Michelin Dining", "Casino", "Fitness Center"}),
				Description: "Rising above Singapore's vibrant Marina Bay, this premier world-class destination redefines luxury hospitality with exceptional service and breathtaking panoramas.",
				PartnerID:   partnerSG.ID,
				Status:      models.HotelStatusActive,
			},
			Rooms: []models.Room{
				{
					Name:             "Deluxe Harbour View Room",
					Type:             "Deluxe",
					Price:            650,
					OriginalPrice:    780,
					MaxGuests:        2,
					BedType:          "1 King Bed",
					SizeSqm:          42,
					RemainingRooms:   8,
					TotalRooms:       20,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"Floor-to-ceiling Bay View", "Deep Soaking Tub", "Complimentary High-speed Wi-Fi", "Nespresso Machine"}),
					Image:            "https://images.unsplash.com/photo-1582719508461-905c673771fd?w=800",
					Status:           "available",
				},
				{
					Name:             "Sands Premier Suite",
					Type:             "Suite",
					Price:            1250,
					OriginalPrice:    1480,
					MaxGuests:        3,
					BedType:          "1 King Bed + Living Lounge",
					SizeSqm:          75,
					RemainingRooms:   3,
					TotalRooms:       8,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"VIP Club Lounge Access", "Dedicated 24/7 Butler Service", "Complimentary Evening Cocktails", "Hermes Bath Amenities"}),
					Image:            "https://images.unsplash.com/photo-1618773928121-c32242e63f39?w=800",
					Status:           "available",
				},
			},
		},
		{
			Hotel: models.Hotel{
				Name:        "The Fullerton Hotel Singapore",
				Tagline:     "Timeless neoclassical elegance by the Singapore River",
				Type:        "boutique",
				City:        "Singapore",
				Area:        "Downtown Core",
				Country:     "SG",
				Address:     "1 Fullerton Square, Singapore 049178",
				StarRating:  5,
				Rating:      4.8,
				ReviewCount: 1640,
				PriceStart:  480,
				Currency:    "SGD",
				Image:       "https://images.unsplash.com/photo-1566073771259-6a8506099945?w=1000&auto=format&fit=crop&q=80",
				Gallery:     toJSON([]string{"https://images.unsplash.com/photo-1566073771259-6a8506099945?w=800"}),
				Facilities:  toJSON([]string{"Riverfront Pool", "Heritage Tours", "The Fullerton Spa", "Post Bar", "Fine Dining"}),
				Description: "A landmark hotel steeped in rich architectural heritage, offering quintessential luxury right at the heart of the Singapore financial and cultural district.",
				PartnerID:   partnerSG.ID,
				Status:      models.HotelStatusActive,
			},
			Rooms: []models.Room{
				{
					Name:             "Heritage River Room",
					Type:             "Deluxe",
					Price:            480,
					OriginalPrice:    560,
					MaxGuests:        2,
					BedType:          "1 King Bed",
					SizeSqm:          38,
					RemainingRooms:   5,
					TotalRooms:       15,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"River Views", "Marble Bathroom", "Balmain Toiletries", "Smart IPTV"}),
					Image:            "https://images.unsplash.com/photo-1591088398332-8a7791972843?w=800",
					Status:           "available",
				},
			},
		},
		{
			Hotel: models.Hotel{
				Name:        "Capella Singapore Sentosa",
				Tagline:     "Secluded island tranquility and colonial grandeur",
				Type:        "resort",
				City:        "Singapore",
				Area:        "Sentosa Island",
				Country:     "SG",
				Address:     "1 The Knolls, Sentosa Island, Singapore 098297",
				StarRating:  5,
				Rating:      4.9,
				ReviewCount: 930,
				PriceStart:  920,
				Currency:    "SGD",
				Image:       "https://images.unsplash.com/photo-1571896349842-33c89424de2d?w=1000&auto=format&fit=crop&q=80",
				Gallery:     toJSON([]string{"https://images.unsplash.com/photo-1571896349842-33c89424de2d?w=800"}),
				Facilities:  toJSON([]string{"Cascading Tiered Pools", "Auriga Spa", "Private Beach Access", "Curated Culturist Experiences"}),
				Description: "Nestled within 30 acres of lush landscaped grounds on Sentosa Island, Capella Singapore blends old-world serenity with ultra-modern design.",
				PartnerID:   partnerSG.ID,
				Status:      models.HotelStatusActive,
			},
			Rooms: []models.Room{
				{
					Name:             "Premier Garden Villa with Private Pool",
					Type:             "Villa",
					Price:            1450,
					OriginalPrice:    1700,
					MaxGuests:        3,
					BedType:          "1 King Bed",
					SizeSqm:          133,
					RemainingRooms:   2,
					TotalRooms:       6,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"Private Plunge Pool", "Outdoor Rainforest Shower", "Personal Culturist Host", "Aesop Amenities"}),
					Image:            "https://images.unsplash.com/photo-1540541338287-41700207dee6?w=800",
					Status:           "available",
				},
			},
		},

		// --- INDONESIA ---
		{
			Hotel: models.Hotel{
				Name:        "AYANA Resort Bali",
				Tagline:     "Majestic cliff-top sanctuary with sunset Rock Bar",
				Type:        "resort",
				City:        "Bali",
				Area:        "Jimbaran Bay",
				Country:     "ID",
				Address:     "Jl. Karang Mas Sejahtera, Jimbaran, Bali 80364",
				StarRating:  5,
				Rating:      4.9,
				ReviewCount: 3120,
				PriceStart:  3500000,
				Currency:    "IDR",
				Image:       "https://images.unsplash.com/photo-1540541338287-41700207dee6?w=1000&auto=format&fit=crop&q=80",
				Gallery:     toJSON([]string{"https://images.unsplash.com/photo-1540541338287-41700207dee6?w=800"}),
				Facilities:  toJSON([]string{"World-famous Rock Bar", "14 Swimming Pools", "Aquatonic Seawater Pool", "Kubu Private Beach", "Spa on the Rocks"}),
				Description: "Perched above Jimbaran Bay on 90 hectares of cliff-front tropical gardens, AYANA is an award-winning integrated luxury resort with iconic ocean sunsets.",
				PartnerID:   partnerID.ID,
				Status:      models.HotelStatusActive,
			},
			Rooms: []models.Room{
				{
					Name:             "Deluxe Ocean View Room",
					Type:             "Deluxe",
					Price:            3500000,
					OriginalPrice:    4200000,
					MaxGuests:        2,
					BedType:          "1 King Bed",
					SizeSqm:          48,
					RemainingRooms:   12,
					TotalRooms:       30,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"Indian Ocean View", "Private Balcony Daybed", "Marble Bathtub", "Sensatia Botanicals Amenities"}),
					Image:            "https://images.unsplash.com/photo-1566665797739-1674de7a421a?w=800",
					Status:           "available",
				},
				{
					Name:             "Cliff Ocean Pool Villa",
					Type:             "Villa",
					Price:            9800000,
					OriginalPrice:    11500000,
					MaxGuests:        3,
					BedType:          "1 King Bed",
					SizeSqm:          298,
					RemainingRooms:   4,
					TotalRooms:       10,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"Private Infinity Edge Pool", "Dedicated 24-hr Butler", "VIP Access to Rock Bar", "Balinese Garden Gazebo"}),
					Image:            "https://images.unsplash.com/photo-1582719508461-905c673771fd?w=800",
					Status:           "available",
				},
			},
		},
		{
			Hotel: models.Hotel{
				Name:        "The Langham Jakarta",
				Tagline:     "Ultra-luxury British craftsmanship in the SCBD heart",
				Type:        "luxury",
				City:        "Jakarta",
				Area:        "SCBD Senayan",
				Country:     "ID",
				Address:     "District 8, SCBD Lot 28, Senayan, Jakarta 12190",
				StarRating:  5,
				Rating:      4.9,
				ReviewCount: 1450,
				PriceStart:  3800000,
				Currency:    "IDR",
				Image:       "https://images.unsplash.com/photo-1551882547-ff40c63fe5fa?w=1000&auto=format&fit=crop&q=80",
				Gallery:     toJSON([]string{"https://images.unsplash.com/photo-1551882547-ff40c63fe5fa?w=800"}),
				Facilities:  toJSON([]string{"Infinity Sky Pool", "Chuan Spa", "Tom's by Tom Aikens", "Langham Club Lounge", "Helipad"}),
				Description: "Soaring majestically in Jakarta's prime commercial precinct, The Langham Jakarta is an oasis of refined elegance, gourmet culinary artistry, and bespoke pampering.",
				PartnerID:   partnerID.ID,
				Status:      models.HotelStatusActive,
			},
			Rooms: []models.Room{
				{
					Name:             "Deluxe Skyline King",
					Type:             "Deluxe",
					Price:            3800000,
					OriginalPrice:    4500000,
					MaxGuests:        2,
					BedType:          "1 King Bed",
					SizeSqm:          42,
					RemainingRooms:   8,
					TotalRooms:       25,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"Floor-to-ceiling Glass Views", "French Marble Bathroom", "Diptyque Bath Amenities", "Cocktail Bar Cabinet"}),
					Image:            "https://images.unsplash.com/photo-1590490360182-c33d57733427?w=800",
					Status:           "available",
				},
			},
		},
		{
			Hotel: models.Hotel{
				Name:        "Amanjiwo Borobudur",
				Tagline:     "Peaceful soul overlooking mystical Borobudur temple",
				Type:        "resort",
				City:        "Yogyakarta",
				Area:        "Magelang / Borobudur",
				Country:     "ID",
				Address:     "Ds. Majaksingi, Borobudur, Central Java 56553",
				StarRating:  5,
				Rating:      5.0,
				ReviewCount: 680,
				PriceStart:  14500000,
				Currency:    "IDR",
				Image:       "https://images.unsplash.com/photo-1571003123894-1f0594d2b5d9?w=1000&auto=format&fit=crop&q=80",
				Gallery:     toJSON([]string{"https://images.unsplash.com/photo-1571003123894-1f0594d2b5d9?w=800"}),
				Facilities:  toJSON([]string{"Borobudur Private Sunrise Tour", "40m Infinity Pool", "Ancient Javanese Healer Spa", "Private Dining Rotunda"}),
				Description: "Framed by the Menoreh Hills and facing UNESCO World Heritage Borobudur, Amanjiwo is built of local limestone and reflects 9th-century Javanese temple architecture.",
				PartnerID:   partnerID.ID,
				Status:      models.HotelStatusActive,
			},
			Rooms: []models.Room{
				{
					Name:             "Borobudur Pool Suite",
					Type:             "Suite",
					Price:            18500000,
					OriginalPrice:    22000000,
					MaxGuests:        2,
					BedType:          "1 King Bed",
					SizeSqm:          215,
					RemainingRooms:   2,
					TotalRooms:       4,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"Direct Borobudur Temple View", "Private Plunge Pool", "Outdoor Sunken Terrazzo Tub", "Artisan Welcome Basket"}),
					Image:            "https://images.unsplash.com/photo-1618773928121-c32242e63f39?w=800",
					Status:           "available",
				},
			},
		},

		// --- MALAYSIA ---
		{
			Hotel: models.Hotel{
				Name:        "The Ritz-Carlton Kuala Lumpur",
				Tagline:     "Legendary five-star hospitality in Bukit Bintang",
				Type:        "luxury",
				City:        "Kuala Lumpur",
				Area:        "Bukit Bintang / Golden Triangle",
				Country:     "MY",
				Address:     "168 Jalan Imbi, Bukit Bintang, Kuala Lumpur 55100",
				StarRating:  5,
				Rating:      4.9,
				ReviewCount: 1890,
				PriceStart:  820,
				Currency:    "MYR",
				Image:       "https://images.unsplash.com/photo-1566665797739-1674de7a421a?w=1000&auto=format&fit=crop&q=80",
				Gallery:     toJSON([]string{"https://images.unsplash.com/photo-1566665797739-1674de7a421a?w=800"}),
				Facilities:  toJSON([]string{"Full Butler Service in All Rooms", "Spa Village Kuala Lumpur", "The Cobalt Room", "Outdoor Lap Pool"}),
				Description: "Set in the heart of Kuala Lumpur's glamorous shopping hub, The Ritz-Carlton delivers prestigious timeless luxury, lavish suites, and bespoke guest care.",
				PartnerID:   partnerMY.ID,
				Status:      models.HotelStatusActive,
			},
			Rooms: []models.Room{
				{
					Name:             "Executive King Room with Butler Service",
					Type:             "Executive",
					Price:            820,
					OriginalPrice:    980,
					MaxGuests:        2,
					BedType:          "1 King Bed",
					SizeSqm:          46,
					RemainingRooms:   9,
					TotalRooms:       25,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"24-hour Dedicated Butler", "Full Marble Bath", "Evening Turndown", "High-speed Fiber Wi-Fi"}),
					Image:            "https://images.unsplash.com/photo-1590490360182-c33d57733427?w=800",
					Status:           "available",
				},
				{
					Name:             "Premier Suite with City Skyline",
					Type:             "Suite",
					Price:            1450,
					OriginalPrice:    1750,
					MaxGuests:        3,
					BedType:          "1 King Bed + Lounge",
					SizeSqm:          82,
					RemainingRooms:   4,
					TotalRooms:       10,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"Separate Living Room", "Club Lounge Access", "Cocktails & Canapes Daily", "Asprey Amenities"}),
					Image:            "https://images.unsplash.com/photo-1582719508461-905c673771fd?w=800",
					Status:           "available",
				},
			},
		},
		{
			Hotel: models.Hotel{
				Name:        "The Datai Langkawi",
				Tagline:     "Ancient 10-million-year rainforest by the Andaman Sea",
				Type:        "resort",
				City:        "Langkawi",
				Area:        "Teluk Datai",
				Country:     "MY",
				Address:     "Jln Teluk Datai, Teluk Datai, Langkawi 07000",
				StarRating:  5,
				Rating:      5.0,
				ReviewCount: 1120,
				PriceStart:  2100,
				Currency:    "MYR",
				Image:       "https://images.unsplash.com/photo-1520250497591-112f2f40a3f4?w=1000&auto=format&fit=crop&q=80",
				Gallery:     toJSON([]string{"https://images.unsplash.com/photo-1520250497591-112f2f40a3f4?w=800"}),
				Facilities:  toJSON([]string{"Private White Sandy Bay", "Resident Naturalist Tours", "The Datai Spa", "Beach Bar", "Kayaking"}),
				Description: "Enfolded in an ancient pristine virgin rainforest and sweeping down to a calm turquoise bay, The Datai Langkawi is universally celebrated as a masterpiece of eco-luxury.",
				PartnerID:   partnerMY.ID,
				Status:      models.HotelStatusActive,
			},
			Rooms: []models.Room{
				{
					Name:             "Rainforest Villa with Plunge Pool",
					Type:             "Villa",
					Price:            3400,
					OriginalPrice:    3900,
					MaxGuests:        2,
					BedType:          "1 King Bed",
					SizeSqm:          123,
					RemainingRooms:   3,
					TotalRooms:       8,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"Private Plunge Pool on Elevated Deck", "Rainforest Canopy Views", "Complimentary Mini Bar", "Deep Forest Quietude"}),
					Image:            "https://images.unsplash.com/photo-1540541338287-41700207dee6?w=800",
					Status:           "available",
				},
			},
		},
		{
			Hotel: models.Hotel{
				Name:        "Eastern & Oriental Hotel Penang",
				Tagline:     "Colonial maritime majesty along George Town seafront",
				Type:        "boutique",
				City:        "Penang",
				Area:        "George Town",
				Country:     "MY",
				Address:     "10 Lebuh Farquhar, George Town, Penang 10200",
				StarRating:  5,
				Rating:      4.8,
				ReviewCount: 1340,
				PriceStart:  750,
				Currency:    "MYR",
				Image:       "https://images.unsplash.com/photo-1571896349842-33c89424de2d?w=1000&auto=format&fit=crop&q=80",
				Gallery:     toJSON([]string{"https://images.unsplash.com/photo-1571896349842-33c89424de2d?w=800"}),
				Facilities:  toJSON([]string{"Seafront Infinity Pool", "Heritage High Tea", "The Farquhar Bar", "UNESCO Heritage Walking Distance"}),
				Description: "Established in 1885 by the Sarkies brothers, the E&O Hotel Penang is a legendary monument of imperial elegance with picturesque views across the Malacca Strait.",
				PartnerID:   partnerMY.ID,
				Status:      models.HotelStatusActive,
			},
			Rooms: []models.Room{
				{
					Name:             "Heritage Wing Strait Suite",
					Type:             "Suite",
					Price:            750,
					OriginalPrice:    890,
					MaxGuests:        2,
					BedType:          "1 King Bed",
					SizeSqm:          58,
					RemainingRooms:   6,
					TotalRooms:       18,
					FreeBreakfast:    true,
					FreeCancellation: true,
					Amenities:        toJSON([]string{"Seafront View Balcony", "Clawfoot Bathtub", "Decadent Afternoon Tea Included", "Evening Cocktails"}),
					Image:            "https://images.unsplash.com/photo-1591088398332-8a7791972843?w=800",
					Status:           "available",
				},
			},
		},
	}

	var seededHotels []models.Hotel
	for _, item := range hotelsData {
		h := item.Hotel
		config.DB.Create(&h)
		for _, r := range item.Rooms {
			r.HotelID = h.ID
			config.DB.Create(&r)
		}
		seededHotels = append(seededHotels, h)
	}

	// 3. Seed Initial Realistic Bookings & Settlements
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
			hotelIndex:  3, // AYANA Bali
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
			hotelIndex:  6, // Ritz Carlton KL
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
				SpecialRequests: "Quiet floor preferred, late check-in",
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
			Country:       "ID",
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
			Country:       "MY",
			City:          "Malacca",
			PropertyCount: 1,
			Status:        "pending",
			Notes:         "Historic 1920s Straits settlement mansion by the Melaka River.",
		},
	}
	for _, app := range mockApps {
		config.DB.Create(&app)
	}

	log.Println("[INFO] Initial database seed successfully completed!")
}

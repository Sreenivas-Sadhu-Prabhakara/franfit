// Package catalog holds the embedded, illustrative directory of Indian
// franchise brands. Every brand here is fictional — invented for the demo —
// and the UI labels the directory as illustrative data.
package catalog

// Brand describes one franchise opportunity. Money fields are in ₹ lakhs.
type Brand struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Category          string   `json:"category"`
	InvestmentMinL    float64  `json:"investmentMinL"`
	InvestmentMaxL    float64  `json:"investmentMaxL"`
	FranchiseFeeL     float64  `json:"franchiseFeeL"`
	AreaSqftMin       int      `json:"areaSqftMin"`
	ModelsSupported   []string `json:"modelsSupported"`
	RoyaltyPct        float64  `json:"royaltyPct"`
	PaybackMonthsEst  int      `json:"paybackMonthsEst"`
	MonthlyProfitEstL float64  `json:"monthlyProfitEstL"`
	CityTiers         []int    `json:"cityTiers"`
	BrandStory        string   `json:"brandStory"`
}

// Categories lists every category present in the directory, in display order.
var Categories = []string{
	"QSR", "Cafe", "Salon", "Pharmacy", "Education", "Grocery", "Fitness", "Courier",
}

// Models lists the four franchise operating models in canonical order.
var Models = []string{"FOFO", "FOCO", "COCO", "FICO"}

// All returns the full illustrative brand directory.
func All() []Brand { return brands }

// ByID returns the brand with the given id, if present.
func ByID(id string) (Brand, bool) {
	for _, b := range brands {
		if b.ID == id {
			return b, true
		}
	}
	return Brand{}, false
}

var brands = []Brand{
	// ---- QSR ----
	{
		ID: "tandoor-tales", Name: "Tandoor Tales", Category: "QSR",
		InvestmentMinL: 22, InvestmentMaxL: 38, FranchiseFeeL: 5, AreaSqftMin: 350,
		ModelsSupported: []string{"FOFO", "FOCO"}, RoyaltyPct: 6, PaybackMonthsEst: 22,
		MonthlyProfitEstL: 1.6, CityTiers: []int{1, 2},
		BrandStory: "Started as a single kebab counter in Lucknow's Hazratganj, Tandoor Tales serves rolls and tikka platters built for the office-lunch crowd. Its 40-item central-kitchen menu keeps outlet kitchens down to two staff.",
	},
	{
		ID: "vadapav-junction", Name: "VadaPav Junction", Category: "QSR",
		InvestmentMinL: 8, InvestmentMaxL: 15, FranchiseFeeL: 2, AreaSqftMin: 120,
		ModelsSupported: []string{"FOFO", "FICO"}, RoyaltyPct: 4, PaybackMonthsEst: 14,
		MonthlyProfitEstL: 0.9, CityTiers: []int{1, 2, 3},
		BrandStory: "A Mumbai-style vada pav kiosk format that fits railway-station forecourts and college gates alike. Standardised chutney packs shipped weekly mean any operator can hit the same taste from day one.",
	},
	{
		ID: "biryani-bandhu", Name: "Biryani Bandhu", Category: "QSR",
		InvestmentMinL: 30, InvestmentMaxL: 55, FranchiseFeeL: 8, AreaSqftMin: 500,
		ModelsSupported: []string{"FOCO", "COCO"}, RoyaltyPct: 8, PaybackMonthsEst: 30,
		MonthlyProfitEstL: 2.8, CityTiers: []int{1, 2},
		BrandStory: "A delivery-first Hyderabadi biryani kitchen with dine-in counters in metro high streets. The brand runs operations end to end, so investors track dashboards, not dum handis.",
	},
	{
		ID: "momo-metro", Name: "Momo Metro", Category: "QSR",
		InvestmentMinL: 10, InvestmentMaxL: 18, FranchiseFeeL: 3, AreaSqftMin: 150,
		ModelsSupported: []string{"FOFO", "FICO", "FOCO"}, RoyaltyPct: 5, PaybackMonthsEst: 16,
		MonthlyProfitEstL: 1.1, CityTiers: []int{1, 2, 3},
		BrandStory: "Darjeeling-style steamed and pan-fried momos out of a six-by-ten kiosk. Momo Metro grew from Siliguri to 90 towns by keeping the menu to twelve SKUs and the equipment to one steamer line.",
	},
	{
		ID: "dosa-dhamaka", Name: "Dosa Dhamaka", Category: "QSR",
		InvestmentMinL: 18, InvestmentMaxL: 32, FranchiseFeeL: 4.5, AreaSqftMin: 300,
		ModelsSupported: []string{"FOFO", "FOCO"}, RoyaltyPct: 6, PaybackMonthsEst: 20,
		MonthlyProfitEstL: 1.5, CityTiers: []int{1, 2, 3},
		BrandStory: "A quick-service dosa house from Coimbatore where every dosa hits the table in under four minutes. Batter arrives fermented from regional hubs, cutting kitchen skill requirements in half.",
	},

	// ---- Cafe ----
	{
		ID: "chai-chaupal", Name: "Chai Chaupal", Category: "Cafe",
		InvestmentMinL: 9, InvestmentMaxL: 16, FranchiseFeeL: 2.5, AreaSqftMin: 200,
		ModelsSupported: []string{"FOFO", "FICO"}, RoyaltyPct: 5, PaybackMonthsEst: 15,
		MonthlyProfitEstL: 0.85, CityTiers: []int{1, 2, 3},
		BrandStory: "Kulhad chai, bun-maska and board games — a village-square feel for city corners. Chai Chaupal outlets become the adda for offices nearby, with evening footfall double the morning's.",
	},
	{
		ID: "kulhad-coffee-co", Name: "Kulhad Coffee Co.", Category: "Cafe",
		InvestmentMinL: 14, InvestmentMaxL: 26, FranchiseFeeL: 4, AreaSqftMin: 300,
		ModelsSupported: []string{"FOFO", "FOCO", "FICO"}, RoyaltyPct: 6, PaybackMonthsEst: 19,
		MonthlyProfitEstL: 1.2, CityTiers: []int{1, 2},
		BrandStory: "South Indian filter coffee served in terracotta kulhads with a modern espresso menu alongside. Beans are roasted in Chikkamagaluru and shipped fresh every fortnight to each cafe.",
	},
	{
		ID: "filter-kaapi-house", Name: "Filter Kaapi House", Category: "Cafe",
		InvestmentMinL: 20, InvestmentMaxL: 35, FranchiseFeeL: 6, AreaSqftMin: 450,
		ModelsSupported: []string{"FOCO", "COCO"}, RoyaltyPct: 7, PaybackMonthsEst: 26,
		MonthlyProfitEstL: 1.9, CityTiers: []int{1},
		BrandStory: "A sit-down kaapi and tiffin house recreating the Mylapore mess experience for metro neighbourhoods. The brand operates every outlet itself so the davara-tumbler ritual never gets diluted.",
	},
	{
		ID: "bakers-gully", Name: "Baker's Gully", Category: "Cafe",
		InvestmentMinL: 12, InvestmentMaxL: 22, FranchiseFeeL: 3.5, AreaSqftMin: 250,
		ModelsSupported: []string{"FOFO", "FICO"}, RoyaltyPct: 5, PaybackMonthsEst: 18,
		MonthlyProfitEstL: 1.0, CityTiers: []int{1, 2, 3},
		BrandStory: "An Irani-bakery revival — mawa cakes, khari and fresh pav from a live oven in the front window. Baker's Gully thrives beside schools and bus stands where fresh-baked aroma does the marketing.",
	},

	// ---- Salon ----
	{
		ID: "roopvati-salon", Name: "Roopvati Salon", Category: "Salon",
		InvestmentMinL: 15, InvestmentMaxL: 28, FranchiseFeeL: 4, AreaSqftMin: 400,
		ModelsSupported: []string{"FOFO", "FOCO"}, RoyaltyPct: 8, PaybackMonthsEst: 24,
		MonthlyProfitEstL: 1.4, CityTiers: []int{1, 2, 3},
		BrandStory: "A women's salon chain from Indore built around bridal packages and membership plans. Roopvati trains every stylist at its own academy and posts them to franchise outlets.",
	},
	{
		ID: "urbantrim-studio", Name: "UrbanTrim Studio", Category: "Salon",
		InvestmentMinL: 10, InvestmentMaxL: 18, FranchiseFeeL: 3, AreaSqftMin: 250,
		ModelsSupported: []string{"FOFO", "FICO"}, RoyaltyPct: 6, PaybackMonthsEst: 18,
		MonthlyProfitEstL: 1.0, CityTiers: []int{1, 2},
		BrandStory: "A men's grooming studio with app-based queue booking, so no one waits on the bench. UrbanTrim's subscription haircut plan keeps chairs 80% booked on weekdays.",
	},
	{
		ID: "mehndi-and-mane", Name: "Mehndi & Mane", Category: "Salon",
		InvestmentMinL: 8, InvestmentMaxL: 14, FranchiseFeeL: 2, AreaSqftMin: 200,
		ModelsSupported: []string{"FOFO", "FICO"}, RoyaltyPct: 5, PaybackMonthsEst: 15,
		MonthlyProfitEstL: 0.8, CityTiers: []int{2, 3},
		BrandStory: "A festival-season powerhouse combining mehndi artistry with everyday hair services for tier-2 and tier-3 towns. Wedding-season bookings alone recover a third of the year's rent.",
	},
	{
		ID: "glowkatha-lounge", Name: "GlowKatha Beauty Lounge", Category: "Salon",
		InvestmentMinL: 25, InvestmentMaxL: 45, FranchiseFeeL: 7, AreaSqftMin: 600,
		ModelsSupported: []string{"FOCO", "COCO"}, RoyaltyPct: 9, PaybackMonthsEst: 32,
		MonthlyProfitEstL: 2.2, CityTiers: []int{1},
		BrandStory: "A premium skin-and-hair lounge with dermatologist tie-ups and machine-led facials. GlowKatha runs outlets with its own therapists, positioning investors as silent capital partners.",
	},

	// ---- Pharmacy ----
	{
		ID: "aarogyamed-pharmacy", Name: "AarogyaMed Pharmacy", Category: "Pharmacy",
		InvestmentMinL: 12, InvestmentMaxL: 20, FranchiseFeeL: 3, AreaSqftMin: 250,
		ModelsSupported: []string{"FOFO", "FICO"}, RoyaltyPct: 3, PaybackMonthsEst: 20,
		MonthlyProfitEstL: 1.0, CityTiers: []int{1, 2, 3},
		BrandStory: "A neighbourhood pharmacy format with generic-medicine counters and free BP checks that build daily footfall. AarogyaMed's supply platform keeps 96% of prescriptions fillable on first visit.",
	},
	{
		ID: "sehatplus-chemists", Name: "SehatPlus Chemists", Category: "Pharmacy",
		InvestmentMinL: 18, InvestmentMaxL: 30, FranchiseFeeL: 5, AreaSqftMin: 350,
		ModelsSupported: []string{"FOFO", "FOCO"}, RoyaltyPct: 4, PaybackMonthsEst: 24,
		MonthlyProfitEstL: 1.5, CityTiers: []int{1, 2},
		BrandStory: "A 24x7 chemist chain anchored near hospitals, with a wellness aisle that lifts average bill values by 30%. Night-shift staffing is handled through the brand's own pharmacist pool.",
	},
	{
		ID: "nirogya-pharmacy", Name: "Nirogya Pharmacy", Category: "Pharmacy",
		InvestmentMinL: 8, InvestmentMaxL: 14, FranchiseFeeL: 2, AreaSqftMin: 180,
		ModelsSupported: []string{"FOFO", "FICO"}, RoyaltyPct: 3, PaybackMonthsEst: 17,
		MonthlyProfitEstL: 0.75, CityTiers: []int{2, 3},
		BrandStory: "A compact pharmacy built for tehsil towns, pairing medicines with ayurvedic staples people already ask for. Nirogya's hub-and-spoke supply van restocks every outlet twice a week.",
	},

	// ---- Education ----
	{
		ID: "gyansetu-coaching", Name: "GyanSetu Coaching", Category: "Education",
		InvestmentMinL: 10, InvestmentMaxL: 20, FranchiseFeeL: 4, AreaSqftMin: 600,
		ModelsSupported: []string{"FOFO", "FOCO"}, RoyaltyPct: 15, PaybackMonthsEst: 18,
		MonthlyProfitEstL: 1.3, CityTiers: []int{1, 2, 3},
		BrandStory: "Board-exam and foundation coaching for classes 8–12, with recorded lectures from Kota faculty and local doubt-solvers in every centre. Parents get weekly progress reports on WhatsApp.",
	},
	{
		ID: "vidyavriksh-learning", Name: "VidyaVriksh Learning", Category: "Education",
		InvestmentMinL: 6, InvestmentMaxL: 12, FranchiseFeeL: 2.5, AreaSqftMin: 400,
		ModelsSupported: []string{"FOFO", "FICO"}, RoyaltyPct: 12, PaybackMonthsEst: 15,
		MonthlyProfitEstL: 0.9, CityTiers: []int{2, 3},
		BrandStory: "An after-school learning centre for tier-2 and tier-3 towns covering spoken English, maths and computer basics. Its Hindi-first pedagogy is the reason retention crosses 85% year on year.",
	},
	{
		ID: "abacusank-academy", Name: "AbacusAnk Academy", Category: "Education",
		InvestmentMinL: 4, InvestmentMaxL: 8, FranchiseFeeL: 1.5, AreaSqftMin: 300,
		ModelsSupported: []string{"FOFO", "FICO"}, RoyaltyPct: 10, PaybackMonthsEst: 12,
		MonthlyProfitEstL: 0.6, CityTiers: []int{1, 2, 3},
		BrandStory: "Abacus and Vedic-maths classes for ages 5–14, run in two-hour weekend batches. AbacusAnk is the lowest-capital brand in most portfolios and often runs from a converted living room.",
	},
	{
		ID: "codechotu-kids", Name: "CodeChotu Kids Coding", Category: "Education",
		InvestmentMinL: 8, InvestmentMaxL: 16, FranchiseFeeL: 3, AreaSqftMin: 350,
		ModelsSupported: []string{"FOFO", "FOCO", "FICO"}, RoyaltyPct: 14, PaybackMonthsEst: 16,
		MonthlyProfitEstL: 1.0, CityTiers: []int{1, 2},
		BrandStory: "Scratch-to-Python coding labs for school kids, taught on the centre's own machines so no parent needs to buy a laptop. Hackathon Saturdays double as the brand's best enrolment funnel.",
	},

	// ---- Grocery ----
	{
		ID: "kirana-junction", Name: "Kirana Junction", Category: "Grocery",
		InvestmentMinL: 25, InvestmentMaxL: 45, FranchiseFeeL: 5, AreaSqftMin: 800,
		ModelsSupported: []string{"FOFO", "FOCO"}, RoyaltyPct: 2, PaybackMonthsEst: 28,
		MonthlyProfitEstL: 2.0, CityTiers: []int{1, 2, 3},
		BrandStory: "A modern-trade makeover for the neighbourhood kirana: barcode billing, loyalty points and a fresh-atta chakki at the back. Kirana Junction converts existing shopkeepers as often as new investors.",
	},
	{
		ID: "sabzisaathi-mart", Name: "SabziSaathi Fresh Mart", Category: "Grocery",
		InvestmentMinL: 15, InvestmentMaxL: 28, FranchiseFeeL: 4, AreaSqftMin: 500,
		ModelsSupported: []string{"FOCO", "COCO"}, RoyaltyPct: 3, PaybackMonthsEst: 26,
		MonthlyProfitEstL: 1.6, CityTiers: []int{1, 2},
		BrandStory: "Farm-to-shelf fruits and vegetables with 4 a.m. mandi procurement handled centrally. The brand operates the cold chain and staffing; partners bring the storefront and capital.",
	},
	{
		ID: "dailydibba-grocery", Name: "DailyDibba Grocery", Category: "Grocery",
		InvestmentMinL: 35, InvestmentMaxL: 60, FranchiseFeeL: 8, AreaSqftMin: 1200,
		ModelsSupported: []string{"FOCO", "COCO"}, RoyaltyPct: 2.5, PaybackMonthsEst: 34,
		MonthlyProfitEstL: 3.2, CityTiers: []int{1},
		BrandStory: "A mini-supermarket format with 4,000 SKUs and 10-minute delivery from the store's own dark corner. DailyDibba runs company-operated stores where investors earn a fixed-plus-share return.",
	},

	// ---- Fitness ----
	{
		ID: "fitbharat-gym", Name: "FitBharat Gym", Category: "Fitness",
		InvestmentMinL: 50, InvestmentMaxL: 90, FranchiseFeeL: 10, AreaSqftMin: 2500,
		ModelsSupported: []string{"FOFO", "FOCO"}, RoyaltyPct: 7, PaybackMonthsEst: 36,
		MonthlyProfitEstL: 3.5, CityTiers: []int{1, 2},
		BrandStory: "A full-floor strength and cardio gym priced for the middle class, with ₹999 monthly memberships that fill 1,500-member rosters. Equipment financing through the brand cuts upfront capital by a fifth.",
	},
	{
		ID: "yogvan-studio", Name: "YogVan Studio", Category: "Fitness",
		InvestmentMinL: 8, InvestmentMaxL: 15, FranchiseFeeL: 2.5, AreaSqftMin: 700,
		ModelsSupported: []string{"FOFO", "FICO"}, RoyaltyPct: 8, PaybackMonthsEst: 16,
		MonthlyProfitEstL: 0.9, CityTiers: []int{1, 2, 3},
		BrandStory: "Sunrise and sunset yoga batches in a mat-and-mirror studio that needs no machines at all. YogVan certifies local instructors, making it one of the fastest formats to open.",
	},
	{
		ID: "akhada45-fitness", Name: "Akhada45 Functional Fitness", Category: "Fitness",
		InvestmentMinL: 30, InvestmentMaxL: 55, FranchiseFeeL: 7, AreaSqftMin: 1500,
		ModelsSupported: []string{"FOCO", "COCO"}, RoyaltyPct: 9, PaybackMonthsEst: 30,
		MonthlyProfitEstL: 2.4, CityTiers: []int{1},
		BrandStory: "Group functional-training classes on a 45-minute clock, drawing on desi akhada conditioning with kettlebells and gada work. Coaches are hired, trained and rotated by the brand.",
	},
	{
		ID: "zumfit-dance", Name: "ZumFit Dance Fitness", Category: "Fitness",
		InvestmentMinL: 12, InvestmentMaxL: 22, FranchiseFeeL: 3.5, AreaSqftMin: 900,
		ModelsSupported: []string{"FOFO", "FOCO", "FICO"}, RoyaltyPct: 8, PaybackMonthsEst: 20,
		MonthlyProfitEstL: 1.2, CityTiers: []int{1, 2},
		BrandStory: "Bollywood-cardio classes where the evening ladies' batch is routinely waitlisted. ZumFit's playlist-and-choreo drops arrive monthly so every studio dances the same steps.",
	},

	// ---- Courier ----
	{
		ID: "swiftdak-couriers", Name: "SwiftDak Couriers", Category: "Courier",
		InvestmentMinL: 5, InvestmentMaxL: 10, FranchiseFeeL: 1.5, AreaSqftMin: 150,
		ModelsSupported: []string{"FOFO", "FICO"}, RoyaltyPct: 10, PaybackMonthsEst: 12,
		MonthlyProfitEstL: 0.65, CityTiers: []int{1, 2, 3},
		BrandStory: "A pin-code-level courier booking counter riding the e-commerce returns wave. SwiftDak's franchisees earn on every parcel booked, picked up or held for customer collection.",
	},
	{
		ID: "gaonexpress-logistics", Name: "GaonExpress Logistics", Category: "Courier",
		InvestmentMinL: 10, InvestmentMaxL: 18, FranchiseFeeL: 3, AreaSqftMin: 400,
		ModelsSupported: []string{"FOFO", "FOCO"}, RoyaltyPct: 8, PaybackMonthsEst: 18,
		MonthlyProfitEstL: 1.1, CityTiers: []int{2, 3},
		BrandStory: "Last-mile delivery hubs for villages and tehsil towns that national carriers skip. A GaonExpress hub aggregates parcels for 30–40 surrounding pin codes and runs two delivery bikes.",
	},
	{
		ID: "citycart-deliveries", Name: "CityCart Deliveries", Category: "Courier",
		InvestmentMinL: 20, InvestmentMaxL: 35, FranchiseFeeL: 5, AreaSqftMin: 600,
		ModelsSupported: []string{"FOCO", "COCO"}, RoyaltyPct: 6, PaybackMonthsEst: 26,
		MonthlyProfitEstL: 1.8, CityTiers: []int{1},
		BrandStory: "Intra-city micro-warehouses powering same-day delivery for local sellers. CityCart staffs and runs each node; partners fund the racking, e-bikes and deposit.",
	},
}

package fitscore

import (
	"reflect"
	"testing"

	"franfit/internal/catalog"
)

func typicalInput() Input {
	return Input{
		BudgetL:     30,
		City:        "Jaipur",
		Tier:        2,
		Involvement: OwnerOperator,
		Risk:        3,
		Categories:  []string{"QSR", "Cafe"},
		SpaceSqft:   500,
	}
}

func TestScoringIsDeterministic(t *testing.T) {
	a := Rank(catalog.All(), typicalInput())
	b := Rank(catalog.All(), typicalInput())
	if !reflect.DeepEqual(a, b) {
		t.Fatal("same input produced different rankings")
	}
	if len(a.Matches) == 0 {
		t.Fatal("expected matches for a ₹30 L budget")
	}
	for i := 1; i < len(a.Matches); i++ {
		if a.Matches[i].Score > a.Matches[i-1].Score {
			t.Fatalf("matches not sorted: %d before %d", a.Matches[i-1].Score, a.Matches[i].Score)
		}
	}
}

func TestWeightBoundaries(t *testing.T) {
	// A brand engineered to fit the typical input perfectly must score 100.
	perfect := catalog.Brand{
		ID: "perfect", Name: "Perfect Brand", Category: "QSR",
		InvestmentMinL: 10, InvestmentMaxL: 20, AreaSqftMin: 100,
		ModelsSupported: []string{"FOFO"}, PaybackMonthsEst: 30, // horizon 0.5 == appetite (risk 3)
		CityTiers: []int{2},
	}
	res := Rank([]catalog.Brand{perfect}, typicalInput())
	if len(res.Matches) != 1 || res.Matches[0].Score != 100 {
		t.Fatalf("perfect brand should score 100, got %+v", res.Matches)
	}
	// Every factor for every real brand must sit within [0, weight] and the
	// factor weights must total exactly 100.
	if WeightCapital+WeightModel+WeightCategory+WeightTier+WeightRisk != 100 {
		t.Fatal("factor weights must sum to 100")
	}
	all := Rank(catalog.All(), typicalInput())
	for _, m := range all.Matches {
		var sum float64
		for _, f := range m.Factors {
			if f.Points < 0 || f.Points > f.Max {
				t.Fatalf("%s factor %q out of range: %v (max %v)", m.Brand.ID, f.Label, f.Points, f.Max)
			}
			sum += f.Points
		}
		if m.Score < 0 || m.Score > 100 {
			t.Fatalf("%s total score out of range: %d", m.Brand.ID, m.Score)
		}
		if diff := float64(m.Score) - sum; diff > 0.5 || diff < -0.5 {
			t.Fatalf("%s score %d does not match factor sum %v", m.Brand.ID, m.Score, sum)
		}
	}
}

func TestBudgetBelowAllBrandsExplains(t *testing.T) {
	in := typicalInput()
	in.BudgetL = 1 // below every brand's minimum investment
	res := Rank(catalog.All(), in)
	if len(res.Matches) != 0 {
		t.Fatalf("expected zero matches, got %d", len(res.Matches))
	}
	if !res.NoMatches {
		t.Fatal("NoMatches flag should be set")
	}
	if res.Explanation == "" {
		t.Fatal("explanation should tell the user why nothing matched")
	}
}

func TestUnaffordableBrandIsExcluded(t *testing.T) {
	in := typicalInput()
	in.BudgetL = 12
	res := Rank(catalog.All(), in)
	for _, m := range res.Matches {
		if m.Brand.InvestmentMinL > in.BudgetL {
			t.Fatalf("brand %s (min %v L) should be excluded at budget %v L",
				m.Brand.ID, m.Brand.InvestmentMinL, in.BudgetL)
		}
	}
}

func TestModelRecommendationPerInvolvement(t *testing.T) {
	allModels := catalog.Brand{
		ID: "omni", Name: "Omni", Category: "QSR",
		InvestmentMinL: 10, InvestmentMaxL: 20, AreaSqftMin: 100,
		ModelsSupported:  []string{"FOFO", "FOCO", "COCO", "FICO"},
		PaybackMonthsEst: 24, CityTiers: []int{1, 2, 3},
	}
	cases := []struct {
		inv  Involvement
		want string
	}{
		{OwnerOperator, "FOFO"},
		{PartTime, "FICO"},
		{Investor, "FOCO"},
	}
	for _, c := range cases {
		in := typicalInput()
		in.Involvement = c.inv
		res := Rank([]catalog.Brand{allModels}, in)
		if len(res.Matches) != 1 {
			t.Fatalf("%s: expected one match", c.inv)
		}
		if got := res.Matches[0].RecommendedModel; got != c.want {
			t.Errorf("%s: recommended %s, want %s", c.inv, got, c.want)
		}
		if res.Matches[0].Reasoning == "" {
			t.Errorf("%s: reasoning paragraph missing", c.inv)
		}
	}
	// A brand that only offers company-operated models must still recommend
	// the best available option for an owner-operator.
	cocoOnly := allModels
	cocoOnly.ModelsSupported = []string{"FOCO", "COCO"}
	in := typicalInput()
	in.Involvement = OwnerOperator
	res := Rank([]catalog.Brand{cocoOnly}, in)
	if got := res.Matches[0].RecommendedModel; got != "FOCO" {
		t.Errorf("owner-operator with FOCO/COCO brand: recommended %s, want FOCO", got)
	}
}

func TestSpacePenaltyHalvesCapitalFactor(t *testing.T) {
	b := catalog.Brand{
		ID: "roomy", Name: "Roomy", Category: "QSR",
		InvestmentMinL: 10, InvestmentMaxL: 20, AreaSqftMin: 1000,
		ModelsSupported: []string{"FOFO"}, PaybackMonthsEst: 24, CityTiers: []int{2},
	}
	withSpace := typicalInput()
	withSpace.SpaceSqft = 1200
	tight := typicalInput()
	tight.SpaceSqft = 300
	a := Rank([]catalog.Brand{b}, withSpace).Matches[0].Factors[0]
	c := Rank([]catalog.Brand{b}, tight).Matches[0].Factors[0]
	if c.Points >= a.Points {
		t.Fatalf("insufficient space should reduce capital factor: %v vs %v", c.Points, a.Points)
	}
	if c.Points != a.Points/2 {
		t.Fatalf("space penalty should halve the factor: got %v, want %v", c.Points, a.Points/2)
	}
}

func TestRiskPaybackAlignment(t *testing.T) {
	short := catalog.Brand{ID: "s", Name: "S", Category: "QSR", InvestmentMinL: 5, InvestmentMaxL: 10,
		ModelsSupported: []string{"FOFO"}, PaybackMonthsEst: 12, CityTiers: []int{2}}
	long := short
	long.ID, long.Name, long.PaybackMonthsEst = "l", "L", 48

	cautious := typicalInput()
	cautious.Risk = 1
	res := Rank([]catalog.Brand{short, long}, cautious)
	if res.Matches[0].Brand.ID != "s" {
		t.Fatal("risk appetite 1 should prefer the 12-month payback brand")
	}
	aggressive := typicalInput()
	aggressive.Risk = 5
	res = Rank([]catalog.Brand{short, long}, aggressive)
	if res.Matches[0].Brand.ID != "l" {
		t.Fatal("risk appetite 5 should prefer the 48-month payback brand")
	}
}

func TestValidateRejectsBadInput(t *testing.T) {
	bad := []Input{
		{BudgetL: 0, Tier: 1, Risk: 3, Involvement: Investor},
		{BudgetL: 10, Tier: 0, Risk: 3, Involvement: Investor},
		{BudgetL: 10, Tier: 4, Risk: 3, Involvement: Investor},
		{BudgetL: 10, Tier: 1, Risk: 0, Involvement: Investor},
		{BudgetL: 10, Tier: 1, Risk: 6, Involvement: Investor},
		{BudgetL: 10, Tier: 1, Risk: 3, Involvement: "ceo"},
		{BudgetL: 10, Tier: 1, Risk: 3, Involvement: Investor, SpaceSqft: -5},
	}
	for i, in := range bad {
		if err := Validate(in); err == nil {
			t.Errorf("case %d: expected validation error for %+v", i, in)
		}
	}
	if err := Validate(typicalInput()); err != nil {
		t.Errorf("valid input rejected: %v", err)
	}
}

func TestCatalogIntegrity(t *testing.T) {
	all := catalog.All()
	if len(all) != 30 {
		t.Fatalf("directory should hold 30 brands, has %d", len(all))
	}
	seen := map[string]bool{}
	validModel := map[string]bool{"FOFO": true, "FOCO": true, "COCO": true, "FICO": true}
	for _, b := range all {
		if seen[b.ID] {
			t.Errorf("duplicate brand id %s", b.ID)
		}
		seen[b.ID] = true
		if b.InvestmentMinL <= 0 || b.InvestmentMaxL < b.InvestmentMinL {
			t.Errorf("%s: bad investment range %v–%v", b.ID, b.InvestmentMinL, b.InvestmentMaxL)
		}
		if len(b.ModelsSupported) == 0 {
			t.Errorf("%s: no models supported", b.ID)
		}
		for _, m := range b.ModelsSupported {
			if !validModel[m] {
				t.Errorf("%s: invalid model %q", b.ID, m)
			}
		}
		if len(b.CityTiers) == 0 {
			t.Errorf("%s: no city tiers", b.ID)
		}
		for _, tier := range b.CityTiers {
			if tier < 1 || tier > 3 {
				t.Errorf("%s: invalid tier %d", b.ID, tier)
			}
		}
		if b.PaybackMonthsEst <= 0 || b.MonthlyProfitEstL <= 0 || b.AreaSqftMin <= 0 {
			t.Errorf("%s: non-positive economics fields", b.ID)
		}
		if b.BrandStory == "" {
			t.Errorf("%s: missing brand story", b.ID)
		}
	}
}

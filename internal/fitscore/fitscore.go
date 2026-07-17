// Package fitscore is the pure, deterministic matching engine.
// It scores every brand 0–100 against a quiz input using weighted factors:
//
//	capital fit 40% · model-vs-involvement 25% · category preference 15%
//	city-tier availability 10% · risk-vs-payback 10%
package fitscore

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"franfit/internal/catalog"
	"franfit/internal/money"
)

// Involvement is how hands-on the aspirant wants to be.
type Involvement string

const (
	OwnerOperator Involvement = "owner-operator"
	PartTime      Involvement = "part-time"
	Investor      Involvement = "investor"
)

// Factor weights (points out of 100).
const (
	WeightCapital  = 40.0
	WeightModel    = 25.0
	WeightCategory = 15.0
	WeightTier     = 10.0
	WeightRisk     = 10.0
)

// Input is one completed quiz.
type Input struct {
	BudgetL     float64     `json:"budgetL"`
	City        string      `json:"city"`
	Tier        int         `json:"tier"`
	Involvement Involvement `json:"involvement"`
	Risk        int         `json:"risk"` // 1 (safety-first) … 5 (aggressive)
	Categories  []string    `json:"categories"`
	SpaceSqft   int         `json:"spaceSqft"`
}

// FactorScore is one factor's contribution to the total.
type FactorScore struct {
	Label  string  `json:"label"`
	Points float64 `json:"points"`
	Max    float64 `json:"max"`
	Note   string  `json:"note"`
}

// Match is one scored brand.
type Match struct {
	Brand            catalog.Brand `json:"brand"`
	Score            int           `json:"score"`
	Factors          []FactorScore `json:"factors"`
	RecommendedModel string        `json:"recommendedModel"`
	Reasoning        string        `json:"reasoning"`
}

// Result is the full ranked response.
type Result struct {
	Matches     []Match `json:"matches"`
	NoMatches   bool    `json:"noMatches"`
	Explanation string  `json:"explanation,omitempty"`
}

// Validate rejects malformed quiz inputs before scoring.
func Validate(in Input) error {
	if in.BudgetL <= 0 {
		return errors.New("budgetL must be a positive number of ₹ lakhs")
	}
	if in.Tier < 1 || in.Tier > 3 {
		return errors.New("tier must be 1, 2 or 3")
	}
	if in.Risk < 1 || in.Risk > 5 {
		return errors.New("risk must be between 1 and 5")
	}
	switch in.Involvement {
	case OwnerOperator, PartTime, Investor:
	default:
		return errors.New("involvement must be owner-operator, part-time or investor")
	}
	if in.SpaceSqft < 0 {
		return errors.New("spaceSqft cannot be negative")
	}
	return nil
}

// modelAffinity maps involvement → operating model → fit in [0,1].
var modelAffinity = map[Involvement]map[string]float64{
	OwnerOperator: {"FOFO": 1.0, "FICO": 0.85, "FOCO": 0.40, "COCO": 0.25},
	PartTime:      {"FICO": 1.0, "FOCO": 0.90, "FOFO": 0.60, "COCO": 0.50},
	Investor:      {"FOCO": 1.0, "COCO": 0.95, "FICO": 0.60, "FOFO": 0.20},
}

var involvementLabel = map[Involvement]string{
	OwnerOperator: "hands-on owner-operator",
	PartTime:      "part-time operator",
	Investor:      "capital investor",
}

var modelExplain = map[string]string{
	"FOFO": "you own the outlet and run daily operations yourself, keeping the largest share of profits",
	"FOCO": "you fund and own the outlet while the brand's trained team runs it day to day",
	"COCO": "the company owns and operates the outlet, and you participate mainly as a capital partner",
	"FICO": "you invest and stay lightly involved while the brand co-manages operations with you",
}

// Rank scores every brand against the input and returns matches sorted by
// score (ties broken by brand name, so output is fully deterministic).
// Brands whose minimum investment exceeds the budget are excluded; if that
// removes everything, NoMatches is set with a plain-English explanation.
func Rank(brands []catalog.Brand, in Input) Result {
	matches := []Match{}
	cheapest := math.Inf(1)
	for _, b := range brands {
		if b.InvestmentMinL < cheapest {
			cheapest = b.InvestmentMinL
		}
		if in.BudgetL < b.InvestmentMinL {
			continue // unaffordable: excluded, not just down-ranked
		}
		matches = append(matches, scoreBrand(b, in))
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}
		return matches[i].Brand.Name < matches[j].Brand.Name
	})
	res := Result{Matches: matches}
	if len(matches) == 0 {
		res.NoMatches = true
		res.Explanation = fmt.Sprintf(
			"Your budget of %s is below the minimum investment of every brand in the directory — the most affordable format starts at %s. Consider raising your budget, adding a co-investor, or exploring bank franchise-financing schemes.",
			money.FormatLakhs(in.BudgetL), money.FormatLakhs(cheapest))
	}
	return res
}

func scoreBrand(b catalog.Brand, in Input) Match {
	capital := capitalFactor(b, in)
	model, recommended := modelFactor(b, in)
	category := categoryFactor(b, in)
	tier := tierFactor(b, in)
	risk := riskFactor(b, in)

	total := capital.Points + model.Points + category.Points + tier.Points + risk.Points
	score := int(math.Round(total))
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return Match{
		Brand:            b,
		Score:            score,
		Factors:          []FactorScore{capital, model, category, tier, risk},
		RecommendedModel: recommended,
		Reasoning:        reasoning(b, in, recommended),
	}
}

// capitalFactor (40%): full marks when the budget covers the top of the
// investment range; scales down toward the minimum. Insufficient floor space
// halves the factor, since capital tied to an unusable site is capital at risk.
func capitalFactor(b catalog.Brand, in Input) FactorScore {
	frac := 1.0
	note := fmt.Sprintf("Budget %s covers the full setup range %s–%s.",
		money.FormatLakhs(in.BudgetL), money.FormatLakhs(b.InvestmentMinL), money.FormatLakhs(b.InvestmentMaxL))
	if in.BudgetL < b.InvestmentMaxL {
		span := b.InvestmentMaxL - b.InvestmentMinL
		pos := 1.0
		if span > 0 {
			pos = (in.BudgetL - b.InvestmentMinL) / span
		}
		frac = 0.6 + 0.4*pos
		note = fmt.Sprintf("Budget %s clears the %s entry point but sits inside the %s–%s range — keep a working-capital buffer.",
			money.FormatLakhs(in.BudgetL), money.FormatLakhs(b.InvestmentMinL), money.FormatLakhs(b.InvestmentMinL), money.FormatLakhs(b.InvestmentMaxL))
	}
	if in.SpaceSqft > 0 && in.SpaceSqft < b.AreaSqftMin {
		frac *= 0.5
		note += fmt.Sprintf(" Your %d sqft is below the %d sqft this format needs, so a new site hunt would be required.", in.SpaceSqft, b.AreaSqftMin)
	}
	return FactorScore{Label: "Capital fit", Points: round1(WeightCapital * frac), Max: WeightCapital, Note: note}
}

// modelFactor (25%): the best affinity between the aspirant's involvement
// style and any operating model the brand supports. Ties resolve in canonical
// model order (FOFO, FOCO, COCO, FICO) so results are deterministic.
func modelFactor(b catalog.Brand, in Input) (FactorScore, string) {
	aff := modelAffinity[in.Involvement]
	best, bestModel := -1.0, ""
	for _, m := range catalog.Models {
		if !contains(b.ModelsSupported, m) {
			continue
		}
		if aff[m] > best {
			best, bestModel = aff[m], m
		}
	}
	note := fmt.Sprintf("As a %s, the %s model is your best route among the models this brand offers (%s).",
		involvementLabel[in.Involvement], bestModel, strings.Join(b.ModelsSupported, ", "))
	return FactorScore{Label: "Model vs involvement", Points: round1(WeightModel * best), Max: WeightModel, Note: note}, bestModel
}

// categoryFactor (15%): full marks for a preferred category, full marks when
// the aspirant is open to everything, a token score otherwise.
func categoryFactor(b catalog.Brand, in Input) FactorScore {
	if len(in.Categories) == 0 {
		return FactorScore{Label: "Category preference", Points: WeightCategory, Max: WeightCategory,
			Note: "You kept all categories open, so " + b.Category + " counts fully."}
	}
	if contains(in.Categories, b.Category) {
		return FactorScore{Label: "Category preference", Points: WeightCategory, Max: WeightCategory,
			Note: b.Category + " is one of your chosen categories."}
	}
	return FactorScore{Label: "Category preference", Points: round1(WeightCategory * 0.15), Max: WeightCategory,
		Note: b.Category + " is outside your chosen categories."}
}

// tierFactor (10%): the brand either franchises in the aspirant's city tier
// or it does not.
func tierFactor(b catalog.Brand, in Input) FactorScore {
	city := in.City
	if city == "" {
		city = "your city"
	}
	for _, t := range b.CityTiers {
		if t == in.Tier {
			return FactorScore{Label: "City-tier availability", Points: WeightTier, Max: WeightTier,
				Note: fmt.Sprintf("Open for tier-%d cities, so %s qualifies.", in.Tier, city)}
		}
	}
	return FactorScore{Label: "City-tier availability", Points: 0, Max: WeightTier,
		Note: fmt.Sprintf("Not currently franchising in tier-%d cities like %s.", in.Tier, city)}
}

// riskFactor (10%): aligns risk appetite with payback horizon. A 12-month
// payback suits appetite 1; a 48-month payback suits appetite 5. Score decays
// with the distance between the two on a 0–1 scale.
func riskFactor(b catalog.Brand, in Input) FactorScore {
	appetite := float64(in.Risk-1) / 4.0
	horizon := clamp((float64(b.PaybackMonthsEst)-12.0)/36.0, 0, 1)
	frac := 1.0 - math.Abs(appetite-horizon)
	note := fmt.Sprintf("Estimated payback of %d months against your risk appetite of %d/5.", b.PaybackMonthsEst, in.Risk)
	return FactorScore{Label: "Risk vs payback", Points: round1(WeightRisk * frac), Max: WeightRisk, Note: note}
}

// reasoning builds the plain-English paragraph explaining the recommended
// model for this specific person and brand.
func reasoning(b catalog.Brand, in Input, model string) string {
	return fmt.Sprintf(
		"%s supports the %s model, which suits your profile as a %s because %s. Your budget of %s stands against a setup cost of %s–%s (including a %s franchise fee), with royalty at %s%% and estimated payback in %d months on projected profits of about %s per month — a pairing that sits comfortably with your risk appetite of %d/5.",
		b.Name, model, involvementLabel[in.Involvement], modelExplain[model],
		money.FormatLakhs(in.BudgetL), money.FormatLakhs(b.InvestmentMinL), money.FormatLakhs(b.InvestmentMaxL),
		money.FormatLakhs(b.FranchiseFeeL), trimPct(b.RoyaltyPct), b.PaybackMonthsEst,
		money.FormatLakhs(b.MonthlyProfitEstL), in.Risk)
}

func trimPct(f float64) string {
	s := fmt.Sprintf("%.1f", f)
	return strings.TrimSuffix(s, ".0")
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if strings.EqualFold(s, v) {
			return true
		}
	}
	return false
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func round1(f float64) float64 { return math.Round(f*10) / 10 }

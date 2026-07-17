# FranFit — Franchise Matchmaker

FranFit matches aspiring Indian franchisees to the right franchise brand *and* the right
operating model (FOFO / FOCO / COCO / FICO). A three-step quiz — money & space, involvement
& risk, categories — feeds a weighted fit-score engine that ranks an embedded directory of
30 illustrative Indian franchise brands, explains every score factor-by-factor, and captures
"request intro" leads.

**Who pays:** franchise brands. Each captured lead arrives pre-qualified with budget, city,
fit score and the model the aspirant should be sold — exactly the lead a franchise sales
team pays per-unit for. The Leads page (with CSV export) is that billable inventory.

> All brand data is **illustrative** — fictional brands invented for the demo, not real
> trademarks.

## Quickstart

```bash
make run          # serves UI + API on http://localhost:8101
make test         # unit tests over scoring, money math, catalog integrity
make build        # single binary at bin/franfit
```

Go 1.25, standard library only. The frontend is embedded in the binary via `go:embed` —
no CDN, fonts or external assets; it runs fully offline. State persists to `./data/store.json`.

## Fit-score engine

Each brand is scored 0–100 with weighted factors (pure, deterministic — see
`internal/fitscore`):

| Factor                  | Weight | Logic |
|-------------------------|-------:|-------|
| Capital fit             |    40% | Full marks when budget covers the top of the investment range; scales down to the entry point; halved if your floor space is under the format's minimum. Brands you can't afford are excluded outright. |
| Model vs involvement    |    25% | owner-operator → FOFO/FICO, part-time → FICO/FOCO, investor → FOCO/COCO; best supported model becomes the recommendation. |
| Category preference     |    15% | Full marks for chosen categories (or when open to all). |
| City-tier availability  |    10% | Brand franchises in your tier, or not. |
| Risk vs payback         |    10% | Aligns risk appetite 1–5 with the payback horizon (12 → 48 months). |

If the budget is below every brand's minimum, the API returns zero matches with
`noMatches: true` and a plain-English `explanation`.

## API summary

| Method | Path                  | Purpose |
|--------|-----------------------|---------|
| GET    | `/api/v1/health`      | `{"status":"ok"}` + provider-mode fields |
| GET    | `/api/v1/brands`      | Full directory (optional `?category=QSR` filter) |
| GET    | `/api/v1/brands/{id}` | One brand |
| POST   | `/api/v1/match`       | Quiz input → ranked matches with factor breakdowns, recommended model, reasoning |
| POST   | `/api/v1/leads`       | Capture a lead (name, phone, email, budgetL, brandId) |
| GET    | `/api/v1/leads`       | All captured leads, newest first |
| GET    | `/api/v1/leads.csv`   | Lead export for brand billing |

Example match call:

```bash
curl -s localhost:8101/api/v1/match -d '{
  "budgetL": 30, "city": "Jaipur", "tier": 2,
  "involvement": "owner-operator", "risk": 3,
  "categories": ["QSR","Cafe"], "spaceSqft": 500
}'
```

## Upgrade to live

Every external integration ships as a deterministic, zero-key mock behind a provider
interface (`internal/notify.Provider`). Message ids are FNV-1a hashes of the input, so the
same lead always produces the same id.

| Env var                    | Integration                | Today (unset)                         | When set |
|----------------------------|----------------------------|---------------------------------------|----------|
| `FRANFIT_WHATSAPP_TOKEN`   | WhatsApp Business Cloud API| `notify.MockWhatsApp` (FNV mock ids)  | Live lead alerts to brand franchise teams |
| `FRANFIT_WHATSAPP_PHONE_ID`| WhatsApp sender number     | unused                                | Live sender identity |
| `PORT`                     | —                          | 8101                                  | Any port |

`GET /api/v1/health` reports the active mode: `{"status":"ok","providerMode":"mock","providers":{"whatsapp":"mock"}}`.

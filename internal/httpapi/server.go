// Package httpapi wires the JSON API and the embedded web UI onto one mux.
package httpapi

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"franfit/internal/catalog"
	"franfit/internal/fitscore"
	"franfit/internal/notify"
	"franfit/internal/store"
	"franfit/web"
)

// Server holds the app's dependencies.
type Server struct {
	store    *store.Store
	notifier notify.Provider
}

// New builds the full HTTP handler: /api/v1/* plus the embedded UI at /.
func New(st *store.Store, notifier notify.Provider) http.Handler {
	s := &Server{store: st, notifier: notifier}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/health", s.handleHealth)
	mux.HandleFunc("GET /api/v1/brands", s.handleBrands)
	mux.HandleFunc("GET /api/v1/brands/{id}", s.handleBrand)
	mux.HandleFunc("POST /api/v1/match", s.handleMatch)
	mux.HandleFunc("POST /api/v1/leads", s.handleCreateLead)
	mux.HandleFunc("GET /api/v1/leads", s.handleLeads)
	mux.HandleFunc("GET /api/v1/leads.csv", s.handleLeadsCSV)

	mux.Handle("GET /", http.FileServerFS(web.Files))
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":       "ok",
		"service":      "franfit",
		"providerMode": s.notifier.Mode(),
		"providers": map[string]string{
			s.notifier.Name(): s.notifier.Mode(),
		},
	})
}

func (s *Server) handleBrands(w http.ResponseWriter, r *http.Request) {
	brands := catalog.All()
	if cat := r.URL.Query().Get("category"); cat != "" {
		filtered := brands[:0:0]
		for _, b := range brands {
			if strings.EqualFold(b.Category, cat) {
				filtered = append(filtered, b)
			}
		}
		brands = filtered
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"brands":     brands,
		"count":      len(brands),
		"categories": catalog.Categories,
	})
}

func (s *Server) handleBrand(w http.ResponseWriter, r *http.Request) {
	b, ok := catalog.ByID(r.PathValue("id"))
	if !ok {
		writeError(w, http.StatusNotFound, "no brand with id "+r.PathValue("id"))
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (s *Server) handleMatch(w http.ResponseWriter, r *http.Request) {
	var in fitscore.Input
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := fitscore.Validate(in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, fitscore.Rank(catalog.All(), in))
}

type leadRequest struct {
	BrandID  string  `json:"brandId"`
	Name     string  `json:"name"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	BudgetL  float64 `json:"budgetL"`
	FitScore int     `json:"fitScore"`
	City     string  `json:"city"`
}

func (s *Server) handleCreateLead(w http.ResponseWriter, r *http.Request) {
	var req leadRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)
	switch {
	case req.Name == "":
		writeError(w, http.StatusBadRequest, "name is required")
		return
	case len(strings.Map(keepDigits, req.Phone)) < 10:
		writeError(w, http.StatusBadRequest, "phone must contain at least 10 digits, e.g. +91-98XXXXXXXX")
		return
	case !strings.Contains(req.Email, "@") || strings.HasPrefix(req.Email, "@") || strings.HasSuffix(req.Email, "@"):
		writeError(w, http.StatusBadRequest, "a valid email is required")
		return
	case req.BudgetL <= 0:
		writeError(w, http.StatusBadRequest, "budgetL must be a positive number of ₹ lakhs")
		return
	}
	brand, ok := catalog.ByID(req.BrandID)
	if !ok {
		writeError(w, http.StatusBadRequest, "unknown brandId "+req.BrandID)
		return
	}
	lead, err := s.store.AddLead(store.Lead{
		BrandID: brand.ID, BrandName: brand.Name,
		Name: req.Name, Phone: req.Phone, Email: req.Email,
		BudgetL: req.BudgetL, FitScore: req.FitScore, City: strings.TrimSpace(req.City),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save lead: "+err.Error())
		return
	}
	msgID, err := s.notifier.NotifyLead(brand.Name, lead.Name, lead.Phone)
	if err != nil {
		log.Printf("notify failed for %s: %v", lead.ID, err)
	} else if err := s.store.SetNotificationID(lead.ID, msgID); err == nil {
		lead.NotificationID = msgID
	}
	writeJSON(w, http.StatusCreated, lead)
}

func (s *Server) handleLeads(w http.ResponseWriter, _ *http.Request) {
	leads := s.store.Leads()
	writeJSON(w, http.StatusOK, map[string]any{"leads": leads, "count": len(leads)})
}

func (s *Server) handleLeadsCSV(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="franfit-leads.csv"`)
	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"id", "createdAt", "name", "phone", "email", "city", "budgetL", "brandId", "brandName", "fitScore", "notificationId"})
	for _, l := range s.store.Leads() {
		_ = cw.Write([]string{
			l.ID, l.CreatedAt.Format("2006-01-02 15:04:05"), l.Name, l.Phone, l.Email, l.City,
			strconv.FormatFloat(l.BudgetL, 'f', -1, 64), l.BrandID, l.BrandName,
			strconv.Itoa(l.FitScore), l.NotificationID,
		})
	}
	cw.Flush()
}

func keepDigits(r rune) rune {
	if r >= '0' && r <= '9' {
		return r
	}
	return -1
}

func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

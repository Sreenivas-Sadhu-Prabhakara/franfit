// Package notify defines the lead-notification provider seam. In production
// a live WhatsApp Business API client would implement Provider (switched on
// via FRANFIT_WHATSAPP_TOKEN); this repo ships only the deterministic mock.
package notify

import (
	"fmt"
	"hash/fnv"
)

// Provider notifies a brand's franchise team about a new qualified lead and
// returns the provider-side message id.
type Provider interface {
	NotifyLead(brandName, leadName, phone string) (messageID string, err error)
	Mode() string // "mock" or "live"
	Name() string // e.g. "whatsapp"
}

// MockWhatsApp is a zero-key stand-in for the WhatsApp Business API. The
// message id is an FNV-1a hash of the inputs, so the same lead always yields
// the same id — deterministic and testable.
type MockWhatsApp struct{}

// NotifyLead returns a stable synthetic WhatsApp message id.
func (MockWhatsApp) NotifyLead(brandName, leadName, phone string) (string, error) {
	h := fnv.New64a()
	fmt.Fprintf(h, "%s|%s|%s", brandName, leadName, phone)
	return fmt.Sprintf("wamid.MOCK-%016X", h.Sum64()), nil
}

// Mode reports that this provider is a mock.
func (MockWhatsApp) Mode() string { return "mock" }

// Name identifies the integration this mock stands in for.
func (MockWhatsApp) Name() string { return "whatsapp" }

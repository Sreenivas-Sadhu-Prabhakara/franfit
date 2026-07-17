// Command server runs FranFit: the JSON API under /api/v1 and the embedded
// web UI on one port (PORT env var, default 8101).
package main

import (
	"log"
	"net/http"
	"os"

	"franfit/internal/httpapi"
	"franfit/internal/notify"
	"franfit/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8101"
	}
	st, err := store.Open("data/store.json")
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	var notifier notify.Provider = notify.MockWhatsApp{}
	// A live WhatsApp Business client would be selected here when
	// FRANFIT_WHATSAPP_TOKEN is set; only the mock ships in this build.

	handler := httpapi.New(st, notifier)
	log.Printf("FranFit listening on http://localhost:%s (notify provider: %s/%s)",
		port, notifier.Name(), notifier.Mode())
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

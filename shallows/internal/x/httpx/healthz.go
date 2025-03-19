package httpx

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/internal/x/cryptox"
	"github.com/retrovibed/retrovibed/internal/x/slicesx"
)

// Healthz
func Healthz(seed string, prob float64, defaultCode int, failures ...int) http.HandlerFunc {
	log.Println("healthz enabled with", defaultCode)
	return func(w http.ResponseWriter, r *http.Request) {
		_, week := time.Now().ISOWeek()
		cha := rand.New(cryptox.NewChaCha8(fmt.Sprintf("%s+%d", seed, week)))
		if chance := cha.Float64(); chance < prob {
			failure := slicesx.FirstOrZero(failures...)
			log.Println("returning alternative healthz code", chance, "<", prob, failure)
			w.WriteHeader(failure)
			return
		}

		w.WriteHeader(defaultCode)
	}
}

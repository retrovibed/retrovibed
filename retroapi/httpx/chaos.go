package httpx

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func ChaosRateLimited(max time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", fmt.Sprintf("%d", int(max.Truncate(time.Second)/time.Second)))
		w.WriteHeader(http.StatusTooManyRequests)
	})
}

func ChaosStatusCodes(codes ...int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := codes[rand.Intn(len(codes))]
		w.WriteHeader(code)
	})
}

// Chaos inject random errors into the application
// enabled only in dev environments. rate is the percentage
// of request to mess with.
func Chaos(rate float64, behavior ...http.Handler) func(http.Handler) http.Handler {
	if rate == 0 || len(behavior) == 0 {
		return func(original http.Handler) http.Handler {
			return original
		}
	}
	log.Println("chaos events enabled", rate)
	return func(original http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if prob := rand.Float64(); rate < prob {
				original.ServeHTTP(resp, req)
				return
			} else {
				log.Println(" generating choas event", prob, "<", rate, prob < rate)
			}

			n := behavior[rand.Intn(len(behavior))]
			n.ServeHTTP(resp, req)
		})
	}
}

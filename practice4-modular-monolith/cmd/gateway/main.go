package gateway

import (
	"log"
	"net/http"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		shared.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "UP",
		})
	})

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

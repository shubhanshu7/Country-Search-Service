package httpapi

import "net/http"

func Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/countries/search", handleSearch)
}
func handleSearch(w http.ResponseWriter, r *http.Request) {

}

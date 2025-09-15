# Country Search Service

A production-ready Go service that fetches country details from the [REST Countries API](https://restcountries.com), with in-memory caching, concurrency safety, and graceful shutdown.

---

##  Features
- **REST API**: `GET /api/countries/search?name={country}`
- **In-memory cache** (thread-safe, custom implementation)
- **Service layer**: validates, caches, dedups requests
- **External API client**: wraps REST Countries v3.1
- **Graceful shutdown**: handles `SIGINT` / `SIGTERM`
- **Configurable timeouts** for HTTP client & handlers
- **Extensive test suite**: >90% coverage, includes concurrency race checks
- **Clean project layout** (cmd/internal split)

---

## 📂 Project Structure
country-search/
├─ cmd/server/ # entrypoint
│ └─ main.go
├─ internal/
│ ├─ cache/ # custom cache
│ ├─ countries/ # REST Countries client
│ ├─ service/ # business logic
│ └─ httpapi/ # HTTP handlers
├─ go.mod / go.sum
└─ README.md


---

### Prerequisites
- Go 1.22+
- Internet access (to call REST Countries API)

### Install dependencies
```bash
go mod tidy

### Run the service
go run ./cmd/server

### API Usage
**Search for a country**
curl "http://localhost:8000/api/countries/search?name=India"


**Example Response**
{
  "name": "India",
  "capital": "New Delhi",
  "currency": "₹",
  "population": 1380004385
}


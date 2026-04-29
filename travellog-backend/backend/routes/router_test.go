package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/travellog/backend/config"
	"github.com/travellog/backend/database"
	"github.com/travellog/backend/models"
	"github.com/travellog/backend/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)

	dsn := "host=localhost port=5432 user=travellog password=travellog123 dbname=travellog sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// reset tables for each test
	err = db.Exec("TRUNCATE TABLE flights, hotels, users RESTART IDENTITY CASCADE").Error
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Hotel{}, &models.Flight{})
	if err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	database.DB = db

	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "travellog",
		DBPassword: "travellog123",
		DBName:     "travellog",
		JWTSecret:  "test-secret",
		ServerPort: "8080",
	}

	return routes.SetupRouter(cfg)
}

func performRequest(r http.Handler, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func registerAndLogin(t *testing.T, r *gin.Engine, email, password string) string {
	t.Helper()

	registerBody := map[string]string{
		"email":    email,
		"password": password,
	}

	w := performRequest(r, http.MethodPost, "/auth/register", registerBody, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("register failed: status=%d body=%s", w.Code, w.Body.String())
	}

	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}

	w = performRequest(r, http.MethodPost, "/auth/login", loginBody, "")
	if w.Code != http.StatusOK {
		t.Fatalf("login failed: status=%d body=%s", w.Code, w.Body.String())
	}

	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}
	if resp.Token == "" {
		t.Fatalf("empty token returned")
	}

	return resp.Token
}

func TestRegisterLoginMe(t *testing.T) {
	r := setupTestRouter(t)

	token := registerAndLogin(t, r, "user1@test.com", "12345678")

	w := performRequest(r, http.MethodGet, "/me", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("/me failed: status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateFlight(t *testing.T) {
	r := setupTestRouter(t)
	token := registerAndLogin(t, r, "user2@test.com", "12345678")

	body := map[string]interface{}{
		"airline":       "Delta",
		"flight_number": "DL123",
		"from_airport":  "JFK",
		"to_airport":    "LAX",
		"depart_time":   "2026-03-01T10:00:00Z",
		"arrive_time":   "2026-03-01T14:00:00Z",
		"price":         320.5,
		"notes":         "window seat",
	}

	w := performRequest(r, http.MethodPost, "/flights", body, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create flight failed: status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestGetFlightsOnlyOwnData(t *testing.T) {
	r := setupTestRouter(t)

	token1 := registerAndLogin(t, r, "user3@test.com", "12345678")
	token2 := registerAndLogin(t, r, "user4@test.com", "12345678")

	body1 := map[string]interface{}{
		"airline":       "AA",
		"flight_number": "AA100",
		"from_airport":  "JFK",
		"to_airport":    "ORD",
		"depart_time":   "2026-03-01T10:00:00Z",
		"arrive_time":   "2026-03-01T12:00:00Z",
		"price":         200,
		"notes":         "user1 flight",
	}
	body2 := map[string]interface{}{
		"airline":       "UA",
		"flight_number": "UA200",
		"from_airport":  "LAX",
		"to_airport":    "SFO",
		"depart_time":   "2026-03-02T10:00:00Z",
		"arrive_time":   "2026-03-02T11:00:00Z",
		"price":         120,
		"notes":         "user2 flight",
	}

	w := performRequest(r, http.MethodPost, "/flights", body1, token1)
	if w.Code != http.StatusCreated {
		t.Fatalf("create flight for user1 failed: status=%d body=%s", w.Code, w.Body.String())
	}

	w = performRequest(r, http.MethodPost, "/flights", body2, token2)
	if w.Code != http.StatusCreated {
		t.Fatalf("create flight for user2 failed: status=%d body=%s", w.Code, w.Body.String())
	}

	w = performRequest(r, http.MethodGet, "/flights", nil, token1)
	if w.Code != http.StatusOK {
		t.Fatalf("get flights failed: status=%d body=%s", w.Code, w.Body.String())
	}

	var resp struct {
		Flights []struct {
			Notes string `json:"notes"`
		} `json:"flights"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse flights response: %v", err)
	}

	if len(resp.Flights) != 1 {
		t.Fatalf("expected 1 flight for user1, got %d, body=%s", len(resp.Flights), w.Body.String())
	}
	if resp.Flights[0].Notes != "user1 flight" {
		t.Fatalf("expected only user1 flight, got body=%s", w.Body.String())
	}
}

func TestDeleteFlight(t *testing.T) {
	r := setupTestRouter(t)
	token := registerAndLogin(t, r, "user5@test.com", "12345678")

	createBody := map[string]interface{}{
		"airline":       "JetBlue",
		"flight_number": "B6123",
		"from_airport":  "BOS",
		"to_airport":    "MCO",
		"depart_time":   "2026-03-03T10:00:00Z",
		"arrive_time":   "2026-03-03T13:00:00Z",
		"price":         180,
		"notes":         "to delete",
	}

	w := performRequest(r, http.MethodPost, "/flights", createBody, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create flight failed: status=%d body=%s", w.Code, w.Body.String())
	}

	var createResp struct {
		Flight struct {
			ID uint `json:"id"`
		} `json:"flight"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to parse create response: %v", err)
	}

	deletePath := "/flights/" + strconv.FormatUint(uint64(createResp.Flight.ID), 10)

	w = performRequest(r, http.MethodDelete, deletePath, nil, token)
	if w.Code != http.StatusOK && w.Code != http.StatusNoContent {
		t.Fatalf("delete flight failed: status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestFlightsRequireAuth(t *testing.T) {
	r := setupTestRouter(t)

	w := performRequest(r, http.MethodGet, "/flights", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unauthenticated request, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateFlightInvalidTime(t *testing.T) {
	r := setupTestRouter(t)
	token := registerAndLogin(t, r, "user6@test.com", "12345678")

	body := map[string]interface{}{
		"airline":       "Delta",
		"flight_number": "DL999",
		"from_airport":  "JFK",
		"to_airport":    "LAX",
		"depart_time":   "not-a-time",
		"arrive_time":   "2026-03-01T14:00:00Z",
		"price":         320.5,
		"notes":         "bad time",
	}

	w := performRequest(r, http.MethodPost, "/flights", body, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid time, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateHotel(t *testing.T) {
	r := setupTestRouter(t)
	token := registerAndLogin(t, r, "hotel1@test.com", "12345678")

	body := map[string]interface{}{
		"hotel_name": "Hilton Orlando",
		"city":       "Orlando",
		"country":    "USA",
		"check_in":   "2026-03-10",
		"check_out":  "2026-03-12",
		"price":      250.0,
		"notes":      "near airport",
	}

	w := performRequest(r, http.MethodPost, "/hotels", body, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create hotel failed: status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestGetHotelsOnlyOwnData(t *testing.T) {
	r := setupTestRouter(t)

	token1 := registerAndLogin(t, r, "hotel2@test.com", "12345678")
	token2 := registerAndLogin(t, r, "hotel3@test.com", "12345678")

	body1 := map[string]interface{}{
		"hotel_name": "Marriott",
		"city":       "Miami",
		"country":    "USA",
		"check_in":   "2026-03-15",
		"check_out":  "2026-03-18",
		"price":      300.0,
		"notes":      "user1 hotel",
	}

	body2 := map[string]interface{}{
		"hotel_name": "Hyatt",
		"city":       "Tampa",
		"country":    "USA",
		"check_in":   "2026-03-20",
		"check_out":  "2026-03-22",
		"price":      220.0,
		"notes":      "user2 hotel",
	}

	w := performRequest(r, http.MethodPost, "/hotels", body1, token1)
	if w.Code != http.StatusCreated {
		t.Fatalf("create hotel for user1 failed: status=%d body=%s", w.Code, w.Body.String())
	}

	w = performRequest(r, http.MethodPost, "/hotels", body2, token2)
	if w.Code != http.StatusCreated {
		t.Fatalf("create hotel for user2 failed: status=%d body=%s", w.Code, w.Body.String())
	}

	w = performRequest(r, http.MethodGet, "/hotels", nil, token1)
	if w.Code != http.StatusOK {
		t.Fatalf("get hotels failed: status=%d body=%s", w.Code, w.Body.String())
	}

	var resp struct {
		Hotels []struct {
			Notes string `json:"notes"`
		} `json:"hotels"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse hotels response: %v", err)
	}

	if len(resp.Hotels) != 1 {
		t.Fatalf("expected 1 hotel for user1, got %d, body=%s", len(resp.Hotels), w.Body.String())
	}
	if resp.Hotels[0].Notes != "user1 hotel" {
		t.Fatalf("expected only user1 hotel, got body=%s", w.Body.String())
	}
}

// Sprint 3: Map & Summary Tests

func TestGetMap_Success(t *testing.T) {
	r := setupTestRouter(t)
	token := registerAndLogin(t, r, "map1@test.com", "12345678")

	// create hotel
	hotelBody := map[string]interface{}{
		"hotel_name": "Test Hotel",
		"city":       "NYC",
		"country":    "USA",
		"check_in":   "2026-03-01",
		"check_out":  "2026-03-03",
		"price":      200,
	}
	w := performRequest(r, http.MethodPost, "/hotels", hotelBody, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create hotel failed: %s", w.Body.String())
	}

	// create flight
	flightBody := map[string]interface{}{
		"airline":       "Delta",
		"flight_number": "DL100",
		"from_airport":  "JFK",
		"to_airport":    "LAX",
		"depart_time":   "2026-03-01T10:00:00Z",
		"arrive_time":   "2026-03-01T14:00:00Z",
		"price":         300,
	}
	w = performRequest(r, http.MethodPost, "/flights", flightBody, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create flight failed: %s", w.Body.String())
	}

	// call /map
	w = performRequest(r, http.MethodGet, "/map", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("map failed: %d %s", w.Code, w.Body.String())
	}

	var resp struct {
		Hotels  []interface{} `json:"hotels"`
		Flights []interface{} `json:"flights"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("parse map response failed: %v", err)
	}

	if len(resp.Hotels) != 1 || len(resp.Flights) != 1 {
		t.Fatalf("unexpected map data: %s", w.Body.String())
	}
}

func TestGetMap_Unauthorized(t *testing.T) {
	r := setupTestRouter(t)

	w := performRequest(r, http.MethodGet, "/map", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestGetSummary_Month(t *testing.T) {
	r := setupTestRouter(t)
	token := registerAndLogin(t, r, "sum1@test.com", "12345678")

	// create flight
	flightBody := map[string]interface{}{
		"airline":       "AA",
		"flight_number": "AA1",
		"from_airport":  "JFK",
		"to_airport":    "LAX",
		"depart_time":   "2026-03-01T10:00:00Z",
		"arrive_time":   "2026-03-01T14:00:00Z",
		"price":         100,
	}
	performRequest(r, http.MethodPost, "/flights", flightBody, token)

	// create hotel
	hotelBody := map[string]interface{}{
		"hotel_name": "Hilton",
		"city":       "LA",
		"country":    "USA",
		"check_in":   "2026-03-01",
		"check_out":  "2026-03-03",
		"price":      200,
	}
	performRequest(r, http.MethodPost, "/hotels", hotelBody, token)

	w := performRequest(r, http.MethodGet, "/summary?type=month", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("summary failed: %d %s", w.Code, w.Body.String())
	}

	var resp struct {
		Type      string `json:"type"`
		Summaries []struct {
			Flights int     `json:"flights"`
			Hotels  int     `json:"hotels"`
			Nights  int     `json:"nights"`
			Spend   float64 `json:"spend"`
		} `json:"summaries"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("parse summary failed: %v", err)
	}

	if resp.Type != "month" {
		t.Fatalf("wrong type")
	}

	if len(resp.Summaries) == 0 {
		t.Fatalf("empty summary")
	}
}

func TestGetSummary_Unauthorized(t *testing.T) {
	r := setupTestRouter(t)

	w := performRequest(r, http.MethodGet, "/summary", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestGetMap_OnlyOwnData(t *testing.T) {
	r := setupTestRouter(t)

	token1 := registerAndLogin(t, r, "mapown1@test.com", "12345678")
	token2 := registerAndLogin(t, r, "mapown2@test.com", "12345678")

	hotel1 := map[string]interface{}{
		"hotel_name": "User1 Hotel",
		"city":       "New York",
		"country":    "USA",
		"check_in":   "2026-03-01",
		"check_out":  "2026-03-03",
		"price":      100,
		"notes":      "u1",
	}
	hotel2 := map[string]interface{}{
		"hotel_name": "User2 Hotel",
		"city":       "Boston",
		"country":    "USA",
		"check_in":   "2026-03-05",
		"check_out":  "2026-03-06",
		"price":      120,
		"notes":      "u2",
	}

	flight1 := map[string]interface{}{
		"airline":       "AA",
		"flight_number": "AA101",
		"from_airport":  "JFK",
		"to_airport":    "LAX",
		"depart_time":   "2026-03-01T10:00:00Z",
		"arrive_time":   "2026-03-01T14:00:00Z",
		"price":         200,
		"notes":         "u1-flight",
	}
	flight2 := map[string]interface{}{
		"airline":       "UA",
		"flight_number": "UA202",
		"from_airport":  "SFO",
		"to_airport":    "SEA",
		"depart_time":   "2026-03-02T10:00:00Z",
		"arrive_time":   "2026-03-02T12:00:00Z",
		"price":         180,
		"notes":         "u2-flight",
	}

	w := performRequest(r, http.MethodPost, "/hotels", hotel1, token1)
	if w.Code != http.StatusCreated {
		t.Fatalf("create user1 hotel failed: %s", w.Body.String())
	}
	w = performRequest(r, http.MethodPost, "/hotels", hotel2, token2)
	if w.Code != http.StatusCreated {
		t.Fatalf("create user2 hotel failed: %s", w.Body.String())
	}
	w = performRequest(r, http.MethodPost, "/flights", flight1, token1)
	if w.Code != http.StatusCreated {
		t.Fatalf("create user1 flight failed: %s", w.Body.String())
	}
	w = performRequest(r, http.MethodPost, "/flights", flight2, token2)
	if w.Code != http.StatusCreated {
		t.Fatalf("create user2 flight failed: %s", w.Body.String())
	}

	w = performRequest(r, http.MethodGet, "/map", nil, token1)
	if w.Code != http.StatusOK {
		t.Fatalf("get map failed: %d %s", w.Code, w.Body.String())
	}

	var resp struct {
		Hotels []struct {
			HotelName string `json:"hotel_name"`
		} `json:"hotels"`
		Flights []struct {
			FlightNumber string `json:"flight_number"`
		} `json:"flights"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("parse map response failed: %v", err)
	}

	if len(resp.Hotels) != 1 || resp.Hotels[0].HotelName != "User1 Hotel" {
		t.Fatalf("expected only user1 hotel, got %s", w.Body.String())
	}
	if len(resp.Flights) != 1 || resp.Flights[0].FlightNumber != "AA101" {
		t.Fatalf("expected only user1 flight, got %s", w.Body.String())
	}
}

func TestGetSummary_InvalidTypeDefaultsToMonth(t *testing.T) {
	r := setupTestRouter(t)
	token := registerAndLogin(t, r, "suminvalid@test.com", "12345678")

	flightBody := map[string]interface{}{
		"airline":       "Delta",
		"flight_number": "DL555",
		"from_airport":  "JFK",
		"to_airport":    "LAX",
		"depart_time":   "2026-03-01T10:00:00Z",
		"arrive_time":   "2026-03-01T14:00:00Z",
		"price":         150,
	}

	w := performRequest(r, http.MethodPost, "/flights", flightBody, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create flight failed: %s", w.Body.String())
	}

	w = performRequest(r, http.MethodGet, "/summary?type=badvalue", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("summary failed: %d %s", w.Code, w.Body.String())
	}

	var resp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("parse summary response failed: %v", err)
	}

	if resp.Type != "month" {
		t.Fatalf("expected default type month, got %s", resp.Type)
	}
}

func TestHotelsRequireAuth(t *testing.T) {
	r := setupTestRouter(t)

	req, _ := http.NewRequest("GET", "/hotels", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestUnknownRoute(t *testing.T) {
	r := setupTestRouter(t)

	req, _ := http.NewRequest("GET", "/unknown_route_123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestSummaryNoAuth(t *testing.T) {
	r := setupTestRouter(t)

	req, _ := http.NewRequest("GET", "/summary?type=month", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestCreateFlightEmptyBody(t *testing.T) {
	r := setupTestRouter(t)

	req, _ := http.NewRequest("POST", "/flights", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

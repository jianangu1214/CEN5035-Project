package handlers

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/travellog/backend/database"
	"github.com/travellog/backend/models"
)

type SummaryHandler struct{}

func NewSummaryHandler() *SummaryHandler {
	return &SummaryHandler{}
}

type SummaryEntry struct {
	Period  string  `json:"period"`
	Flights int     `json:"flights"`
	Hotels  int     `json:"hotels"`
	Nights  int     `json:"nights"`
	Spend   float64 `json:"spend"`
}

func (h *SummaryHandler) GetSummary(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	summaryType := c.DefaultQuery("type", "month")
	if summaryType != "month" && summaryType != "quarter" && summaryType != "year" {
		summaryType = "month"
	}

	db := database.GetDB()

	var flights []models.Flight
	if err := db.Where("user_id = ?", userID).Find(&flights).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flights"})
		return
	}

	var hotels []models.Hotel
	if err := db.Where("user_id = ?", userID).Find(&hotels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hotels"})
		return
	}

	buckets := map[string]*SummaryEntry{}

	for _, f := range flights {
		key := periodKey(f.DepartTime, summaryType)
		ensureBucket(buckets, key)
		b := buckets[key]
		b.Flights++
		b.Spend += f.Price
	}

	for _, h := range hotels {
		key := periodKey(h.CheckIn, summaryType)
		ensureBucket(buckets, key)
		b := buckets[key]
		b.Hotels++
		nights := int(h.CheckOut.Sub(h.CheckIn).Hours() / 24)
		if nights < 0 {
			nights = 0
		}
		b.Nights += nights
		b.Spend += h.Price
	}

	// sort periods ascending
	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]SummaryEntry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, *buckets[k])
	}

	c.JSON(http.StatusOK, gin.H{
		"type":      summaryType,
		"summaries": entries,
	})
}

func periodKey(t time.Time, summaryType string) string {
	switch summaryType {
	case "year":
		return t.Format("2006")
	case "quarter":
		q := (int(t.Month())-1)/3 + 1
		return t.Format("2006") + "-Q" + string('0'+q)
	default: // month
		return t.Format("2006-01")
	}
}

func ensureBucket(b map[string]*SummaryEntry, key string) {
	if _, ok := b[key]; !ok {
		b[key] = &SummaryEntry{Period: key}
	}
}

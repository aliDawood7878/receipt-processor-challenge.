package main

import (
    "fmt"
    "log"
    "net/http"
	"encoding/json"
    "github.com/google/uuid"
	"math"
    "regexp"
    "strconv"
    "strings"
	"time"
)

type Item struct {
    ShortDescription string `json:"shortDescription"`
    Price            string `json:"price"`
}

type Receipt struct {
    Retailer     string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"` 
    PurchaseTime string `json:"purchaseTime"` 
    Items        []Item `json:"items"`
    Total        string `json:"total"`
}

var receipts = make(map[string]int)

func processReceiptsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var receipt Receipt
    if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

	if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || len(receipt.Items) == 0 || receipt.Total == "" {
		http.Error(w, "Missing required fields.", http.StatusBadRequest)
		return
	}

	if !validDateFormat(receipt.PurchaseDate) {
		http.Error(w, "Invalid purchaseDate format. Expected YYYY-MM-DD.", http.StatusBadRequest)
		return
	}
	
	if !validTimeFormat(receipt.PurchaseTime) {
		http.Error(w, "Invalid purchaseTime format. Expected HH:MM in 24-hour format.", http.StatusBadRequest)
		return
	}

    
    points := calculatePoints(receipt)

    
    id := uuid.New().String()

    receipts[id] = points

    
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"id":"%s"}`, id)
}

func validDateFormat(dateStr string) bool {
    _, err := time.Parse("2006-01-02", dateStr)
    return err == nil
}

func validTimeFormat(timeStr string) bool {
    _, err := time.Parse("15:04", timeStr)
    return err == nil
}

func calculatePoints(receipt Receipt) int {
    points := 0

    
    retailerAlnum := regexp.MustCompile("[A-Za-z0-9]")
    alnumCount := len(retailerAlnum.FindAllString(receipt.Retailer, -1))
    points += alnumCount

    
    totalVal, _ := strconv.ParseFloat(receipt.Total, 64)

    
    if math.Mod(totalVal, 1.0) == 0 {
        points += 50
    }

    
    quarter := 0.25
    if math.Mod(totalVal, quarter) == 0 {
        points += 25
    }

    
    itemCount := len(receipt.Items)
    points += (itemCount / 2) * 5

    
    for _, item := range receipt.Items {
        descLen := len(strings.TrimSpace(item.ShortDescription))
        itemPrice, _ := strconv.ParseFloat(item.Price, 64)
        if descLen%3 == 0 {
            bonus := math.Ceil(itemPrice * 0.2)
            points += int(bonus)
        }
    }

    
    dateParts := strings.Split(receipt.PurchaseDate, "-")
    if len(dateParts) == 3 {
        day, _ := strconv.Atoi(dateParts[2])
        if day%2 == 1 {
            points += 6
        }
    }

    
    timeParts := strings.Split(receipt.PurchaseTime, ":")
    if len(timeParts) == 2 {
        hour, _ := strconv.Atoi(timeParts[0])
        if hour >= 14 && hour < 16 {
            points += 10
        }
    }

    return points
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    
    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 3 {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    id := pathParts[2] 
    points, exists := receipts[id]
    if !exists {
        http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"points": %d}`, points)
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Receipt Processor is running!")
    })

	http.HandleFunc("/receipts/process", processReceiptsHandler)

	http.HandleFunc("/receipts/", func(w http.ResponseWriter, r *http.Request) {
		
		if strings.HasSuffix(r.URL.Path, "/points") {
			getPointsHandler(w, r)
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})


    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
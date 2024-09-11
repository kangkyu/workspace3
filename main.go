package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ReceiptProcessResponse struct {
	ID string `json:"id"`
}

type ReceiptRewardResponse struct {
	Points int `json:"points"`
}

type Reward struct {
	id      string
	points  int
	receipt Receipt
}

var rewards = []Reward{}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /receipts/process", receiptProcess)
	mux.HandleFunc("GET /receipts/{id}/points", receiptReward)

	log.Print("starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func receiptReward(w http.ResponseWriter, r *http.Request) {

	receiptID := r.PathValue("id")
	w.Header().Set("Content-Type", "application/json")

	reward, ok := findReceipt(receiptID)

	if !ok {
		http.Error(w, "receipt not found", http.StatusNotFound)
		return
	}

	resp := ReceiptRewardResponse{
		Points: reward.points,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "error building the response", http.StatusInternalServerError)
		return
	}
}

func receiptProcess(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var receipt = Receipt{}
	err := d.Decode(&receipt)
	if err != nil {
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	points := rewardCalculate(receipt)

	rewards = append(rewards, Reward{receipt: receipt, points: points, id: id})

	resp := ReceiptProcessResponse{
		ID: id,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "error building the response", http.StatusInternalServerError)
		return
	}
}

func findReceipt(id string) (Reward, bool) {
	var reward Reward
	for _, v := range rewards {
		if v.id == id {
			reward = v
		}
	}
	return reward, reward.id != ""
}

func rewardCalculate(receipt Receipt) int {
	var points int
	/*
		One point for every alphanumeric character in the retailer name.
		50 points if the total is a round dollar amount with no cents.
		25 points if the total is a multiple of 0.25.
		5 points for every two items on the receipt.
		If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
		6 points if the day in the purchase date is odd.
		10 points if the time of purchase is after 2:00pm and before 4:00pm.
	*/

	totalFloat, err := strconv.ParseFloat(strings.TrimSpace(receipt.Total), 64)
	if err == nil {
		if float64(int(totalFloat)) == totalFloat {
			points += 50
		}
		if float64(int(totalFloat/0.25)) == totalFloat/0.25 {
			points += 25
		}
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			priceFloat, err := strconv.ParseFloat(strings.TrimSpace(item.Price), 64)
			if err == nil {
				points += int(math.Ceil(priceFloat * 0.2))
			}
		}
	}

	purchaseDate, err := time.Parse("2006-01-02", strings.TrimSpace(receipt.PurchaseDate))
	if err == nil && purchaseDate.Day()%2 == 1 {
		points += 6
	}

	purchaseTime, err := time.Parse("2006-01-02 15:04", strings.TrimSpace(receipt.PurchaseDate)+" "+strings.TrimSpace(receipt.PurchaseTime))
	if err == nil {

		twoPM, _ := time.Parse("2006-01-02 15:04", strings.TrimSpace(receipt.PurchaseDate)+" "+"14:00")
		fourPM, _ := time.Parse("2006-01-02 15:04", strings.TrimSpace(receipt.PurchaseDate)+" "+"16:00")

		if purchaseTime.After(twoPM) && purchaseTime.Before(fourPM) {
			points += 10
		}
	}

	points += countAlphanumericRunes([]rune(receipt.Retailer)) * 1

	return points
}

func countAlphanumericRunes(chars []rune) int {
	var count int
	for _, r := range chars {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			count += 1
		}
	}
	return count
}

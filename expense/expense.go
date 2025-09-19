package expense

import (
	"errors"
	"regexp"
	"time"
)

type expense struct {
	title        string
	amount       float64
	dateIncurred string
	category     string
	notes        string
	payeeID      int
	receiptURI   string
}

var (
	ErrInvalidTitle      = errors.New("payoutmanagementsystem.NewExpense: title should not be empty")
	ErrInvalidAmount     = errors.New("payoutmanagementsystem.NewExpense: amount must be greater than 0")
	ErrInvalidDate       = errors.New("payoutmanagementsystem.NewExpense: invalid date values or format (YYYY-MM-DD)")
	ErrInvalidCategory   = errors.New("payoutmanagementsystem.NewExpense: category should not be empty")
	ErrInvalidPayeeID    = errors.New("payoutmanagementsystem.NewExpense: payeeID must be positive")
	ErrInvalidReceiptURI = errors.New("payoutmanagementsystem.NewExpense: invalid receipt URI - must be file path")
)

func NewExpense(title string, amount float64, dateIncurred string, category string, notes string, payeeID int, receiptURI string) (*expense, error) {
	if title == "" {
		return nil, ErrInvalidTitle
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if !checkDate(dateIncurred) {
		return nil, ErrInvalidDate
	}
	if category == "" {
		return nil, ErrInvalidCategory
	}
	if payeeID <= 0 {
		return nil, ErrInvalidPayeeID
	}
	if !checkReceiptURI(receiptURI) {
		return nil, ErrInvalidReceiptURI
	}
	return &expense{
		title:        title,
		amount:       amount,
		dateIncurred: dateIncurred,
		category:     category,
		notes:        notes,
		payeeID:      payeeID,
		receiptURI:   receiptURI,
	}, nil
}

func checkDate(dateStr string) bool {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	year := date.Year()
	if year < 2025 || year > 2050 {
		return false
	}

	today := time.Now().Truncate(24 * time.Hour)
	return !date.Before(today)
}

func checkReceiptURI(uri string) bool {
	isMatching, _ := regexp.MatchString(`^/`, uri)
	return isMatching
}

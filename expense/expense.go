package expense

import (
	"errors"
	"regexp"
)

type expense struct {
	title        string
	amount       float64
	dateIncurred string
	category     string
	notes        string
	payeeID      int
	receiptURI   string
	status       string
}

var ErrInvalidTitle = errors.New("payoutmanagementsystem.NewExpense: title should not be empty")
var ErrInvalidAmount = errors.New("payoutmanagementsystem.NewExpense: amount must be greater than 0")
var ErrInvalidDate = errors.New("payoutmanagementsystem.NewExpense: invalid date values or format (YYYY-MM-DD)")
var ErrInvalidCategory = errors.New("payoutmanagementsystem.NewExpense: category should not be empty")
var ErrInvalidPayeeID = errors.New("payoutmanagementsystem.NewExpense: payeeID must be positive")
var ErrInvalidReceiptURI = errors.New("payoutmanagementsystem.NewExpense: invalid receipt URI - must be file path")

func NewExpense(title string, amount float64, dateIncurred string, category string, notes string, payeeID int, receiptURI string) (*expense, error) {
	if title == "" {
		return nil, ErrInvalidTitle
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if !checkDateFormat(dateIncurred) {
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
		status:       "Pending",
	}, nil
}

func checkDateFormat(date string) bool {
	pattern := `^(202[5-9]|20(3\d|4\d|50))-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`
	match, _ := regexp.MatchString(pattern, date)
	return match
}

func checkReceiptURI(uri string) bool {
	match2, _ := regexp.MatchString(`^/`, uri)
	return match2
}

func (e *expense) Title() string {
	return e.title
}

func (e *expense) Amount() float64 {
	return e.amount
}

func (e *expense) DateIncurred() string {
	return e.dateIncurred
}

func (e *expense) Category() string {
	return e.category
}

func (e *expense) Notes() string {
	return e.notes
}

func (e *expense) PayeeID() int {
	return e.payeeID
}

func (e *expense) ReceiptURI() string {
	return e.receiptURI
}

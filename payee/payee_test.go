package payee

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var errTests = []struct {
	testName        string
	errMsg          string
	beneficiaryName string
	beneficiaryCode string
	accNo           int
	ifsc            string
	bankName        string
	email           string
	mobile          int
	payeeCategory   string
	expectedErr     error
}{
	{"TestAccountNumberWith9DigitsIsInvalid", "invalid account number: 9 digits is used", "abc", "123", 678000234, "CBIN0564891", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"TestAccountNumberWith11DigitsIsInvalid", "invalid account number: 11 digits is used", "abc", "123", 67800023445, "CBIN0564891", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"TestAccountNumberWith15DigitsIsInvalid", "invalid account number: 15 digits is used", "abc", "123", 678000234576543, "CBIN0564891", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"TestInvalidAccountNumberOfLength17", "invalid account number: 17 digits is used", "abc", "123", 67800023457654324, "CBIN0564891", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"TestInvalidMobileNumberOfLength9", "invalid mobile number: 9 digits is used", "abc", "123", 6780002345765432, "CBIN0564891", "cbi", "abc@gmail.com", 912345678, "Employee", ErrInvalidMobileNumber},
	{"TestInvalidMobileNumberOfLength11", "invalid mobile number: 11 digits is used", "abc", "123", 6780002345765432, "CBIN0564891", "cbi", "abc@gmail.com", 91234567891, "Employee", ErrInvalidMobileNumber},
	{"TestInvalidEmailMissingAtSymbol", "invalid email: @ symbol is not used", "abc", "123", 6780002345765432, "CBIN0564891", "cbi", "abc.com", 9123456782, "Employee", ErrInvalidEmail},
	{"TestNameEmptyIsInvalid", "invalid name: name is empty", "", "123", 6780002345765432, "CBIN0564891", "cbi", "abc@gmail.com", 9123456782, "Employee", ErrEmptyName},
	{"TestCodeEmptyReturnsErrEmptyCode", "invalid code: code is empty", "abc", "", 6780002345765432, "CBIN0564891", "cbi", "abc@gmail.com", 9123456782, "Employee", ErrEmptyCode},
	{"TestIFSCMissingNumeralsReturnsErrInvalidIFSC", "invalid ifsc: the numerals are missed", "abc", "123", 6700345678, "CBIN0789", "cbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidIFSC},
	{"TestIFSCContainsLowercaseLettersReturnsErrInvalidIFSC", "invalid ifsc: there are lowercase alphabets", "abc", "123", 6700345678, "cbin0456671", "cbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidIFSC},
	{"TestIFSCBranchCodeContainsAlphabetsReturnsErrInvalidIFSC", "invalid ifsc: the alphabets is used as part of branch code", "abc", "123", 6700345678, "CBIN0456ab1", "cbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidIFSC},
	{"TestInvalidBankNameOfLengthGreaterThan50", "invalid bank name: bank name exceeds 50 characters", "abc", "123", 6700345678, "CBIN0456671", "cbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidBankName},
	{"TestInvalidPayeeCategory", "invalid payee category: given category is not listed", "abc", "123", 6700543678, "CBIN0123451", "cbi", "abc@gmail.com", 9123456780, "Student", ErrInvalidCategory},
}

func TestInvalidPayee(t *testing.T) {
	for _, tt := range errTests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := NewPayee(tt.beneficiaryName, tt.beneficiaryCode, tt.accNo, tt.ifsc, tt.bankName, tt.email, tt.mobile, tt.payeeCategory)
			assert.ErrorIs(t, err, tt.expectedErr, "Error Test Case: %v", tt.errMsg)

		})
	}
}

func TestValidPayee(t *testing.T) {
	name := "abc"
	code := "123"
	accNo := 1234567890
	bankIFSC := "CBIN0123456"
	bankName := "cbi"
	emailID := "abc@gmail.com"
	mobile := 9876543210
	category := "Employee"
	p, err := NewPayee(name, code, accNo, bankIFSC, bankName, emailID, mobile, category)
	assert.NoError(t, err, "payee should be created")

	assert.Equal(t, name, p.beneficiaryName, "name mismatch")
	assert.Equal(t, code, p.beneficiaryCode, "beneficiary code mismatch")
	assert.Equal(t, accNo, p.accNo, "account number mismatch")
	assert.Equal(t, bankIFSC, p.ifsc, "IFSC mismatch")
	assert.Equal(t, bankName, p.bankName, "bank name mismatch")
	assert.Equal(t, emailID, p.email, "email mismatch")
	assert.Equal(t, mobile, p.mobile, "mobile mismatch")
	assert.Equal(t, category, p.payeeCategory, "category mismatch")
}

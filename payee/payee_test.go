package payoutmanagementsystem

import (
	"testing"
)

var errTests = []struct {
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
	{"abc", "123", 678000234, "CBIN056489", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"abc", "123", 67800023445, "CBIN056489", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"abc", "123", 678000234576543, "CBIN056489", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"abc", "123", 67800023457654324, "CBIN056489", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"abc", "123", 6780002345765432, "CBIN056489", "cbi", "abc@gmail.com", 912345678, "Employee", ErrInvalidMobileNumber},
	{"abc", "123", 6780002345765432, "CBIN056489", "cbi", "abc@gmail.com", 91234567891, "Employee", ErrInvalidMobileNumber},
	{"abc", "123", 6780002345765432, "CBIN056489", "cbi", "abc.com", 9123456782, "Employee", ErrInvalidEmail},
	{"", "123", 6780002345765432, "CBIN056489", "cbi", "abc@gmail.com", 9123456782, "Employee", ErrEmptyName},
	{"abc", "", 6780002345765432, "CBIN056489", "cbi", "abc@gmail.com", 9123456782, "Employee", ErrEmptyCode},
	{"abc", "123", 6700345678, "CBIN0789", "cbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidIFSC},
	{"abc", "123", 6700345678, "cbin045667", "cbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidIFSC},
	{"abc", "123", 6700345678, "CBIN0456ab", "cbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidIFSC},
	{"abc", "123", 6700345678, "CBIN045667", "cbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidBankName},
}

func TestValidatePayee(t *testing.T) {
	for _, tt := range errTests {
		_, err := NewPayee(tt.beneficiaryName, tt.beneficiaryCode, tt.accNo, tt.ifsc, tt.bankName, tt.email, tt.mobile, tt.payeeCategory)
		if err != tt.expectedErr {
			t.Fatalf("Expected Error: %v but Actual Error: %v", tt.expectedErr, err)
		}
	}
}

var validPayeeFieldsTests = []struct {
	beneficiaryName string
	beneficiaryCode string
	accNo           int
	ifsc            string
	bankName        string
	email           string
	mobile          int
	payeeCategory   string
}{
	{"abc", "123", 6780003045, "CBIN056489", "cbi", "abc@gmail.com", 9123456789, "Employee"},
	{"abc", "123", 6780002345765432, "CBIN056489", "cbi", "abc@gmail.com", 9123456789, "Employee"},
}

func TestPayeeCanBeCreatedWithValidValues(t *testing.T) {
	for _, tt := range validPayeeFieldsTests {
		_, err := NewPayee(tt.beneficiaryName, tt.beneficiaryCode, tt.accNo, tt.ifsc, tt.bankName, tt.email, tt.mobile, tt.payeeCategory)
		if err != nil {
			t.Fatalf("Payee should be created but got error: %v", err)
		}
	}
}

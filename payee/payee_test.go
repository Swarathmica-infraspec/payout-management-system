package payee

import (
	"testing"
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
	{"TestAccountNumberWith9DigitsIsInvalid", "invalid account number: 9 digits is used", "abc", "123", 678000234, "CBIN056489", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"TestAccountNumberWith11DigitsIsInvalid", "invalid account number: 11 digits is used", "abc", "123", 67800023445, "CBIN056489", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"TestAccountNumberWith15DigitsIsInvalid", "invalid account number: 15 digits is used", "abc", "123", 678000234576543, "CBIN056489", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"TestInvalidAccountNumberOfLength17", "invalid account number: 17 digits is used", "abc", "123", 67800023457654324, "CBIN056489", "cbi", "abc@gmail.com", 9123456789, "Employee", ErrInvalidAccountNumber},
	{"TestInvalidMobileNumberOfLength9", "invalid mobile number: 9 digits is used", "abc", "123", 6780002345765432, "CBIN056489", "cbi", "abc@gmail.com", 912345678, "Employee", ErrInvalidMobileNumber},
	{"TestInvalidMobileNumberOfLength11", "invalid mobile number: 11 digits is used", "abc", "123", 6780002345765432, "CBIN056489", "cbi", "abc@gmail.com", 91234567891, "Employee", ErrInvalidMobileNumber},
	{"TestInvalidEmailMissingAtSymbol", "invalid email: @ symbol is not used", "abc", "123", 6780002345765432, "CBIN056489", "cbi", "abc.com", 9123456782, "Employee", ErrInvalidEmail},
	{"TestNameEmptyIsInvalid", "invalid name: name is empty", "", "123", 6780002345765432, "CBIN056489", "cbi", "abc@gmail.com", 9123456782, "Employee", ErrEmptyName},
	{"TestCodeEmptyReturnsErrEmptyCode", "invalid code: code is empty", "abc", "", 6780002345765432, "CBIN056489", "cbi", "abc@gmail.com", 9123456782, "Employee", ErrEmptyCode},
	{"TestIFSCMissingNumeralsReturnsErrInvalidIFSC", "invalid ifsc: the numerals are missed", "abc", "123", 6700345678, "CBIN0789", "cbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidIFSC},
	{"TestIFSCContainsLowercaseLettersReturnsErrInvalidIFSC", "invalid ifsc: there are lowercase alphabets", "abc", "123", 6700345678, "cbin045667", "cbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidIFSC},
	{"TestIFSCBranchCodeContainsAlphabetsReturnsErrInvalidIFSC", "invalid ifsc: the alphabets is used as part of branch code", "abc", "123", 6700345678, "CBIN0456ab", "cbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidIFSC},
	{"TestInvalidBankNameOfLengthGreaterThan50", "invalid bank name: bank name exceeds 50 characters", "abc", "123", 6700345678, "CBIN045667", "cbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbicbi", "abc@gmail.com", 9123456666, "Employee", ErrInvalidBankName},
}

func TestValidatePayee(t *testing.T) {
	for _, tt := range errTests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := NewPayee(tt.beneficiaryName, tt.beneficiaryCode, tt.accNo, tt.ifsc, tt.bankName, tt.email, tt.mobile, tt.payeeCategory)
			if err != tt.expectedErr {
				t.Fatalf("Error Test Case: %v , Expected Error: %v but Actual Error: %v", tt.errMsg, tt.expectedErr, err)
			}
		})
	}
}

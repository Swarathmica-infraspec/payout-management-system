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
			if err != tt.expectedErr {
				t.Fatalf("Error Test Case: %v , Expected Error: %v but Actual Error: %v", tt.errMsg, tt.expectedErr, err)
			}
		})
	}
}

func TestValidPayee(t *testing.T) {

	_, err := NewPayee("abc", "123", 1234567890, "CBIN0123456", "cbi", "abc@gmail.com", 9876543210, "Employee")
	if err != nil {
		t.Fatalf("payee should be created but got error: %v", err)
	}

	name := "abc"
	code := "123"
	accNo := 1234567890
	bankIFSC := "CBIN0123456"
	bankName := "cbi"
	emailID := "abc@gmail.com"
	mobile := 9876543210
	category := "Employee"
	p, err := NewPayee(name, code, accNo, bankIFSC, bankName, emailID, mobile, category)
	if err != nil {
		t.Fatalf("payee should be created but got error: %v", err)
	}

	if p.beneficiaryName != name {
		t.Errorf("expected name: %v but stored name: %v", name, p.beneficiaryName)
	}

	if p.beneficiaryCode != code {
		t.Errorf("expected beneficiary code: %v but stored code: %v", code, p.beneficiaryCode)
	}

	if p.accNo != accNo {
		t.Errorf("expected account Number: %v but stored account number: %v", accNo, p.accNo)
	}

	if p.ifsc != bankIFSC {
		t.Errorf("expected IFSC: %v but stored IFSC: %v", bankIFSC, p.ifsc)
	}

	if p.bankName != bankName {
		t.Errorf("expected bank name: %v but stored bank name: %v", bankName, p.bankName)
	}

	if p.email != emailID {
		t.Errorf("expected emailID: %v but stored email ID: %v", emailID, p.email)
	}

	if p.mobile != mobile {
		t.Errorf("expected mobile number: %v but stored mobile number: %v", mobile, p.mobile)
	}

	if p.payeeCategory != category {
		t.Errorf("expected category: %v but stored category: %v", category, p.payeeCategory)
	}

}

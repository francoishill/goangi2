package emailUtils

import (
	"fmt"
)

type EmailRecipient struct {
	FullName     string
	EmailAddress string
}

func (this *EmailRecipient) AsArray() []*EmailRecipient {
	return []*EmailRecipient{this}
}

func (this *EmailRecipient) formatFullnameAndAddress() string {
	return fmt.Sprintf("%s<%s>", this.FullName, this.EmailAddress)
}

func CreateEmailRecipient(fullName, emailAddress string) *EmailRecipient {
	return &EmailRecipient{
		FullName:     fullName,
		EmailAddress: emailAddress,
	}
}

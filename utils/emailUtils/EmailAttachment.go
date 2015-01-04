package emailUtils

type EmailAttachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

func CreateNewEmailAttachment(fileName string, data []byte) *EmailAttachment {
	return &EmailAttachment{
		Filename: fileName,
		Data:     data,
		Inline:   false, //Inline ALWAYS false as it does not work as expected, gets messed up in gmail and outlook shows as attachment but ends up being empty
	}
}

package emailUtils

import (
	"database/sql"
	"strings"
	"time"

	. "github.com/francoishill/goangi2/utils/entityUtils"
)

const (
	MAIL_QUEUE__TABLE_NAME                   = "queued_email"
	MAIL_QUEUE__COLUMN__SUCCESSFULLY_SENT_ON = "SuccessfullySentOn"
	MAIL_QUEUE__COLUMN__SEND_ERROR           = "SendError"
)

func CreateFieldsToLoadInQueuedEmail() *RelatedFieldsToLoad {
	relatedFields := CreateRelatedFieldsToLoad()
	return relatedFields
}

func CreateQueuedEmail_FromEmailMessage(emailMessage *EmailMessage) *QueuedEmail {
	onlyEmailAddress := true
	return &QueuedEmail{
		FromEmail:          emailMessage.From.EmailAddress,
		RecipientEmailsCsv: emailMessage.getCommaSeparatedFormattedToAddresses(onlyEmailAddress),
		Content:            emailMessage.GetContent(),
		SendDueTime:        emailMessage.SendDueTime,
		// SuccessfullySentOn: nil, //Leave blank to keep as null
		SendError: sql.NullString{Valid: false}, //Set NULL
		DebugInfo: emailMessage.DebugInfo,
	}
}

type QueuedEmail struct {
	Id                 int64
	FromEmail          string         `orm:"type(text)"` //Only the email address not the person name
	RecipientEmailsCsv string         `orm:"type(text)"` //Only their email addresses not the person names
	Content            string         `orm:"type(text)"`
	SendDueTime        time.Time      `orm:"type(datetime)"`
	SuccessfullySentOn time.Time      `orm:"null;type(datetime)"` //Default null
	SendError          sql.NullString `orm:"null"`                //When an error occurs in sending
	DebugInfo          string
}

func (this *QueuedEmail) ReadUsingID(ormContext *OrmContext, id int64, loadRelatedSettings *RelatedFieldsToLoad) {
	this.Id = id
	OrmRepo.BaseReadEntityUsingPK(ormContext, this, loadRelatedSettings)
}

func (this *QueuedEmail) Insert(ormContext *OrmContext) {
	OrmRepo.BaseInsertEntity(ormContext, this)
}

func (this *QueuedEmail) update(ormContext *OrmContext, modifiedEntity interface{}, onlyAllowTheseFieldsToSave ...string) []ChangedField {
	return OrmRepo.BaseUpdateEntity(ormContext, this, modifiedEntity, onlyAllowTheseFieldsToSave...)
}

func (this *QueuedEmail) Update_MarkAsSuccessfullySent(ormContext *OrmContext) []ChangedField {
	modifiedEntity := &QueuedEmail{
		Id:                 this.Id,
		SuccessfullySentOn: time.Now(),
		SendError:          sql.NullString{Valid: false}, //Set NULL
	}
	return this.update(ormContext, modifiedEntity, MAIL_QUEUE__COLUMN__SUCCESSFULLY_SENT_ON, MAIL_QUEUE__COLUMN__SEND_ERROR)
}

func (this *QueuedEmail) Update_SetSendError(ormContext *OrmContext, errorString string) []ChangedField {
	modifiedEntity := &QueuedEmail{
		Id:        this.Id,
		SendError: sql.NullString{Valid: true, String: errorString},
	}
	return this.update(ormContext, modifiedEntity, MAIL_QUEUE__COLUMN__SEND_ERROR)
}

/*func (this *QueuedEmail) Delete(ormContext *OrmContext) {
	OrmRepo.BaseDeleteEntity(ormContext, this)
}*/

func (this *QueuedEmail) List_Unsent(ormContext *OrmContext, loadRelatedSettings *RelatedFieldsToLoad, orderByFields ...string) []*QueuedEmail {
	fieldFilters := []map[string]interface{}{
		map[string]interface{}{
			MAIL_QUEUE__COLUMN__SUCCESSFULLY_SENT_ON + "__isnull": true,
		},
	}

	list := []*QueuedEmail{}
	OrmRepo.BaseListEntities_ANDFilters_OrderBy(ormContext, MAIL_QUEUE__TABLE_NAME, fieldFilters, orderByFields, loadRelatedSettings, &list)
	return list
}

func (this *QueuedEmail) GetSplittedRecipientEmailAddresses() []string {
	return strings.Split(this.RecipientEmailsCsv, ",")
}

func (u *QueuedEmail) TableEngine() string     { return "INNODB" }
func (u *QueuedEmail) TableName() string       { return MAIL_QUEUE__TABLE_NAME }
func (u *QueuedEmail) TableIndex() [][]string  { return [][]string{} }
func (u *QueuedEmail) TableUnique() [][]string { return [][]string{} }

func init() {
	DefaultRegisterModel(new(QueuedEmail))
}

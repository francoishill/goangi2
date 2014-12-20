package oauth2Utils

import (
	"time"

	. "github.com/francoishill/goangi2/utils/entityUtils"
)

const (
	OAUTH2_AUTHORIZE_TABLE_NAME         = "oauth2_authorize"
	OAUTH2_AUTHORIZE_CLIENT_COLUMN_NAME = "Client"
	OAUTH2_AUTHORIZE_CODE_COLUMN_NAME   = "Code"
)

type OAuth2Authorize struct {
	Id          int64
	ClientId    int64  //Client        *OAuth2Client    `orm:"rel(fk)"`
	Code        string //The authorization code
	ExpiresIn   int32
	Scope       string `orm:"type(text)"`
	RedirectUri string
	State       string
	CreatedAt   time.Time
}

func (this *OAuth2Authorize) ReadUsingID(ormContext *OrmContext, id int64, loadRelatedSettings *RelatedFieldsToLoad) {
	this.Id = id
	OrmRepo.BaseReadEntityUsingPK(ormContext, this, loadRelatedSettings)
}

func (this *OAuth2Authorize) ReadUsingCode(ormContext *OrmContext, code string, loadRelatedSettings *RelatedFieldsToLoad) bool {
	this.Code = code
	return OrmRepo.BaseReadEntityUsingFields(ormContext, this, loadRelatedSettings, OAUTH2_AUTHORIZE_CODE_COLUMN_NAME)
}

func (this *OAuth2Authorize) Insert(ormContext *OrmContext) {
	OrmRepo.BaseInsertEntity(ormContext, this)
}

func (this *OAuth2Authorize) Delete(ormContext *OrmContext) {
	OrmRepo.BaseDeleteEntity(ormContext, this)
}

func (u *OAuth2Authorize) TableEngine() string {
	return "INNODB"
}

func (u *OAuth2Authorize) TableName() string {
	return OAUTH2_AUTHORIZE_TABLE_NAME
}

func (u *OAuth2Authorize) TableIndex() [][]string {
	return [][]string{}
}

func (u *OAuth2Authorize) TableUnique() [][]string {
	return [][]string{
		[]string{OAUTH2_AUTHORIZE_CODE_COLUMN_NAME},
	}
}

func init() {
	DefaultRegisterModel(new(OAuth2Authorize))
}

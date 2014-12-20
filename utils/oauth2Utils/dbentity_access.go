package oauth2Utils

import (
	"time"

	. "github.com/francoishill/goangi2/utils/entityUtils"
)

const (
	OAUTH2_ACCESS_TABLE_NAME                 = "oauth2_access"
	OAUTH2_ACCESS_AUTHORIZE_DATA_COLUMN_NAME = "AuthorizeData"
	OAUTH2_ACCESS_ACCESS_DATA_COLUMN_NAME    = "AccessData"
	OAUTH2_ACCESS_ACCESS_TOKEN_COLUMN_NAME   = "AccessToken"
	OAUTH2_ACCESS_REFRESH_TOKEN_COLUMN_NAME  = "RefreshToken"
)

var (
	relatedAuthorizeDataInOAuth2Access *RelatedField = CreateRelatedField(OAUTH2_ACCESS_AUTHORIZE_DATA_COLUMN_NAME, true)
	relatedAccessDataInOAuth2Access    *RelatedField = CreateRelatedField(OAUTH2_ACCESS_ACCESS_DATA_COLUMN_NAME, true)
)

func CreateFieldsToLoadInOAuth2Access(loadAuthorizeData, loadAccessData bool) *RelatedFieldsToLoad {
	relatedFields := CreateRelatedFieldsToLoad()
	relatedFields.AppendIfTrue(loadAuthorizeData, relatedAuthorizeDataInOAuth2Access)
	relatedFields.AppendIfTrue(loadAccessData, relatedAccessDataInOAuth2Access)
	return relatedFields
}

type OAuth2Access struct {
	Id            int64
	ClientId      int64            //Client        *OAuth2Client    `orm:"rel(fk)"`
	UserId        int64            //Store the ID so we can keep it generic
	AuthorizeData *OAuth2Authorize `orm:"rel(fk);null"`
	AccessData    *OAuth2Access    `orm:"rel(fk);on_delete(set_null);null"` //Previous access data, for refresh token. It must only set_null on_delete otherwise the workflow fails
	AccessToken   string
	RefreshToken  string
	ExpiresIn     int32
	Scope         string `orm:"type(text)"`
	RedirectUri   string
	CreatedAt     time.Time

	User IExpectedUser `orm:"-"` //For now lets not allow NULL users as if we want to create our own 'bot' accessing the data, we can also create a user for it. And perhaps have a boolean flag for user entities called 'bot'?
}

func (this *OAuth2Access) ReadUsingID(ormContext *OrmContext, id int64, loadRelatedSettings *RelatedFieldsToLoad) {
	this.Id = id
	OrmRepo.BaseReadEntityUsingPK(ormContext, this, loadRelatedSettings)
}

func (this *OAuth2Access) ReadUsingAccessToken(ormContext *OrmContext, accessToken string, loadRelatedSettings *RelatedFieldsToLoad) bool {
	this.AccessToken = accessToken
	return OrmRepo.BaseReadEntityUsingFields(ormContext, this, loadRelatedSettings, OAUTH2_ACCESS_ACCESS_TOKEN_COLUMN_NAME)
}

func (this *OAuth2Access) ReadUsingRefreshToken(ormContext *OrmContext, refreshToken string, loadRelatedSettings *RelatedFieldsToLoad) bool {
	this.RefreshToken = refreshToken
	return OrmRepo.BaseReadEntityUsingFields(ormContext, this, loadRelatedSettings, OAUTH2_ACCESS_REFRESH_TOKEN_COLUMN_NAME)
}

func (this *OAuth2Access) Insert(ormContext *OrmContext) {
	if this.UserId == 0 && this.User != nil {
		this.UserId = this.User.GetId()
	}
	OrmRepo.BaseInsertEntity(ormContext, this)
}

func (this *OAuth2Access) Delete(ormContext *OrmContext) {
	OrmRepo.BaseDeleteEntity(ormContext, this)
}

func (u *OAuth2Access) TableEngine() string {
	return "INNODB"
}

func (u *OAuth2Access) TableName() string {
	return OAUTH2_ACCESS_TABLE_NAME
}

func (u *OAuth2Access) TableIndex() [][]string {
	return [][]string{}
}

func (u *OAuth2Access) TableUnique() [][]string {
	return [][]string{
		[]string{OAUTH2_ACCESS_ACCESS_TOKEN_COLUMN_NAME},
		[]string{OAUTH2_ACCESS_REFRESH_TOKEN_COLUMN_NAME},
	}
}

func init() {
	DefaultRegisterModel(new(OAuth2Access))
}

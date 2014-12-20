package oauth2Utils

var predefinedOauthClients_ByClientId = map[string]*OAuth2Client_NonDb{}
var predefinedOauthClients_ById = map[int64]*OAuth2Client_NonDb{}

type OAuth2Client_NonDb struct {
	Id                int64 //Do not get confused with the ClientId field which is a string
	ClientId          string
	ClientSecret      string
	RedirectUri       string
	ClientDisplayName string

	//Created time.Time `orm:"auto_now_add"`
	//Updated time.Time `orm:"auto_now"`
}

func GetClientUsingClientId(clientId string) (*OAuth2Client_NonDb, bool) {
	if client, ok := predefinedOauthClients_ByClientId[clientId]; ok {
		return client, true
	} else {
		return nil, false
	}
}

func GetClientUsingId(id int64) (*OAuth2Client_NonDb, bool) {
	if client, ok := predefinedOauthClients_ById[id]; ok {
		return client, true
	} else {
		return nil, false
	}
}

var tmpIdIncrement int64 = 1

func AddPredefinedOAuthClient(client *OAuth2Client_NonDb) {
	client.Id = tmpIdIncrement
	predefinedOauthClients_ByClientId[client.ClientId] = client
	predefinedOauthClients_ById[client.Id] = client
	tmpIdIncrement++
}

/*func init() {
	DefaultRegisterModel(new(OAuth2Client_NonDb))
}*/

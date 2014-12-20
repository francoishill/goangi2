package oauth2Utils

//This file mainly deals with translating our OAuth2Authorize, OAuth2Access and OAuth2Client_NonDb into
//osin.AuthorizeData, osin.AccessData and osin.Client
import (
	"github.com/RangelReale/osin"
)

func ConvertIntoOsinClient(client *OAuth2Client_NonDb) *osin.Client {
	if client == nil {
		return nil
	}

	return &osin.Client{
		Id:          client.ClientId,
		Secret:      client.ClientSecret,
		RedirectUri: client.RedirectUri,
		//UserData:    &client, NOT using this
	}
}

func ConvertFromOsinAuthorize(osinAuthorize *osin.AuthorizeData) *OAuth2Authorize {
	var client *OAuth2Client_NonDb = nil
	if osinAuthorize.Client != nil && osinAuthorize.Client.Id != "" {
		if tmpClient, found := GetClientUsingClientId(osinAuthorize.Client.Id); found {
			client = tmpClient
		}
	}

	return &OAuth2Authorize{
		ClientId:    client.Id,
		Code:        osinAuthorize.Code,
		ExpiresIn:   osinAuthorize.ExpiresIn,
		Scope:       osinAuthorize.Scope,
		RedirectUri: osinAuthorize.RedirectUri,
		State:       osinAuthorize.State,
		CreatedAt:   osinAuthorize.CreatedAt,
	}
}

func ConvertIntoOsinAuthorize(authorize *OAuth2Authorize) *osin.AuthorizeData {
	if authorize == nil {
		return nil
	}

	client, foundClient := GetClientUsingId(authorize.ClientId)
	if !foundClient {
		return nil
	}

	return &osin.AuthorizeData{
		Client:      ConvertIntoOsinClient(client),
		Code:        authorize.Code,
		ExpiresIn:   authorize.ExpiresIn,
		Scope:       authorize.Scope,
		RedirectUri: authorize.RedirectUri,
		State:       authorize.State,
		CreatedAt:   authorize.CreatedAt,
		//UserData:    authorize, NOT using this
	}
}

func ConvertFromOsinAccess(osinAccess *osin.AccessData) *OAuth2Access {
	var client *OAuth2Client_NonDb = nil
	if osinAccess.Client != nil && osinAccess.Client.Id != "" {
		if tmpClient, found := GetClientUsingClientId(osinAccess.Client.Id); found {
			client = tmpClient
		}
	}

	var authorize *OAuth2Authorize = nil
	if osinAccess.AuthorizeData != nil && osinAccess.AuthorizeData.Code != "" {
		authorize = &OAuth2Authorize{}
		if existed := authorize.ReadUsingCode(nil, osinAccess.AuthorizeData.Code, nil); !existed {
			authorize = nil
		}
	}

	var previousAccessPointer *OAuth2Access = nil
	if osinAccess.AccessData != nil && osinAccess.AccessData.AccessToken != "" {
		previousAccessPointer = &OAuth2Access{}
		if existed := previousAccessPointer.ReadUsingAccessToken(nil, osinAccess.AccessData.AccessToken, nil); !existed {
			previousAccessPointer = nil
		}
	}

	return &OAuth2Access{
		ClientId:      client.Id,
		AuthorizeData: authorize,
		AccessData:    previousAccessPointer,
		AccessToken:   osinAccess.AccessToken,
		RefreshToken:  osinAccess.RefreshToken,
		ExpiresIn:     osinAccess.ExpiresIn,
		Scope:         osinAccess.Scope,
		RedirectUri:   osinAccess.RedirectUri,
		CreatedAt:     osinAccess.CreatedAt,
	}
}

func recuresiveConvertIntoOsinAccess(authorize *OAuth2Access, currentLevel int) *osin.AccessData {
	if authorize == nil {
		return nil
	}

	var previousAccess *osin.AccessData = nil
	if currentLevel == 0 && authorize.AccessData != nil && authorize.AccessData.AccessToken != "" {
		//Check for currentLevel == 0 because we do not want to retrieve possible nested AccessTokens
		previousAccess = recuresiveConvertIntoOsinAccess(authorize.AccessData, currentLevel+1)
	}

	client, foundClient := GetClientUsingId(authorize.ClientId)
	if !foundClient {
		return nil
	}

	return &osin.AccessData{
		Client:        ConvertIntoOsinClient(client),
		AuthorizeData: ConvertIntoOsinAuthorize(authorize.AuthorizeData),
		AccessData:    previousAccess,
		AccessToken:   authorize.AccessToken,
		RefreshToken:  authorize.RefreshToken,
		ExpiresIn:     authorize.ExpiresIn,
		Scope:         authorize.Scope,
		RedirectUri:   authorize.RedirectUri,
		CreatedAt:     authorize.CreatedAt,
		UserData:      authorize.User,
	}
}

func ConvertIntoOsinAccess(authorize *OAuth2Access) *osin.AccessData {
	return recuresiveConvertIntoOsinAccess(authorize, 0)
}

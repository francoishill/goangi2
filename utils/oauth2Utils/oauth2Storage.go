package oauth2Utils

import (
	"errors"
	"fmt"
	"github.com/RangelReale/osin"
)

type OAuth2Storage struct {
	// clients map[string]*osin.DefaultClient
	// authorize map[string]*osin.AuthorizeData
	// access    map[string]*osin.AccessData
	// refresh   map[string]string
}

func NewOAuth2Storage() *OAuth2Storage {
	return &OAuth2Storage{}
}

func (s *OAuth2Storage) Clone() osin.Storage {
	return s
}

func (s *OAuth2Storage) Close() {
	//Do nothing
}

func (s *OAuth2Storage) GetClient(clientId string) (returnClient osin.Client, returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			returnClient = nil
			returnErr = extractErrorFromRecoveryObject(r)
		}
	}()

	if client, found := GetClientUsingClientId(clientId); !found {
		panic("Unable to find OAuth Client with ID: " + clientId)
	} else {
		returnClient = ConvertIntoOsinClient(client)
		returnErr = nil
		return
	}
}

func (s *OAuth2Storage) SaveAuthorize(data *osin.AuthorizeData) (returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			returnErr = extractErrorFromRecoveryObject(r)
		}
	}()

	authorize := ConvertFromOsinAuthorize(data)

	authorize.Insert(nil)
	returnErr = nil
	return
}

func (s *OAuth2Storage) LoadAuthorize(code string) (data *osin.AuthorizeData, returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			data = nil
			returnErr = extractErrorFromRecoveryObject(r)
		}
	}()

	authorize := &OAuth2Authorize{}
	if foundAuthorize := authorize.ReadUsingCode(nil, code, nil); !foundAuthorize {
		panic("Unable to find OAuth Authorize with code: " + code)
	}

	data = ConvertIntoOsinAuthorize(authorize)
	returnErr = nil
	return
}

func (s *OAuth2Storage) RemoveAuthorize(code string) (returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			returnErr = extractErrorFromRecoveryObject(r)
		}
	}()

	authorize := &OAuth2Authorize{}
	if foundAuthorize := authorize.ReadUsingCode(nil, code, nil); !foundAuthorize {
		return nil //It has already been removed
	}

	authorize.Delete(nil)
	returnErr = nil
	return
}

func (s *OAuth2Storage) SaveAccess(data *osin.AccessData) (returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			returnErr = extractErrorFromRecoveryObject(r)
		}
	}()

	if data.UserData == nil {
		panic(errorMapKeys[E_ACCESS_EMPTY_USER_ERROR])
	}

	usr, ok := data.UserData.(IExpectedUser)
	if !ok {
		panic(errorMapKeys[E_ACCESS_EMPTY_USER_ERROR])
	}

	access := ConvertFromOsinAccess(data)

	access.UserId = usr.GetId()
	access.User = usr
	access.Insert(nil)

	returnErr = nil
	return
}

type tmpUserStruct struct {
	Id int64
}

func (this *tmpUserStruct) IAmAUser() {
}
func (this *tmpUserStruct) GetRands() string {
	//We have not implemented this yet, perhaps we should never trigger this workflow? Otherwise this must be implemented
	panic("Internal server error 600002")
}
func (this *tmpUserStruct) GetPassword() string {
	//We have not implemented this yet, perhaps we should never trigger this workflow? Otherwise this must be implemented
	panic("Internal server error 600003")
}
func (this *tmpUserStruct) GetId() int64 {
	return this.Id
}

func (s *OAuth2Storage) LoadAccess(accessToken string) (data *osin.AccessData, returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			data = nil
			returnErr = extractErrorFromRecoveryObject(r)
		}
	}()

	access := &OAuth2Access{}
	loadAuthorizeData := false
	loadAccessData := false
	loadRelatedSettings := CreateFieldsToLoadInOAuth2Access(loadAuthorizeData, loadAccessData)
	if foundAccessToken := access.ReadUsingAccessToken(nil, accessToken, loadRelatedSettings); !foundAccessToken {
		panic("Unable to find OAuth AccessToken: " + accessToken)
	}

	data = ConvertIntoOsinAccess(access)
	data.UserData = IExpectedUser(&tmpUserStruct{Id: access.UserId})
	returnErr = nil
	return
}

func (s *OAuth2Storage) RemoveAccess(accessToken string) (returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			returnErr = extractErrorFromRecoveryObject(r)
		}
	}()

	access := &OAuth2Access{}
	loadAuthorizeData := false
	loadAccessData := false
	loadRelatedSettings := CreateFieldsToLoadInOAuth2Access(loadAuthorizeData, loadAccessData)
	if foundAuthorize := access.ReadUsingAccessToken(nil, accessToken, loadRelatedSettings); !foundAuthorize {
		return nil //It has already been removed
	}
	access.Delete(nil)
	returnErr = nil
	return
}

func (s *OAuth2Storage) LoadRefresh(refreshToken string) (data *osin.AccessData, returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			data = nil
			returnErr = extractErrorFromRecoveryObject(r)
		}
	}()

	access := &OAuth2Access{}
	loadAuthorizeData := false
	loadAccessData := false
	loadRelatedSettings := CreateFieldsToLoadInOAuth2Access(loadAuthorizeData, loadAccessData)
	if foundAccessToken := access.ReadUsingRefreshToken(nil, refreshToken, loadRelatedSettings); !foundAccessToken {
		panic("Unable to find OAuth RefreshToken: " + refreshToken)
	}

	data = ConvertIntoOsinAccess(access)
	returnErr = nil
	return
}

func (s *OAuth2Storage) RemoveRefresh(refreshToken string) (returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			returnErr = extractErrorFromRecoveryObject(r)
		}
	}()

	access := &OAuth2Access{}
	loadAuthorizeData := false
	loadAccessData := false
	loadRelatedSettings := CreateFieldsToLoadInOAuth2Access(loadAuthorizeData, loadAccessData)
	if foundAuthorize := access.ReadUsingRefreshToken(nil, refreshToken, loadRelatedSettings); !foundAuthorize {
		return nil //It has already been removed
	}
	access.Delete(nil)

	returnErr = nil
	return
}

func extractErrorFromRecoveryObject(recoveryObj interface{}) error {
	switch e := recoveryObj.(type) {
	case string:
		return errors.New(e)
	case error:
		return e
	default:
		return fmt.Errorf("%+v", e)
	}
}

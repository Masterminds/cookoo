package auth

import (
	//"fmt"
	"github.com/Masterminds/cookoo"
	"bitbucket.org/mobiplug/sugarmill/model"
) 


// Authenticate a username/password pair.
//
// This expects a username and a password as strings.
// A boolean `true` indicates that the user has been authenticated. A `false`
// indicates that the user/password combo has failed to auth. This is not
// necessarily an error. An error should only be returned when an unexpected
// condition has obtained during authentication.
type UserDatasource interface {
	AuthUser(username, password string) (bool, error)
}

type Admin struct {
	Cxt cookoo.Context
}

func (ad *Admin) AuthUser(username, password string) (bool, error) {
	acct := model.NewApiKey(ad.Cxt)
	acct.Account = username
	
	// Check if account under username exists
	if err := acct.LoadAccount(); err != nil {
		ad.Cxt.Logf("warn", "AuthUser could not find %s. Error %s", username, err)
		return false, err
	}
	
	// Check if account under username is in Admin group
	if "admin" != acct.AccountType {
		ad.Cxt.Logf("debug", "%s is not a member of the group Admin", acct.AccountType)
		return false, nil	
	}

	return acct.Valid(password), nil
}

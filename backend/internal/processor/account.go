package processor

import (
	"backend/constant"
	"backend/model"
	"fmt"
	"net/http"

	"github.com/free-ran-ue/util"
)

func (p *Processor) Login(req *model.RequestLogin) (*model.ResponseLogin, *model.ErrorDetail) {
	p.ProcLog.Debugf("Processing login for username: %s", req.Username)

	account, err := p.DmContext.Db().LoadAll(constant.COLL_ACCOUNT)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to load account: %v", err),
		}
	}
	p.ProcLog.Debugf("account get from db: %v", account)

	if len(account) != 0 { // user had modified the username or password
		hashValue, found := account[req.Username]
		if !found {
			return nil, &model.ErrorDetail{
				HttpStatus: http.StatusUnauthorized,
				Detail:     "Invalid username or incorrect password",
			}
		}

		verifyResult, err := p.DmContext.Password().Verify(req.Password, hashValue)
		if err != nil {
			return nil, &model.ErrorDetail{
				HttpStatus: http.StatusInternalServerError,
				Detail:     fmt.Sprintf("failed to verify account: %v", err),
			}
		}

		if !verifyResult {
			return nil, &model.ErrorDetail{
				HttpStatus: http.StatusUnauthorized,
				Detail:     "Invalid username or incorrect password",
			}
		}
	} else {
		p.ProcLog.Debugf("account get from db: %v", account)
		if req.Username != p.username || req.Password != p.password {
			return nil, &model.ErrorDetail{
				HttpStatus: http.StatusUnauthorized,
				Detail:     "Invalid username or incorrect password",
			}
		}
	}

	token, err := util.CreateJWT(p.jwtSecret, req.Username, p.jwtExpiresIn, nil)
	if err != nil {
		p.ProcLog.Errorf("Failed to create JWT for username %s: %v", req.Username, err)
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     "Failed to create JWT",
		}
	}

	return &model.ResponseLogin{
		Message: "Login successful",
		Token:   token,
	}, nil
}

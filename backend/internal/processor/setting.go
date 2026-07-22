package processor

import (
	"backend/constant"
	"backend/model"
	"fmt"
	"net/http"
)

func (p *Processor) SettingAccount(req *model.RequestSettingAccount) *model.ErrorDetail {
	account, err := p.DmContext.LoadAll(constant.COLL_ACCOUNT)
	if err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to load exist account: %v", err),
		}
	}

	existAccount, existHashValue := "", ""
	for k, v := range account {
		existAccount, existHashValue = k, v
		p.ProcLog.Debugf("Load exist account: %s", k)
	}

	hashValue, err := p.DmContext.Password().Hash(req.Password)
	if err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to hash the new password: %v", err),
		}
	}

	if err := p.DmContext.Db().RemoveAll(constant.COLL_ACCOUNT); err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to drop exist account: %v", err),
		}
	}

	if err := p.DmContext.Db().Save(constant.COLL_ACCOUNT, req.Username, hashValue); err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_ACCOUNT, existAccount, existHashValue); err != nil {
			p.ProcLog.Errorf("failed to rollback exist user: %v", err)
		}

		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to set account: %v", err),
		}
	}

	return nil
}

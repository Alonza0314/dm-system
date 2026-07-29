package processor

import (
	"backend/constant"
	"backend/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func (p *Processor) GetHistory(cate, dev string) (*model.ResponseGetHistory, *model.ErrorDetail) {
	// unlock at line #43
	p.DmContext.Db().MainLock().RLock()
	defer p.DmContext.Db().MainLock().RUnlock()

	existCate, err := p.DmContext.Db().Exist(constant.COLL_CATEGORY, cate)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to check if category exists: %v", err),
		}
	}
	if !existCate {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusNotFound,
			Detail:     fmt.Sprintf("category %s not found", cate),
		}
	}

	existDev, err := p.DmContext.Db().Exist(constant.COLL_CATEGORY_TAG+cate, dev)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to check if device exists: %v", err),
		}
	}
	if !existDev {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusNotFound,
			Detail:     fmt.Sprintf("device %s not found", dev),
		}
	}

	p.DmContext.Db().HistoryLock().RLock()
	defer p.DmContext.Db().HistoryLock().RUnlock()

	key := p.DmContext.Db().GetHistoryKey(cate, dev)

	response := &model.ResponseGetHistory{
		Histories: make(model.Histories, 0),
	}

	history, err := p.DmContext.Db().Load(constant.COLL_HISTORY, key)
	if err != nil {
		if err.Error() == "bucket not found" {
			p.ProcLog.Debugln("history is empty for category %s, device %s", cate, dev)
			return response, nil
		}

		p.ProcLog.Errorf("failed to load history bucket: %v", err)
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to load history bucket: %v", err),
		}
	}

	var historyUnmarshal model.Histories
	if err := json.Unmarshal([]byte(history), &historyUnmarshal); err != nil {
		p.ProcLog.Errorf("failed to unmarshal history get from db: %v", err)
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to unmarshal history for device %s: %v", dev, err),
		}
	}

	response.Histories = historyUnmarshal

	return response, nil
}

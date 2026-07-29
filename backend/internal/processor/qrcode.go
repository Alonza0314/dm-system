package processor

import (
	"backend/constant"
	"backend/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func taiwanNow() *time.Time {
	now := time.Now().In(time.FixedZone("CST", 8*60*60))
	return &now
}

func (p *Processor) Borrow(cate, dev string, req *model.RequestQrcodeBorrow) *model.ErrorDetail {
	p.DmContext.Db().MainLock().Lock()
	defer p.DmContext.Db().MainLock().Unlock()

	existCate, err := p.DmContext.Db().Exist(constant.COLL_CATEGORY, cate)
	if err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to check if category exists: %v", err),
		}
	}
	if !existCate {
		return &model.ErrorDetail{
			HttpStatus: http.StatusNotFound,
			Detail:     fmt.Sprintf("category %s not found", cate),
		}
	}

	existDev, err := p.DmContext.Db().Exist(constant.COLL_CATEGORY_TAG+cate, dev)
	if err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to check if device exists: %v", err),
		}
	}
	if !existDev {
		return &model.ErrorDetail{
			HttpStatus: http.StatusNotFound,
			Detail:     fmt.Sprintf("device %s not found", dev),
		}
	}

	device, err := p.DmContext.Load(constant.COLL_CATEGORY_TAG+cate, dev)
	if err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get device: %v", err),
		}
	}

	var deviceUnmarshal model.Device
	if err := json.Unmarshal([]byte(device), &deviceUnmarshal); err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to unmarshal device: %v", err),
		}
	}

	if deviceUnmarshal.Status == constant.STATUS_USING {
		return &model.ErrorDetail{
			HttpStatus: http.StatusConflict,
			Detail:     fmt.Sprintf("device %s is still using by %s", dev, deviceUnmarshal.User),
		}
	}

	newDeviceUnmarshal := model.Device{
		Id:         deviceUnmarshal.Id,
		Category:   deviceUnmarshal.Category,
		Name:       deviceUnmarshal.Name,
		Status:     constant.STATUS_USING,
		User:       req.User,
		LastBorrow: taiwanNow(),
		LastReturn: nil,
		Owner:      deviceUnmarshal.Owner,
		Note:       deviceUnmarshal.Note,
	}

	newDeviceMarshal, err := json.Marshal(newDeviceUnmarshal)
	if err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to marshal json: %v", err),
		}
	}

	if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+newDeviceUnmarshal.Category, newDeviceUnmarshal.Name, string(newDeviceMarshal)); err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to save device: %v", err),
		}
	}

	category, err := p.DmContext.Db().Load(constant.COLL_CATEGORY, newDeviceUnmarshal.Category)
	if err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+deviceUnmarshal.Category, deviceUnmarshal.Name, string(device)); err != nil {
			p.ProcLog.Errorf("failed to rollback device: %v", err)
		}
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get category: %v", err),
		}
	}

	var categoryUnmarshal model.Category
	if err := json.Unmarshal([]byte(category), &categoryUnmarshal); err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+deviceUnmarshal.Category, deviceUnmarshal.Name, string(device)); err != nil {
			p.ProcLog.Errorf("failed to rollback device: %v", err)
		}
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to unmarshal category: %v", err),
		}
	}

	categoryUnmarshal.IdleDevice, categoryUnmarshal.UsingDevice = categoryUnmarshal.IdleDevice-1, categoryUnmarshal.UsingDevice+1

	categoryMarshal, err := json.Marshal(categoryUnmarshal)
	if err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+deviceUnmarshal.Category, deviceUnmarshal.Name, string(device)); err != nil {
			p.ProcLog.Errorf("failed to rollback device: %v", err)
		}
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to marshal json: %v", err),
		}
	}

	if err := p.DmContext.Db().Save(constant.COLL_CATEGORY, newDeviceUnmarshal.Category, string(categoryMarshal)); err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+deviceUnmarshal.Category, deviceUnmarshal.Name, string(device)); err != nil {
			p.ProcLog.Errorf("failed to rollback device: %v", err)
		}
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to save category: %v", err),
		}
	}

	return nil
}

func (p *Processor) Return(cate, dev string, req *model.RequestQrcodeReturn) *model.ErrorDetail {
	p.DmContext.Db().MainLock().Lock()
	defer p.DmContext.Db().MainLock().Unlock()

	existCate, err := p.DmContext.Db().Exist(constant.COLL_CATEGORY, cate)
	if err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to check if category exists: %v", err),
		}
	}
	if !existCate {
		return &model.ErrorDetail{
			HttpStatus: http.StatusNotFound,
			Detail:     fmt.Sprintf("category %s not found", cate),
		}
	}

	existDev, err := p.DmContext.Db().Exist(constant.COLL_CATEGORY_TAG+cate, dev)
	if err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to check if device exists: %v", err),
		}
	}
	if !existDev {
		return &model.ErrorDetail{
			HttpStatus: http.StatusNotFound,
			Detail:     fmt.Sprintf("device %s not found", dev),
		}
	}

	device, err := p.DmContext.Load(constant.COLL_CATEGORY_TAG+cate, dev)
	if err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get device: %v", err),
		}
	}

	var deviceUnmarshal model.Device
	if err := json.Unmarshal([]byte(device), &deviceUnmarshal); err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to unmarshal device: %v", err),
		}
	}

	if deviceUnmarshal.Status == constant.STATUS_IDLE {
		return &model.ErrorDetail{
			HttpStatus: http.StatusConflict,
			Detail:     fmt.Sprintf("device %s is idle", dev),
		}
	}

	newDeviceUnmarshal := model.Device{
		Id:         deviceUnmarshal.Id,
		Category:   deviceUnmarshal.Category,
		Name:       deviceUnmarshal.Name,
		Status:     constant.STATUS_IDLE,
		User:       deviceUnmarshal.User,
		LastBorrow: deviceUnmarshal.LastBorrow,
		LastReturn: taiwanNow(),
		Owner:      deviceUnmarshal.Owner,
		Note:       deviceUnmarshal.Note,
	}

	newDeviceMarshal, err := json.Marshal(newDeviceUnmarshal)
	if err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to marshal json: %v", err),
		}
	}

	if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+newDeviceUnmarshal.Category, newDeviceUnmarshal.Name, string(newDeviceMarshal)); err != nil {
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to save device: %v", err),
		}
	}

	category, err := p.DmContext.Db().Load(constant.COLL_CATEGORY, newDeviceUnmarshal.Category)
	if err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+deviceUnmarshal.Category, deviceUnmarshal.Name, string(device)); err != nil {
			p.ProcLog.Errorf("failed to rollback device: %v", err)
		}
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get category: %v", err),
		}
	}

	var categoryUnmarshal model.Category
	if err := json.Unmarshal([]byte(category), &categoryUnmarshal); err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+deviceUnmarshal.Category, deviceUnmarshal.Name, string(device)); err != nil {
			p.ProcLog.Errorf("failed to rollback device: %v", err)
		}
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to unmarshal category: %v", err),
		}
	}

	categoryUnmarshal.IdleDevice, categoryUnmarshal.UsingDevice = categoryUnmarshal.IdleDevice+1, categoryUnmarshal.UsingDevice-1

	categoryMarshal, err := json.Marshal(categoryUnmarshal)
	if err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+deviceUnmarshal.Category, deviceUnmarshal.Name, string(device)); err != nil {
			p.ProcLog.Errorf("failed to rollback device: %v", err)
		}
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to marshal json: %v", err),
		}
	}

	if err := p.DmContext.Db().Save(constant.COLL_CATEGORY, newDeviceUnmarshal.Category, string(categoryMarshal)); err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+deviceUnmarshal.Category, deviceUnmarshal.Name, string(device)); err != nil {
			p.ProcLog.Errorf("failed to rollback device: %v", err)
		}
		return &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to save category: %v", err),
		}
	}

	p.DmContext.Db().HistoryLock().Lock()
	defer p.DmContext.Db().HistoryLock().Unlock()

	key, targetHistory := p.DmContext.Db().GetHistoryKey(cate, dev), model.History{
		User:       newDeviceUnmarshal.User,
		LastBorrow: newDeviceUnmarshal.LastBorrow,
		LastReturn: newDeviceUnmarshal.LastReturn,
	}

	history, err := p.DmContext.Db().Load(constant.COLL_HISTORY, key)
	if err != nil && err.Error() != "bucket not found" {
		p.ProcLog.Errorf("failed to load history bucket: %v", err)
		return nil
	}

	if history == "" {
		p.ProcLog.Debugln("history is empty, create new one")

		history := make(model.Histories, 0)
		history = append(history, targetHistory)

		historyMarshal, err := json.Marshal(history)
		if err != nil {
			p.ProcLog.Errorf("failed to marshal history: %v", err)
			return nil
		}

		if err := p.DmContext.Db().Save(constant.COLL_HISTORY, key, string(historyMarshal)); err != nil {
			p.ProcLog.Errorf("failed to save marshal history to db: %v", err)
		}

		return nil
	}

	var historyUnmarshal model.Histories
	if err := json.Unmarshal([]byte(history), &historyUnmarshal); err != nil {
		p.ProcLog.Errorf("failed to unmarshal history get from db: %v", err)
		return nil
	}

	historyUnmarshal = append([]model.History{targetHistory}, historyUnmarshal...)
	if len(historyUnmarshal) > p.maxHistory {
		historyUnmarshal = historyUnmarshal[:p.maxHistory]
	}

	historyMarshal, err := json.Marshal(historyUnmarshal)
	if err != nil {
		p.ProcLog.Errorf("failed to marshal back history: %v", err)
		return nil
	}

	if err := p.DmContext.Db().Save(constant.COLL_HISTORY, key, string(historyMarshal)); err != nil {
		p.ProcLog.Errorf("failed to save marshal history to db: %v", err)
	}

	return nil
}

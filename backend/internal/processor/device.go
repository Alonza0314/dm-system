package processor

import (
	"backend/constant"
	"backend/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func (p *Processor) GetDevices(cate string) (*model.ResponseGetDevices, *model.ErrorDetail) {
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

	devices, err := p.DmContext.Db().LoadAll(constant.COLL_CATEGORY_TAG + cate)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get category devices: %v", err),
		}
	}

	response := &model.ResponseGetDevices{
		Devices: make([]model.DeviceShort, 0, len(devices)),
	}
	for _, v := range devices {
		var device model.DeviceShort
		if err := json.Unmarshal([]byte(v), &device); err != nil {
			return nil, &model.ErrorDetail{
				HttpStatus: http.StatusInternalServerError,
				Detail:     fmt.Sprintf("failed to unmarshal category devices: %v", err),
			}
		}

		response.Devices = append(response.Devices, device)
	}

	return response, nil
}

func (p *Processor) GetDevice(cate, dev string) (*model.ResponseGetDevice, *model.ErrorDetail) {
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

	device, err := p.DmContext.Db().Load(constant.COLL_CATEGORY_TAG+cate, dev)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get device: %v", err),
		}
	}

	var deviceUnmarshal model.ResponseGetDevice
	if err := json.Unmarshal([]byte(device), &deviceUnmarshal); err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to unmarshal device: %v", err),
		}
	}

	return &deviceUnmarshal, nil
}

func (p *Processor) CreateDevice(req *model.RequestCreateDevice) (*model.ResponseCreateDevice, *model.ErrorDetail) {
	existCate, err := p.DmContext.Db().Exist(constant.COLL_CATEGORY, req.Category)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to check if category exists: %v", err),
		}
	}
	if !existCate {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusNotFound,
			Detail:     fmt.Sprintf("category %s not found", req.Category),
		}
	}

	existDev, err := p.DmContext.Db().Exist(constant.COLL_CATEGORY_TAG+req.Category, req.Name)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to check if device exists: %v", err),
		}
	}
	if existDev {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusConflict,
			Detail:     "device already exists",
		}
	}

	devId, err := p.DmContext.RequestId(constant.ID_KEY_DEVICE)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to request device id: %v", err),
		}
	}
	req.Id, req.Status = devId, constant.STATUS_IDLE

	device, err := json.Marshal(req)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to marshal json: %v", err),
		}
	}

	if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+req.Category, req.Name, string(device)); err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to save device: %v", err),
		}
	}

	category, err := p.DmContext.Db().Load(constant.COLL_CATEGORY, req.Category)
	if err != nil {
		if err := p.DmContext.Db().Remove(constant.COLL_CATEGORY_TAG+req.Category, req.Name); err != nil {
			p.ProcLog.Errorf("failed to rollback already saved device: %v", err)
		}
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get category: %v", err),
		}
	}

	var categoryUnmarshal model.Category
	if err := json.Unmarshal([]byte(category), &categoryUnmarshal); err != nil {
		if err := p.DmContext.Db().Remove(constant.COLL_CATEGORY_TAG+req.Category, req.Name); err != nil {
			p.ProcLog.Errorf("failed to rollback already saved device: %v", err)
		}
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to unmarshal category: %v", err),
		}
	}

	categoryUnmarshal.IdleDevice += 1

	categoryMarshal, err := json.Marshal(categoryUnmarshal)
	if err != nil {
		if err := p.DmContext.Db().Remove(constant.COLL_CATEGORY_TAG+req.Category, req.Name); err != nil {
			p.ProcLog.Errorf("failed to rollback already saved device: %v", err)
		}
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to marshal json: %v", err),
		}
	}

	if err := p.DmContext.Db().Save(constant.COLL_CATEGORY, req.Category, string(categoryMarshal)); err != nil {
		if err := p.DmContext.Db().Remove(constant.COLL_CATEGORY_TAG+req.Category, req.Name); err != nil {
			p.ProcLog.Errorf("failed to rollback already saved device: %v", err)
		}
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to save category: %v", err),
		}
	}

	return &model.ResponseCreateDevice{
		Message: "Created successful",
	}, nil
}

func (p *Processor) DeleteDevice(cate, dev string) (*model.ResponseDeleteDevice, *model.ErrorDetail) {
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

	device, err := p.DmContext.Load(constant.COLL_CATEGORY_TAG+cate, dev)
	if err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get device: %v", err),
		}
	}

	var deviceUnmarshal model.Device
	if err := json.Unmarshal([]byte(device), &deviceUnmarshal); err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to unmarshal device: %v", err),
		}
	}

	if deviceUnmarshal.Status == constant.STATUS_USING {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusConflict,
			Detail:     fmt.Sprintf("device %s is still using by %s", dev, deviceUnmarshal.User),
		}
	}

	if err := p.DmContext.Remove(constant.COLL_CATEGORY_TAG+cate, dev); err != nil {
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to delete device: %v", err),
		}
	}

	category, err := p.DmContext.Db().Load(constant.COLL_CATEGORY, cate)
	if err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+cate, dev, device); err != nil {
			p.ProcLog.Errorf("failed to rollback already deleted device: %v", err)
		}
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to get category: %v", err),
		}
	}

	var categoryUnmarshal model.Category
	if err := json.Unmarshal([]byte(category), &categoryUnmarshal); err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+cate, dev, device); err != nil {
			p.ProcLog.Errorf("failed to rollback already deleted device: %v", err)
		}
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to unmarshal category: %v", err),
		}
	}

	categoryUnmarshal.IdleDevice -= 1

	categoryMarshal, err := json.Marshal(categoryUnmarshal)
	if err != nil {
		if err := p.DmContext.Db().Save(constant.COLL_CATEGORY_TAG+cate, dev, device); err != nil {
			p.ProcLog.Errorf("failed to rollback already saved device: %v", err)
		}
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to marshal json: %v", err),
		}
	}

	if err := p.DmContext.Db().Save(constant.COLL_CATEGORY, cate, string(categoryMarshal)); err != nil {
		if err := p.DmContext.Db().Remove(constant.COLL_CATEGORY_TAG+cate, dev); err != nil {
			p.ProcLog.Errorf("failed to rollback already saved device: %v", err)
		}
		return nil, &model.ErrorDetail{
			HttpStatus: http.StatusInternalServerError,
			Detail:     fmt.Sprintf("failed to save category: %v", err),
		}
	}

	return &model.ResponseDeleteDevice{
		Message: "Delete successful",
	}, nil
}

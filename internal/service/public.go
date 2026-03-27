package service

import (
	"context"
	"fmt"
	"time"

	"go-admin/server/internal/model"
	"go-admin/server/internal/repository"
)

type PublicService struct {
	repo *repository.Repository
}

func NewPublicService(repo *repository.Repository) *PublicService {
	return &PublicService{repo: repo}
}

type PublicDeviceQuery struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	DeviceNo   string `form:"deviceNo"`
	MerchantID string `form:"merchantId"`
	Status     *uint8 `form:"status"`
}

type UpdatePublicDeviceStatusPayload struct {
	Status uint8 `json:"status" binding:"required"`
}

type RegisterDevicePayload struct {
	DeviceNo   string `json:"deviceNo" binding:"required"`
	MerchantID string `json:"merchantId" binding:"required"`
}

func (q PublicDeviceQuery) Normalize() PublicDeviceQuery {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	return q
}

func (s *PublicService) ListDevices(ctx context.Context, query PublicDeviceQuery) (*PageResult, error) {
	query = query.Normalize()
	devices, total, err := s.repo.ListDevicesByQuery(ctx, query.Page, query.PageSize, repository.DeviceQuery{
		DeviceNo:   query.DeviceNo,
		MerchantID: query.MerchantID,
		Status:     query.Status,
	})
	if err != nil {
		return nil, err
	}
	return &PageResult{List: devices, Total: total, Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *PublicService) GetDevice(ctx context.Context, deviceNo string) (*model.Device, error) {
	return s.repo.GetDeviceByDeviceNo(ctx, deviceNo)
}

func (s *PublicService) UpdateDeviceStatus(ctx context.Context, deviceNo string, status uint8) error {
	if status != 0 && status != 1 {
		return fmt.Errorf("invalid status")
	}
	device, err := s.repo.GetDeviceByDeviceNo(ctx, deviceNo)
	if err != nil {
		return err
	}
	device.Status = status
	return s.repo.SaveDevice(ctx, device)
}

func (s *PublicService) RegisterDevice(ctx context.Context, payload RegisterDevicePayload, clientIP string) (*model.Device, error) {
	device := &model.Device{
		DeviceNo:   payload.DeviceNo,
		MerchantID: payload.MerchantID,
		Status:     0,
		IP:         clientIP,
		CreateT:    uint64(time.Now().UnixMilli()),
	}
	if err := s.repo.CreateDevice(ctx, device); err != nil {
		return nil, err
	}
	return device, nil
}

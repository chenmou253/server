package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-admin/server/internal/model"
	"go-admin/server/internal/service"
	"go-admin/server/pkg/response"
)

type PublicHandler struct {
	public *service.PublicService
}

type registerDeviceSwaggerResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    model.Device `json:"data"`
}

type errorSwaggerResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewPublicHandler(public *service.PublicService) *PublicHandler {
	return &PublicHandler{public: public}
}

func (h *PublicHandler) ListDevices(c *gin.Context) {
	var query service.PublicDeviceQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.public.ListDevices(c.Request.Context(), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// RegisterDevice godoc
// @Summary 设备注册
// @Description 向 device 表插入一条设备记录。
// @Description
// @Description 请求头鉴权说明：
// @Description - X-API-Key: 开放接口访问标识
// @Description - X-Timestamp: 秒级 Unix 时间戳
// @Description - X-Nonce: 每次请求唯一随机串
// @Description - X-Signature: 按 apiKey\ntimestamp\nnonce\nMETHOD\n/path\nsha256(body)\nsecret 规则生成的 SHA256 签名
// @Description
// @Description 业务说明：
// @Description - deviceNo 重复时禁止注册
// @Tags 开放设备接口
// @Accept json
// @Produce json
// @Param X-API-Key header string true "开放接口访问标识"
// @Param X-Timestamp header string true "秒级 Unix 时间戳"
// @Param X-Nonce header string true "请求唯一随机串"
// @Param X-Signature header string true "SHA256 签名值"
// @Param request body service.RegisterDevicePayload true "设备注册请求参数"
// @Success 200 {object} registerDeviceSwaggerResponse "注册成功"
// @Failure 400 {object} errorSwaggerResponse "参数错误或设备号重复"
// @Failure 401 {object} errorSwaggerResponse "鉴权失败"
// @Router /api/v1/open/devices [post]
func (h *PublicHandler) RegisterDevice(c *gin.Context) {
	var payload service.RegisterDevicePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	device, err := h.public.RegisterDevice(c.Request.Context(), payload, c.ClientIP())
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, device)
}

func (h *PublicHandler) GetDevice(c *gin.Context) {
	deviceNo := c.Param("deviceNo")
	device, err := h.public.GetDevice(c.Request.Context(), deviceNo)
	if err != nil {
		response.Error(c, http.StatusNotFound, "device not found")
		return
	}
	response.Success(c, device)
}

func (h *PublicHandler) UpdateDeviceStatus(c *gin.Context) {
	deviceNo := c.Param("deviceNo")
	var payload service.UpdatePublicDeviceStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.public.UpdateDeviceStatus(c.Request.Context(), deviceNo, payload.Status); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{})
}

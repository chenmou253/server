package handler

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-admin/server/internal/service"
	"go-admin/server/pkg/response"
)

type PublicHandler struct {
	public *service.PublicService
}

type deviceSwaggerData struct {
	// 设备主键 ID
	ID uint `json:"id" example:"1"`

	// 设备编号
	DeviceNo string `json:"deviceNo" example:"DVC202603290001"`

	// 商户编号
	MerchantID string `json:"merchantId" example:"M10001"`

	// 设备状态。0 表示待机，1 表示接单中
	Status uint8 `json:"status" example:"1" enums:"0,1"`

	// 设备上报的客户端 IP
	IP string `json:"ip" example:"192.168.1.10"`

	// 设备创建时间，Unix 毫秒时间戳
	CreateT uint64 `json:"createT" example:"1711699200000"`
}

type deviceSwaggerResponse struct {
	Code    int               `json:"code" example:"0"`
	Message string            `json:"message" example:"ok"`
	Data    deviceSwaggerData `json:"data"`
}

type errorSwaggerData struct {
	// 错误原因标识
	Reason string `json:"reason" example:"device_not_found"`

	// 详细错误说明
	Detail string `json:"detail" example:"device not found"`

	// 相关设备编号。仅在设备相关错误时返回
	DeviceNo string `json:"deviceNo,omitempty" example:"DVC202603290001"`
}

type errorSwaggerResponse struct {
	Code    int              `json:"code" example:"404"`
	Message string           `json:"message" example:"device not found"`
	Data    errorSwaggerData `json:"data"`
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
// @Description - deviceNo 重复时默认返回已存在记录，视为注册成功
// @Tags 开放设备接口
// @Accept json
// @Produce json
// @Param X-API-Key header string true "开放接口访问标识"
// @Param X-Timestamp header string true "秒级 Unix 时间戳"
// @Param X-Nonce header string true "请求唯一随机串"
// @Param X-Signature header string true "SHA256 签名值"
// @Param request body service.RegisterDevicePayload true "设备注册请求参数"
// @Success 200 {object} deviceSwaggerResponse "注册成功"
// @Failure 400 {object} errorSwaggerResponse "参数错误"
// @Failure 401 {object} errorSwaggerResponse "鉴权失败"
// @Router /api/v1/open/devices [post]
func (h *PublicHandler) RegisterDevice(c *gin.Context) {
	var payload service.RegisterDevicePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.ErrorWithData(c, http.StatusBadRequest, err.Error(), gin.H{
			"reason": "invalid_request",
			"detail": err.Error(),
		})
		return
	}
	device, err := h.public.RegisterDevice(c.Request.Context(), payload, resolveClientIP(c))
	if err != nil {
		response.ErrorWithData(c, http.StatusBadRequest, err.Error(), gin.H{
			"reason": "register_device_failed",
			"detail": err.Error(),
		})
		return
	}
	response.Success(c, device)
}

func resolveClientIP(c *gin.Context) string {
	candidates := []string{
		c.GetHeader("X-Forwarded-For"),
		c.GetHeader("X-Real-IP"),
		c.Request.Header.Get("CF-Connecting-IP"),
		c.Request.RemoteAddr,
		c.ClientIP(),
	}
	for _, candidate := range candidates {
		if ip := firstRoutableIP(candidate); ip != "" {
			return ip
		}
	}
	for _, candidate := range candidates {
		if ip := firstParsedIP(candidate, true); ip != "" {
			return ip
		}
	}
	return ""
}

func firstRoutableIP(raw string) string {
	return firstParsedIP(raw, false)
}

func firstParsedIP(raw string, allowLoopback bool) string {
	for _, part := range strings.Split(raw, ",") {
		candidate := strings.TrimSpace(part)
		if candidate == "" {
			continue
		}
		if host, _, err := net.SplitHostPort(candidate); err == nil {
			candidate = host
		}
		ip := net.ParseIP(candidate)
		if ip == nil || ip.IsUnspecified() {
			continue
		}
		if !allowLoopback && ip.IsLoopback() {
			continue
		}
		return ip.String()
	}
	return ""
}

// GetDevice godoc
// @Summary 查询设备
// @Description 根据 deviceNo 查询单个设备信息。
// @Tags 开放设备接口
// @Accept json
// @Produce json
// @Param X-API-Key header string true "开放接口访问标识"
// @Param X-Timestamp header string true "秒级 Unix 时间戳"
// @Param X-Nonce header string true "请求唯一随机串"
// @Param X-Signature header string true "SHA256 签名值"
// @Param deviceNo path string true "设备编号"
// @Success 200 {object} deviceSwaggerResponse "查询成功"
// @Failure 401 {object} errorSwaggerResponse "鉴权失败"
// @Failure 404 {object} errorSwaggerResponse "设备不存在"
// @Router /api/v1/open/devices/{deviceNo} [get]
func (h *PublicHandler) GetDevice(c *gin.Context) {
	deviceNo := c.Param("deviceNo")
	device, err := h.public.GetDevice(c.Request.Context(), deviceNo)
	if err != nil {
		response.ErrorWithData(c, http.StatusNotFound, "device not found", gin.H{
			"reason":   "device_not_found",
			"detail":   "device not found",
			"deviceNo": deviceNo,
		})
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

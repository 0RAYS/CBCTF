package dto

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type UpdateSettingForm struct {
	Host *string `form:"host" json:"host" binding:"omitempty,url"`
	Path *string `form:"path" json:"path" binding:"omitempty,dir"`

	LogLevel *string `form:"log_level" json:"log_level" binding:"omitempty,oneof=DEBUG INFO WARNING ERROR"`
	LogSave  *bool   `form:"log_save" json:"log_save"`

	AsyncQLogLevel    *string `form:"asyncq_log_level" json:"asyncq_log_level" binding:"omitempty,oneof=DEBUG INFO WARNING ERROR"`
	AsyncQConcurrency *int    `form:"asyncq_concurrency" json:"asyncq_concurrency" binding:"omitempty,gte=1"`

	GinMode               *string   `form:"gin_mode" json:"gin_mode" binding:"omitempty,oneof=debug test release"`
	GinHost               *string   `form:"gin_host" json:"gin_host" binding:"omitempty,ip"`
	GinPort               *uint     `form:"gin_port" json:"gin_port" binding:"omitempty,port"`
	GinUploadMax          *int      `form:"gin_upload_max" json:"gin_upload_max" binding:"omitempty,gte=1"`
	GinProxies            *[]string `form:"gin_proxies" json:"gin_proxies" binding:"omitempty,dive,ip|cidr"`
	GinRateLimitGlobal    *int      `form:"gin_ratelimit_global" json:"gin_ratelimit_global" binding:"omitempty,gte=1"`
	GinRateLimitWhitelist *[]string `form:"gin_ratelimit_whitelist" json:"gin_ratelimit_whitelist" binding:"omitempty,dive,ip|cidr"`
	GinCORS               *[]string `form:"gin_cors" json:"gin_cors" binding:"omitempty,dive,url"`
	GinLogWhitelist       *[]string `form:"gin_log_whitelist" json:"gin_log_whitelist" binding:"omitempty,dive,uri"`

	GormMySQLHost   *string `form:"gorm_mysql_host" json:"gorm_mysql_host" binding:"omitempty,ip"`
	GormMySQLPort   *uint   `form:"gorm_mysql_port" json:"gorm_mysql_port" binding:"omitempty,port"`
	GormMySQLUser   *string `form:"gorm_mysql_user" json:"gorm_mysql_user" binding:"omitempty,min=1,alphanum"`
	GormMySQLPwd    *string `form:"gorm_mysql_pwd" json:"gorm_mysql_pwd" binding:"omitempty,min=1,ascii"`
	GormMySQLDB     *string `form:"gorm_mysql_db" json:"gorm_mysql_db" binding:"omitempty,min=1,alphanum"`
	GormMySQLMXOpen *int    `form:"gorm_mysql_mxopen" json:"gorm_mysql_mxopen" binding:"omitempty,min=1"`
	GormMySQLMXIdle *int    `form:"gorm_mysql_mxidle" json:"gorm_mysql_mxidle" binding:"omitempty,min=1"`
	GormLogLevel    *string `form:"gorm_log_level" json:"gorm_log_level" binding:"omitempty,oneof=SILENT INFO WARNING ERROR"`

	RedisHost *string `form:"redis_host" json:"redis_host" binding:"omitempty,ip"`
	RedisPort *uint   `form:"redis_port" json:"redis_port" binding:"omitempty,port"`
	RedisPwd  *string `form:"redis_pwd" json:"redis_pwd" binding:"omitempty,min=1,ascii"`

	K8SConfig                    *string   `form:"k8s_config" json:"k8s_config" binding:"omitempty,file"`
	K8SNamespace                 *string   `form:"k8s_namespace" json:"k8s_namespace" binding:"omitempty,min=1,alphanum"`
	K8SExternalNetworkCIDR       *string   `form:"k8s_external_network_cidr" json:"k8s_external_network_cidr" binding:"omitempty,cidr"`
	K8SExternalNetworkGateway    *string   `form:"k8s_external_network_gateway" json:"k8s_external_network_gateway" binding:"omitempty,ip"`
	K8SExternalNetworkInterface  *string   `form:"k8s_external_network_interface" json:"k8s_external_network_interface" binding:"omitempty,min=1,alphanum"`
	K8SExternalNetworkExcludeIPs *[]string `form:"k8s_external_network_exclude_ips" json:"k8s_external_network_exclude_ips" binding:"omitempty,dive,ip|cidr"`
	K8STCPDumpImage              *string   `form:"k8s_tcpdump" json:"k8s_tcpdump" binding:"omitempty,min=1"`

	K8SFrpOn         *bool                `form:"k8s_frp_on" json:"k8s_frp_on"`
	K8SFrpFrpcImage  *string              `form:"k8s_frp_frpc" json:"k8s_frp_frpc" binding:"omitempty,min=1"`
	K8SFrpNginxImage *string              `form:"k8s_frp_nginx" json:"k8s_frp_nginx" binding:"omitempty,min=1"`
	K8SFrpFrps       *[]config.FrpsConfig `form:"k8s_frp_frps" json:"k8s_frp_frps"`

	K8SGeneratorWorker *int `form:"k8s_generator_worker" json:"k8s_generator_worker" binding:"omitempty,gte=1"`

	NFSServer  *string `form:"nfs_server" json:"nfs_server" binding:"omitempty,ip|hostname"`
	NFSPath    *string `form:"nfs_path" json:"nfs_path" binding:"omitempty,dirpath|filepath"`
	NFSStorage *string `form:"nfs_storage" json:"nfs_storage"`

	CheatIPWhitelist *[]string `form:"cheat_ip_whitelist" json:"cheat_ip_whitelist" binding:"omitempty,dive,ip|cidr"`
}

func (f *UpdateSettingForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

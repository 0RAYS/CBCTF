/**
 * Builds update payload from nested config structure to flat structure expected by backend
 * @param {Object} config - Nested configuration object
 * @returns {Object} - Flat configuration payload
 */
export function buildPayload(config) {
  return {
    host: config.host,
    asyncq_log_level: config.asyncq.log.level,
    asyncq_victim_concurrency: config.asyncq.queues.victim,
    asyncq_traffic_concurrency: config.asyncq.queues.traffic,
    asyncq_generator_concurrency: config.asyncq.queues.generator,
    asyncq_attachment_concurrency: config.asyncq.queues.attachment,
    asyncq_email_concurrency: config.asyncq.queues.email,
    asyncq_webhook_concurrency: config.asyncq.queues.webhook,
    asyncq_image_concurrency: config.asyncq.queues.image,
    gin_mode: config.gin.mode,
    gin_upload_picture: config.gin.upload.picture,
    gin_upload_challenge: config.gin.upload.challenge,
    gin_upload_writeup: config.gin.upload.writeup,
    gin_proxies: config.gin.proxies,
    gin_ratelimit_global: config.gin.ratelimit.global,
    gin_ratelimit_whitelist: config.gin.ratelimit.whitelist,
    gin_cors: config.gin.cors,
    gin_log_whitelist: config.gin.log.whitelist,
    gin_jwt_secret: config.gin.jwt.secret || undefined,
    gin_metrics_whitelist: config.gin.metrics.whitelist,
    gorm_postgres_mxopen: config.gorm.postgres.mxopen,
    gorm_postgres_mxidle: config.gorm.postgres.mxidle,
    gorm_log_level: config.gorm.log.level,
    k8s_namespace: config.k8s.namespace,
    k8s_capture: config.k8s.capture,
    k8s_frp_on: config.k8s.frp.on,
    k8s_frp_frpc: config.k8s.frp.frpc,
    k8s_frp_nginx: config.k8s.frp.nginx,
    k8s_frp_frps: config.k8s.frp.frps,
    cheat_ip_whitelist: config.cheat.ip_whitelist,
    webhook_whitelist: config.webhook.whitelist,
    registration_enabled: config.registration.enabled,
    registration_default_group: config.registration.default_group,
  };
}

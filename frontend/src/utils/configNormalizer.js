/**
 * Normalizes incoming config from backend to nested structure
 * Handles both flat keys (gorm_mysql_host) and nested objects (gorm.mysql.host)
 */

const fallback = (value, defaultValue) => (value !== undefined && value !== null ? value : defaultValue);

const normalizeSecret = (value) => {
  if (!value || value === '******') {
    return '';
  }
  return value;
};

export function normalizeConfig(source) {
  if (!source) {
    return null;
  }

  const logLevel = fallback(source?.log_level, fallback(source?.log?.level, ''));
  const logSave = fallback(source?.log_save, fallback(source?.log?.save, false));
  const asyncqConcurrency = fallback(source?.asyncq_concurrency, fallback(source?.asynq?.concurrency, 0));
  const asyncqLogLevel = fallback(source?.asyncq_log_level, fallback(source?.asynq?.log?.level, ''));

  return {
    host: fallback(source?.host, ''),
    path: fallback(source?.path, ''),
    log: {
      level: logLevel,
      save: logSave,
    },
    asyncq: {
      concurrency: asyncqConcurrency,
      log: {
        level: asyncqLogLevel,
      },
    },
    gin: {
      host: fallback(source?.gin_host, fallback(source?.gin?.host, '')),
      mode: fallback(source?.gin_mode, fallback(source?.gin?.mode, '')),
      port: fallback(source?.gin_port, fallback(source?.gin?.port, 0)),
      upload: {
        max: fallback(source?.gin_upload_max, fallback(source?.gin?.upload?.max, 0)),
      },
      ratelimit: {
        global: fallback(source?.gin_ratelimit_global, fallback(source?.gin?.ratelimit?.global, 0)),
        whitelist: fallback(source?.gin_ratelimit_whitelist, fallback(source?.gin?.ratelimit?.whitelist, [])),
      },
      proxies: fallback(source?.gin_proxies, fallback(source?.gin?.proxies, [])),
      log: {
        whitelist: fallback(source?.gin_log_whitelist, fallback(source?.gin?.log?.whitelist, [])),
      },
      cors: fallback(source?.gin_cors, fallback(source?.gin?.cors, [])),
      jwt: {
        secret: fallback(source?.gin_jwt_secret, fallback(source?.gin?.jwt?.secret, '')),
        static: fallback(source?.gin_jwt_static, fallback(source?.gin?.jwt?.static, false)),
      },
    },
    gorm: {
      log: {
        level: fallback(source?.gorm_log_level, fallback(source?.gorm?.log?.level, '')),
      },
      mysql: {
        host: fallback(source?.gorm_mysql_host, fallback(source?.gorm?.mysql?.host, '')),
        port: fallback(source?.gorm_mysql_port, fallback(source?.gorm?.mysql?.port, 0)),
        db: fallback(source?.gorm_mysql_db, fallback(source?.gorm?.mysql?.db, '')),
        user: fallback(source?.gorm_mysql_user, fallback(source?.gorm?.mysql?.user, '')),
        mxidle: fallback(source?.gorm_mysql_mxidle, fallback(source?.gorm?.mysql?.mxidle, 0)),
        mxopen: fallback(source?.gorm_mysql_mxopen, fallback(source?.gorm?.mysql?.mxopen, 0)),
        pwd: normalizeSecret(fallback(source?.gorm_mysql_pwd, fallback(source?.gorm?.mysql?.pwd, ''))),
      },
    },
    redis: {
      host: fallback(source?.redis_host, fallback(source?.redis?.host, '')),
      port: fallback(source?.redis_port, fallback(source?.redis?.port, 0)),
      pwd: normalizeSecret(fallback(source?.redis_pwd, fallback(source?.redis?.pwd, ''))),
    },
    k8s: {
      config: fallback(source?.k8s_config, fallback(source?.k8s?.config, '')),
      namespace: fallback(source?.k8s_namespace, fallback(source?.k8s?.namespace, '')),
      tcpdump: fallback(source?.k8s_tcpdump, fallback(source?.k8s?.tcpdump, '')),
      generator_worker: fallback(source?.k8s_generator_worker, fallback(source?.k8s?.generator_worker, 0)),
      external_network: {
        cidr: fallback(source?.k8s_external_network_cidr, fallback(source?.k8s?.external_network?.cidr, '')),
        gateway: fallback(source?.k8s_external_network_gateway, fallback(source?.k8s?.external_network?.gateway, '')),
        interface: fallback(
          source?.k8s_external_network_interface,
          fallback(source?.k8s?.external_network?.interface, '')
        ),
        exclude_ips: fallback(
          source?.k8s_external_network_exclude_ips,
          fallback(source?.k8s?.external_network?.exclude_ips, [])
        ),
      },
      frp: {
        frpc: fallback(source?.k8s_frp_frpc, fallback(source?.k8s?.frp?.frpc, '')),
        nginx: fallback(source?.k8s_frp_nginx, fallback(source?.k8s?.frp?.nginx, '')),
        on: fallback(source?.k8s_frp_on, fallback(source?.k8s?.frp?.on, false)),
        frps: fallback(source?.k8s_frp_frps, fallback(source?.k8s?.frp?.frps, [])),
      },
    },
    nfs: {
      server: fallback(source?.nfs_server, fallback(source?.nfs?.server, '')),
      path: fallback(source?.nfs_path, fallback(source?.nfs?.path, '')),
      storage: fallback(source?.nfs_storage, fallback(source?.nfs?.storage, '')),
    },
    cheat: {
      ip_whitelist: fallback(source?.cheat_ip_whitelist, fallback(source?.cheat?.ip_whitelist, [])),
    },
    webhook: {
      blacklist: fallback(source?.webhook_blacklist, fallback(source?.webhook?.blacklist, [])),
    },
    registration: {
      enabled: fallback(source?.registration_enabled, fallback(source?.registration?.enabled, true)),
      default_group: fallback(source?.registration_default_group, fallback(source?.registration?.default_group, 0)),
    },
    geocity_db: fallback(source?.geocity_db, ''),
  };
}

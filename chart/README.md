# CBCTF

Helm Charts

The chart renders `.Values.cbctf` into `/app/config.yaml`; sensitive settings are stored in the ConfigMap. Change default passwords before production use.

These fields are deployment-only:

- PostgreSQL connection: `gorm.postgres.*`
- Redis connection: `redis.*`
- Data path: `path`
- Gin listen address and port: `gin.host`, `gin.port`

`cbctf.gorm.postgres.sslmode` is a boolean value.

- `false` maps to PostgreSQL DSN `sslmode=disable`
- `true` maps to PostgreSQL DSN `sslmode=require`

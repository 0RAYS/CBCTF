# CBCTF
预期中，是一个CTF平台

## Config
初次运行会创建`config.json`，配置后重新启动
```yaml
log:
    # DEBUG INFO WARNING ERROR
    level: info
    save: true
gin:
    mode: release
    host: 127.0.0.1
    port: 8000
    upload:
        path: ./uploads
        max: 8
gorm:
    # mysql sqlite
    type: sqlite
    sqlite:
        file: data.db
    mysql:
        host: 127.0.0.1
        port: 3306
        user: cbctf
        pwd: password
        db: cbctf
    log:
        # INFO WARNING ERROR SILENT
        level: silent
redis:
    addr: 127.0.0.1:6379
    pwd: ""
    # millisecond
    timeout: 10
frontend: http://127.0.0.1:3000
backend: http://127.0.0.1:8000
```
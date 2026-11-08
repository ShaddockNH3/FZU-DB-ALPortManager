package constants

import (
	"log"
	"os"
	"strconv"
)

// 数据库连接的默认值取自 docker-compose.yml，本地 `docker compose up -d` 之后不设环境变量也能跑。
// 线上换 openGauss 时用环境变量覆盖，不要改代码——之前这里写死了华为云的公网 IP 和密码。
var (
	PostgresHost     = envStr("POSTGRES_HOST", "localhost")
	PostgresPort     = envInt("POSTGRES_PORT", 5432)
	PostgresUser     = envStr("POSTGRES_USER", "commander")
	PostgresPassword = envStr("POSTGRES_PASSWORD", "mysecretpassword")
	PostgresDBName   = envStr("POSTGRES_DB", "al_port_db")
	PostgresSSLMode  = envStr("POSTGRES_SSLMODE", "disable")
	PostgresTimeZone = envStr("POSTGRES_TIMEZONE", "Asia/Shanghai")
)

func envStr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// envInt 端口这类数值给个能读的报错，别让手滑写成 "5432 " 那样悄悄失效。
func envInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("[CONFIG] %s=%q 不是整数，回退到 %d", key, raw, fallback)
		return fallback
	}
	return value
}

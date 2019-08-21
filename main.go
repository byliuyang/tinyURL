package main

import (
	"os"
	"short/cmd"
)

func main() {
	host := getEnv("DB_HOST", "localhost")
	portStr := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "short")
	recaptchaSecret := getEnv("RECAPTCHA_SECRET", "")
	githubClientId := getEnv("GITHUB_CLIENT_ID", "")
	githubClientSecret := getEnv("GITHUB_CLIENT_SECRET", "")
	jwtSecret := getEnv("JWT_SECRET", "")

	cmd.Execute(host, portStr, user, password, dbName, recaptchaSecret, githubClientId, githubClientSecret, jwtSecret)
}

func getEnv(varName string, defaultVal string) string {
	val := os.Getenv(varName)

	if val == "" {
		return defaultVal
	}

	return val
}

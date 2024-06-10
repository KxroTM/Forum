package forum

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var GoogleClientID string
var GoogleClientSecret string
var GithubClientID string
var GithubClientSecret string
var EmailPassword string

func LoadEnvFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("ligne mal formatée: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	return scanner.Err()
}

func GetEnvData() error {
	err := LoadEnvFile("./.env")
	if err != nil {
		return err
	}
	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	GithubClientID = os.Getenv("GITHUB_CLIENT_ID")
	GithubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	EmailPassword = os.Getenv("EMAIL_PASSWORD_MAILING")

	return nil
}

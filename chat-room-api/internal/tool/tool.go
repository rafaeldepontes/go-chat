package tool

import "os"

func ChecksEnvFile(s *string) {
	_, err := os.Stat(*s)
	if err != nil {
		*s = ".env.example"
	}
}

func FileExists(s string) bool {
	if _, err := os.Stat(s); err != nil {
		return false
	}
	return true
}
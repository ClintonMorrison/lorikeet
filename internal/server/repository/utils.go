package repository

import (
	"fmt"
	"strings"

	"github.com/ClintonMorrison/lorikeet/internal/model"
)

func pathForUserFolder(dataPath string, auth model.Auth) string {
	// Sanitize username to prevent path traversal
	username := strings.ReplaceAll(auth.Username, "/", "")
	username = strings.ReplaceAll(username, "\\", "")
	return fmt.Sprintf("%s/%s", dataPath, username)
}

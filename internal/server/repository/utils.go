package repository

import (
	"fmt"

	"github.com/ClintonMorrison/lorikeet/internal/model"
)

func pathForUserFolder(dataPath string, auth model.Auth) string {
	return fmt.Sprintf("%s/%s", dataPath, auth.Username)
}

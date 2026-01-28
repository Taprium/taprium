package migrations

import (
	"os"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/security"
)

func init() {
	m.Register(func(app core.App) error {
		defaultSettings, _ := app.FindFirstRecordByFilter("settings", "")

		if defaultSettings == nil {
			return nil
		}
		// 2. Determine initial secret
		secret := os.Getenv("CLUSTER_SECRET")
		if secret == "" {
			secret = security.RandomString(32)
		}

		defaultSettings.Set("auth_secret", secret)

		return app.Save(defaultSettings)
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}

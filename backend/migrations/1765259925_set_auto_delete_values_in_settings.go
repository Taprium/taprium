package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		defaultSettings, _ := app.FindFirstRecordByFilter("settings", "")

		if defaultSettings == nil {
			return nil
		}

		defaultSettings.Set("delete_finished_image_queues", false)
		defaultSettings.Set("delete_finished_text_queues", false)
		app.Save(defaultSettings)

		return nil
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}

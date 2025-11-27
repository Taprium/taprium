package migrations

import (
	"log"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		generationSettingsRecord, _ := app.FindFirstRecordByFilter("settings", "")

		if generationSettingsRecord == nil {
			return nil
		}

		log.Println("Generating default settings")

		if generationSettingsRecord.Get("upscale_model") == "realesr-animevideov3-x2" {
			generationSettingsRecord.Set("upscale_model", "realesr-animevideov3")
			app.Save(generationSettingsRecord)
		}

		return nil
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}

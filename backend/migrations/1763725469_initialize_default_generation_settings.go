package migrations

import (
	"log"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		generationSettingsRecord, _ := app.FindFirstRecordByFilter("settings", "")

		if generationSettingsRecord != nil {
			return nil
		}

		log.Println("Initializing default generation record")

		settingsCollection, err := app.FindCollectionByNameOrId("settings")
		if err != nil {
			return err
		}

		generationSettingsRecord = core.NewRecord(settingsCollection)
		generationSettingsRecord.Set("img_width", 512)
		generationSettingsRecord.Set("img_height", 512)
		generationSettingsRecord.Set("default_queue_count", 4)
		generationSettingsRecord.Set("upscale_times", 2)
		generationSettingsRecord.Set("upscale_timeout_in_second", 120)
		app.Save(generationSettingsRecord)

		return nil
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}

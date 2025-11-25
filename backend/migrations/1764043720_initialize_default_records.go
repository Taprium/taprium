package migrations

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// add up queries...
		generationSettingsRecord, _ := app.FindFirstRecordByFilter("settings", "")

		if generationSettingsRecord != nil {
			return nil
		}

		log.Println("Generating default settings")

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
		generationSettingsRecord.Set("upscale_model", "realesr-animevideov3-x2")
		app.Save(generationSettingsRecord)

		usersColleciton, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		userRecord := core.NewRecord(usersColleciton)
		userRecord.SetEmail(fmt.Sprintf("%s@tap.rium", uuid.New().String()))
		userRecord.Set("name", "taprium")
		userRecord.SetPassword("tapriumpassword")
		userRecord.SetVerified(true)
		err = app.Save(userRecord)
		if err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}

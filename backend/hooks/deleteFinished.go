package hooks

import (
	"log"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

func DeleteFinishedQueuesCheck(app *pocketbase.PocketBase) {
	settingsRecord, err := app.FindFirstRecordByFilter("settings", "")
	if err != nil {
		log.Printf("Failed to find default settings record: %v", err)
		return
	}

	if settingsRecord.GetBool("delete_finished_image_queues") {
		imageQueueRecords, err := app.FindRecordsByFilter("image_queues", "user_confirmed_upscale=true", "", 0, 0)
		if err != nil {
			log.Printf("Failed to find finished image queue records: %v", err)
			return
		}

		for _, q := range imageQueueRecords {
			images, err := app.FindRecordsByFilter("generated_images", "queue='{:queueId}'", "", 0, 0, dbx.Params{
				"queueId": q.Id,
			})
			if err != nil {
				log.Printf("Failed to find generated images from queue: %v", err)
				continue
			}

			for _, i := range images {
				app.Delete(i)
			}
		}
	}

	if settingsRecord.GetBool("delete_finished_text_queues") {
		textQueueRecords, err := app.FindRecordsByFilter("text_queues", "result!=''", "", 0, 0)
		if err != nil {
			log.Printf("Failed to find finished text queues: %v", err)
			return
		}

		for _, q := range textQueueRecords {
			app.Delete(q)
		}

	}

}

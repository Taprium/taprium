package hooks

import (
	"log"
	"time"

	"github.com/gofrs/flock"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

func ResetUpscaleTimeout(app *pocketbase.PocketBase) {
	fileLock := flock.New("/var/lock/timeout-reset.lock")

	locked, err := fileLock.TryLock()
	if err != nil || !locked {
		log.Printf("Reset processing.")
		return
	}

	defer func() {
		if locked {
			fileLock.Unlock()
		}
	}()

	settingsRecord, err := app.FindFirstRecordByFilter("settings", "")
	if err != nil {
		log.Printf("Failed to find settings record: %v", err)
		return
	}
	timeoutRecords, err := app.FindRecordsByFilter("generated_images", "selected=true && upscaled=true && runner !='' updated<{:timeoutTimestap} ", "", 0, 0,
		dbx.Params{
			"timeoutTimestap": time.Now().Add(-time.Duration(settingsRecord.GetInt("upscale_timeout_in_second") * int(time.Second))),
		})
	if err != nil {
		log.Printf("Failed to find upscale timeout records: %v", err)
		return
	}

	for _, r := range timeoutRecords {
		log.Printf("Upscale for image %s was reset.", r.Id)
		r.Set("runner", "")
		app.Save(r)
	}

	if locked {
		fileLock.Unlock()
	}
}

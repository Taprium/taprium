package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"taprium/backend/hooks"
	_ "taprium/backend/migrations"
)

func main() {
	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	app.OnRecordAfterCreateSuccess("image_queues").BindFunc(func(e *core.RecordEvent) error {
		go func() {
			hooks.GenerateImage(app, e.Record)
			hooks.ImageGenerationRecover(app)
		}()
		return e.Next()
	})

	app.OnRecordAfterCreateSuccess("text_queues").BindFunc(func(e *core.RecordEvent) error {
		go func() {
			hooks.GenerateText(app, e.Record)
			hooks.TextGenerationRecover(app)
		}()
		return e.Next()
	})

	// update runner pinged_at
	app.OnRecordAuthRequest("upscale_runners").BindFunc(func(e *core.RecordAuthRequestEvent) error {
		go func() {
			e.Record.Set("pinged_at", time.Now())
			app.Save(e.Record)
		}()

		return e.Next()
	})

	app.Cron().MustAdd("img-gen-recovery", "* * * * *", func() {
		hooks.ImageGenerationRecover(app)
	})

	app.Cron().MustAdd("text-gen-recovery", "* * * * *", func() {
		hooks.TextGenerationRecover(app)
	})

	app.Cron().MustAdd("delete-finished-check", "* * * * *", func() {
		hooks.DeleteFinishedQueuesCheck(app)
	})

	app.Cron().MustAdd("upscale-timeout-check", "* * * * *", func() {
		hooks.ResetUpscaleTimeout(app)
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		app.Settings().Batch.Enabled = true
		app.Settings().Batch.MaxRequests = 100
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		se.Router.POST("/api/cluster/auth", func(e *core.RequestEvent) error {
			data := struct {
				Secret    string `json:"secret"`
				MachineID string `json:"machine_id"`
				Hostname  string `json:"hostname"`
			}{}

			if err := e.BindBody(&data); err != nil {
				return e.BadRequestError("Invalid request data", err)
			}

			settingsRecord, err := app.FindFirstRecordByFilter("settings", "")
			if err != nil || data.Secret != settingsRecord.GetString("auth_secret") {
				return e.ForbiddenError("Invalid cluster secret", nil)
			}

			device, err := app.FindFirstRecordByData("upscale_runners", "device_id", data.MachineID)
			if err != nil {
				// 3. AUTO-CREATE if device doesn't exist
				collection, err := app.FindCollectionByNameOrId("upscale_runners")
				if err != nil {
					return e.InternalServerError("Collection not found", err)
				}

				newDevice := core.NewRecord(collection)
				newDevice.Set("device_id", data.MachineID)
				newDevice.Set("name", data.Hostname)
				newDevice.Set("verified", false) // Default locked
				newDevice.SetRandomPassword()

				if err := app.Save(newDevice); err != nil {
					return e.InternalServerError("Failed to auto-register device", err)
				}

				return e.JSON(http.StatusAccepted, map[string]string{
					"message": "Device registered. Awaiting admin verification.",
				})
			}

			// 4. Check if Admin has verified the device
			if !device.GetBool("verified") {
				return e.ForbiddenError("Device is not yet verified by an admin", nil)
			}

			token, err := device.NewAuthToken()
			if err != nil {
				return e.InternalServerError("Failed to generate token", err)
			}

			return e.JSON(http.StatusOK, map[string]string{
				"token": token,
			})
		})

		go func() {
			hooks.ImageGenerationRecover(app)
			hooks.TextGenerationRecover(app)
		}()

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

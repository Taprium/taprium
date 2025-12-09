package main

import (
	"log"
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

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		app.Settings().Batch.Enabled = true
		app.Settings().Batch.MaxRequests = 100
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

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

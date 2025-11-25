package main

import (
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"ai-shared/cf-image-request/hooks"
	_ "ai-shared/cf-image-request/migrations"
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

	app.OnRecordAfterUpdateSuccess("generate_queues").BindFunc(func(e *core.RecordEvent) error {
		go func() {
			if e.Record.GetString("status") == "queue" || e.Record.GetString("status") == "processing" {
				hooks.GenerateImage(app, e.Record)
			}
		}()

		return e.Next()
	})

	app.OnRecordAfterCreateSuccess("generate_queues").BindFunc(func(e *core.RecordEvent) error {
		go func() {
			if e.Record.GetString("status") == "queue" || e.Record.GetString("status") == "processing" {
				hooks.GenerateImage(app, e.Record)
			}
		}()
		return e.Next()
	})

	app.Cron().MustAdd("img-gen-recovery", "* * * * *", func() {
		hooks.GenerationRecover(app)
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		app.Settings().Batch.Enabled = true
		app.Settings().Batch.MaxRequests = 100
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		hooks.GenerationRecover(app)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

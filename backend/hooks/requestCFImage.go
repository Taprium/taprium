package hooks

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofrs/flock"
	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// WorkerAIPayload defines the JSON structure for the API request body.
type WorkerAIPayload struct {
	Prompt            string `json:"prompt"`
	NegativePrompt    string `json:"negative_prompt"`
	NumInferenceSteps int    `json:"num_inference_steps"`
	Width             int    `json:"width"`
	Height            int    `json:"height"`
}

// WorkerAIResult is nested inside the main response and holds the Base64 image.
type WorkerAIResult struct {
	Image string `json:"image"`
}

// WorkerAIResponse defines the top-level structure of the JSON response.
type WorkerAIResponse struct {
	Result  WorkerAIResult `json:"result"`
	Success bool           `json:"success"`
	Errors  []any          `json:"errors"`
}

var CloudFlareAIAccID = os.Getenv("CF_ACCID")
var CloudFlareAIAPIToken = os.Getenv("CF_API_TOKEN")

func requestQueue(record *core.Record) ([]byte, error) {

	const modelName = "@cf/black-forest-labs/flux-1-schnell"

	// Define the generation parameters
	payload := WorkerAIPayload{
		Prompt:            record.GetString("positive_prompt"),
		NegativePrompt:    record.GetString("negative_prompt"),
		NumInferenceSteps: record.GetInt("number"),
		Width:             record.GetInt("width"),
		Height:            record.GetInt("height"),
	}

	// --- 1. Setup API Call ---
	apiURL := fmt.Sprintf(
		"https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s",
		CloudFlareAIAccID,
		modelName,
	)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CloudFlareAIAPIToken))
	req.Header.Set("Content-Type", "application/json")

	// --- 2. Execute API Call ---
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request to Cloudflare: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// --- 3. Handle JSON Response ---
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("API request failed with status code %d. Response: %s\n", resp.StatusCode, string(bodyBytes))
		return nil, err
	}

	var aiResponse WorkerAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResponse); err != nil {
		fmt.Printf("Error decoding JSON response: %v\n", err)
		return nil, err
	}

	if !aiResponse.Success {
		fmt.Printf("API reported an error: %v\n", aiResponse.Errors)
		return nil, err
	}

	// --- 4. Decode Base64 and Save File ---
	// The image string starts with the image type header (e.g., "/9j/...")
	// We use the standard base64 decoder.
	imgBytes, err := base64.StdEncoding.DecodeString(aiResponse.Result.Image)
	if err != nil {
		fmt.Printf("Error decoding Base64 image string: %v\n", err)
		return nil, err
	}
	return imgBytes, nil

}

func generateImage(app *pocketbase.PocketBase, queueRecord *core.Record) {
	if CloudFlareAIAPIToken == "" || CloudFlareAIAccID == "" {
		log.Println("Cloudflare credentials not set!")
		return
	}

	queueRecord.Set("status", "processing")
	app.Save(queueRecord)

	generatedImages, err := app.FindRecordsByFilter("generated_images", "queue={:queueId}", "", 0, 0, dbx.Params{
		"queueId": queueRecord.Id,
	})
	if err != nil {
		log.Printf("Failed to find generated images from prompt [%s]: %v", queueRecord.Id, err)
	}

	toGenCount := queueRecord.GetInt("number") - len(generatedImages)

	log.Printf("Will request %d images for prompt %s", toGenCount, queueRecord.Id)
	generatedImageCollection, err := app.FindCollectionByNameOrId("generated_images")
	if err != nil {
		log.Printf("Failed to find generated_image collection: %v", err)
		return
	}

	for range toGenCount {
		imgBytes, err := requestQueue(queueRecord)
		if err != nil {
			log.Printf("Failed to send image generation request: %v", err)
			return
		}
		newRecord := core.NewRecord(generatedImageCollection)
		newRecord.Set("queue", queueRecord.Id)
		imgFile, _ := filesystem.NewFileFromBytes(imgBytes, fmt.Sprintf("%s.png", uuid.New().String()))
		newRecord.Set("image", imgFile)
		err = app.Save(newRecord)
		if err != nil {
			log.Printf("Failed to add create new image record: %v", err)
		}
	}

	queueRecord.Set("status", "finished")
	app.Save(queueRecord)
}

func GenerationRecover(app *pocketbase.PocketBase) {
	fileLock := flock.New("/var/lock/generate.lock")

	locked, err := fileLock.TryLock()
	if err != nil || !locked {
		log.Printf("Generation processing.")
		return
	}

	unfinishedQueues, err := app.FindRecordsByFilter("generate_queues", "status='queue' || status='processing'", "", 0, 0)
	if err != nil {
		log.Printf("Failed to find unfinished queues: %v", err)
	}
	for _, uq := range unfinishedQueues {
		generateImage(app, uq)
	}

	if locked {
		fileLock.Unlock()
	}
}

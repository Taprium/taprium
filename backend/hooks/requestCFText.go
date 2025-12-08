package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func requestText(prompt string) (string, error) {
	const modelName = "@cf/meta/llama-3.1-8b-instruct-fast"
	apiURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s", CloudFlareAIAccID, modelName)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": strings.ReplaceAll(prompt, `"`, `\"`)},
		},
		"max_tokens": 256, // Limit output to keep Neuron cost low
	})

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+CloudFlareAIAPIToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return "", fmt.Errorf("cloudflare AI API request failed with status code %d: %s", resp.StatusCode, string(body))

	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Extract the response text
	if res, ok := result["result"].(map[string]interface{}); ok {
		if response, ok := res["response"].(string); ok {
			// log.Printf("Failed to generate Cloudflare Text: %v", err)
			return response, nil
		}
	}

	return "", fmt.Errorf("failed to parse response from Cloudflare AI")
}

func GenerateText(app *pocketbase.PocketBase, queueRecord *core.Record) {
	if CloudFlareAIAPIToken == "" || CloudFlareAIAccID == "" {
		log.Println("Cloudflare credentials not set!")
		return
	}

	fileLock := flock.New("/var/lock/text.lock")

	locked, err := fileLock.TryLock()
	if err != nil || !locked {
		log.Printf("Generation processing.")
		return
	}
	defer func() {
		if locked {
			fileLock.Unlock()
		}
	}()

	queueRecord.Set("status", "processing")
	app.Save(queueRecord)

	if len(queueRecord.GetString("result")) != 0 {
		log.Printf("Text generated")
		return
	}

	result, err := requestText(queueRecord.GetString("user_prompt"))
	if err != nil {
		log.Printf("Failed to request cloudflare for text generation: %v", err)
		return
	}

	queueRecord.Set("result", result)
	queueRecord.Set("status", "finished")
	app.Save(queueRecord)
}

func TextGenerationRecover(app *pocketbase.PocketBase) {

	unfinishedQueues, err := app.FindRecordsByFilter("text_queues", " status='' || status='queue' || status='processing'", "", 0, 0)
	if err != nil {
		log.Printf("Failed to find unfinished queues: %v", err)
	}
	for _, uq := range unfinishedQueues {
		GenerateText(app, uq)
	}

}

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"bytes"

	"github.com/JxavierP/competitor-chat/internal/models"
	"github.com/JxavierP/competitor-chat/web/components"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	datastar "github.com/starfederation/datastar/sdk/go"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"google.golang.org/genai"
)

type PromptRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

func PromptHandler(c *gin.Context) {
	var request PromptRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	message := models.Message{
		Role:    "user",
		Content: request.Prompt,
	}

	sse := datastar.NewSSE(c.Writer, c.Request)
	fragment := components.ChatMessage(string(message.Role), message.Content, fmt.Sprintf("prompt-%s", uuid.New().String()))
	sse.MergeFragmentTempl(fragment, datastar.WithMergeMode("append"), datastar.WithSelector("#chat-messages"))
}

func ResponseHandler(gemini *genai.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JSON string from the 'datastar' query param
		jsonPromptPayload := c.Query("datastar")
		if jsonPromptPayload == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing datastar payload"})
			return
		}

		// Parse it into a map
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(jsonPromptPayload), &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid datastar JSON"})
			return
		}

		// Extract the prompt from the parsed map
		prompt := data["prompt"]
		if prompt == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt is required"})
			return
		}
		if prompt == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt is required"})
			return
		}
		sse := datastar.NewSSE(c.Writer, c.Request)
		jsonIDPayload := c.Query("id")
		jsonIDPayload = strings.TrimPrefix(jsonIDPayload, "prompt-")

		flush, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
			return
		}

		sse.MergeFragmentTempl(components.ChatMessage("assistant", "Thinking...", fmt.Sprintf("reponse-%s", jsonIDPayload)), datastar.WithMergeMode("append"), datastar.WithSelector("#chat-messages"))
		flush.Flush()

		contents := []*genai.Content{
			genai.NewContentFromText(prompt.(string), "user"),
		}

		int32ThinkingBudget := int32(0)
		var fullResponse string

		for response, err := range gemini.Models.GenerateContentStream(
			c.Request.Context(),
			"gemini-2.5-flash",
			contents,
			&genai.GenerateContentConfig{
				ThinkingConfig: &genai.ThinkingConfig{
					ThinkingBudget: &int32ThinkingBudget, // Disables thinking
				},
			}) {
			if err != nil {
				log.Fatal("Error generating content:", err)
			}

			chunk := response.Candidates[0].Content.Parts[0].Text
			fullResponse += chunk

			var md = goldmark.New(
				goldmark.WithExtensions(
					extension.GFM,
				),
			)

			var buf bytes.Buffer
			if err := md.Convert([]byte(fullResponse), &buf); err != nil {
				log.Println("markdown render error:", err)
				continue
			}

			htmlOutput := buf.String()

			sse.MergeFragmentTempl(
				components.ChatMessage("assistant", htmlOutput, fmt.Sprintf("reponse-%s", jsonIDPayload)),
				datastar.WithMergeMode("morph"),
				datastar.WithSelector(fmt.Sprintf("#reponse-%s", jsonIDPayload)),
			)
			flush.Flush()
		}

	}
}

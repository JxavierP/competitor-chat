package handlers

import (
	"fmt"
	"log"
	"net/http"

	"bytes"
	"github.com/JxavierP/competitor-chat/internal/models"
	"github.com/JxavierP/competitor-chat/web/components"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	datastar "github.com/starfederation/datastar/sdk/go"
	"github.com/yuin/goldmark"
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
	fragment := components.ChatMessage(string(message.Role), message.Content, uuid.New().String())
	sse.MergeFragmentTempl(fragment, datastar.WithMergeMode("append"), datastar.WithSelector("#chat-messages"))
}

func ResponseHandler(gemini *genai.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sse := datastar.NewSSE(c.Writer, c.Request)
		id := uuid.New().String()

		flush, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
			return
		}

		sse.MergeFragmentTempl(components.ChatMessage("assistant", "Thinking...", fmt.Sprintf("assistant-%s", id)), datastar.WithMergeMode("append"), datastar.WithSelector("#chat-messages"))
		flush.Flush()

		contents := []*genai.Content{
			genai.NewContentFromText("Send a response that uses all the features of markdown", "user"),
		}

		int32ThinkingBudget := int32(0) // Set to 0 to disable thinking
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

			var buf bytes.Buffer
			if err := goldmark.Convert([]byte(fullResponse), &buf); err != nil {
				log.Println("markdown render error:", err)
				continue
			}

			htmlOutput := buf.String()

			sse.MergeFragmentTempl(
				components.ChatMessage("assistant", htmlOutput, fmt.Sprintf("assistant-%s", id)),
				datastar.WithMergeMode("morph"),
				datastar.WithSelector(fmt.Sprintf("#assistant-%s", id)),
			)
			flush.Flush()
		}

	}
}

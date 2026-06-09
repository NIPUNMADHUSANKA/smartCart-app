package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"smartCart-app/database"
	"smartCart-app/models"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ChatMessage struct {
	role    string
	content string
}

type aiGeneratedItem struct {
	Item     string      `json:"item"`
	Quantity interface{} `json:"quantity"`
	Unit     string      `json:"unit"`
}

type aiGeneratedPayload struct {
	Prompt       string            `json:"prompt"`
	Intent       string            `json:"intent"`
	Category     string            `json:"category"`
	ShoppingList []aiGeneratedItem `json:"shopping_list"`
}

var OPEN_API_KEY string = os.Getenv("OPENAI_API_KEY")
var MODEL_NAME = openai.ChatModelGPT4_1Mini

func generateAIPrompt(userPrompt string) *ChatMessage {
	prompt := fmt.Sprintf(`
        You are an assistant for a smart grocery and shopping app called SmartCart.

        Your task is to convert a user's shopping intention into a structured shopping list.

        The intention may relate to cooking, cleaning, travel, daily needs, or general shopping.

        Rules:
        - Only include items that must be purchased.
        - Do NOT include instructions or explanations.
        - Return valid JSON only.
        - Do not include markdown or extra text.
        - Quantity must be numeric when possible; otherwise use "1".
        - Unit must be one of:
          kg, piece, pack, dozen, box, g, l, ml,
          bottle, can, cup, other

        Example:

        User input:
        I want to cook fried rice for 2 people

        Assistant output:
        {
          "prompt": "I want to cook fried rice for 2 people",
          "intent": "Cook Fried Rice",
          "category": "Cooking",
          "shopping_list": [
            { "item": "Rice", "quantity": "2", "unit": "cup" },
            { "item": "Eggs", "quantity": "2", "unit": "piece" },
            { "item": "Carrots", "quantity": "1", "unit": "cup" },
            { "item": "Green onions", "quantity": "1", "unit": "bunch" },
            { "item": "Soy sauce", "quantity": "2", "unit": "cup" },
            { "item": "Cooking oil", "quantity": "1", "unit": "bottle" }
          ]
        }

        User input:
        %s
        `, userPrompt)

	messages := &ChatMessage{
		role:    "system",
		content: prompt,
	}

	return messages
}

/*
String or []byte already available → json.Unmarshal
Common use cases
JSON stored in a string
JSON stored in a byte slice
Small API responses already loaded into memory

File, HTTP body, stream → json.NewDecoder().Decode()
Common use cases
HTTP requests
Files
Large JSON documents
Streaming JSON
*/
func parseAIGeneratedData(aiText string) (*aiGeneratedPayload, error) {
	var payload aiGeneratedPayload
	decoder := json.NewDecoder(strings.NewReader(aiText))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to parse AI response JSON: %w", err)
	}

	if payload.Prompt == "" || payload.Category == "" || len(payload.ShoppingList) == 0 {
		return nil, errors.New("AI response JSON is missing required fields")
	}

	return &payload, nil
}

func parseQuantity(value interface{}) (float32, error) {
	switch q := value.(type) {
	case float64:
		return float32(q), nil
	case string:
		if q == "" {
			return 1, nil
		}
		f, err := strconv.ParseFloat(q, 32)
		if err != nil {
			return 0, err
		}
		return float32(f), nil
	case json.Number:
		f, err := q.Float64()
		if err != nil {
			return 0, err
		}
		return float32(f), nil
	default:
		return 0, fmt.Errorf("unsupported quantity type %T", value)
	}
}

func normalizeUnit(unit string) models.UnitStatus {
	normalized := strings.ToLower(strings.TrimSpace(unit))
	switch normalized {
	case "kg":
		return models.Kg
	case "g", "gram", "grams":
		return models.G
	case "l", "liter", "litre", "liters", "litres":
		return models.L
	case "ml", "milliliter", "millilitre", "milliliters", "millilitres":
		return models.Ml
	case "piece", "pieces":
		return models.Piece
	case "pack", "packs":
		return models.Pack
	case "dozen":
		return models.Dozen
	case "box", "boxes":
		return models.Box
	case "bottle", "bottles":
		return models.Bottle
	case "can", "cans":
		return models.Can
	case "cup", "cups":
		return models.Cup
	default:
		return models.Other
	}
}

func saveAIGeneratedData(ctx context.Context, userId string, payload *aiGeneratedPayload) error {
	tx, err := database.DBPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	suggestionId := uuid.New()
	if _, err := tx.Exec(ctx,
		`INSERT INTO "AiSuggestion" ("id", "userId", "prompt") VALUES ($1, $2, $3)`,
		suggestionId,
		userId,
		payload.Prompt,
	); err != nil {
		return err
	}

	categoryId := uuid.New()
	if _, err := tx.Exec(ctx,
		`INSERT INTO "AICategory" ("id", "suggestionId", "categoryName", "priority") VALUES ($1, $2, $3, $4)`,
		categoryId,
		suggestionId,
		payload.Category,
		models.PriorityStatusNormal,
	); err != nil {
		return err
	}

	for _, item := range payload.ShoppingList {
		quantity, err := parseQuantity(item.Quantity)
		if err != nil {
			return fmt.Errorf("invalid quantity for item %q: %w", item.Item, err)
		}

		if item.Item == "" {
			return errors.New("AI response contains an item with empty name")
		}

		unit := normalizeUnit(item.Unit)
		itemId := uuid.New()
		if _, err := tx.Exec(ctx,
			`INSERT INTO "AIItem" ("id", "categoryId", "itemName", "quantity", "unit", "priority") VALUES ($1, $2, $3, $4, $5, $6)`,
			itemId,
			categoryId,
			item.Item,
			quantity,
			unit,
			models.PriorityStatusNormal,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func processAIGenratedData(ctx context.Context, userId string, result []openai.ChatCompletionChoice) error {
	if len(result) == 0 {
		return errors.New("no AI response returned")
	}

	aiText := result[0].Message.Content
	if aiText == "" {
		return errors.New("AI response content is empty")
	}

	parsed, err := parseAIGeneratedData(aiText)
	if err != nil {
		return err
	}

	return saveAIGeneratedData(ctx, userId, parsed)
}

func GenerateAI(ctx context.Context, userPrompt, userId string) error {
	ChatMessage := generateAIPrompt(userPrompt)

	if OPEN_API_KEY == "" {
		return errors.New("The OPENAI_API_KEY is not configured on the server")
	}

	if ChatMessage == nil {
		return errors.New("prompt cannot be generated, please contact the system administrator")
	}

	client := openai.NewClient(
		option.WithAPIKey(OPEN_API_KEY),
	)

	resp, err := client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: MODEL_NAME,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(ChatMessage.content),
			},
		},
	)

	if err != nil {
		return err
	}

	return processAIGenratedData(ctx, userId, resp.Choices)
}

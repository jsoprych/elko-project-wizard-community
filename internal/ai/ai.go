// Package ai provides AI-assisted directive generation.
// Community edition: builds a structured "prompt sandwich" for the user to
// paste into Claude/ChatGPT/Grok/DeepSeek.
// Private edition: shells out to the Claude API (set ANTHROPIC_API_KEY).
package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

// PromptRequest describes what the user wants generated.
type PromptRequest struct {
	Requirement string `json:"requirement"` // e.g. "enforce 12-factor app principles"
	ProjectType string `json:"project_type"`
	TechStack   string `json:"tech_stack"`
	Category    string `json:"category"` // directive category to target
}

// PromptSandwich is the structured JSON the user pastes into any LLM.
type PromptSandwich struct {
	Instructions string          `json:"instructions"`
	Context      PromptRequest   `json:"context"`
	OutputFormat OutputFormat    `json:"output_format"`
	Example      json.RawMessage `json:"example"`
}

// OutputFormat tells the LLM what shape to return.
type OutputFormat struct {
	Type   string          `json:"type"`
	Schema json.RawMessage `json:"schema"`
}

const sandwichTmpl = `You are an expert in AI agent development and project configuration.
You specialize in writing AGENTS.md directives that guide AI coding agents
(Claude Code, OpenCode, Cursor, Zed) to follow specific rules and policies.

Task: Generate a custom directive JSON object for the following requirement.

Context:
  - Requirement : {{.Requirement}}
  - Project type: {{.ProjectType}}
  - Tech stack  : {{.TechStack}}
  - Category    : {{.Category}}

Return ONLY a valid JSON object matching this schema (no markdown fences):
{
  "id":          "kebab-case-unique-id",
  "name":        "Human readable name",
  "category":    "{{.Category}}",
  "description": "One sentence description",
  "content":     "Markdown content — actionable rules the AI agent will follow"
}

The content field must be markdown-formatted and start with a ## heading.
Be specific, actionable, and concise. Avoid vague advice.`

// BuildPromptSandwich returns a structured prompt the user can paste into any LLM.
func BuildPromptSandwich(req *PromptRequest) (string, error) {
	t, err := template.New("sandwich").Parse(sandwichTmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, req); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateViaAPI calls the Claude API if ANTHROPIC_API_KEY is set.
// Returns the directive JSON string or an error.
func GenerateViaAPI(req *PromptRequest) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY not set — use prompt sandwich mode instead")
	}

	prompt, err := BuildPromptSandwich(req)
	if err != nil {
		return "", err
	}

	payload := map[string]interface{}{
		"model":      "claude-sonnet-4-6",
		"max_tokens": 1024,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	body, _ := json.Marshal(payload)

	return callClaudeAPI(apiKey, body)
}

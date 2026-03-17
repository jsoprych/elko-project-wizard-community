package ai

import (
	"strings"
	"testing"
)

func TestBuildPromptSandwich(t *testing.T) {
	req := &PromptRequest{
		Requirement: "enforce 12-factor app principles",
		ProjectType: "API / Microservice",
		TechStack:   "go",
		Category:    "workflow",
	}
	prompt, err := BuildPromptSandwich(req)
	if err != nil {
		t.Fatalf("BuildPromptSandwich: %v", err)
	}
	if !strings.Contains(prompt, "12-factor app principles") {
		t.Error("expected requirement in prompt")
	}
	if !strings.Contains(prompt, "workflow") {
		t.Error("expected category in prompt")
	}
	if !strings.Contains(prompt, `"content"`) {
		t.Error("expected output schema field 'content' in prompt")
	}
}

func TestGenerateViaAPINoKey(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")
	_, err := GenerateViaAPI(&PromptRequest{Requirement: "test"})
	if err == nil {
		t.Error("expected error with no API key")
	}
	if !strings.Contains(err.Error(), "ANTHROPIC_API_KEY") {
		t.Errorf("expected key hint in error, got: %v", err)
	}
}

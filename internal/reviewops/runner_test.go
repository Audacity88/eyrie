package reviewops

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestLocalDraftRunnerUsesSourceContext(t *testing.T) {
	runner := LocalDraftRunner{}
	result, err := runner.RunTask(context.Background(), DraftInput{
		Task: Task{
			ProjectID:    "p1",
			Kind:         "triage_issue",
			Repo:         "owner/repo",
			TargetNumber: 42,
		},
		SourceContext: &SourceContext{
			Repo:         "owner/repo",
			TargetNumber: 42,
			TargetType:   "issue",
			Title:        "Bug: review ops context",
			State:        "open",
			Author:       "alice",
			URL:          "https://github.com/owner/repo/issues/42",
			Labels:       []string{"bug", "triage"},
			BodyExcerpt:  "The issue body.",
			Comments: []ContextComment{
				{Author: "bob", BodyExcerpt: "I can reproduce this.", CreatedAt: "2026-04-29T00:00:00Z"},
			},
			FetchedAt: time.Now(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || result.Content == "" {
		t.Fatal("expected draft content")
	}
	for _, want := range []string{
		"Bug: review ops context",
		"labels: bug, triage",
		"bob: I can reproduce this.",
		"No GitHub mutations were performed",
	} {
		if !strings.Contains(result.Content, want) {
			t.Fatalf("draft missing %q:\n%s", want, result.Content)
		}
	}
}

func TestLocalDraftRunnerNotesMissingSourceContext(t *testing.T) {
	runner := LocalDraftRunner{}
	result, err := runner.RunTask(context.Background(), DraftInput{
		Task: Task{
			ProjectID:    "p1",
			Kind:         "review_pr",
			Repo:         "owner/repo",
			TargetNumber: 99,
		},
		SourceFetchError: "GitHub API returned 404",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Content, "Source context was unavailable") {
		t.Fatalf("draft should explain missing context:\n%s", result.Content)
	}
	if !strings.Contains(result.Content, "GitHub API returned 404") {
		t.Fatalf("draft should include fetch note:\n%s", result.Content)
	}
}

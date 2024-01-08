package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Distributed-Lab-Testing/integration-tests/config"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func updateNote(config config.TestsConfig, noteID, newContent string) (response *http.Response, err error) {
	url := config.FromExampleSvc().WithSuffix("notes/" + noteID).MustBuild()

	var updatedData = struct {
		Content string `json:"content"`
	}{
		Content: newContent,
	}

	updatedDataJSON, err := json.Marshal(updatedData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal updated data")
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(updatedDataJSON))
	if err != nil {
		return nil, errors.Wrap(err, "unable to form the new request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}

	return resp, nil
}

func testNoteUpdating(t *testing.T) {
	dsn := "host=example user=example password=example dbname=example sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err)

	var originalContent string
	var createdAt time.Time
	err = db.QueryRow(`SELECT content, created_at FROM notes WHERE id = $1`, globalNoteID).Scan(&originalContent, &createdAt)
	require.NoError(t, err)
	t.Logf("Original note content: %s, created at: %s", originalContent, createdAt.Format(time.RFC3339))

	startTime := time.Now()
	fmt.Printf("Test started at: %s\n", startTime)

	newContent := "Updated Content"
	cfg := config.NewTestsEnvConfig()

	resp, err := updateNote(cfg, globalNoteID, newContent)
	if err != nil {
		t.Fatalf("Failed to update note: %v", err)
	}

	fmt.Printf("Operation Status: %s\n", resp.Status)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status NoContent, got %v", resp.Status)
	}

	endTime := time.Now()
	fmt.Printf("Test ended at: %s\n", endTime)

	var updatedContent string
	var updatedAt time.Time
	err = db.QueryRow(`SELECT content, created_at FROM notes WHERE id = $1`, globalNoteID).Scan(&updatedContent, &updatedAt)
	require.NoError(t, err)
	t.Logf("Updated note content: %s, updated at: %s", updatedContent, updatedAt.Format(time.RFC3339))

}

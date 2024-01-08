package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/Distributed-Lab-Testing/integration-tests/config"
	"github.com/Distributed-Lab-Testing/integration-tests/resources"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var globalNoteID string

func createNote(config config.TestsConfig, content string, created_at string) (list []string, response *http.Response, err error) {
	var credentials resources.Note

	credentials.Data.Attributes.Content = content
	credentials.Data.Attributes.CreatedAt = created_at

	credentialsJSON, err := json.Marshal(credentials)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to marshal json with credentials")
	}

	url := config.FromExampleSvc().WithSuffix("notes").MustBuild()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(credentialsJSON))
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to form the new request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to read response body")
	}

	var noteResp resources.NoteIDResponse
	if err = json.Unmarshal(body, &noteResp); err != nil {
		return nil, nil, errors.Wrap(err, "failed to unmarshal response")
	}

	list = append(list, noteResp.Data.ID)

	return list, resp, nil
}

func TestNotes(t *testing.T) {
	t.Run("Inserting", testNoteCreation)
	t.Run("Updating db", testNoteUpdating)
	t.Run("Deleting from db", testNoteDeletion)
}

func testNoteCreation(t *testing.T) {
	var (
		content    = "test note"
		created_at = time.Now().Format(time.RFC3339)
		cfg        = config.NewTestsEnvConfig()
	)

	list, resp, err := createNote(cfg, content, created_at)
	require.NoError(t, err)

	t.Log("Response Status:", resp.Status)
	t.Log("Created note ID:", list)

	dsn := "host=upstream port=5433 user=example password=example dbname=example sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err)

	globalNoteID = list[0]

	err = db.QueryRow(`SELECT content, created_at FROM notes WHERE id = $1`, globalNoteID).Scan(&content, &created_at)
	require.NoError(t, err)
	t.Logf("Retrieved note content: %s, created at: %s", content, created_at)

	require.NotEmpty(t, list, "Note should not be empty")
}

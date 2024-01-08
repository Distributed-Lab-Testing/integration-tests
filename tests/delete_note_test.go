package tests

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/Distributed-Lab-Testing/integration-tests/config"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func deleteNote(config config.TestsConfig, list string) (response *http.Response, err error) {
	url := config.FromExampleSvc().WithSuffix("notes/" + list).MustBuild() // Предполагается, что URL для удаления записи выглядит как /notes/{noteID}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to form the new request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()

	// Дополнительные проверки можно добавить здесь, например, проверка статуса ответа

	return resp, nil
}

func TestNoteDeletion(t *testing.T) {
	cfg := config.NewTestsEnvConfig()

	preDeletionTime := time.Now()
	t.Logf("Attempting to delete note at: %v", preDeletionTime)

	resp, err := deleteNote(cfg, globalNoteID)
	require.NoError(t, err)

	t.Log("Response Status on Deletion:", resp.Status)

	postDeletionTime := time.Now()
	t.Logf("Note deletion attempted at: %v", postDeletionTime)

	// Проверка на отсутствие записи в БД
	dsn := "host=localhost user=example password=example dbname=example sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	// Пингуем базу данных, чтобы убедиться, что подключение установлено
	err = db.Ping()
	require.NoError(t, err)

	// Пытаемся получить запись по ID
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM notes WHERE id = $1", globalNoteID).Scan(&count)
	require.NoError(t, err)

	t.Logf("Number of records with ID %s: %d", globalNoteID, count)
	// Проверяем, что запись была удалена (счетчик равен 0)
	require.Equal(t, 0, count, "The note should be deleted, but it still exists.")
}
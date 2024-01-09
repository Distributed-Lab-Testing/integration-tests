package tests

import "testing"

func TestNotes(t *testing.T) {
	t.Run("Inserting", testNoteCreation)
	t.Run("Updating db", testNoteUpdating)
	t.Run("Deleting from db", testNoteDeletion)
}

package repository

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeArtifactPathStrict(t *testing.T) {
	t.Run("accepts nested paths", func(t *testing.T) {
		got, err := normalizeArtifactPathStrict("group/artifact/1.0/artifact.jar")
		require.NoError(t, err)
		assert.Equal(t, "group/artifact/1.0/artifact.jar", got)
	})

	t.Run("rejects traversal", func(t *testing.T) {
		_, err := normalizeArtifactPathStrict("../evil.jar")
		require.ErrorIs(t, err, ErrInvalidArtifactPath)
	})

	t.Run("rejects absolute paths", func(t *testing.T) {
		_, err := normalizeArtifactPathStrict("/evil.jar")
		require.ErrorIs(t, err, ErrInvalidArtifactPath)
	})
}

func TestArtifactStorageSaveOpenAndDelete(t *testing.T) {
	root := t.TempDir()
	store, err := newArtifactStorage(root)
	require.NoError(t, err)

	repoName := "demo-repo"
	artifactPath := "com/example/app/1.0.0/app.jar"
	content := []byte("artifact-bytes")

	fullPath, size, checksums, err := store.saveArtifact(repoName, artifactPath, strings.NewReader(string(content)))
	require.NoError(t, err)
	assert.EqualValues(t, len(content), size)
	assert.FileExists(t, fullPath)
	assert.NotEmpty(t, checksums.MD5)
	assert.NotEmpty(t, checksums.SHA1)
	assert.NotEmpty(t, checksums.SHA256)

	reader, err := store.openArtifact(repoName, artifactPath)
	require.NoError(t, err)
	defer reader.Close()

	loaded, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, content, loaded)

	require.NoError(t, store.deleteArtifact(repoName, artifactPath))
	_, statErr := os.Stat(filepath.Clean(fullPath))
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestArtifactStorageRejectsTraversal(t *testing.T) {
	store, err := newArtifactStorage(t.TempDir())
	require.NoError(t, err)

	_, _, _, err = store.saveArtifact("demo-repo", "../evil.jar", strings.NewReader("x"))
	require.ErrorIs(t, err, ErrInvalidArtifactPath)
}

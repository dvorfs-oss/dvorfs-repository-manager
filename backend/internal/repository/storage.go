package repository

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"hash"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/datatypes"
)

var (
	ErrRepositoryNotFound       = errors.New("repository not found")
	ErrArtifactNotFound         = errors.New("artifact not found")
	ErrInvalidRepositoryName    = errors.New("invalid repository name")
	ErrInvalidArtifactPath      = errors.New("invalid artifact path")
	ErrUnsupportedRepositoryType = errors.New("unsupported repository type")
)

type artifactStorage struct {
	root string
}

type artifactChecksums struct {
	MD5    string `json:"md5"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
}

func (c artifactChecksums) toJSON() datatypes.JSON {
	payload, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return datatypes.JSON(payload)
}

func newArtifactStorage(root string) (*artifactStorage, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(absRoot, 0o755); err != nil {
		return nil, err
	}

	return &artifactStorage{root: absRoot}, nil
}

func defaultArtifactStorageRoot() string {
	if envRoot := strings.TrimSpace(os.Getenv("ARTIFACT_STORAGE_DIR")); envRoot != "" {
		return envRoot
	}
	if envRoot := strings.TrimSpace(os.Getenv("REPOSITORY_ARTIFACT_STORAGE_DIR")); envRoot != "" {
		return envRoot
	}
	return filepath.Join(".", "storage", "artifacts")
}

func (s *artifactStorage) ensureRepositoryRoot(repoName string) error {
	root, err := s.repoRoot(repoName)
	if err != nil {
		return err
	}
	return os.MkdirAll(root, 0o755)
}

func (s *artifactStorage) deleteRepositoryRoot(repoName string) error {
	root, err := s.repoRoot(repoName)
	if err != nil {
		return err
	}
	return os.RemoveAll(root)
}

func (s *artifactStorage) saveArtifact(repoName, artifactPath string, body io.Reader) (string, int64, artifactChecksums, error) {
	target, err := s.artifactFilePath(repoName, artifactPath)
	if err != nil {
		return "", 0, artifactChecksums{}, err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", 0, artifactChecksums{}, err
	}

	tempFile, err := os.CreateTemp(filepath.Dir(target), ".artifact-*")
	if err != nil {
		return "", 0, artifactChecksums{}, err
	}
	cleanupTemp := func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
	}

	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()

	size, copyErr := io.Copy(io.MultiWriter(tempFile, md5Hash, sha1Hash, sha256Hash), body)
	if copyErr != nil {
		cleanupTemp()
		return "", 0, artifactChecksums{}, copyErr
	}

	if err := tempFile.Sync(); err != nil {
		cleanupTemp()
		return "", 0, artifactChecksums{}, err
	}

	if err := tempFile.Close(); err != nil {
		cleanupTemp()
		return "", 0, artifactChecksums{}, err
	}

	if err := os.Remove(target); err != nil && !errors.Is(err, os.ErrNotExist) {
		cleanupTemp()
		return "", 0, artifactChecksums{}, err
	}

	if err := os.Rename(tempFile.Name(), target); err != nil {
		cleanupTemp()
		return "", 0, artifactChecksums{}, err
	}

	return target, size, artifactChecksums{
		MD5:    hex.EncodeToString(md5Hash.Sum(nil)),
		SHA1:   hex.EncodeToString(sha1Hash.Sum(nil)),
		SHA256: hex.EncodeToString(sha256Hash.Sum(nil)),
	}, nil
}

func (s *artifactStorage) openArtifact(repoName, artifactPath string) (io.ReadCloser, error) {
	target, err := s.artifactFilePath(repoName, artifactPath)
	if err != nil {
		return nil, err
	}

	return os.Open(target)
}

func (s *artifactStorage) deleteArtifact(repoName, artifactPath string) error {
	target, err := s.artifactFilePath(repoName, artifactPath)
	if err != nil {
		return err
	}

	if err := os.Remove(target); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return err
		}
		return err
	}
	return nil
}

func (s *artifactStorage) repoRoot(repoName string) (string, error) {
	if err := validateRepositoryName(repoName); err != nil {
		return "", err
	}

	root := filepath.Clean(filepath.Join(s.root, "repositories", repoName))
	rel, err := filepath.Rel(s.root, root)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", ErrInvalidRepositoryName
	}
	return root, nil
}

func (s *artifactStorage) artifactFilePath(repoName, artifactPath string) (string, error) {
	root, err := s.repoRoot(repoName)
	if err != nil {
		return "", err
	}

	relativePath, err := normalizeArtifactPathStrict(artifactPath)
	if err != nil {
		return "", err
	}

	target := filepath.Clean(filepath.Join(root, filepath.FromSlash(relativePath)))
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", ErrInvalidArtifactPath
	}
	return target, nil
}

func validateRepositoryName(repoName string) error {
	cleaned := strings.TrimSpace(repoName)
	if cleaned == "" {
		return ErrInvalidRepositoryName
	}
	if strings.ContainsAny(cleaned, `\/:*?"<>|`) || cleaned == "." || cleaned == ".." {
		return ErrInvalidRepositoryName
	}
	return nil
}

func normalizeArtifactPath(artifactPath string) string {
	relativePath, err := normalizeArtifactPathStrict(artifactPath)
	if err != nil {
		return strings.TrimSpace(strings.ReplaceAll(artifactPath, `\`, "/"))
	}
	return relativePath
}

func normalizeArtifactPathStrict(artifactPath string) (string, error) {
	cleaned := strings.TrimSpace(strings.ReplaceAll(artifactPath, `\`, "/"))
	if cleaned == "" {
		return "", ErrInvalidArtifactPath
	}

	if strings.HasPrefix(cleaned, "/") || strings.HasPrefix(cleaned, "./") || strings.HasPrefix(cleaned, "../") {
		return "", ErrInvalidArtifactPath
	}

	parts := strings.Split(cleaned, "/")
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return "", ErrInvalidArtifactPath
		}
	}

	return path.Clean(cleaned), nil
}

func currentTime() time.Time {
	return time.Now().UTC()
}

func checksumWriterSet() (hash.Hash, hash.Hash, hash.Hash) {
	return md5.New(), sha1.New(), sha256.New()
}

func hashStrings(h hash.Hash) string {
	return hex.EncodeToString(h.Sum(nil))
}

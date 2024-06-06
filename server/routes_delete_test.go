package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/ollama/ollama/api"
)

func TestDelete(t *testing.T) {
	p := t.TempDir()
	t.Setenv("OLLAMA_MODELS", p)
	var s Server

	w := createRequest(t, s.CreateModelHandler, api.CreateRequest{
		Name:      "test",
		Modelfile: fmt.Sprintf("FROM %s", createBinFile(t)),
	})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code 200, actual %d", w.Code)
	}

	w = createRequest(t, s.CreateModelHandler, api.CreateRequest{
		Name:      "test2",
		Modelfile: fmt.Sprintf("FROM %s\nTEMPLATE {{ .System }} {{ .Prompt }}", createBinFile(t)),
	})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code 200, actual %d", w.Code)
	}

	checkFileExists(t, filepath.Join(p, "manifests", "*", "*", "*", "*"), []string{
		filepath.Join(p, "manifests", "registry.ollama.ai", "library", "test", "latest"),
		filepath.Join(p, "manifests", "registry.ollama.ai", "library", "test2", "latest"),
	})

	checkFileExists(t, filepath.Join(p, "blobs", "*"), []string{
		filepath.Join(p, "blobs", "sha256-8f2c2167d789c6b2302dff965160fa5029f6a24096d262c1cbb469f21a045382"),
		filepath.Join(p, "blobs", "sha256-a4e5e156ddec27e286f75328784d7106b60a4eb1d246e950a001a3f944fbda99"),
		filepath.Join(p, "blobs", "sha256-ca239d7bd8ea90e4a5d2e6bf88f8d74a47b14336e73eb4e18bed4dd325018116"),
		filepath.Join(p, "blobs", "sha256-fe7ac77b725cda2ccad03f88a880ecdfd7a33192d6cae08fce2c0ee1455991ed"),
	})

	w = createRequest(t, s.DeleteModelHandler, api.DeleteRequest{Name: "test"})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code 200, actual %d", w.Code)
	}

	checkFileExists(t, filepath.Join(p, "manifests", "*", "*", "*", "*"), []string{
		filepath.Join(p, "manifests", "registry.ollama.ai", "library", "test2", "latest"),
	})

	checkFileExists(t, filepath.Join(p, "blobs", "*"), []string{
		filepath.Join(p, "blobs", "sha256-8f2c2167d789c6b2302dff965160fa5029f6a24096d262c1cbb469f21a045382"),
		filepath.Join(p, "blobs", "sha256-a4e5e156ddec27e286f75328784d7106b60a4eb1d246e950a001a3f944fbda99"),
		filepath.Join(p, "blobs", "sha256-fe7ac77b725cda2ccad03f88a880ecdfd7a33192d6cae08fce2c0ee1455991ed"),
	})

	w = createRequest(t, s.DeleteModelHandler, api.DeleteRequest{Name: "test2"})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code 200, actual %d", w.Code)
	}

	checkFileExists(t, filepath.Join(p, "manifests", "*", "*", "*", "*"), []string{})
	checkFileExists(t, filepath.Join(p, "blobs", "*"), []string{})
}

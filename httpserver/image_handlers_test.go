package httpserver

import (
	"bytes"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ericls/imgdd/domainmodels"
)

func makePNGBytes(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func makeMultipartRequest(t *testing.T, fieldName, filename string, content []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, err := mw.CreateFormFile(fieldName, filename)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatalf("failed to write form file: %v", err)
	}
	mw.Close()
	req := httptest.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}

func TestUploadHandlerRejectsOversizedFile(t *testing.T) {
	pngBytes := makePNGBytes(10, 10)
	conf := &HttpServerConfigDef{
		AllowUpload:         true,
		ImageMaxUploadBytes: int64(len(pngBytes) - 1),
	}

	req := makeMultipartRequest(t, "image", "test.png", pngBytes)
	res := httptest.NewRecorder()

	// nil deps are safe here: the size check fires before auth or storage are touched
	handler := makeUploadHandler(conf, nil, nil, nil, nil)
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), "File too large") {
		t.Fatalf("expected 'File too large' in response, got: %s", res.Body.String())
	}
}

func TestUploadHandlerAcceptsFileWithinLimit(t *testing.T) {
	pngBytes := makePNGBytes(10, 10)
	storageDefRepo := newTestUploadStorageDefRepo()
	conf := &HttpServerConfigDef{
		AllowUpload:         true,
		ImageMaxUploadBytes: int64(len(pngBytes) + 1),
	}

	req := makeMultipartRequest(t, "image", "test.png", pngBytes)
	res := httptest.NewRecorder()

	handler := makeUploadHandler(conf, nil, storageDefRepo, nil, nil)
	handler.ServeHTTP(res, req)

	// Size check passes; failure at storage (no enabled storage definitions) means 500, not 400
	if res.Code == http.StatusBadRequest && strings.Contains(res.Body.String(), "File too large") {
		t.Fatalf("file within limit was wrongly rejected as too large")
	}
}

// testUploadStorageDefRepo returns an empty list so the handler fails at "no storage definitions"
// rather than panicking on a nil repo.
type testUploadStorageDefRepo struct{}

func newTestUploadStorageDefRepo() *testUploadStorageDefRepo {
	return &testUploadStorageDefRepo{}
}

func (r *testUploadStorageDefRepo) ListStorageDefinitions() ([]*domainmodels.StorageDefinition, error) {
	return nil, nil
}

func (r *testUploadStorageDefRepo) GetStorageDefinitionById(id string) (*domainmodels.StorageDefinition, error) {
	return nil, nil
}

func (r *testUploadStorageDefRepo) GetStorageDefinitionsByIds(ids []string) ([]*domainmodels.StorageDefinition, error) {
	return nil, nil
}

func (r *testUploadStorageDefRepo) GetStorageDefinitionByIdentifier(id string) (*domainmodels.StorageDefinition, error) {
	return nil, nil
}

func (r *testUploadStorageDefRepo) CreateStorageDefinition(storage_type string, config string, identifier string, isEnabled bool, priority int64) (*domainmodels.StorageDefinition, error) {
	return nil, nil
}

func (r *testUploadStorageDefRepo) UpdateStorageDefinition(identifier string, storage_type *string, config *string, isEnabled *bool, priority *int64) (*domainmodels.StorageDefinition, error) {
	return nil, nil
}

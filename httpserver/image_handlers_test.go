package httpserver

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/identity"
)

// guestContextUserManager always returns nil (unauthenticated guest).
type guestContextUserManager struct{}

func (m *guestContextUserManager) GetAuthenticationInfo(_ context.Context) *identity.AuthenticationInfo {
	return nil
}
func (m *guestContextUserManager) WithAuthenticationInfo(c context.Context, _ *identity.AuthenticationInfo) context.Context {
	return c
}
func (m *guestContextUserManager) SetAuthenticationInfo(_ context.Context, _ *identity.AuthenticationInfo) {
}

// authenticatedContextUserManager returns a fixed OrganizationUser.
type authenticatedContextUserManager struct {
	orgUser *domainmodels.OrganizationUser
}

func (m *authenticatedContextUserManager) GetAuthenticationInfo(_ context.Context) *identity.AuthenticationInfo {
	return &identity.AuthenticationInfo{
		AuthenticatedUser: &identity.AuthenticatedUser{User: m.orgUser.User},
		AuthorizedUser:    &identity.AuthorizedUser{OrganizationUser: m.orgUser},
	}
}
func (m *authenticatedContextUserManager) WithAuthenticationInfo(c context.Context, _ *identity.AuthenticationInfo) context.Context {
	return c
}
func (m *authenticatedContextUserManager) SetAuthenticationInfo(_ context.Context, _ *identity.AuthenticationInfo) {
}

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

func guestIdentityManager() *IdentityManager {
	return &IdentityManager{ContextUserManager: &guestContextUserManager{}}
}

func authenticatedIdentityManager(orgUser *domainmodels.OrganizationUser) *IdentityManager {
	return &IdentityManager{ContextUserManager: &authenticatedContextUserManager{orgUser: orgUser}}
}

// testUploadStorageDefRepo returns an empty list so the handler fails at "no storage definitions"
// rather than panicking on a nil repo.
type testUploadStorageDefRepo struct{}

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

// assertFileTooLarge checks the response indicates the file size was rejected.
func assertFileTooLarge(t *testing.T, res *httptest.ResponseRecorder) {
	t.Helper()
	if res.Code != http.StatusBadRequest || !strings.Contains(res.Body.String(), "File too large") {
		t.Fatalf("expected 400 File too large, got %d: %s", res.Code, res.Body.String())
	}
}

// assertNotFileTooLarge checks the response was NOT rejected due to size.
func assertNotFileTooLarge(t *testing.T, res *httptest.ResponseRecorder) {
	t.Helper()
	if res.Code == http.StatusBadRequest && strings.Contains(res.Body.String(), "File too large") {
		t.Fatalf("file was wrongly rejected as too large")
	}
}

// --- Global max ---

func TestUploadHandlerRejectsOversizedFile(t *testing.T) {
	pngBytes := makePNGBytes(10, 10)
	conf := &HttpServerConfigDef{
		AllowUpload:         true,
		ImageMaxUploadBytes: int64(len(pngBytes) - 1),
	}
	req := makeMultipartRequest(t, "image", "test.png", pngBytes)
	res := httptest.NewRecorder()
	makeUploadHandler(conf, guestIdentityManager(), nil, nil, nil).ServeHTTP(res, req)
	assertFileTooLarge(t, res)
}

func TestUploadHandlerAcceptsFileWithinGlobalLimit(t *testing.T) {
	pngBytes := makePNGBytes(10, 10)
	conf := &HttpServerConfigDef{
		AllowUpload:         true,
		ImageMaxUploadBytes: int64(len(pngBytes) + 1),
	}
	req := makeMultipartRequest(t, "image", "test.png", pngBytes)
	res := httptest.NewRecorder()
	makeUploadHandler(conf, guestIdentityManager(), &testUploadStorageDefRepo{}, nil, nil).ServeHTTP(res, req)
	assertNotFileTooLarge(t, res)
}

// --- Guest limit ---

func TestUploadHandlerRejectsGuestOverGuestLimit(t *testing.T) {
	pngBytes := makePNGBytes(10, 10)
	conf := &HttpServerConfigDef{
		AllowUpload:              true,
		ImageMaxUploadBytes:      int64(len(pngBytes) + 100),
		GuestImageMaxUploadBytes: int64(len(pngBytes) - 1),
	}
	req := makeMultipartRequest(t, "image", "test.png", pngBytes)
	res := httptest.NewRecorder()
	makeUploadHandler(conf, guestIdentityManager(), nil, nil, nil).ServeHTTP(res, req)
	assertFileTooLarge(t, res)
}

func TestUploadHandlerAcceptsGuestWithinGuestLimit(t *testing.T) {
	pngBytes := makePNGBytes(10, 10)
	conf := &HttpServerConfigDef{
		AllowUpload:              true,
		ImageMaxUploadBytes:      int64(len(pngBytes) + 100),
		GuestImageMaxUploadBytes: int64(len(pngBytes) + 1),
	}
	req := makeMultipartRequest(t, "image", "test.png", pngBytes)
	res := httptest.NewRecorder()
	makeUploadHandler(conf, guestIdentityManager(), &testUploadStorageDefRepo{}, nil, nil).ServeHTTP(res, req)
	assertNotFileTooLarge(t, res)
}

func TestUploadHandlerAuthenticatedUserFallsBackToGuestLimit(t *testing.T) {
	// Authenticated users with no per-user override use the guest limit as their base.
	pngBytes := makePNGBytes(10, 10)
	conf := &HttpServerConfigDef{
		AllowUpload:              true,
		ImageMaxUploadBytes:      int64(len(pngBytes) + 100),
		GuestImageMaxUploadBytes: int64(len(pngBytes) - 1),
	}
	orgUser := &domainmodels.OrganizationUser{
		Id:   "user-1",
		User: &domainmodels.User{Id: "user-1"},
	}
	req := makeMultipartRequest(t, "image", "test.png", pngBytes)
	res := httptest.NewRecorder()
	makeUploadHandler(conf, authenticatedIdentityManager(orgUser), nil, nil, nil).ServeHTTP(res, req)
	assertFileTooLarge(t, res)
}

func TestUploadHandlerPerUserLimitCanExceedGuestLimit(t *testing.T) {
	// A per-user override allows uploading above the guest limit (up to global max).
	pngBytes := makePNGBytes(10, 10)
	perUserLimit := int64(len(pngBytes) + 1)
	conf := &HttpServerConfigDef{
		AllowUpload:              true,
		ImageMaxUploadBytes:      int64(len(pngBytes) + 100),
		GuestImageMaxUploadBytes: int64(len(pngBytes) - 1), // would reject without override
	}
	orgUser := &domainmodels.OrganizationUser{
		Id:               "user-1",
		User:             &domainmodels.User{Id: "user-1"},
		UploadLimitBytes: &perUserLimit,
	}
	req := makeMultipartRequest(t, "image", "test.png", pngBytes)
	res := httptest.NewRecorder()
	makeUploadHandler(conf, authenticatedIdentityManager(orgUser), &testUploadStorageDefRepo{}, nil, nil).ServeHTTP(res, req)
	assertNotFileTooLarge(t, res)
}

// --- Per-user limit ---

func TestUploadHandlerRejectsUserOverPerUserLimit(t *testing.T) {
	pngBytes := makePNGBytes(10, 10)
	perUserLimit := int64(len(pngBytes) - 1)
	conf := &HttpServerConfigDef{
		AllowUpload:         true,
		ImageMaxUploadBytes: int64(len(pngBytes) + 100),
	}
	orgUser := &domainmodels.OrganizationUser{
		Id:               "user-1",
		User:             &domainmodels.User{Id: "user-1"},
		UploadLimitBytes: &perUserLimit,
	}
	req := makeMultipartRequest(t, "image", "test.png", pngBytes)
	res := httptest.NewRecorder()
	makeUploadHandler(conf, authenticatedIdentityManager(orgUser), nil, nil, nil).ServeHTTP(res, req)
	assertFileTooLarge(t, res)
}

func TestUploadHandlerAcceptsUserWithinPerUserLimit(t *testing.T) {
	pngBytes := makePNGBytes(10, 10)
	perUserLimit := int64(len(pngBytes) + 1)
	conf := &HttpServerConfigDef{
		AllowUpload:         true,
		ImageMaxUploadBytes: int64(len(pngBytes) + 100),
	}
	orgUser := &domainmodels.OrganizationUser{
		Id:               "user-1",
		User:             &domainmodels.User{Id: "user-1"},
		UploadLimitBytes: &perUserLimit,
	}
	req := makeMultipartRequest(t, "image", "test.png", pngBytes)
	res := httptest.NewRecorder()
	makeUploadHandler(conf, authenticatedIdentityManager(orgUser), &testUploadStorageDefRepo{}, nil, nil).ServeHTTP(res, req)
	assertNotFileTooLarge(t, res)
}

func TestUploadHandlerPerUserLimitCappedByGlobalMax(t *testing.T) {
	pngBytes := makePNGBytes(10, 10)
	// Per-user limit set higher than global max — global max should still win.
	perUserLimit := int64(len(pngBytes) + 100)
	conf := &HttpServerConfigDef{
		AllowUpload:         true,
		ImageMaxUploadBytes: int64(len(pngBytes) - 1),
	}
	orgUser := &domainmodels.OrganizationUser{
		Id:               "user-1",
		User:             &domainmodels.User{Id: "user-1"},
		UploadLimitBytes: &perUserLimit,
	}
	req := makeMultipartRequest(t, "image", "test.png", pngBytes)
	res := httptest.NewRecorder()
	makeUploadHandler(conf, authenticatedIdentityManager(orgUser), nil, nil, nil).ServeHTTP(res, req)
	assertFileTooLarge(t, res)
}

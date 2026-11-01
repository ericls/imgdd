package httpserver

import (
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/storage"
	"github.com/gorilla/mux"
)

const redirectTestImageId = "00000000-0000-0000-0000-000000000003"

func makeRedirectTestRouter(webUIHost string, storedImage *domainmodels.StoredImage) *mux.Router {
	storageDefRepo := storage.NewInMemoryStorageDefRepo()
	storageDefRepo.AddStorageDefinition(&domainmodels.StorageDefinition{
		Id:          "00000000-0000-0000-0000-000000000001",
		Identifier:  "fs1",
		StorageType: domainmodels.FSStorageType,
		IsEnabled:   true,
	})
	storedImageRepo := &testStoredImageRepo{storedImage: storedImage}

	rootRouter := mux.NewRouter()
	mountImageRedirectRoutes(rootRouter, webUIHost, storageDefRepo, storedImageRepo)
	return rootRouter
}

func TestImageRedirectHandlers(t *testing.T) {
	storedImage := &domainmodels.StoredImage{
		Id:                  "00000000-0000-0000-0000-000000000002",
		FileIdentifier:      "stored.png",
		StorageDefinitionId: "00000000-0000-0000-0000-000000000001",
		Image:               &domainmodels.Image{Id: redirectTestImageId},
	}
	cases := []struct {
		name         string
		storedImage  *domainmodels.StoredImage
		path         string
		expectedCode int
	}{
		{"canonical", storedImage, "/i/image/abc.png", http.StatusFound},
		{"canonical unknown extension", storedImage, "/i/image/abc.notanext", http.StatusNotFound},
		{"canonical missing image", nil, "/i/image/abc.png", http.StatusNotFound},
		{"direct", storedImage, "/i/direct/fs1.stored.png", http.StatusFound},
		{"direct malformed", storedImage, "/i/direct/fs1", http.StatusNotFound},
		{"direct unknown storage", storedImage, "/i/direct/fs2.stored.png", http.StatusNotFound},
		{"direct unknown file", storedImage, "/i/direct/fs1.other.png", http.StatusNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := makeRedirectTestRouter("", tc.storedImage)
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			if res.Code != tc.expectedCode {
				t.Fatalf("expected status %d, got %d", tc.expectedCode, res.Code)
			}
			if tc.expectedCode == http.StatusFound {
				if loc := res.Header().Get("Location"); loc != imageDetailPathPrefix+redirectTestImageId {
					t.Fatalf("expected redirect to image detail page, got %q", loc)
				}
			}
		})
	}
}

func TestImageRedirectUsesWebUIHost(t *testing.T) {
	storedImage := &domainmodels.StoredImage{
		Id:                  "00000000-0000-0000-0000-000000000002",
		FileIdentifier:      "stored.png",
		StorageDefinitionId: "00000000-0000-0000-0000-000000000001",
		Image:               &domainmodels.Image{Id: redirectTestImageId},
	}
	cases := []struct {
		name             string
		webUIHost        string
		path             string
		forwardedProto   string
		expectedLocation string
	}{
		{"canonical http", "imgdd.example", "/i/image/abc.png", "", "http://imgdd.example/images/" + redirectTestImageId},
		{"canonical https", "imgdd.example", "/i/image/abc.png", "https", "https://imgdd.example/images/" + redirectTestImageId},
		{"direct https", "imgdd.example:8080", "/i/direct/fs1.stored.png", "https", "https://imgdd.example:8080/images/" + redirectTestImageId},
		{"no host", "", "/i/direct/fs1.stored.png", "https", imageDetailPathPrefix + redirectTestImageId},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := makeRedirectTestRouter(tc.webUIHost, storedImage)
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			if tc.forwardedProto != "" {
				req.Header.Set("X-Forwarded-Proto", tc.forwardedProto)
			}
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			if res.Code != http.StatusFound {
				t.Fatalf("expected status %d, got %d", http.StatusFound, res.Code)
			}
			if loc := res.Header().Get("Location"); loc != tc.expectedLocation {
				t.Fatalf("expected redirect to %q, got %q", tc.expectedLocation, loc)
			}
		})
	}
}

// The redirect destination must be the web client's public image detail route. Hopefully.
func TestImageDetailPathMatchesWebClientRoute(t *testing.T) {
	routesSrc, err := os.ReadFile("../web_client/src/routes.ts")
	if err != nil {
		t.Fatalf("failed to read web client routes: %v", err)
	}
	segment := regexp.MustCompile(`(?s)routeSegments = \{.*?\bimages: "([^"]+)"`).FindSubmatch(routesSrc)
	if segment == nil {
		t.Fatal("routeSegments.images not found in web_client/src/routes.ts")
	}
	if expected := "/" + string(segment[1]) + "/"; imageDetailPathPrefix != expected {
		t.Fatalf("imageDetailPathPrefix is %q, but the web client's images segment gives %q", imageDetailPathPrefix, expected)
	}
	if !regexp.MustCompile(`\bimage: \(imageId: string\) =>\s*` + "`" + `/\$\{routeSegments\.images\}/\$\{imageId\}` + "`").Match(routesSrc) {
		t.Fatal("routes.image in web_client/src/routes.ts is no longer /${routeSegments.images}/${imageId}")
	}

	entrySrc, err := os.ReadFile("../web_client/src/entry.tsx")
	if err != nil {
		t.Fatalf("failed to read web client entry: %v", err)
	}
	publicRoute := regexp.MustCompile(`path: ` + "`" + `\$\{routeSegments\.images\}/:imageId` + "`" + `,\s*lazy: async \(\) => \{\s*const \{ PublicImageDetail \}`)
	if !publicRoute.Match(entrySrc) {
		t.Fatal("web_client/src/entry.tsx no longer serves PublicImageDetail at ${routeSegments.images}/:imageId")
	}
}

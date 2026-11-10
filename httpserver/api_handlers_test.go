package httpserver_test

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/httpserver"
	"github.com/ericls/imgdd/identity"
	imgddimage "github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/ratelimit"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/test_support"

	"github.com/stretchr/testify/require"
)

type apiTestFile struct {
	field    string
	filename string
	content  []byte
}

func pngBytes(t *testing.T, w, h int) []byte {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, w, h))))
	return buf.Bytes()
}

func newCreateImagesRequest(t *testing.T, files ...apiTestFile) *http.Request {
	t.Helper()
	return newCreateImagesRequestWithFields(t, nil, files...)
}

func newCreateImagesRequestWithFields(t *testing.T, fields [][2]string, files ...apiTestFile) *http.Request {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	for _, field := range fields {
		require.NoError(t, mw.WriteField(field[0], field[1]))
	}
	for _, f := range files {
		fw, err := mw.CreateFormFile(f.field, f.filename)
		require.NoError(t, err)
		_, err = fw.Write(f.content)
		require.NoError(t, err)
	}
	require.NoError(t, mw.Close())
	req := httptest.NewRequest(http.MethodPost, "/api/v1/createImages", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}

type createImagesResponse struct {
	Results []struct {
		Filename string
		Image    *struct {
			ID        string
			Filename  string
			URL       string
			MIMEType  string
			Width     int
			Height    int
			ByteSize  int
			ExpiresAt *time.Time
		}
		Error *struct {
			Code    string
			Message string
		}
	}
	Error *struct {
		Code    string
		Message string
	}
}

func decodeCreateImages(t *testing.T, res *httptest.ResponseRecorder) createImagesResponse {
	t.Helper()
	require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	var out createImagesResponse
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &out), res.Body.String())
	return out
}

func TestAPICreateImages(t *testing.T) {
	dbConf := TestServiceMan.GetDBConfig()
	test_support.ResetDatabase(dbConf)
	conn := db.GetConnection(dbConf)
	identityRepo := identity.NewDBIdentityRepo(conn)
	imageRepo := imgddimage.NewDBImageRepo(conn)
	storageDefRepo := storage.NewInMemoryStorageDefRepo()
	_, err := storageDefRepo.CreateStorageDefinition("fs", `{"mediaRoot": "`+t.TempDir()+`"}`, "apitest", true, 0)
	require.NoError(t, err)

	orgUser, err := identityRepo.CreateUserWithOrganization("api@home.arpa", "api", "pw")
	require.NoError(t, err)
	_, secret, err := identityRepo.CreateAPIKey(orgUser.Id, "test")
	require.NoError(t, err)

	img := pngBytes(t, 7, 5)
	conf := &httpserver.HttpServerConfigDef{
		AllowUpload:         true,
		ImageMaxUploadBytes: 1 << 20,
		DefaultURLFormat:    "canonical",
	}
	serve := func(conf *httpserver.HttpServerConfigDef, limiter *ratelimit.RateLimiter, req *http.Request) *httptest.ResponseRecorder {
		if limiter == nil {
			limiter = ratelimit.NewRateLimiter(100, 100)
		}
		res := httptest.NewRecorder()
		httpserver.MakeAPICreateImagesHandler(conf, identityRepo, storageDefRepo, imageRepo, limiter).ServeHTTP(res, req)
		return res
	}

	t.Run("anonymous upload", func(t *testing.T) {
		res := serve(conf, nil, newCreateImagesRequest(t, apiTestFile{"image", "a.png", img}))
		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		out := decodeCreateImages(t, res)
		require.Len(t, out.Results, 1)
		got := out.Results[0].Image
		require.NotNil(t, got)
		require.Equal(t, "a.png", got.Filename)
		require.Equal(t, "image/png", got.MIMEType)
		require.Equal(t, 7, got.Width)
		require.Equal(t, 5, got.Height)
		require.NotContains(t, res.Body.String(), `"identifier"`)

		stored, err := imageRepo.GetImageById(got.ID)
		require.NoError(t, err)
		require.Equal(t, "http://example.com/image/"+stored.Identifier+".png", got.URL)
		require.NotEqual(t, orgUser.Id, stored.CreatedById)
		require.Nil(t, got.ExpiresAt)
		require.Nil(t, stored.ExpiresAt)
		require.Contains(t, res.Body.String(), `"expiresAt":null`)
	})

	t.Run("expires_in applies to every image in the request", func(t *testing.T) {
		before := time.Now()
		res := serve(conf, nil, newCreateImagesRequestWithFields(t,
			[][2]string{{"expires_in", "3600"}},
			apiTestFile{"image", "a.png", img},
			apiTestFile{"image", "b.png", img},
		))
		after := time.Now()
		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		out := decodeCreateImages(t, res)
		require.Len(t, out.Results, 2)
		for _, result := range out.Results {
			require.NotNil(t, result.Image)
			require.NotNil(t, result.Image.ExpiresAt)
			require.WithinRange(t, *result.Image.ExpiresAt, before.Add(time.Hour), after.Add(time.Hour))
			stored, err := imageRepo.GetImageById(result.Image.ID)
			require.NoError(t, err)
			require.NotNil(t, stored.ExpiresAt)
			require.WithinDuration(t, *result.Image.ExpiresAt, *stored.ExpiresAt, time.Second)
		}
	})

	t.Run("expires_in bounds are inclusive", func(t *testing.T) {
		for _, value := range []string{"60", "31536000"} {
			res := serve(conf, nil, newCreateImagesRequestWithFields(t,
				[][2]string{{"expires_in", value}},
				apiTestFile{"image", "a.png", img},
			))
			require.Equal(t, http.StatusOK, res.Code, value)
		}
	})

	t.Run("invalid expires_in rejects the whole request", func(t *testing.T) {
		invalid := [][][2]string{
			{{"expires_in", ""}},
			{{"expires_in", "abc"}},
			{{"expires_in", "1.5"}},
			{{"expires_in", "-60"}},
			{{"expires_in", "0"}},
			{{"expires_in", "59"}},
			{{"expires_in", "31536001"}},
			{{"expires_in", "99999999999999999999"}},
			{{"expires_in", "60"}, {"expires_in", "120"}},
		}
		for _, fields := range invalid {
			res := serve(conf, nil, newCreateImagesRequestWithFields(t, fields, apiTestFile{"image", "a.png", img}))
			require.Equal(t, http.StatusBadRequest, res.Code, fields)
			out := decodeCreateImages(t, res)
			require.Equal(t, "invalid_expires_in", out.Error.Code, fields)
			require.Empty(t, out.Results)
		}
	})

	t.Run("authenticated upload is owned by the key's user", func(t *testing.T) {
		req := newCreateImagesRequest(t,
			apiTestFile{"image", "a.png", img},
			apiTestFile{"image", "b.png", img},
		)
		req.Header.Set("Authorization", "Bearer "+secret)
		res := serve(conf, nil, req)
		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		out := decodeCreateImages(t, res)
		require.Len(t, out.Results, 2)
		for i, name := range []string{"a.png", "b.png"} {
			require.Equal(t, name, out.Results[i].Filename)
			require.NotNil(t, out.Results[i].Image)
			stored, err := imageRepo.GetImageById(out.Results[i].Image.ID)
			require.NoError(t, err)
			require.Equal(t, orgUser.Id, stored.CreatedById)
		}
	})

	t.Run("authenticated uploads are not rate limited", func(t *testing.T) {
		limiter := ratelimit.NewRateLimiter(0, 0)
		for range 3 {
			req := newCreateImagesRequest(t, apiTestFile{"image", "a.png", img})
			req.Header.Set("Authorization", "Bearer "+secret)
			res := serve(conf, limiter, req)
			require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		}
	})

	t.Run("anonymous uploads are rate limited", func(t *testing.T) {
		limiter := ratelimit.NewRateLimiter(0, 0)
		// The limiter admits the first request it sees from an IP.
		res := serve(conf, limiter, newCreateImagesRequest(t, apiTestFile{"image", "a.png", img}))
		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		res = serve(conf, limiter, newCreateImagesRequest(t, apiTestFile{"image", "a.png", img}))
		require.Equal(t, http.StatusTooManyRequests, res.Code, res.Body.String())
		out := decodeCreateImages(t, res)
		require.Equal(t, "rate_limited", out.Results[0].Error.Code)
	})

	t.Run("invalid credentials are rejected, not downgraded to anonymous", func(t *testing.T) {
		for _, header := range []string{
			"Bearer imgdd_wrong",
			"Bearer ",
			"Basic " + secret,
			secret,
		} {
			req := newCreateImagesRequest(t, apiTestFile{"image", "a.png", img})
			req.Header.Set("Authorization", header)
			res := serve(conf, nil, req)
			require.Equal(t, http.StatusUnauthorized, res.Code, header)
			require.Equal(t, "invalid_api_key", decodeCreateImages(t, res).Error.Code)
		}
	})

	t.Run("revoked key is rejected", func(t *testing.T) {
		key, revokedSecret, err := identityRepo.CreateAPIKey(orgUser.Id, "revoked")
		require.NoError(t, err)
		require.NoError(t, identityRepo.RevokeAPIKey(orgUser.Id, key.Id))
		req := newCreateImagesRequest(t, apiTestFile{"image", "a.png", img})
		req.Header.Set("Authorization", "Bearer "+revokedSecret)
		res := serve(conf, nil, req)
		require.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("partial failure reports per-image errors", func(t *testing.T) {
		res := serve(conf, nil, newCreateImagesRequest(t,
			apiTestFile{"image", "bad.jpg", img},
			apiTestFile{"image", "good.png", img},
		))
		require.Equal(t, http.StatusOK, res.Code)
		out := decodeCreateImages(t, res)
		require.Equal(t, "mime_type_mismatch", out.Results[0].Error.Code)
		require.Nil(t, out.Results[0].Image)
		require.NotNil(t, out.Results[1].Image)
	})

	t.Run("all failed uses first failure status", func(t *testing.T) {
		small := *conf
		small.ImageMaxUploadBytes = int64(len(img) - 1)
		res := serve(&small, nil, newCreateImagesRequest(t, apiTestFile{"image", "a.png", img}))
		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, "file_too_large", decodeCreateImages(t, res).Results[0].Error.Code)
	})

	t.Run("undecodable images are rejected as invalid", func(t *testing.T) {
		for name, file := range map[string]apiTestFile{
			"truncated image": {"image", "a.png", img[:16]},
			"not an image":    {"image", "a.pdf", []byte("%PDF-1.7")},
		} {
			res := serve(conf, nil, newCreateImagesRequest(t, file))
			require.Equal(t, http.StatusBadRequest, res.Code, name)
			out := decodeCreateImages(t, res)
			require.Nil(t, out.Results[0].Image, name)
			require.Equal(t, "invalid_image", out.Results[0].Error.Code, name)
		}
	})

	t.Run("temporary files are removed when the query string is malformed", func(t *testing.T) {
		tmp := t.TempDir()
		t.Setenv("TMPDIR", tmp)
		small := *conf
		small.ImageMaxUploadBytes = 1024
		// Larger than ImageMaxUploadBytes, so the part is spooled to disk.
		req := newCreateImagesRequest(t, apiTestFile{"image", "a.png", make([]byte, 4096)})
		req.URL.RawQuery = "a=%zz"
		res := serve(&small, nil, req)
		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, "invalid_request", decodeCreateImages(t, res).Error.Code)
		entries, err := os.ReadDir(tmp)
		require.NoError(t, err)
		require.Empty(t, entries)
	})

	t.Run("request errors", func(t *testing.T) {
		res := serve(conf, nil, newCreateImagesRequest(t))
		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, "no_images", decodeCreateImages(t, res).Error.Code)

		res = serve(conf, nil, newCreateImagesRequest(t, apiTestFile{"file", "a.png", img}))
		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, "no_images", decodeCreateImages(t, res).Error.Code)

		many := make([]apiTestFile, 6)
		for i := range many {
			many[i] = apiTestFile{"image", "a.png", img}
		}
		res = serve(conf, nil, newCreateImagesRequest(t, many...))
		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, "too_many_images", decodeCreateImages(t, res).Error.Code)

		res = serve(conf, nil, httptest.NewRequest(http.MethodPost, "/api/v1/createImages", bytes.NewBufferString("{}")))
		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, "invalid_request", decodeCreateImages(t, res).Error.Code)

		res = serve(conf, nil, httptest.NewRequest(http.MethodGet, "/api/v1/createImages", nil))
		require.Equal(t, http.StatusMethodNotAllowed, res.Code)
		require.Equal(t, http.MethodPost, res.Header().Get("Allow"))

		disabled := *conf
		disabled.AllowUpload = false
		res = serve(&disabled, nil, newCreateImagesRequest(t, apiTestFile{"image", "a.png", img}))
		require.Equal(t, http.StatusForbidden, res.Code)
		require.Equal(t, "upload_disabled", decodeCreateImages(t, res).Error.Code)

		tiny := *conf
		tiny.ImageMaxUploadBytes = 10
		huge := make([]byte, 2<<20)
		res = serve(&tiny, nil, newCreateImagesRequest(t, apiTestFile{"image", "a.png", huge}))
		require.Equal(t, http.StatusRequestEntityTooLarge, res.Code)
		require.Equal(t, "request_too_large", decodeCreateImages(t, res).Error.Code)
	})
}

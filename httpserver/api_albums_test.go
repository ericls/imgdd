package httpserver_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ericls/imgdd/db"
	dm "github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/httpserver"
	"github.com/ericls/imgdd/identity"
	imgddimage "github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/ratelimit"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/test_support"

	"github.com/stretchr/testify/require"
)

type createAlbumResponse struct {
	Album *struct {
		ID    string
		URL   string
		Title string
	}
	Error *struct {
		Code    string
		Message string
	}
}

func decodeCreateAlbum(t *testing.T, res *httptest.ResponseRecorder) createAlbumResponse {
	t.Helper()
	require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	var out createAlbumResponse
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &out), res.Body.String())
	return out
}

func TestAPIAlbums(t *testing.T) {
	dbConf := TestServiceMan.GetDBConfig()
	test_support.ResetDatabase(dbConf)
	conn := db.GetConnection(dbConf)
	identityRepo := identity.NewDBIdentityRepo(conn)
	imageRepo := imgddimage.NewDBImageRepo(conn)
	storageDefRepo := storage.NewInMemoryStorageDefRepo()
	_, err := storageDefRepo.CreateStorageDefinition("fs", `{"mediaRoot": "`+t.TempDir()+`"}`, "albumtest", true, 0)
	require.NoError(t, err)

	owner, err := identityRepo.CreateUserWithOrganization("owner@home.arpa", "owner", "pw")
	require.NoError(t, err)
	_, ownerSecret, err := identityRepo.CreateAPIKey(owner.Id, "owner")
	require.NoError(t, err)
	// The site owner can manage everyone's images, but must still not be able
	// to add to someone else's album.
	siteOwner, err := identityRepo.CreateUserWithOrganization("site-owner@home.arpa", "site", "pw")
	require.NoError(t, err)
	require.NoError(t, identityRepo.AddRoleToOrganizationUser(siteOwner.Id, "site_owner"))
	require.True(t, identityRepo.GetOrganizationUserById(siteOwner.Id).IsSiteOwner())
	_, siteOwnerSecret, err := identityRepo.CreateAPIKey(siteOwner.Id, "site owner")
	require.NoError(t, err)

	img := pngBytes(t, 4, 3)
	conf := &httpserver.HttpServerConfigDef{
		AllowUpload:         true,
		ImageMaxUploadBytes: 1 << 20,
		DefaultURLFormat:    "canonical",
	}

	createAlbum := func(limiter *ratelimit.RateLimiter, secret string, contentType string, body string) *httptest.ResponseRecorder {
		t.Helper()
		if limiter == nil {
			limiter = ratelimit.NewRateLimiter(100, 100)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/v1/createAlbum", bytes.NewBufferString(body))
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		if secret != "" {
			req.Header.Set("Authorization", "Bearer "+secret)
		}
		res := httptest.NewRecorder()
		httpserver.MakeAPICreateAlbumHandler(conf, identityRepo, imageRepo, limiter).ServeHTTP(res, req)
		return res
	}
	newAlbum := func(title string) string {
		t.Helper()
		res := createAlbum(nil, ownerSecret, "application/json", `{"title":"`+title+`"}`)
		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		return decodeCreateAlbum(t, res).Album.ID
	}
	createImages := func(secret string, fields [][2]string, files ...apiTestFile) *httptest.ResponseRecorder {
		t.Helper()
		req := newCreateImagesRequestWithFields(t, fields, files...)
		if secret != "" {
			req.Header.Set("Authorization", "Bearer "+secret)
		}
		res := httptest.NewRecorder()
		httpserver.MakeAPICreateImagesHandler(conf, identityRepo, storageDefRepo, imageRepo, ratelimit.NewRateLimiter(100, 100)).ServeHTTP(res, req)
		return res
	}

	t.Run("create album from JSON, form, or empty body", func(t *testing.T) {
		for name, tc := range map[string]struct {
			contentType string
			body        string
			title       string
		}{
			"json":       {"application/json", `{"title":"  Fix login layout "}`, "Fix login layout"},
			"urlencoded": {"application/x-www-form-urlencoded", "title=" + url.QueryEscape("Fix login"), "Fix login"},
			"empty":      {"", "", ""},
			"empty json": {"application/json", "", ""},
		} {
			res := createAlbum(nil, ownerSecret, tc.contentType, tc.body)
			require.Equal(t, http.StatusOK, res.Code, name+": "+res.Body.String())
			out := decodeCreateAlbum(t, res)
			require.Equal(t, tc.title, out.Album.Title, name)
			require.Len(t, out.Album.ID, 22, name)
			require.Equal(t, "http://example.com/albums/"+out.Album.ID, out.Album.URL, name)

			stored, err := imageRepo.GetAlbumByIdentifier(out.Album.ID)
			require.NoError(t, err)
			require.Equal(t, owner.Id, stored.CreatedById, name)
			// The internal id is never exposed.
			require.NotContains(t, res.Body.String(), stored.Id, name)
		}
	})

	t.Run("page URLs use the configured web UI host", func(t *testing.T) {
		webUIConf := *conf
		webUIConf.WebUIHost = "app.example.net"
		req := httptest.NewRequest(http.MethodPost, "/api/v1/createAlbum", nil)
		req.Header.Set("Authorization", "Bearer "+ownerSecret)
		res := httptest.NewRecorder()
		httpserver.MakeAPICreateAlbumHandler(&webUIConf, identityRepo, imageRepo, ratelimit.NewRateLimiter(100, 100)).ServeHTTP(res, req)
		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		albumID := decodeCreateAlbum(t, res).Album.ID
		require.Equal(t, "http://app.example.net/albums/"+albumID, decodeCreateAlbum(t, res).Album.URL)

		req = newCreateImagesRequestWithFields(t, [][2]string{{"album", albumID}}, apiTestFile{"image", "shot.png", img})
		req.Header.Set("Authorization", "Bearer "+ownerSecret)
		req.Header.Set("X-Forwarded-Proto", "https")
		res = httptest.NewRecorder()
		httpserver.MakeAPICreateImagesHandler(&webUIConf, identityRepo, storageDefRepo, imageRepo, ratelimit.NewRateLimiter(100, 100)).ServeHTTP(res, req)
		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		got := decodeCreateImages(t, res).Results[0].Image
		require.NotNil(t, got)
		require.Equal(t, "https://app.example.net/albums/"+albumID, *got.AlbumURL)
		require.Equal(t, "https://app.example.net/images/"+got.ID, got.PageURL)
	})

	t.Run("create album requires an API key", func(t *testing.T) {
		for _, secret := range []string{"", "imgdd_wrong"} {
			res := createAlbum(nil, secret, "application/json", `{"title":"x"}`)
			require.Equal(t, http.StatusUnauthorized, res.Code)
			require.Equal(t, "invalid_api_key", decodeCreateAlbum(t, res).Error.Code)
		}
	})

	t.Run("create album request errors", func(t *testing.T) {
		res := createAlbum(nil, ownerSecret, "application/json", `{"title":`+`"`+strings.Repeat("a", 256)+`"}`)
		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, "invalid_title", decodeCreateAlbum(t, res).Error.Code)

		for _, tc := range [][2]string{
			{"application/json", `{"title":`},
			{"application/json", `{"title":1}`},
			{"application/json", `{"name":"misspelled field"}`},
			{"application/json", `{"title":"a"}{"title":"b"}`},
			{"text/plain", "title"},
			{"application/x-www-form-urlencoded", "title=a&title=b"},
		} {
			res := createAlbum(nil, ownerSecret, tc[0], tc[1])
			require.Equal(t, http.StatusBadRequest, res.Code, tc)
			require.Equal(t, "invalid_request", decodeCreateAlbum(t, res).Error.Code, tc)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/createAlbum", nil)
		res = httptest.NewRecorder()
		httpserver.MakeAPICreateAlbumHandler(conf, identityRepo, imageRepo, ratelimit.NewRateLimiter(100, 100)).ServeHTTP(res, req)
		require.Equal(t, http.StatusMethodNotAllowed, res.Code)
	})

	t.Run("create album is rate limited per user", func(t *testing.T) {
		limiter := ratelimit.NewRateLimiter(0, 0)
		// The limiter admits the first request it sees from a key.
		require.Equal(t, http.StatusOK, createAlbum(limiter, ownerSecret, "", "").Code)
		res := createAlbum(limiter, ownerSecret, "", "")
		require.Equal(t, http.StatusTooManyRequests, res.Code)
		require.Equal(t, "rate_limited", decodeCreateAlbum(t, res).Error.Code)
		// Other users have their own budget.
		require.Equal(t, http.StatusOK, createAlbum(limiter, siteOwnerSecret, "", "").Code)
	})

	t.Run("owner uploads into album with one caption per image", func(t *testing.T) {
		albumID := newAlbum("captions")
		res := createImages(ownerSecret,
			[][2]string{{"album", albumID}, {"caption", "before fix"}, {"caption", " after fix "}},
			apiTestFile{"image", "before.png", img},
			apiTestFile{"image", "after.png", img},
		)
		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		out := decodeCreateImages(t, res)
		require.Len(t, out.Results, 2)
		album := mustAlbum(t, imageRepo, albumID)
		for i, caption := range []string{"before fix", "after fix"} {
			got := out.Results[i].Image
			require.NotNil(t, got)
			require.Equal(t, caption, *got.Caption)
			require.Equal(t, "http://example.com/albums/"+albumID, *got.AlbumURL)
			require.Equal(t, "http://example.com/images/"+got.ID, got.PageURL)
		}
		images, err := imageRepo.ListAlbumImages(album.Id, 0)
		require.NoError(t, err)
		require.Len(t, images, 2)
		for i, caption := range []string{"before fix", "after fix"} {
			require.Equal(t, out.Results[i].Image.ID, images[i].Image.Id)
			require.Equal(t, caption, images[i].Caption)
		}
	})

	t.Run("one caption applies to every image", func(t *testing.T) {
		albumID := newAlbum("same caption")
		res := createImages(ownerSecret,
			[][2]string{{"album", albumID}, {"caption", "same"}},
			apiTestFile{"image", "a.png", img},
			apiTestFile{"image", "b.png", img},
		)
		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		for _, result := range decodeCreateImages(t, res).Results {
			require.Equal(t, "same", *result.Image.Caption)
		}
	})

	t.Run("captions need an album", func(t *testing.T) {
		for _, secret := range []string{"", ownerSecret} {
			res := createImages(secret, [][2]string{{"caption", "no album"}}, apiTestFile{"image", "a.png", img})
			require.Equal(t, http.StatusBadRequest, res.Code, res.Body.String())
			require.Equal(t, "invalid_caption", decodeCreateImages(t, res).Error.Code)
		}
		res := createImages(ownerSecret, nil, apiTestFile{"image", "a.png", img})
		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		got := decodeCreateImages(t, res).Results[0].Image
		require.Nil(t, got.Caption)
		require.Nil(t, got.AlbumURL)
	})

	t.Run("invalid captions reject the whole request", func(t *testing.T) {
		files := []apiTestFile{{"image", "a.png", img}, {"image", "b.png", img}, {"image", "c.png", img}}
		for name, fields := range map[string][][2]string{
			"too few":  {{"caption", "a"}, {"caption", "b"}},
			"too many": {{"caption", "a"}, {"caption", "b"}, {"caption", "c"}, {"caption", "d"}},
			"too long": {{"caption", strings.Repeat("a", 2001)}},
		} {
			res := createImages(ownerSecret, append([][2]string{{"album", newAlbum("bad captions")}}, fields...), files...)
			require.Equal(t, http.StatusBadRequest, res.Code, name)
			out := decodeCreateImages(t, res)
			require.Equal(t, "invalid_caption", out.Error.Code, name)
			require.Empty(t, out.Results, name)
		}
	})

	t.Run("only the owner can add to an album", func(t *testing.T) {
		albumID := newAlbum("private")
		for name, tc := range map[string]struct {
			secret string
			album  string
		}{
			"anonymous":        {"", albumID},
			"site owner":       {siteOwnerSecret, albumID},
			"unknown album":    {ownerSecret, "does-not-exist"},
			"empty album":      {ownerSecret, ""},
			"internal uuid id": {ownerSecret, mustAlbum(t, imageRepo, albumID).Id},
		} {
			res := createImages(tc.secret, [][2]string{{"album", tc.album}}, apiTestFile{"image", "a.png", img})
			require.Equal(t, http.StatusNotFound, res.Code, name)
			out := decodeCreateImages(t, res)
			require.Equal(t, "album_not_found", out.Error.Code, name)
			require.Empty(t, out.Results, name)
		}
		images, err := imageRepo.ListAlbumImages(mustAlbum(t, imageRepo, albumID).Id, 0)
		require.NoError(t, err)
		require.Empty(t, images)
	})

	t.Run("more than one album is rejected", func(t *testing.T) {
		albumID := newAlbum("one")
		res := createImages(ownerSecret, [][2]string{{"album", albumID}, {"album", albumID}}, apiTestFile{"image", "a.png", img})
		require.Equal(t, http.StatusBadRequest, res.Code)
		require.Equal(t, "invalid_album", decodeCreateImages(t, res).Error.Code)
	})

	type albumEditResponse struct {
		Album struct {
			ID  string `json:"id"`
			URL string `json:"url"`
		} `json:"album"`
		Added   int `json:"added"`
		Removed int `json:"removed"`
		Error   struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	albumEdit := func(handler func(*httpserver.HttpServerConfigDef, identity.IdentityRepo, imgddimage.ImageRepo, *ratelimit.RateLimiter) http.HandlerFunc, secret string, body string) (int, albumEditResponse) {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/albumEdit", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		if secret != "" {
			req.Header.Set("Authorization", "Bearer "+secret)
		}
		res := httptest.NewRecorder()
		handler(conf, identityRepo, imageRepo, ratelimit.NewRateLimiter(100, 100)).ServeHTTP(res, req)
		var out albumEditResponse
		require.NoError(t, json.Unmarshal(res.Body.Bytes(), &out), res.Body.String())
		return res.Code, out
	}
	addImages := func(secret string, body string) (int, albumEditResponse) {
		t.Helper()
		return albumEdit(httpserver.MakeAPIAddImagesToAlbumHandler, secret, body)
	}
	removeImages := func(secret string, body string) (int, albumEditResponse) {
		t.Helper()
		return albumEdit(httpserver.MakeAPIRemoveImagesFromAlbumHandler, secret, body)
	}
	upload := func(secret string) string {
		t.Helper()
		res := createImages(secret, nil, apiTestFile{"image", "a.png", img})
		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		return decodeCreateImages(t, res).Results[0].Image.ID
	}
	albumContents := func(albumID string) map[string]string {
		t.Helper()
		images, err := imageRepo.ListAlbumImages(mustAlbum(t, imageRepo, albumID).Id, 0)
		require.NoError(t, err)
		out := map[string]string{}
		for _, img := range images {
			out[img.Image.Id] = img.Caption
		}
		return out
	}

	t.Run("owner adds existing images to several albums", func(t *testing.T) {
		run := newAlbum("run")
		showcase := newAlbum("showcase")
		shot := upload(ownerSecret)
		code, out := addImages(ownerSecret, `{"album":"`+run+`","images":[{"id":"`+shot+`","caption":" step 3 "}]}`)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, 1, out.Added)
		require.Equal(t, run, out.Album.ID)
		require.Equal(t, "http://example.com/albums/"+run, out.Album.URL)
		code, out = addImages(ownerSecret, `{"album":"`+showcase+`","images":[{"id":"`+shot+`"}]}`)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, 1, out.Added)
		require.Equal(t, map[string]string{shot: "step 3"}, albumContents(run))
		require.Equal(t, map[string]string{shot: ""}, albumContents(showcase))

		// Adding again is a no-op that keeps the first caption.
		code, out = addImages(ownerSecret, `{"album":"`+run+`","images":[{"id":"`+shot+`","caption":"changed"}]}`)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, 0, out.Added)
		require.Equal(t, map[string]string{shot: "step 3"}, albumContents(run))
	})

	t.Run("adding requires owning the album and every image", func(t *testing.T) {
		albumID := newAlbum("guarded")
		mine := upload(ownerSecret)
		theirs := upload(siteOwnerSecret)
		guest := upload("")
		expired := upload(ownerSecret)
		_, err := conn.Exec("UPDATE image_table SET expires_at = now() - interval '1 minute' WHERE id = $1", expired)
		require.NoError(t, err)
		for name, tc := range map[string]struct {
			secret string
			body   string
			status int
			code   string
		}{
			"anonymous":            {"", `{"album":"` + albumID + `","images":[{"id":"` + mine + `"}]}`, http.StatusUnauthorized, "invalid_api_key"},
			"not the album owner":  {siteOwnerSecret, `{"album":"` + albumID + `","images":[{"id":"` + theirs + `"}]}`, http.StatusNotFound, "album_not_found"},
			"someone else's image": {ownerSecret, `{"album":"` + albumID + `","images":[{"id":"` + mine + `"},{"id":"` + theirs + `"}]}`, http.StatusNotFound, "image_not_found"},
			"guest image":          {ownerSecret, `{"album":"` + albumID + `","images":[{"id":"` + guest + `"}]}`, http.StatusNotFound, "image_not_found"},
			"expired image":        {ownerSecret, `{"album":"` + albumID + `","images":[{"id":"` + expired + `"}]}`, http.StatusNotFound, "image_not_found"},
			"unknown image":        {ownerSecret, `{"album":"` + albumID + `","images":[{"id":"00000000-0000-0000-0000-000000000009"}]}`, http.StatusNotFound, "image_not_found"},
			"no images":            {ownerSecret, `{"album":"` + albumID + `","images":[]}`, http.StatusBadRequest, "invalid_images"},
			"duplicate images":     {ownerSecret, `{"album":"` + albumID + `","images":[{"id":"` + mine + `"},{"id":"` + mine + `"}]}`, http.StatusBadRequest, "invalid_images"},
			"malformed image id":   {ownerSecret, `{"album":"` + albumID + `","images":[{"id":"nope"}]}`, http.StatusBadRequest, "invalid_images"},
			"caption too long":     {ownerSecret, `{"album":"` + albumID + `","images":[{"id":"` + mine + `","caption":"` + strings.Repeat("a", 2001) + `"}]}`, http.StatusBadRequest, "invalid_caption"},
			"misspelled field":     {ownerSecret, `{"album":"` + albumID + `","image":[{"id":"` + mine + `"}]}`, http.StatusBadRequest, "invalid_request"},
			"images as plain ids":  {ownerSecret, `{"album":"` + albumID + `","images":["` + mine + `"]}`, http.StatusBadRequest, "invalid_request"},
		} {
			code, out := addImages(tc.secret, tc.body)
			require.Equal(t, tc.status, code, name)
			require.Equal(t, tc.code, out.Error.Code, name)
		}
		// Nothing was added by the failed requests.
		require.Empty(t, albumContents(albumID))
	})

	t.Run("owner removes images without deleting them", func(t *testing.T) {
		albumID := newAlbum("remove")
		keep := newAlbum("keep")
		shot := upload(ownerSecret)
		other := upload(ownerSecret)
		for _, a := range []string{albumID, keep} {
			code, _ := addImages(ownerSecret, `{"album":"`+a+`","images":[{"id":"`+shot+`"},{"id":"`+other+`"}]}`)
			require.Equal(t, http.StatusOK, code)
		}

		code, out := removeImages(siteOwnerSecret, `{"album":"`+albumID+`","images":["`+shot+`"]}`)
		require.Equal(t, http.StatusNotFound, code)
		require.Equal(t, "album_not_found", out.Error.Code)

		code, out = removeImages(ownerSecret, `{"album":"`+albumID+`","images":["`+shot+`","00000000-0000-0000-0000-000000000009"]}`)
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, 1, out.Removed)
		require.Equal(t, map[string]string{other: ""}, albumContents(albumID))
		require.Len(t, albumContents(keep), 2)
		stillThere, err := imageRepo.GetImageById(shot)
		require.NoError(t, err)
		require.NotNil(t, stillThere)

		code, out = removeImages(ownerSecret, `{"album":"`+albumID+`","images":[]}`)
		require.Equal(t, http.StatusBadRequest, code)
		require.Equal(t, "invalid_images", out.Error.Code)
	})
}

func mustAlbum(t *testing.T, imageRepo *imgddimage.DBImageRepo, identifier string) *dm.Album {
	t.Helper()
	album, err := imageRepo.GetAlbumByIdentifier(identifier)
	require.NoError(t, err)
	require.NotNil(t, album)
	return album
}

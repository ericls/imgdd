package httpserver

// Exposes internals to the external httpserver_test package.
var MakeAPICreateImagesHandler = makeAPICreateImagesHandler
var MakeAPICreateAlbumHandler = makeAPICreateAlbumHandler
var MakeAPIAddImagesToAlbumHandler = makeAPIAddImagesToAlbumHandler
var MakeAPIRemoveImagesFromAlbumHandler = makeAPIRemoveImagesFromAlbumHandler
var MakeAppHandler = makeAppHandler
var WithTemplateFS = withTemplateFS
var WithGoogleAnalyticsID = withGoogleAnalyticsID

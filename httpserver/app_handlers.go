package httpserver

import (
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/ericls/imgdd/buildflag"
	"github.com/ericls/imgdd/captcha"
)

type appHandlerOptions struct {
	siteName           string
	siteTitle          string
	templatesFS        fs.FS
	sessionHeaderName  string
	sessionUseCookie   bool
	captchaProvider    captcha.CaptchaProvider
	recaptchaClientKey string
	turnstileSiteKey   string
	customCSS          string
	customJS           string
	googleAnalyticsID  string
}

type appHandlerOption func(*appHandlerOptions)

func withSessionHeaderName(name string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.sessionHeaderName = name
	}
}

func withSessionUseCookie(useCookie bool) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.sessionUseCookie = useCookie
	}
}

func withSiteName(name string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.siteName = name
	}
}

func withSiteTitle(title string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.siteTitle = title
	}
}

func withTemplateFS(fs fs.FS) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.templatesFS = fs
	}
}

func withCaptchaProvider(provider captcha.CaptchaProvider) func(*appHandlerOptions) {
	if !provider.IsValid() {
		panic("Invalid captcha provider")
	}
	return func(o *appHandlerOptions) {
		o.captchaProvider = provider
	}
}

func withRecaptchaClientKey(key string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.recaptchaClientKey = key
	}
}

func withTurnstileSiteKey(key string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.turnstileSiteKey = key
	}
}

func withCustomCSS(css string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.customCSS = css
	}
}

func withCustomJS(js string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.customJS = js
	}
}

func withGoogleAnalyticsID(id string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.googleAnalyticsID = id
	}
}

// albumPathPrefix is where the web client serves album pages, at
// albumPathPrefix + <album identifier>.
const albumPathPrefix = "/albums/"

// isPrivateLinkPage reports whether path is a page whose URL itself grants
// access (an album link). Such pages must not be indexed, and their URLs must
// not be sent to third parties such as analytics.
func isPrivateLinkPage(path string) bool {
	return strings.HasPrefix(path, albumPathPrefix)
}

func makeAppHandler(
	options ...appHandlerOption,
) http.HandlerFunc {
	opts := appHandlerOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		parsedTemplate, err := template.New("app.gotmpl").Funcs(template.FuncMap{
			"safe": func(s string) template.HTML { return template.HTML(s) },
		}).ParseFS(opts.templatesFS, "*.gotmpl")
		if err != nil {
			httpLogger.Err(err).Msg("Error parsing template")
			w.Write([]byte("Error parsing template"))
			return
		}
		var sessionHeaderName string
		if opts.sessionUseCookie {
			sessionHeaderName = ""
		} else {
			sessionHeaderName = opts.sessionHeaderName
		}
		privateLink := isPrivateLinkPage(r.URL.Path)
		if privateLink {
			w.Header().Set("X-Robots-Tag", "noindex, nofollow")
			// Pages opened from here (an image's page in a new tab, say) must
			// not learn the album URL from the Referer, nor pass it on to
			// analytics.
			w.Header().Set("Referrer-Policy", "no-referrer")
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = parsedTemplate.Execute(w, struct {
			Version            string
			VersionHash        string
			SiteName           string
			SiteTitle          string
			Debug              bool
			SessionHeaderName  string
			CaptchaProvider    captcha.CaptchaProvider
			RecaptchaClientKey string
			TurnstileSiteKey   string
			CustomCSS          template.CSS
			CustomJS           template.JS
			GoogleAnalyticsID  string
			PrivateLink        bool
		}{
			Version:            buildflag.Version,
			Debug:              buildflag.IsDebug,
			SiteName:           opts.siteName,
			SiteTitle:          opts.siteTitle,
			VersionHash:        buildflag.VersionHash,
			SessionHeaderName:  sessionHeaderName,
			CaptchaProvider:    opts.captchaProvider,
			RecaptchaClientKey: opts.recaptchaClientKey,
			TurnstileSiteKey:   opts.turnstileSiteKey,
			CustomCSS:          template.CSS(opts.customCSS),
			CustomJS:           template.JS(opts.customJS),
			GoogleAnalyticsID:  opts.googleAnalyticsID,
			PrivateLink:        privateLink,
		})
		if err != nil {
			httpLogger.Err(err).Msg("Error rendering template")
			w.Write([]byte("Error rendering template"))
		}
	}
}

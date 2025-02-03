package httpserver

import (
	"html/template"
	"io/fs"
	"net/http"

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

func makeAppHandler(
	options ...appHandlerOption,
) http.HandlerFunc {
	opts := appHandlerOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		parsedTemplate, err := template.ParseFS(opts.templatesFS, "*.gotmpl")
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
		}{
			Version:            buildflag.VersionHash,
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
		})
		if err != nil {
			w.Write([]byte("Error rendering template"))
		}
	}
}

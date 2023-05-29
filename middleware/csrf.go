package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"

	"github.com/heffcodex/braid/response"
	"github.com/heffcodex/braid/session"
	"github.com/heffcodex/braid/status"
	"github.com/heffcodex/braid/vars"
)

type CSRFConfigExpose struct {
	Header bool `json:"header" yaml:"header" mapstructure:"header"`
}

type CSRFConfig struct {
	Cookie       session.CookieConfig                `json:"cookie" yaml:"cookie" mapstructure:"cookie"`
	Expose       CSRFConfigExpose                    `json:"expose" yaml:"expose" mapstructure:"expose"`
	Storage      fiber.Storage                       `json:"-" yaml:"-" mapstructure:"-"`
	ErrorHandler func(c *fiber.Ctx, err error) error `json:"-" yaml:"-" mapstructure:"-"`
	Extractor    func(c *fiber.Ctx) (string, error)  `json:"-" yaml:"-" mapstructure:"-"`
}

func CSRFConfigDefault(cookie session.CookieConfig, storage ...fiber.Storage) CSRFConfig {
	var _storage fiber.Storage

	if len(storage) > 1 {
		panic("too many arguments")
	} else if len(storage) == 1 {
		_storage = storage[0]
	}

	return CSRFConfig{
		Cookie:       cookie,
		Storage:      _storage,
		ErrorHandler: func(c *fiber.Ctx, err error) error { return response.EBadRequest(status.CSRFTokenMismatch) },
		Extractor: func(c *fiber.Ctx) (string, error) {
			token, err := csrf.CsrfFromHeader(csrf.HeaderName)(c)
			if err != nil {
				return csrf.CsrfFromCookie(csrfCookieName(cookie.Name))(c)
			}

			return token, nil
		},
	}
}

func CSRF(config CSRFConfig) fiber.Handler {
	handler := csrf.New(csrf.Config{
		CookieName:     csrfCookieName(config.Cookie.Name),
		CookieDomain:   config.Cookie.Domain,
		CookiePath:     config.Cookie.Path,
		CookieSecure:   config.Cookie.Secure,
		CookieHTTPOnly: config.Cookie.HTTPOnly,
		CookieSameSite: config.Cookie.SameSite,
		Expiration:     config.Cookie.Expiration(),
		ContextKey:     vars.LocalCSRFToken,
		ErrorHandler:   config.ErrorHandler,
		Extractor:      config.Extractor,
	})

	return func(c *fiber.Ctx) error {
		err := handler(c)

		if config.Expose.Header {
			token, ok := c.Locals(vars.LocalCSRFToken).(string)
			if ok && token != "" {
				c.Set(csrf.HeaderName, token)
			}
		}

		return err
	}
}

func csrfCookieName(name string) string {
	if name != "" {
		return name
	}

	return "csrf_"
}

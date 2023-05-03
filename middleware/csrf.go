package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"

	"github.com/heffcodex/braid/response"
	"github.com/heffcodex/braid/session"
	"github.com/heffcodex/braid/status"
	"github.com/heffcodex/braid/vars"
)

type CSRFConfig struct {
	Cookie       session.CookieConfig                `json:"cookie" yaml:"cookie" mapstructure:"cookie"`
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
		Extractor:    csrf.CsrfFromHeader(csrf.HeaderName),
	}
}

func CSRF(config CSRFConfig) fiber.Handler {
	cookieName := config.Cookie.Name
	if cookieName == "" {
		cookieName = "csrf_"
	}

	handler := csrf.New(csrf.Config{
		CookieName:     cookieName,
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

		token, ok := c.Locals(vars.LocalCSRFToken).(string)
		if ok && token != "" {
			c.Set(csrf.HeaderName, token)
		}

		return err
	}
}

package session

import "time"

type CookieConfig struct {
	Name     string `json:"name" yaml:"name" mapstructure:"name"`
	Domain   string `json:"domain" yaml:"domain" mapstructure:"domain"`
	Path     string `json:"path" yaml:"path" mapstructure:"path"`
	Secure   bool   `json:"secure" yaml:"secure" mapstructure:"secure"`
	HTTPOnly bool   `json:"httpOnly" yaml:"httpOnly" mapstructure:"httpOnly"`
	SameSite string `json:"sameSite" yaml:"sameSite" mapstructure:"sameSite"`
	Expires  uint64 `json:"expires" yaml:"expires" mapstructure:"expires"`
}

func (c *CookieConfig) Expiration() time.Duration {
	return time.Duration(c.Expires) * time.Second
}

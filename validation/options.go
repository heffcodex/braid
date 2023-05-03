package validation

type Option func(v *Validator)

func WithRules(rules ...Rule) Option {
	return func(v *Validator) {
		if err := v.registerRules(rules...); err != nil {
			panic(err)
		}
	}
}

func WithMods(mods ...Mod) Option {
	return func(v *Validator) {
		v.registerModifiers(mods...)
	}
}

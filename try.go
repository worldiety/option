package option

// Try is best used with defer and an according function pointer, which is evaluated when the defer runs the try.
// It works best with [io.Closer]. Hopefully, this will get solved at the language level one day.
func Try(f func() error, err *error) {
	newErr := f()
	if *err == nil {
		*err = newErr
	}
}

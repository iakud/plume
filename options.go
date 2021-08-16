package plume

type Options struct {

}

type Option func(*Options)

func WithServices() Option {
	return func(opt *Options) {

	}
}
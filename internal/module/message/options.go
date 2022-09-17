package message

type Options struct {
	// choose the storage of messages.
	UseStorage string

	// set the maximum number for the history message of per session.
	MaxHistorySize int64
}

func NewOptions() *Options {
	return &Options{}
}

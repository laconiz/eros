package log

type Formatter interface {
	Format(*Entry)
}

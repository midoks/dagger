package conf

var (
	App struct {
		Version string `ini:"-"`
		Name    string

		RunMode string
	}
)

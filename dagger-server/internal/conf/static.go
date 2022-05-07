package conf

// CustomConf returns the absolute path of custom configuration file that is used.
var CustomConf string

var (
	App struct {
		Version string `ini:"-"`
		Name    string
		RunUser string
		RunMode string
		Node    string
	}

	// log
	Log struct {
		Format   string
		RootPath string
	}

	// http settings
	Http struct {
		Port    int    `ini:"port"`
		Path    string `ini:"path"`
		Timeout int    `ini:"timeout"`
	}

	User struct {
		Enable bool
	}

	// Security settings
	Security struct {
		InstallLock             bool
		SecretKey               string
		LoginRememberDays       int
		CookieRememberName      string
		CookieUsername          string
		CookieSecure            bool
		EnableLoginStatusCookie bool
		LoginStatusCookieName   string
	}
)

type DatabaseOpts struct {
	Type         string
	Host         string
	Port         string
	Name         string
	User         string
	Password     string
	SslMode      string `ini:"ssl_mode"`
	Path         string
	Prefix       string
	Charset      string
	MaxOpenConns int
	MaxIdleConns int
}

// Database settings
var Database DatabaseOpts

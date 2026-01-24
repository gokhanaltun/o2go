package o2go

type Config struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	AuthURLParams map[string]string
	TokenParams   map[string]string
}

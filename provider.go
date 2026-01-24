package o2go

type Provider interface {
	AuthURL() string
	TokenURL() string
}

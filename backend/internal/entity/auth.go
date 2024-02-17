package entity

type AuthManager interface {
	MakeAuthn(userID uint) (string, error)
	FetchAuthn(tknString string) (*map[string]string, error)
}

package test

type FakeTokenGenerator struct {
	RawToken    string
	HashedToken string
	GenerateErr error
}

func (f FakeTokenGenerator) Generate() (string, error) {
	if f.GenerateErr != nil {
		return "", f.GenerateErr
	}

	return f.RawToken, nil
}

func (f FakeTokenGenerator) Hash(token string) string {
	return f.HashedToken
}

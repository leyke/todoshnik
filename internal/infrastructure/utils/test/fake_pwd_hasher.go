package test

type FakePwdHasher struct {
	HashResult    string
	HashErr       error
	HashCalls     int
	LastHashInput string
	CompareResult bool
	CompareErr    error
}

func NewFakePwdHasher() *FakePwdHasher {
	return &FakePwdHasher{
		HashResult: "hashed_password",
	}
}

func (f *FakePwdHasher) Hash(password string) (string, error) {
	f.HashCalls++
	f.LastHashInput = password

	return f.HashResult, f.HashErr
}

func (f *FakePwdHasher) Compare(hash string, input string) (bool, error) {
	return f.CompareResult, f.CompareErr
}

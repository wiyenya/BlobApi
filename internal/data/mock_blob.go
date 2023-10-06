package data

type MockBlobModel struct{}

func (m *MockBlobModel) Insert(userID int, data string) (int, error) {
	return 1, nil
}

func (m *MockBlobModel) Get(id int) (*Blob, error) {
	return &Blob{ID: id, UserID: nil, Data: "mock data"}, nil
}

func (m *MockBlobModel) GetBlobList() ([]*Blob, error) {
	return []*Blob{
		{ID: 1, UserID: nil, Data: "mock data1"},
		{ID: 2, UserID: nil, Data: "mock data2"},
	}, nil
}

func (m *MockBlobModel) Delete(id int) error {
	return nil
}

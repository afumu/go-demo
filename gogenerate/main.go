package test

//go:generate mockgen --source=$GOFILE --package=$GOPACKAGE --destination=mock_test.go dependency

type MyInter interface {
	GetName(id int) string
}

func GetUser(m MyInter, id int) string {
	user := m.GetName(id)
	return user
}

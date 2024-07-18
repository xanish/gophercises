package phone

type Phone struct {
	Contacts *Contacts
}

func NewPhone(driverName, dataSource string) (*Phone, error) {
	contacts, err := Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}

	return &Phone{Contacts: contacts}, nil
}

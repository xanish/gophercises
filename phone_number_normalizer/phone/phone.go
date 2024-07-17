package phone

type Phone struct {
	contacts *Contacts
}

func NewPhone(driverName, dataSource string) (*Phone, error) {
	contacts, err := Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}

	return &Phone{contacts: contacts}, nil
}

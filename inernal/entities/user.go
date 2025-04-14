package entities

type User struct {
	ID      int
	Balance int
	Sheba   Sheba
}

func (u *User) CheckImmutable(n *User) bool {
	return u.ID == n.ID && u.Sheba == n.Sheba
}

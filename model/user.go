package model

type User struct {
	Username string
	Password string
	Roles    []string
}

func (u User) ToID() string {
	concatRoles := ""
	for _, r := range u.Roles {
		concatRoles += r + ","
	}
	if len(concatRoles) > 0 {
		concatRoles = concatRoles[:len(concatRoles)-1]
	}
	return u.Username + "#" + u.Password + "#" + concatRoles
}

package model

const (
	PERM_NO    = 1
	PERM_YES   = 2
	PERM_NEVER = 3

	//
	PERM_BROWSE = 1
	PERM_READ   = 2
	PERM_EDIT   = 4
	PERM_ADD    = 8
	PERM_DELETE = 16

	PERM_ALL = PERM_BROWSE | PERM_READ | PERM_EDIT | PERM_ADD | PERM_DELETE
)

type Role struct {
	ID        int
	Name      string
	OrgID     int
	IsDefault bool
	IsBase    bool
}

type UserRole struct {
	UserID int
	RoleID int
}

type Org struct {
	ID        int
	Name      string
	IsDefault bool
	IsBase    bool
}

type Permission struct {
	ID           int
	Action       int
	Grant        int
	ResourceID   int
	ResourceRole int
	ResourceOrg  int
	RoleID       int
	OrgID        int
	UserID       int
}

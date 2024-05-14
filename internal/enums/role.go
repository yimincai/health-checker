package enums

type RoleType int

const (
	RoleType_Normal RoleType = iota
	RoleType_Admin
)

func (r RoleType) String() string {
	switch r {
	case RoleType_Normal:
		return "Normal"
	case RoleType_Admin:
		return "Admin"
	default:
		return "Unknown"
	}
}

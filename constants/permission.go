package constants

type (
	// PermissionCommand : định nghĩ các permission
	PermissionCommand string
)

// Employee permissions
const (
	CreateUser PermissionCommand = "CREATE_USER"
	UpdateUser PermissionCommand = "UPDATE_USER"
	GetUser    PermissionCommand = "GET_USER"
	DeleteUser PermissionCommand = "DELETE_USER"
)

// User permissions
const (
	Test PermissionCommand = "TEST"
)

var (
	// UserPermissions : định nghĩa các quyền của người dùng
	UserPermissions = map[PermissionCommand]bool{
		Test: true,
	}

	// EmployeePermissions : định nghĩa các quyền của nhân viên
	EmployeePermissions = map[PermissionCommand]bool{
		CreateUser: true,
		UpdateUser: true,
		GetUser:    true,
		DeleteUser: true,
	}
)

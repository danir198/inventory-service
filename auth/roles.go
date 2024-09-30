package middleware

type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
)

var rolePermissions = map[Role][]string{
	Admin: {
		"add_product",
		"update_product",
		"delete_product",
		"view_all_products",
		"manage_roles",
	},
	User: {
		"view_product",
		"check_availability",
	},
}

func HasPermission(role Role, permission string) bool {
	for _, p := range rolePermissions[role] {
		if p == permission {
			return true
		}
	}
	return false
}
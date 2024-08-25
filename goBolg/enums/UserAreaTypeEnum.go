package enums

// UserAreaTypeEnum represents the user area type.
type UserAreaTypeEnum struct {
	Type int
	Desc string
}

// Constants for UserAreaTypeEnum.
var (
	USERTYPE = UserAreaTypeEnum{Type: 1, Desc: "用户"}
	VISITOR  = UserAreaTypeEnum{Type: 2, Desc: "游客"}
)

// GetUserAreaType returns the UserAreaTypeEnum based on the provided type.
func GetUserAreaType(t int) *UserAreaTypeEnum {
	switch t {
	case USERTYPE.Type:
		return &USERTYPE
	case VISITOR.Type:
		return &VISITOR
	default:
		return nil
	}
}

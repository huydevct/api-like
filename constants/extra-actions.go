package constants

func ListExtraAction() (actions []string) {
	actions = []string{"ChangePassword", "ChangeSecretkey", "ChangeCover", "ChangeAvatar", "AddMail", "VeryMail"}
	return
}

func CheckInArray(value string, values []string) (res bool) {
	for _, val := range values {
		if val == value {
			res = true
		}
	}
	return
}

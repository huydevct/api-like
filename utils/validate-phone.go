package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// PhoneValidate : Kiểm tra các số điện thoại có hợp lệ hay không ?
// Kiểm tra độ dải
// Kiểm tra đầu số
// Chuẩn format phone input
func PhoneValidate(phone string) (phoneNum string, err error) {

	validHomephoneHeaders := []int{24, 28, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 218, 219, 220, 221, 222, 225, 226, 227, 228, 229, 232, 233, 234, 235, 236, 237, 238, 239, 251, 252, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 269, 270, 271, 272, 273, 274, 275, 276, 277, 290, 291, 292, 293, 294, 296, 297, 299}
	validMobilephoneHeaders := []int{
		80, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 161, 162, 163, 164, 165, 166, 167, 168, 169, 186, 188, 199, 868, 32, 33, 34, 35, 36, 37, 38, 39, 70, 71, 72, 73, 79, 77, 76, 75, 78, 83, 84, 85, 81, 82, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	}
	// Remove white space
	phone = strings.Replace(phone, " ", "", -1)
	phoneNum, _, errCheck := getPhoneNumber(phone)
	if errCheck != nil {
		err = errCheck
		return
	}

	_, err = strconv.Atoi(phoneNum)
	if err != nil {
		err = fmt.Errorf("Phone must only contain number")
		return
	}

	// Kiểm tra chiều dài
	if len(phoneNum) > 10 {
		err = fmt.Errorf("Phone length is invalid")
		return
	}

	// Kiểm tra số điện thoại di đông (số điện thoại di động 10 số)
	headerMobilePhone := phoneNum[0:2]
	isMobile := false
	for _, h := range validMobilephoneHeaders {
		if strings.HasPrefix(headerMobilePhone, strconv.Itoa(h)) {
			isMobile = true
		}
	}
	if isMobile && len(phoneNum[2:]) == 7 {
		phoneNum = fmt.Sprintf("0%s", phoneNum)
		return
	}

	// kiểm tra số điện thoại bàn (số điện thoại bàn 11 số)
	isHome := false
	headerHomePhone := ""
	for _, h := range validHomephoneHeaders {
		if strings.HasPrefix(phoneNum, strconv.Itoa(h)) {
			headerHomePhone = strconv.Itoa(h)
			isHome = true
		}
	}

	if isHome {
		// nếu đầu số điện thoại là 2 số thì chiều dài còn lại phải 8 số
		if len(headerHomePhone) == 2 && len(phoneNum[2:]) == 8 {
			phoneNum = fmt.Sprintf("0%s", phoneNum)
			return
		}
		// nếu đầu số điện thoại là 3 số thì chiều dài còn lại phải 7 số
		if len(headerHomePhone) == 3 && len(phoneNum[3:]) == 7 {
			phoneNum = fmt.Sprintf("0%s", phoneNum)
			return
		}
	}

	err = fmt.Errorf("Phone number is invalid")
	return
}

func getPhoneNumber(phone string) (phoneNum, phoneHeader string, err error) {

	if existed := regexp.MustCompile(`^\+84`).MatchString(phone); existed {
		phoneNum = phone[3:]
		phoneHeader = phone[3:5]
		return
	}

	if existed := regexp.MustCompile(`^84`).MatchString(phone); existed {
		phoneNum = phone[2:]
		phoneHeader = phone[2:4]
		return
	}

	if existed := regexp.MustCompile(`^0`).MatchString(phone); existed {
		phoneNum = phone[1:]
		phoneHeader = phone[1:3]
		return
	}

	err = fmt.Errorf("Phone must be start by +84, 84, 0")
	return
}

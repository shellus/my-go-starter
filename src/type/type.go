package main

import "fmt"

/**
性别类型
女 Female
男 Male
 */
type GenderType struct {
	// 该项目的名称。
	name string
}

const (
	female = "Female"
	male = "Male"
)

/**
人类
 */
type Person struct {
	// 姓。 在美国，一个人的姓氏。 这可以与givenName一起使用，而不是name属性。
	familyName string
	// 给定名字. 在美国，一个人的名字。 这可以与familyName而不是name属性一起使用。
	givenName  string

	// 人的性别。 虽然可以使用http://schema.org/Male和http://schema.org/Female，但是对于不识别为二进制性别的人也可以接受文本字符串。
	gender     GenderType

	// 电子邮件地址
	email      string
	// 传真号码。
	faxNumber  int

	name       string
}

func main() {
	fmt.Print(female, male)
}

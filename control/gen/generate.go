package main

import (
	"fmt"
	serial "github.com/Myriad-Dreamin/minimum-template/lib/serial"
	"github.com/Myriad-Dreamin/minimum-template/model"
	"github.com/Myriad-Dreamin/minimum-template/types"
	"os"
)

type UserCategories struct {
	serial.VirtualService
	List           *serial.Category
	Login          *serial.Category
	Register       *serial.Category
	ChangePassword *serial.Category
	IdGroup        *serial.Category
}

var codeField = serial.Param("code", *new(types.CodeRawType))
var required = serial.Tag("binding", "required")

func DescribeUserService(cat *serial.Category) serial.ProposingService {
	var userModel = new(model.User)
	svc := &UserCategories{
		List: serial.Ink().
			Path("user-list").
			Method(serial.POST, "List",
				serial.Request(
					serial.Transfer(model.Filter{}),
				),
				serial.Reply(
					codeField,
					serial.ArrayParam(serial.Param("users", serial.Object(
						"ListUserReply",
						serial.Param("nick_name", userModel.NickName),
						serial.Param("last_login", userModel.LastLogin),
					))),
				),
			),
		Login: serial.Ink().
			Path("login").
			Method(serial.POST, "Login",
				serial.Request(
					serial.Param("id", userModel.ID),
					serial.Param("nick_name", userModel.NickName),
					serial.Param("phone", userModel.Phone),
					serial.Param("password", serial.String, required),
				),
				serial.Reply(
					codeField,
					serial.Param("id", userModel.ID),
					serial.Param("identity", serial.Strings),
					serial.Param("phone", userModel.Phone),
					serial.Param("nick_name", userModel.NickName),
					serial.Param("name", userModel.Name),
					serial.Param("token", serial.String),
					serial.Param("refresh_token", serial.String),
				),
			),
		Register: serial.Ink().
			Path("register").
			Method(serial.POST, "Register",
				serial.Request(
					serial.Param("name", serial.String, required),
					serial.Param("password", serial.String, required),
					serial.Param("nick_name", serial.String, required),
					serial.Param("phone", serial.String, required),
				),
				serial.Reply(
					codeField,
					serial.Param("id", userModel.ID)),
			),
		ChangePassword: serial.Ink().
			Path("user/:id/password").
			Method(serial.PUT, "ChangePassword",
				serial.Request(
					serial.Param("old_password", serial.String, required),
					serial.Param("new_password", serial.String, required),
				),
			),
		IdGroup: serial.Ink().
			Method(serial.GET, serial.PUT, serial.DELETE).
			Path("user/:id").
			Method(serial.GET, "Get",
				serial.Reply(
					codeField,
					serial.Param("nick_name", userModel.NickName),
					serial.Param("last_login", userModel.LastLogin),
				)).
			Method(serial.PUT, "Put",
				serial.Request(
					codeField,
					serial.Param("phone", userModel.Phone),
				)).
			Method(serial.DELETE, "Delete"),
	}
	svc.Name("UserService").CateOf(cat).UseModel(userModel)
	return svc
}

func main() {
	V1Cate := serial.Ink().Path("v1")

	userCate := DescribeUserService(V1Cate)
	userCate.ToFile("user.go")
	fmt.Println(os.Getwd())
	err := serial.NewService(
		userCate).Publish()
	if err != nil {
		fmt.Println(err)
	}
}

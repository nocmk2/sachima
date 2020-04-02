package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

func main() {
	e, err := casbin.NewEnforcer("./rbac.conf", "./rbac.csv")

	sub := "alice" // the user that wants to access a resource.
	obj := "data1" // the resource that is going to be accessed.
	act := "read"  // the operation that the user performs on the resource.

	ok, err := e.Enforce(sub, obj, act)
	e.EnableAutoSave(true)

	if err != nil {
		// handle err
	}

	if ok == true {
		// permit alice to read data1
	} else {
		// deny the request, show an error
	}
	hasrole, _ := e.HasRoleForUser("wenyin", "Supseradmin")
	e.AddRoleForUser("wenyin", "Uxdd")
	roles, _ := e.GetRolesForUser("wenyin")

	fmt.Println(hasrole)
	fmt.Println(roles)
}

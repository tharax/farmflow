package main

import (
	"fmt"

	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type DataStore_ContainersDB struct {
	Things []interface{}
}
type ProfileKeys struct {
	profiles []string
}
type Global struct{}

type Characters map[string]interface{}

func DoLUAThing() {
	fmt.Println("\nLUA THING\n==================")
	L := lua.NewState()
	defer L.Close()
	if err := L.DoString(`print("hello")`); err != nil {
		panic(err)
	}
	if err := L.DoFile("data/DataStore_Containers.lua"); err != nil {
		panic(err)
	}
	// type Role struct {
	// 	Name string
	// }

	// type Person struct {
	// 	Name      string
	// 	Age       int
	// 	WorkPlace string
	// 	Role      []*Role
	// }

	// L := lua.NewState()
	// if err := L.DoString(`
	// person = {
	//   name = "Michel",
	//   age  = "31", -- weakly input
	//   work_place = "San Jose",
	//   role = {
	// 	{
	// 	  name = "Administrator"
	// 	},
	// 	{
	// 	  name = "Operator"
	// 	}
	//   }
	// }
	// `); err != nil {
	// 	panic(err)
	// }

	_, err := L.LoadFile("data/DataStore_Containers.lua")
	if err != nil {
		panic(err)
	}
	// L.DoString("a")
	// fmt.Printf("%+v", file)

	// fmt.Printf("%+v", file.Proto)
	var db Characters
	if err := gluamapper.Map(L.GetGlobal(`DataStore_ContainersDB`).(*lua.LTable), &db); err != nil {
		panic(err)
	}

	// if err := gluamapper.Map(file.Env, &db); err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("DB:%+v<---", db)

	for k, v := range db {
		fmt.Println(k)
		if k == "Global" {
			x := innerValue(v.(map[interface{}]interface{}))
			// fmt.Println(k, x["Characters"])
			fmt.Printf("%+v", x.Val("Characters").Val("Default.Frostmourne.Humphrëy").Val("Containers"))

			// fmt.Println(x["Characters"]["Default.Frostmourne.Humphrëy"])
		}

	}
	// fmt.Printf("%s %d", person.Name, person.Age)

}

type innerValue map[interface{}]interface{}

func (i *innerValue) Val(value string) *innerValue {
	var x innerValue
	for k, v := range *i {
		if k == value {
			x = innerValue(v.(map[interface{}]interface{}))
		}
	}
	return &x
}

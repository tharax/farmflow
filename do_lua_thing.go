package main

import (
	"fmt"
	"regexp"
	"strings"

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
	items := Inventory{}

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
	var bag4 *innerValue
	for k, v := range db {
		fmt.Println(k)
		if k == "Global" {
			x := innerValue(v.(map[interface{}]interface{}))
			// fmt.Println(k, x["Characters"])

			bag4 = x.Val("Characters").Val("Default.Frostmourne.Humphrëy").Val("Containers").Val("Bag4")

			fmt.Printf("%+v\n", bag4)

			// fmt.Println(x["Characters"]["Default.Frostmourne.Humphrëy"])
		}

	}

	// fmt.Printf("links interface-->%+v\n", bag4.ValAsInterface("Links"))
	// fmt.Printf("links array    -->%+v\n", bag4.ValAsArray("Links"))
	// fmt.Printf("links interface-->%+v\n", bag4.Val("Links"))
	// fmt.Printf("links maps     -->%+v\n", len(bag4.ValAsInterface("Links")))
	// fmt.Printf("links-->%+v\n", bag4["Counts"])
	// fmt.Printf("links-->%+v\n", len(bag4))
	// b4 := innerValue2(bag4)
	// fmt.Printf("b4----->%+v\n", b4)
	// fmt.Printf("%s %d", person.Name, person.Age)
	// v, ok :=

	var re = regexp.MustCompile(`(?:\|h\[)(.*)(?:\]\|)`)

	for _, v := range bag4.ValAsArray("Links") {
		// fmt.Println(v.(string))
		str := v.(string)
		x := re.FindAllString(str, -1)
		// match := re.MatchString(str)
		// fmt.Println(match, len(x), fmt.Sprintf("--->%+v<---", x[0]))

		a := strings.TrimLeft(x[0], "|h[")
		aa := strings.TrimRight(a, "]|")
		// fmt.Println(aa)
		items = append(items, &struct {
			Name     string `json:"name"`
			Quantity int    `json:"quantity"`
		}{
			Name: aa,
		})
		// for k, vv := range x {
		// 	if k == 2 {
		// 		fmt.Println(k, vv)
		// 	}
		// }
		// fmt.Println(match, len(x), fmt.Sprintf("--->%+v<---", x), fmt.Sprintf("--->%+v<---", str))
		// x := re.FindAllString(v.(string), -1)
		// fmt.Println(x, len(x), x[0], x[1], x[len(x)-1])
		// "|cff1eff00|Hitem:10153::::::::120:267:::1:1680:::|h[Mighty Spaulders of the Quickblade]|h|r"
		// 		regexp.Match("|.*)
	}
	for _, v := range items {
		fmt.Println(v.Name, v.Quantity)
	}
}

type innerValue map[interface{}]interface{}
type innerValue2 map[string]interface{}

func (i *innerValue) Val(value string) *innerValue {
	var x innerValue
	for k, v := range *i {
		if k == value {
			x = innerValue(v.(map[interface{}]interface{}))
		}
	}
	return &x
}
func (i *innerValue) ValAsMap(value string) map[string]interface{} {
	for k, v := range *i {
		if k == value {
			return map[string]interface{}(v.(map[string]interface{}))
		}
	}
	return nil
}

func (i *innerValue) ValAsInterface(value string) interface{} {
	var x interface{}
	for k, v := range *i {
		if k == value {
			x = interface{}(v.(interface{}))
		}
	}
	return x
}

func (i *innerValue) ValAsArray(value string) []interface{} {
	var x []interface{}
	for k, v := range *i {
		if k == value {
			x = []interface{}(v.([]interface{}))
		}
	}
	return x
}

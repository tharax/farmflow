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

type Characters map[string]interface{}

type InvItem struct {
	ID       int
	Name     string
	Quantity int
}

type Bag struct {
	Size  int
	Items []BagItem
}

type BagItem struct {
	Slot     int
	ID       int
	Name     string
	Quantity int
}

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

	_, err := L.LoadFile("data/DataStore_Containers.lua")
	if err != nil {
		panic(err)
	}

	var db Characters
	if err := gluamapper.Map(L.GetGlobal(`DataStore_ContainersDB`).(*lua.LTable), &db); err != nil {
		panic(err)
	}

	// if err := gluamapper.Map(file.Env, &db); err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("DB:%+v<---", db)
	// var bag4 *innerValue
	for k, v := range db {
		fmt.Println(k)
		if k == "Global" {
			x := innerValue(v.(map[interface{}]interface{}))
			// fmt.Println(k, x["Characters"])

			// bag4 = x.Val("Characters").Val("Default.Frostmourne.Humphrëy").Val("Containers").Val("Bag0")

			bags := x.Val("Characters").Val("Default.Frostmourne.Petër").Val("Containers")
			fmt.Println(len(*bags))
			for _, value := range *bags {
				x := innerValue(value.(map[interface{}]interface{}))
				for _, v := range x.GetItemsFromBags() {
					if v.Quantity > 0 {
						fmt.Println(v.ID, v.Name, v.Quantity)
					}
				}
			}
			// fmt.Printf("%+v\n", bag4)

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

	// array := bag4.GetItemsFromBags()

	// for _, v := range array {
	// 	fmt.Println(v.ID, v.Name, v.Quantity)
	// }
}

func (i *innerValue) GetItemsFromBags() []InvItem {
	lengthOfBag := i.ValAsInt("Size")
	res := make([]InvItem, lengthOfBag)
	var re = regexp.MustCompile(`(?:\|h\[)(.*)(?:\]\|)`)
	for k, v := range i.ValAsArray("Links") {
		strArray := re.FindAllString(v.(string), -1)
		str := strings.TrimLeft(strArray[0], "|h[")
		str = strings.TrimRight(str, "]|")

		res[k].Name = str
	}
	for k, v := range i.ValAsArray("Counts") {
		if v != nil {
			res[k].Quantity = int(v.(float64))
		} else {
			res[k].Quantity = 1
		}
	}
	for k, v := range i.ValAsArray("Ids") {
		if v != nil {
			res[k].ID = int(v.(float64))
		}
	}
	return res
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

func (i *innerValue) ValAsInt(value string) int {
	var x int
	for k, v := range *i {
		if k == value {
			x = int(v.(float64))
		}
	}
	return x
}

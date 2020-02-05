package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tharax/farmflow/inventory"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type DataStore_ContainersDB struct {
	ProfileKeys map[string]string
	Global      Global
}

type Global struct {
	Characters map[string]Character
	Guilds     map[string]interface{}
	Version    int
}

type Character struct {
	Containers map[string]Container
}

type Container struct {
	Size   int
	Ids    []int
	Links  []string
	Counts []int
}

func DoLUAThing() {

	fmt.Println("\nLUA THING\n==================")
	L := lua.NewState()
	defer L.Close()

	if err := L.DoFile("data/DataStore_Containers.lua"); err != nil {
		panic(err)
	}

	var db DataStore_ContainersDB
	if err := gluamapper.Map(L.GetGlobal(`DataStore_ContainersDB`).(*lua.LTable), &db); err != nil {
		panic(err)
	}

	var inv inventory.Inventory

	for _, character := range db.Global.Characters {
		for _, container := range character.Containers {
			counts := make([]int, container.Size)
			for position, count := range container.Counts {
				if count == 0 {
					counts[position] = 1
				}
				counts[position] = count
			}
			for position, item := range container.Links {
				inv.Add(inventory.Item{
					ID:       container.Ids[position],
					Name:     TrimItemDescription(item),
					Quantity: counts[position],
				})
			}
		}
	}

	fmt.Println(inv)
}

func GetInvFromLUAFile() Inventory {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoFile("data/DataStore_Containers.lua"); err != nil {
		panic(err)
	}

	var db DataStore_ContainersDB
	if err := gluamapper.Map(L.GetGlobal(`DataStore_ContainersDB`).(*lua.LTable), &db); err != nil {
		panic(err)
	}

	var inv inventory.Inventory

	for _, character := range db.Global.Characters {
		for _, container := range character.Containers {
			counts := make([]int, container.Size)
			for position, count := range container.Counts {
				if count == 0 {
					counts[position] = 1
				}
				counts[position] = count
			}
			for position, item := range container.Links {
				inv.Add(inventory.Item{
					ID:       container.Ids[position],
					Name:     TrimItemDescription(item),
					Quantity: counts[position],
				})
			}
		}
	}
	var res Inventory
	for _, v := range inv {
		res = append(res, InvItem{
			Item: Item{
				ID:   v.ID,
				Name: v.Name,
			},
			Quantity: v.Quantity,
		})
	}
	return res
}

func TrimItemDescription(item string) string {
	var re = regexp.MustCompile(`(?:\|h\[)(.*)(?:\]\|)`)
	strArray := re.FindAllString(item, -1)
	var str string
	if len(strArray) > 0 {
		str = strings.TrimLeft(strArray[0], "|h[")
		str = strings.TrimRight(str, "]|")
	}
	if str != "" {
		return str
	}
	return item
}

func (i *innerValue) GetItemsFromBags() []InvItem {
	lengthOfBag := i.ValAsInt("Size")
	res := make([]InvItem, lengthOfBag)
	var re = regexp.MustCompile(`(?:\|h\[)(.*)(?:\]\|)`)
	for k, v := range i.ValAsArray("Links") {
		if v != nil {
			strArray := re.FindAllString(v.(string), -1)
			str := strings.TrimLeft(strArray[0], "|h[")
			str = strings.TrimRight(str, "]|")

			res[k].Name = str
		}
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

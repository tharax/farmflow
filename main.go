package main

import (
	"encoding/json"
	"fmt"
	"github.com/tharax/farmflow/items"
	"github.com/tharax/farmflow/lua"
	"io/ioutil"
	"os"
	"sort"
)

type Items []Item

type Item struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	IsTransmog bool           `json:"is_transmog"`
	Inputs     map[string]int `json:"inputs"`
}

type InvItem struct {
	Item
	Quantity int
}

type Inventory []InvItem

func (i Inventory) Use(name string, amount int) (leftover int) {
	leftover = amount
	for k, v := range i {

		if v.Name == name {
			if amount >= v.Quantity {
				leftover = amount - v.Quantity
				i[k].Quantity = 0
			} else {
				i[k].Quantity = v.Quantity - amount
				leftover = 0
			}
		}
	}
	return
}

func main() {
	itemdb, err := items.LoadItems()
	check(err)

	// m := itemdb.GetInventoryType()
	// for count := 0; count < len(m); count++ {
	// 	fmt.Println(m[count], " = ", count)
	// }

	item, err := itemdb.GetItemByName("Robe of Winter Night")
	check(err)
	fmt.Printf("%+v\n", item)
	fmt.Println(itemdb.GetQualitySums())

	validTransmog := itemdb.GetThingsThatMightBeTransmog()

	fmt.Println(len(itemdb.Items), len(validTransmog))

	crafts, err := lua.LoadDataStoreItems()
	check(err)
	RECIPES, err := crafts.GetDirectedRecipeGraph()

	check(err)
	fmt.Println(len(RECIPES))
	fmt.Println(RECIPES.GetDeepest())
	err = RECIPES.PrintRecipeFor(item.ItemID)
	check(err)

	b, err := ioutil.ReadFile("data/items.json")
	check(err)
	var items Items
	err = json.Unmarshal(b, &items)
	check(err)

	inventory := GetInvFromLUAFile()

	rm := CalculateRawMaterials(&items, &inventory)
	keys := make([]string, 0, len(rm))
	for k := range rm {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	PrintRecipesToCreate(items)
	// PrintFarmListForInventory(Inventory{})         //blank inv
	PrintFarmListForInventory(GetInvFromLUAFile()) //loaded from test LUA file
	// PrintLeftOverInventory(inventory)

}

func PrintRecipesToCreate(items Items) {
	fmt.Println("\nRecipes to Create\n=================")
	r := []string{}
	for _, item := range items {
		if item.IsTransmog {
			r = append(r, item.Name)
		}
	}
	sort.Strings(r)
	for _, item := range r {
		fmt.Println(item)
	}
}

func PrintLeftOverInventory(inv Inventory) {
	fmt.Println("\nLeftover Inventory\n==================")
	for _, k := range inv {
		if k.Quantity > 0 {
			fmt.Printf("%-30v %5v\n", k.Name, k.Quantity)
		}
	}
}

func PrintFarmListForInventory(inv2 Inventory) {
	b2, err := ioutil.ReadFile("data/items.json")
	check(err)
	var items2 Items
	err = json.Unmarshal(b2, &items2)
	rm2 := CalculateRawMaterials(&items2, &inv2)
	keys2 := make([]string, 0, len(rm2))
	for k := range rm2 {
		keys2 = append(keys2, k)
	}
	sort.Strings(keys2)
	r2 := []string{}
	for _, item2 := range items2 {
		if item2.IsTransmog {
			r2 = append(r2, item2.Name)
		}
	}
	sort.Strings(r2)
	fmt.Println("\nRemaining to Farm\n=================")
	for _, k := range keys2 {
		if rm2[k] > 0 {
			fmt.Printf("%-30v %5v\n", k, rm2[k])
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CalculateRawMaterials(items *Items, inventory *Inventory) map[string]int {
	res := map[string]int{}
	var itemsToMake []string
	for _, item := range *items {
		if item.IsTransmog {
			itemsToMake = append(itemsToMake, item.Name)
		}
	}

	loop := func(name string, count int) {}
	loop = func(name string, count int) {
		missing := true
		for _, v := range *items {
			if v.Name == name {
				if len(v.Inputs) != 0 {
					for kk, vv := range v.Inputs {
						loop(kk, inventory.Use(kk, (count*vv)))
					}
				} else {
					res[name] = res[name] + count
				}
				missing = false
			}
		}
		if missing {
			fmt.Printf("missing %s\n", name)
			os.Exit(1)
		}

	}

	for _, item := range itemsToMake {
		loop(item, 1)
	}

	return res
}

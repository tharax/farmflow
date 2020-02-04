package main

import (
	"encoding/json"
	"fmt"
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
	b, err := ioutil.ReadFile("data/items.json")
	check(err)
	var items Items
	err = json.Unmarshal(b, &items)
	check(err)

	b, err = ioutil.ReadFile("data/inventory.json")
	check(err)
	var inventory Inventory
	err = json.Unmarshal(b, &inventory)
	check(err)

	rm := CalculateRawMaterials(&items, &inventory)
	keys := make([]string, 0, len(rm))
	for k := range rm {
		keys = append(keys, k)
	}
	sort.Strings(keys)

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

	fmt.Println("\nRemaining to Farm\n=================")
	for _, k := range keys {
		if rm[k] > 0 {
			fmt.Printf("%-30v %5v\n", k, rm[k])
		}
	}

	fmt.Println("\nLeftover Inventory\n==================")
	for _, k := range inventory {
		if k.Quantity > 0 {
			fmt.Printf("%-30v %5v\n", k.Name, k.Quantity)
		}
	}

	DoLUAThing()
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

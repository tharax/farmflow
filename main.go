package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Data struct {
	Recipes Recipes `json:"recipes"`
	Items   Items   `json:"items"`
}

type Recipes []Recipe

type Recipe struct {
	Name    string         `json:"name"`
	Source  string         `json:"source"`
	Inputs  map[string]int `json:"inputs"`
	Outputs map[string]int `json:"outputs"`
}

type Items []Item

type Item struct {
	Name       string `json:"name"`
	Source     string `json:"source"`
	IsTransmog bool   `json:"is_transmog"`
}

const (
	SourceRecipe = "recipe"
)

func main() {
	fmt.Println("farmflow")

	b, err := ioutil.ReadFile("data.json")
	if err != nil {
		panic(err)
	}
	var data Data
	err = json.Unmarshal(b, &data)
	if err != nil {
		panic(err)
	}

	for _, recipe := range data.Recipes {
		for k, _ := range recipe.Inputs {
			if data.Items.Exists(k) == false {
				log.Printf("Missing Item: %s", k)
			}
		}
	}
	for _, item := range data.Items {
		if item.Source == SourceRecipe {
			if data.Recipes.Exists(item.Name) == false {
				log.Printf("Missing Recipe: %s", item.Name)
			}
		}
	}

	fullRawMaterialsList := map[string]int{}
	for _, item := range data.Items {
		if item.IsTransmog {
			log.Println(item.Name)
			for k, v := range item.RawMaterials(data.Recipes, data.Items) {
				fullRawMaterialsList[k] = fullRawMaterialsList[k] + v
				// log.Printf("%30s %4v\n", k, v)
			}

		}
	}
	for k, v := range fullRawMaterialsList {
		log.Printf("%30s %4v\n", k, v)
	}
}

func (i Item) RawMaterials(recipes Recipes, items Items) map[string]int {
	res := make(map[string]int, 0)
	if i.Source != SourceRecipe {
		res[i.Name] = 1
	}
	if i.Source == SourceRecipe {
		for _, recipe := range recipes {
			if i.Name == recipe.Name {
				for inputKey, inputValue := range recipe.Inputs {
					componentItem := items.Get(inputKey)
					rm := componentItem.RawMaterials(recipes, items)
					for k, v := range rm {
						res[k] = res[k] + (v * inputValue)
					}
				}
			}
		}
	}
	return res
}

// 		for ingredient, count := range recipe.Inputs {
// 			fmt.Println(ingredient, count)
// 			for _, ingred := range i {
// 				if ingred.Name == ingredient {
// 					// var parts map[string]int
// 					if ingred.Source == SourceRecipe {
// 						for kk, _ := range recipe.Inputs {
// 							next := i.Get(kk)
// 							fmt.Println(next.Name)
// 							parts := next.RawMaterials(r, i)
// 							fmt.Println(parts)

// 						}
// 						fmt.Println(recipe)
// 					} else {
// 						fmt.Println(ingred.Source)
// 					}

// 					// 	next := i.Get(ingredient)
// 					// 	parts = next.RawMaterials(r, i)
// 					// }
// 					// for k, v := range parts {
// 					// 	rawParts[k] = rawParts[k] + (v * count)
// 					// }
// 				}
// 			}
// 		}
// 	}
// }
// fmt.Println(rawParts)
// return rawParts

func (i Items) Get(name string) Item {
	for _, ii := range i {
		if ii.Name == name {
			return ii
		}
	}
	return Item{}
}

func (i Items) Exists(name string) bool {
	for _, ii := range i {
		if ii.Name == name {
			return true
		}
	}
	return false
}

func (r Recipes) Exists(name string) bool {
	for _, rr := range r {
		if rr.Name == name {
			return true
		}
	}
	return false
}

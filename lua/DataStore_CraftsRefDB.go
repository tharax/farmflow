package lua

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"regexp"

	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type DataStore_CraftsRefDB struct {
	ProfileKeys map[string]string
	Global      DataStore_CraftsRefDBGlobal
}

type DataStore_CraftsRefDBGlobal struct {
	ReagentsActualUseThisOne map[int]string
	Reagents                 []string
	Recipes                  interface{}
	ResultItems              interface{}
	RecipeCategoryNames      interface{}
}

func LoadDataStoreItems() (*DataStore_CraftsRefDB, error) {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoFile("data/DataStore_Crafts.lua"); err != nil {
		return nil, err
	}

	var db DataStore_CraftsRefDB
	if err := gluamapper.Map(L.GetGlobal(`DataStore_CraftsRefDB`).(*lua.LTable), &db); err != nil {
		return nil, err
	}
	Actual := map[int]string{}
	for k, recipe := range db.Global.Reagents {
		if recipe != "" {
			Actual[k] = recipe
		}
	}
	db.Global.ReagentsActualUseThisOne = Actual
	return &db, nil
}

func (db *DataStore_CraftsRefDB) GetDirectedRecipeGraph() (Items, error) {

	var re = regexp.MustCompile(`(?m)[0-9]*,[0-9]*`)
	var listOfRecipes Items
	fmt.Println("list", len(db.Global.ReagentsActualUseThisOne))
	for recipe, components := range db.Global.ReagentsActualUseThisOne {
		res := re.FindAllString(components, -1)
		m := Inputs{}
		for _, part := range res {
			pair := strings.Split(part, ",")
			if len(pair) < 2 {
				fmt.Println(pair, part, recipe, components)
				return nil, errors.New("fuck")
			}
			first, err := strconv.Atoi(pair[0])
			if err != nil {
				return nil, err
			}
			second, err := strconv.Atoi(pair[1])
			if err != nil {
				return nil, err
			}
			m[first] = second
		}
		err := listOfRecipes.AddItem(Item{
			ID:     recipe,
			Recipe: []Inputs{m},
		})
		if err != nil {
			return nil, err
		}
		// fmt.Println(recipe, len(res))
	}
	return listOfRecipes, nil
}

type Items []Item

func (items *Items) AddItem(item Item) error {
	for k, v := range *items {
		if v.ID == item.ID {
			return fmt.Errorf("Deal With This! %v, %v", k, v)
		}
	}
	*items = append(*items, item)
	return nil
}

type Item struct {
	ID         int
	Name       string
	IsTransmog bool
	Recipe     []Inputs
}

type Inputs map[int]int

func (items *Items) GetDeepest() (ID int, depth int) {
	GoDeeper := func(inputs Inputs, curDepth int) (deeperID, deeperDepth int) {
		return -1, -1
	}
	GoDeeper = func(inputs Inputs, curDepth int) (deeperID, deeperDepth int) {
		for k, _ := range inputs {
			for _, item := range *items {
				if item.ID == k {
					if len(item.Recipe) == 0 {
						return item.ID, curDepth
					} else {
						for _, vv := range item.Recipe {
							return GoDeeper(vv, curDepth+1)
						}
					}
				}
			}
		}
		return -2, -2
	}
	for _, v := range *items {
		for _, vv := range v.Recipe {
			i, d := GoDeeper(vv, 1)
			if i > ID {
				ID, depth = i, d
			}
		}
	}
	return ID, depth
}

func (items *Items) PrintFullList() {
	fmt.Println("-----PrintFullList-----")
	fmt.Println(len(*items))
	for _, item := range *items {
		if len(item.Recipe) != 1 {
			fmt.Println(item.ID, len(item.Recipe))
		} else {
			fmt.Println(item.ID, item.Recipe)
		}
	}
}

func (items *Items) PrintRecipeFor(id int) error {
	sort.Stable(items)
	fmt.Println((*items)[0])
	fmt.Println("-----PrintFullList-----")
	fmt.Println(len(*items))
	for _, item := range *items {
		if item.ID == id {
			fmt.Println(item.ID, len(item.Recipe))
			return nil
		}
	}
	return errors.New("Not Found")
}

func (items *Items) Len() int {
	return len(*items)
}
func (items *Items) Less(i, j int) bool {
	return (*items)[i].ID > (*items)[j].ID
}
func (items *Items) Swap(i, j int) {
	(*items)[i], (*items)[j] = (*items)[j], (*items)[i]
}

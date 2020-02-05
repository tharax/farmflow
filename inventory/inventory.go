package inventory

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Inventory []Item

type Item struct {
	ID       int
	Name     string
	Quantity int
}

func NewInventory() Inventory {
	return Inventory{}
}

func (i Inventory) String() string {
	sort.Stable(&i)
	var res string
	for _, v := range i {
		res += fmt.Sprintf("%10d %-40s %5d\n", v.ID, v.Name, v.Quantity)
	}
	return res
}

func (i *Inventory) Add(item Item) error {
	if item.ID <= 0 {
		return fmt.Errorf("trying to add an item with id less than 1: %d", item.ID)
	}
	if item.Name == "" {
		return fmt.Errorf("trying to add a item with no name")
	}
	if item.Quantity <= 0 {
		return fmt.Errorf("trying to add a item with quantity less than 1: %d", item.Quantity)
	}
	added := false
	for k, v := range *i {
		if v.ID == item.ID {
			if v.Name != item.Name {
				(*i)[k].Name = GetBestName(v.Name, item.Name)
				// fmt.Printf("item names don't match: %+v, %+v, newName:%s\n", v.Name, item.Name, GetBestName(v.Name, item.Name))
			}
			(*i)[k].Quantity += item.Quantity
			added = true
		}
	}
	if !added {
		*i = append(*i, item)
	}
	return nil
}

func (inv *Inventory) Len() int {
	return len(*inv)
}

func (inv *Inventory) Less(i, j int) bool {
	if (*inv)[i].ID < (*inv)[j].ID {
		return true
	}
	return false
}

func (inv *Inventory) Swap(i, j int) {
	(*inv)[i], (*inv)[j] = (*inv)[j], (*inv)[i]
}

func GetBestName(old, new string) string {
	var re = regexp.MustCompile(`(?:\|h\[)(.*)(?:\]\|)`)
	strArray := re.FindAllString(old, -1)
	var str string
	if len(strArray) > 0 {
		str = strings.TrimLeft(strArray[0], "|h[")
		str = strings.TrimRight(str, "]|")
	} else {
		return old
	}
	if str != "" {
		return str
	}
	return new
}

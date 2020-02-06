package items

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/labstack/gommon/log"
)

type DB struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name          string `json:"name"`
	ItemID        int    `json:"itemID"`
	Equippable    bool   `json:"equippable"`
	Quality       int    `json:"quality"`
	Bind          int    `json:"bind"`
	Class         int    `json:"class"`
	Subclass      int    `json:"subclass"`
	ILvl          int    `json:"iLvl"`
	InventoryType int    `json:"inventoryType"`
	Lvl           int    `json:"Lvl"`
	Classes       []int  `json:"classes"`
}

type Bind int

const (
	BindNever    Bind = 0
	BindOnEquip  Bind = 1
	BindOnPickup Bind = 2
	BindOnUse    Bind = 3
	QuestItem    Bind = 4
)

func LoadItems() (*DB, error) {
	b, err := ioutil.ReadFile("data/itemDB.json")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var db DB
	err = json.Unmarshal(b, &db)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &db, nil
}

func (db *DB) GetItem(id int) (*Item, error) {
	for _, item := range db.Items {
		if item.ItemID == id {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("item with id %d not found", id)
}

func (db *DB) GetBind() map[int]int {
	res := map[int]int{}
	for _, item := range db.Items {
		res[item.Bind] += 1
	}
	return res
}

func (db *DB) GetInventoryType() map[int]int {
	res := map[int]int{}
	for _, item := range db.Items {
		res[item.InventoryType] += 1
	}
	return res
}

type InventoryType int

const (
	FoodOrMaybeNotEquipable InventoryType = 0
	Head                    InventoryType = 1
	Neck                    InventoryType = 2
	Shoulder                InventoryType = 3
	Shirt                   InventoryType = 4
	Chest                   InventoryType = 5
	Belt                    InventoryType = 6
	Legs                    InventoryType = 7
	Boots                   InventoryType = 8
	Bracers                 InventoryType = 9
	Gloves                  InventoryType = 10
	Ring                    InventoryType = 11
	Trinket                 InventoryType = 12
	OneHandedSword          InventoryType = 13
	Shield                  InventoryType = 14
	Bow                     InventoryType = 15
	Cloak                   InventoryType = 16
	TwoHandedStaff          InventoryType = 17
	Bag                     InventoryType = 18
	Tabard                  InventoryType = 19
	ChestRobe               InventoryType = 20
	MainhandWeapon          InventoryType = 21
	OffhandWeapon           InventoryType = 22
	OffhandBook             InventoryType = 23
	Arrow                   InventoryType = 24
	ThrowingWeapon          InventoryType = 25
	Gun                     InventoryType = 26
)

func (db *DB) GetItemByName(name string) (*Item, error) {
	for _, item := range db.Items {
		if item.Name == name {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("item with name %s not found", name)
}

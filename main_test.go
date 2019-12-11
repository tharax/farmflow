package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func Test_NoMissingRecipes(t *testing.T) {
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
				t.Errorf("Missing Item: %s", k)
			}
		}
	}
}

func Test_NoMissingItems(t *testing.T) {
	b, err := ioutil.ReadFile("data.json")
	if err != nil {
		panic(err)
	}
	var data Data
	err = json.Unmarshal(b, &data)
	if err != nil {
		panic(err)
	}

	for _, item := range data.Items {
		if item.Source == SourceRecipe {
			if data.Recipes.Exists(item.Name) == false {
				t.Errorf("Missing Recipe: %s", item.Name)
			}
		}
	}
}

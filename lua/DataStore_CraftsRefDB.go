package lua

type DataStore_CraftsRefDB struct {
	ProfileKeys map[string]string
	Global      Global
}

type Global struct {
	Reagents map[int]string
}

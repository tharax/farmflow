package main

type Global struct {
	Characters Character
	Guild      Guild
}

type Character map[string]struct {
	Containers map[string]Container
}

type Container struct {
	IDs    []int
	Counts []int
}

type Guild struct {
}

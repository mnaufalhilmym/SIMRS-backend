package pg

type Where struct {
	Query interface{}
	Args  []interface{}
}

type IncludeTables struct {
	Query string
	Args  []interface{}
}

type CountOption struct {
	Where *[]Where
}

type FindOneOption struct {
	Where         *[]Where
	Order         *[]interface{}
	IncludeTables *[]IncludeTables
}

type FindAllOption struct {
	Where         *[]Where
	Order         *[]interface{}
	Limit         *int
	Offset        *int
	IncludeTables *[]IncludeTables
}

type CreateOption struct {
	IsUpsert bool
}

type UpdateOption struct {
	Where *[]Where
}

type ReplaceOption struct {
	Where *[]Where
}

type DestroyOption struct {
	Where *[]Where
}

type Pagination struct {
	Count int
	Limit int
	Total int
}

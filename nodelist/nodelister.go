package nodelist

// Lister is any type that holds a nodelist and can give out nodelist information
type Lister interface {
	IsSelfOnList() bool
}

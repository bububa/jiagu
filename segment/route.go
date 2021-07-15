package segment

type Route struct {
	I int
	V float64
}

var DefaultRoute = NewRoute(0, 0)

func NewRoute(i int, v float64) Route {
	return Route{
		I: i,
		V: v,
	}
}

type RouteSlice []Route

func NewRouteSlice(cap int) RouteSlice {
	return RouteSlice(make([]Route, 0, cap))
}

func (s RouteSlice) Len() int      { return len(s) }
func (s RouteSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s RouteSlice) Less(i, j int) bool {
	if s[i].V == s[j].V {
		return s[i].I < s[j].I
	}
	return s[i].V < s[j].V
}

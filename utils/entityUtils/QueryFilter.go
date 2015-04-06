package entityUtils

type OrQueryFilterContainer struct {
	IsNot      bool
	Expression string
	Arguments  []interface{}
}

func OrQueryFilterContainer__Constructor(isNot bool, expression string, arguments ...interface{}) *OrQueryFilterContainer {
	return &OrQueryFilterContainer{
		IsNot:      isNot,
		Expression: expression,
		Arguments:  arguments,
	}
}

type AndQueryFilterContainer struct {
	OrList []*OrQueryFilterContainer
}

func AndQueryFilterContainer__Constructor(orList []*OrQueryFilterContainer) *AndQueryFilterContainer {
	return &AndQueryFilterContainer{
		OrList: orList,
	}
}

type QueryFilter struct {
	AndList []*AndQueryFilterContainer
}

func QueryFilter__Constructor(andList []*AndQueryFilterContainer) *QueryFilter {
	return &QueryFilter{
		AndList: andList,
	}
}

func QueryFilter__Constructor__SingleCheck(isNot bool, expression string, arguments ...interface{}) *QueryFilter {
	return QueryFilter__Constructor([]*AndQueryFilterContainer{AndQueryFilterContainer__Constructor([]*OrQueryFilterContainer{
		OrQueryFilterContainer__Constructor(isNot, expression, arguments...),
	})})
}

func QueryFilter__Constructor__Empty() *QueryFilter {
	return QueryFilter__Constructor([]*AndQueryFilterContainer{})
}

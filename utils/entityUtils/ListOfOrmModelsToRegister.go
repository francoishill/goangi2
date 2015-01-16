package entityUtils

type iEntity interface {
	TableName() string
}

var listOfOrmModelsToRegister []iEntity

func AddToListOfOrmModelsToRegister(entityModelToRegister iEntity) {
	listOfOrmModelsToRegister = append(listOfOrmModelsToRegister, entityModelToRegister)
}

func ForeachOrmModelToRegister(function func(interface{})) {
	for _, entityModel := range listOfOrmModelsToRegister {
		function(entityModel)
	}
}

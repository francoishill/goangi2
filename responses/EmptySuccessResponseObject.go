package responses

type EmptySuccessResponseObject struct {
	*BaseResponseObject
}

func CreateEmptySuccessResponseObject() *EmptySuccessResponseObject {
	return &EmptySuccessResponseObject{
		BaseResponseObject: CreateBaseResponseObject(true, ""),
	}
}

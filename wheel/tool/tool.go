package wheel

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SprintfStruct ... 格式化输出
func SprintfStruct(s interface{}) string {
	return fmt.Sprintf("%+v", s)
}

// ConObjToStr	...	objectId 转 string
func ConObjToStr(obj primitive.ObjectID) string {
	u := strings.TrimLeft(obj.Hex(), "0")
	return u
}

// string 转  objectId
func ConStrToObj(s string) (primitive.ObjectID, error) {
	sObj, err := primitive.ObjectIDFromHex(fmt.Sprintf("%024s", s))
	if err != nil {
		return primitive.NewObjectID(), err
	}
	return sObj, nil

}

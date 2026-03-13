package frame

import "github.com/afret0/wheel/database"

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// type PageTag struct {
// 	PageTag        string `json:"pageTag"`
// 	PrevPageTag string `json:"forwardPageTag"`
// 	//IsForward      bool   `json:"isForward"`
// 	IsLastPage bool  `json:"isLastPage"`
// 	Count      int64 `json:"count"`
// }

// type PageTag struct {
// 	PageTag        string `json:"pageTag"`
// 	PrevPageTag string `json:"forwardPageTag"`
// }

// Deprecated, use database.Page
type Page = database.Page

//type Page struct {
//	Count       int64  `json:"count"`
//	IsLastPage  bool   `json:"isLastPage"`
//	PageTag     string `json:"pageTag"`
//	PrevPageTag string `json:"prevPageTag"`
//}

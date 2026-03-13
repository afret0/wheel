package frame

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

type Page struct {
	Count       int64  `json:"count"`
	IsLastPage  bool   `json:"isLastPage"`
	PageTag     string `json:"pageTag"`
	PrevPageTag string `json:"prevPageTag"`
}

const DirectionForward = -1
const DirectionBackward = 1

func (p *Page) Direction() (int, string) {
	if p == nil {
		return DirectionBackward, ""
	}

	if p.PrevPageTag != "" {
		return DirectionForward, p.PrevPageTag
	}
	return DirectionBackward, p.PageTag
}

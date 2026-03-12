package frame

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// type PageTag struct {
// 	PageTag        string `json:"pageTag"`
// 	ForwardPageTag string `json:"forwardPageTag"`
// 	//IsForward      bool   `json:"isForward"`
// 	IsLastPage bool  `json:"isLastPage"`
// 	Count      int64 `json:"count"`
// }

type PageTag struct {
	PageTag        string `json:"pageTag"`
	ForwardPageTag string `json:"forwardPageTag"`
}

type PageTagResp struct {
	PageTag    *PageTag `json:"pageTag"`
	Count      int64    `json:"count"`
	IsLastPage bool     `json:"isLastPage"`
}

const DirectionForward = "forward"
const DirectionBackward = "backward"

func (p *PageTag) Direction() string {
	if p == nil {
		return DirectionBackward
	}

	if p.ForwardPageTag == "" {
		return DirectionBackward
	}
	return DirectionForward
}

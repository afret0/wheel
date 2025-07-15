package frame

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PageTag struct {
	PageTag        string `json:"pageTag"`
	ForwardPageTag string `json:"forwardPageTag"`
	//IsForward      bool   `json:"isForward"`
	IsLastPage bool  `json:"isLastPage"`
	Count      int64 `json:"count"`
}

func (p *PageTag) Direction() int64 {
	if p.PageTag == "" && p.ForwardPageTag == "" {
		return -1
	}

	if p.PageTag != "" {
		return -1
	}

	return 1
}

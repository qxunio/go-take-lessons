package comm

type PageVo struct {
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"totalCount"`
}

type PageParam struct {
	Page  int `json:"page" valid:"required~分页page不能为空"`
	Limit int `json:"limit" valid:"required~分页limit不能为空"`
}

func (p *PageParam) GetOffset() int {
	if p.Page-1 <= 0 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

func (p *PageParam) GetLimit() int {
	if p.Limit < 10 {
		return 10
	}
	return p.Limit
}

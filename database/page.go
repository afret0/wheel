package database

import (
	"context"
	"fmt"

	"github.com/afret0/wheel/constant"
	"github.com/afret0/wheel/log"
	"github.com/afret0/wheel/tool"
	"github.com/samber/lo/mutable"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Page struct {
	Count       int64  `json:"count"`
	IsLastPage  bool   `json:"isLastPage"`
	PageTag     string `json:"pageTag"`
	PrevPageTag string `json:"prevPageTag"`
}

//func ConvPage(pt interface{}) *Page {
//	pt1 := &Page{}
//	tool.CopyByJson(pt, pt1)
//	return pt1
//}

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

type AggrListPage[T any] struct {
	L    []T   `json:"l"`
	Page *Page `json:"page"`
}

func FindWithPage[T any](
	ctx context.Context,
	repo *Repository,
	filter bson.M,
	sortField string,
	ptI interface{},
	optChain ...*options.FindOptions,
) (*AggrListPage[T], error) {

	lg := log.CtxLogger(ctx)

	// 默认排序：倒序（最新在前）
	defaultOpt := &options.FindOptions{
		Sort:  bson.M{sortField: -1},
		Limit: tool.Int64Ptr(constant.FindListOffset),
	}

	// 逐字段合并而非整体覆盖, 调用方只需传关心的字段, 其余保留默认值
	opt := options.MergeFindOptions(append([]*options.FindOptions{defaultOpt}, optChain...)...)

	limit := int64(constant.FindListOffset)
	if opt.Limit != nil && *opt.Limit > 0 {
		limit = *opt.Limit
	}

	// -----------------------------
	// 1. 解析分页方向
	// -----------------------------
	ptS, err := tool.Marshal(ptI)
	if err != nil {
		return nil, err
	}

	pt := &Page{}
	err = tool.Unmarshal(ptS, pt)
	if err != nil {
		return nil, err
	}

	direction, tag := pt.Direction()
	if tag != "" {
		ts := tool.ConStringToInt64WithoutErr(tag)
		if ts > 0 {
			// 游标值取自文档 sortField 的原始值(毫秒时间戳), 比较时必须保持同类型,
			// 否则 BSON 跨类型比较会让条件恒真/恒假, 导致翻页失效
			if direction == DirectionBackward {
				// 往后翻页（下一页）
				filter[sortField] = bson.M{"$lt": ts}
			} else {
				// 往前翻页（上一页）
				opt.Sort = bson.M{sortField: 1}
				filter[sortField] = bson.M{"$gt": ts}
			}
		}
	}

	// -----------------------------
	// 2. 查询
	// -----------------------------
	list := make([]T, 0)
	if err := repo.Find(ctx, &list, filter, opt); err != nil {
		lg.Errorf("mongo find error: %v", err)
		return nil, err
	}

	// forward 查询结果需要反转
	if direction == DirectionForward {
		mutable.Reverse(list)
	}

	// -----------------------------
	// 3. 构造分页信息
	// -----------------------------
	nextPage := &Page{
		Count:      int64(len(list)),
		IsLastPage: int64(len(list)) < limit,
	}

	if len(list) == 0 {
		//return list, nextPage, nil
		return &AggrListPage[T]{
			L:    list,
			Page: nextPage,
		}, nil
	}

	// -----------------------------
	// 4. 自动提取 PageTag / PrevPageTag
	// -----------------------------
	// 最后一条记录（下一页用）
	last := list[len(list)-1]
	if val, ok := tool.ExtractFieldValueByBSONTag(last, sortField); ok {
		nextPage.PageTag = fmt.Sprintf("%v", val)
	}

	// 第一条记录（上一页用）
	first := list[0]
	if val, ok := tool.ExtractFieldValueByBSONTag(first, sortField); ok {
		nextPage.PrevPageTag = fmt.Sprintf("%v", val)
	}

	//return list, nextPage, nil
	return &AggrListPage[T]{
		L:    list,
		Page: nextPage,
	}, nil
}

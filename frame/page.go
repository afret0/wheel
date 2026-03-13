package frame

import (
	"context"
	"fmt"

	"github.com/afret0/wheel/constant"
	"github.com/afret0/wheel/database"
	"github.com/afret0/wheel/log"
	"github.com/afret0/wheel/tool"
	"github.com/afret0/wheel/tool/timeTool"
	"github.com/samber/lo/mutable"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//func WithFind[T any](ctx context.Context, repo *database.Repository, filter bson.M, SortField string, pt *Page, optChain ...*options.FindOptions) ([]T, *Page, error) {
//	lg := log.CtxLogger(ctx).WithFields(logrus.Fields{})
//
//	opt := &options.FindOptions{Sort: bson.M{SortField: -1}, Limit: tool.Int64Ptr(constant.FindListOffset)}
//	for _, v := range optChain {
//		opt = v
//	}
//
//	direction, ptS := pt.Direction()
//	if ptS != "" && tool.ConStringToInt64WithoutErr(ptS) >= 0 {
//		filter[SortField] = bson.M{"$lt": timeTool.ParseMillisecond(tool.ConStringToInt64WithoutErr(ptS))}
//		if direction == DirectionForward {
//			opt.Sort = bson.M{SortField: 1}
//			filter[SortField] = bson.M{"$gt": timeTool.ParseMillisecond(tool.ConStringToInt64WithoutErr(ptS))}
//		}
//	}
//
//	l := make([]T, 0)
//	err := repo.Find(ctx, &l, filter, opt)
//	if err != nil {
//		lg.Errorf("err: %d", err)
//		return nil, nil, err
//	}
//
//	if direction == DirectionForward {
//		mutable.Reverse(l)
//	}
//
//	pt1 := &Page{
//		IsLastPage: len(l) < constant.FindListOffset,
//		Count:      int64(len(l)),
//	}
//
//	if len(l) <= 0 {
//		return l, pt1, nil
//	}
//
//	cur, ok := tool.ExtractFieldValueByBSONTag(l[pt1.Count-1], SortField)
//
//	return l, pt1, nil
//
//}

func FindWithPage[T any](
	ctx context.Context,
	repo *database.Repository,
	filter bson.M,
	sortField string,
	pt *Page,
	optChain ...*options.FindOptions,
) ([]T, *Page, error) {

	lg := log.CtxLogger(ctx)

	// 默认排序：倒序（最新在前）
	opt := &options.FindOptions{
		Sort:  bson.M{sortField: -1},
		Limit: tool.Int64Ptr(constant.FindListOffset),
	}

	for _, v := range optChain {
		if v != nil {
			opt = v
		}
	}

	// -----------------------------
	// 1. 解析分页方向
	// -----------------------------
	direction, tag := pt.Direction()
	if tag != "" {
		ts := tool.ConStringToInt64WithoutErr(tag)
		if ts > 0 {
			if direction == DirectionBackward {
				// 往后翻页（下一页）
				filter[sortField] = bson.M{"$lt": timeTool.ParseMillisecond(ts)}
			} else {
				// 往前翻页（上一页）
				opt.Sort = bson.M{sortField: 1}
				filter[sortField] = bson.M{"$gt": timeTool.ParseMillisecond(ts)}
			}
		}
	}

	// -----------------------------
	// 2. 查询
	// -----------------------------
	list := make([]T, 0)
	if err := repo.Find(ctx, &list, filter, opt); err != nil {
		lg.Errorf("mongo find error: %v", err)
		return nil, nil, err
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
		IsLastPage: len(list) < constant.FindListOffset,
	}

	if len(list) == 0 {
		return list, nextPage, nil
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

	return list, nextPage, nil
}

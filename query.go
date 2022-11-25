// @Author xiaozhaofu 2022/11/25 14:41:00
package goes

import (
	"context"
	"errors"

	"github.com/olivere/elastic/v7"
)

// bool query 条件
type EsSearch struct {
	MustQuery    []elastic.Query
	MustNotQuery []elastic.Query
	ShouldQuery  []elastic.Query
	FilterQuery  []elastic.Query
	Sorters      []elastic.Sorter
	From         int // 分页
	Size         int // 每页的数量
}

// 检查索引是否存在
func CheckIndex(index string) error {
	ctx := context.Background()
	exists, err := esclient.IndexExists(index).Do(ctx)
	if err != nil {
		esLog.Printf("userEs init exist failed err is %s\n", err)
		return err
	}
	if !exists {
		return errors.New("查询的索引不存在")
	}
	return nil
}

func BoolQuery(filter *EsSearch) *elastic.BoolQuery {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(filter.MustQuery...)
	boolQuery.MustNot(filter.MustNotQuery...)
	boolQuery.Should(filter.ShouldQuery...)
	boolQuery.Filter(filter.FilterQuery...)

	// 当should不为空时，保证至少匹配should中的一项
	if len(filter.MustQuery) == 0 && len(filter.MustNotQuery) == 0 && len(filter.ShouldQuery) > 0 {
		boolQuery.MinimumShouldMatch("1")
	}
	return boolQuery
}

// Package elasticInitd
// @author:WXZ
// @date:2022/12/1
// @note
package elasticInit

import (
	"context"
	"errors"
	"fmt"
	"sample_week_aggs/conststat"

	"github.com/olivere/elastic/v6"
)

type elasticer struct {
	c *elastic.Client
}

var e *elasticer

//	openEs
//	@Author WXZ
//	@Description:
//	@return *elasticModel.Client
func openEs(host string) (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(host))
	if err != nil {
		return nil, err
	}
	_, code, err := client.Ping(host).Do(context.Background())

	if err != nil {
		return nil, errors.New(fmt.Sprintf("elasticInitd initd pint error:%v,code:%d\n", err, code))
	}

	return client, nil
}

// Client
// @Author WXZ
// @Description: //TODO
// @return *elastic.Client
// @return error
func Client() (*elastic.Client, error) {
	if e != nil && e.c != nil {
		return e.c, nil
	}

	err := New()
	if err != nil {
		return nil, err
	}

	return e.c, nil
}

//	GetEs
//	@Author WXZ
//	@Description:
//	@return *elasticModel.Client
func New() error {
	client, err := openEs(
		conststat.ELASTIC_HOST,
	)

	if err != nil {
		return err
	}

	e = &elasticer{
		c: client,
	}
	return nil
}

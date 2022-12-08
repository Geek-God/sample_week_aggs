// Package mysqlModel
// @author:WXZ
// @date:2022/12/2
// @note
// Package SampleWeekAggs
// @Description:
package SampleWeekAggs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v6"
	"os"
	"path"
	"sample_week_aggs/conststat"
	"sample_week_aggs/model/elasticModel"
	"time"
)

type aggregation struct {
	DocCount           int64              `json:"doc_count"`
	WeekVirusTypeTop   WeekVirusTypeTop   `json:"week_virus_type_top"`
	MiningTotal        MiningTotal        `json:"mining_total"`
	HighRiskTotal      HighRiskTotal      `json:"high_risk_total"`
	WeekWormSpreadTop  WeekWormSpreadTop  `json:"week_worm_spread_top"`
	ExtortionTotal     ExtortionTotal     `json:"extortion_total"`
	WeekVirusSpreadTop WeekVirusSpreadTop `json:"week_virus_spread_top"`
	WeekVirusNameTop   WeekVirusNameTop   `json:"week_virus_name_top"`
}
type Buckets struct {
	Key      string `json:"key"`
	DocCount int64  `json:"doc_count"`
}
type WeekVirusTypeTop struct {
	DocCountErrorUpperBound int64     `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int64     `json:"sum_other_doc_count"`
	Buckets                 []Buckets `json:"buckets"`
}
type MiningTotal struct {
	DocCount int64 `json:"doc_count"`
}
type HighRiskTotal struct {
	DocCount int64 `json:"doc_count"`
}
type WeekWormSpreadTop struct {
	DocCount int64 `json:"doc_count"`
}
type ExtortionTotal struct {
	DocCount int64 `json:"doc_count"`
}
type WeekVirusSpreadTop struct {
	DocCount int64 `json:"doc_count"`
}
type WeekVirusNameTop struct {
	DocCountErrorUpperBound int64     `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int64     `json:"sum_other_doc_count"`
	Buckets                 []Buckets `json:"buckets"`
}
type option func(a *Aggs)

// task
// @Author WXZ
// @Description: //TODO
type Aggs struct {
	start_time string
	end_time   string
	//样本总数
	total int64
	//报毒数量
	virus_total int64
	//本周新收集样本
	week_total int64
	//本周新收集黑样本
	week_virus_total int64
	//本周病毒样本TOP,报毒数量前3的主类型
	week_virus_type_top string
	//高危病毒数量
	high_risk_total int64
	//勒索病毒数量
	extortion_total int64
	//挖矿病毒数量
	mining_total int64
	//病毒家族排名
	week_virus_name_top string
	//本周病毒传播数量
	week_virus_spread_top int64
	week_worm_spread_top  int64
}

// New
// @Author WXZ
// @Description: //TODO
// @param options ...option
// @return *Aggs
func New(options ...option) *Aggs {
	a := &Aggs{}
	for _, v := range options {
		v(a)
	}
	return a
}

// WitchTime
// @Author WXZ
// @Description: //TODO
// @param start string
// @param end string
// @return option
func WitchTime(start, end string) option {
	return func(a *Aggs) {
		a.start_time = start
		a.end_time = end
	}
}

// elasticExport
// @Author WXZ
// @Description: //TODO es查询导出
// @param s mysqlModel.SampleSearchExport
// @return error
func (a *Aggs) searchTotal() error {
	ctx := context.Background()
	sample := elasticModel.Samples{}
	aggs := make(map[string]elastic.Aggregation, 1)

	aggs["aggs_data"] = elastic.NewFilterAggregation().Filter(
		elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("hr_scan_name", "")),
	)

	result, err := elasticModel.GetListAggs(ctx, sample, nil, aggs)
	if err != nil {
		return err
	}

	if result == nil || result.Hits == nil {
		return nil
	}
	a.total = result.Hits.TotalHits

	b, err := result.Aggregations["aggs_data"].MarshalJSON()
	if err != nil {
		return err
	}
	data := make(map[string]int)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	a.virus_total = int64(data["doc_count"])
	return nil
}

// ElasticTimeRange
// @Author WXZ
// @Description: //TODO 查询时间段之内的数据
// @receiver a *Aggs
// @return error
func (a *Aggs) searchTimeRange() error {
	ctx := context.Background()
	sample := elasticModel.Samples{}
	query := elastic.NewBoolQuery()
	aggs := make(map[string]elastic.Aggregation, 1)
	sub_aggs := make(map[string]elastic.Aggregation, 1)
	query.Must(
		elastic.NewRangeQuery("addtime").Gte(a.start_time).Lte(a.end_time),
	)
	//报毒数量前3的主类型
	sub_aggs["week_virus_type_top"] = elastic.NewTermsAggregation().Field("hr_scan_virus_type").Size(3)
	//Virus、Worm、Exploit 这三种主类型的报毒数量的"总和"
	sub_aggs["high_risk_total"] = elastic.NewFilterAggregation().Filter(
		elastic.NewBoolQuery().MustNot(elastic.NewTermsQuery("hr_scan_virus_type", "Virus", "Worm", "Exploit")),
	)
	//Ransom勒索病毒数量
	sub_aggs["extortion_total"] = elastic.NewFilterAggregation().Filter(
		elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("hr_scan_virus_type", "Ransom")),
	)
	//CoinMiner的报毒数量
	sub_aggs["mining_total"] = elastic.NewFilterAggregation().Filter(
		elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("hr_scan_virus_name", "CoinMiner")),
	)
	//报毒数量最多的前三个家族名称
	sub_aggs["week_virus_name_top"] = elastic.NewTermsAggregation().Field("hr_scan_virus_name").Size(3)
	//单独列出  Virus 、 Worm 本周数量
	sub_aggs["week_virus_spread_top"] = elastic.NewFilterAggregation().Filter(
		elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("hr_scan_virus_type", "Virus")),
	)
	sub_aggs["week_worm_spread_top"] = elastic.NewFilterAggregation().Filter(
		elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("hr_scan_virus_type", "Worm")),
	)

	//本周新收集样本XX个，其中黑样本XX个
	filter_aggs := elastic.NewFilterAggregation().Filter(
		elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("hr_scan_name", "")),
	)
	for k, v := range sub_aggs {
		filter_aggs = filter_aggs.SubAggregation(k, v)
	}
	aggs["week_virus_total"] = filter_aggs
	result, err := elasticModel.GetListAggs(ctx, sample, query, aggs)
	if err != nil {
		return err
	}

	if result == nil || result.Hits == nil {
		return nil
	}

	b, err := result.Aggregations["week_virus_total"].MarshalJSON()
	if err != nil {
		return err
	}

	data := &aggregation{}
	err = json.Unmarshal(b, data)
	if err != nil {
		return err
	}

	a.week_total = result.Hits.TotalHits
	a.week_virus_total = data.DocCount

	for _, v := range data.WeekVirusTypeTop.Buckets {
		a.week_virus_type_top += v.Key + " "
	}
	for _, v := range data.WeekVirusNameTop.Buckets {
		a.week_virus_name_top += v.Key + " "
	}
	a.high_risk_total = data.HighRiskTotal.DocCount
	a.extortion_total = data.ExtortionTotal.DocCount
	a.mining_total = data.MiningTotal.DocCount
	a.week_virus_spread_top = data.WeekVirusSpreadTop.DocCount
	a.week_worm_spread_top = data.WeekWormSpreadTop.DocCount
	return nil
}

// MakeWord
// @Author WXZ
// @Description: //TODO 生成word
// @receiver a *Aggs
func (a *Aggs) MakeWord() error {
	err := a.searchTotal()
	if err != nil {
		return err
	}
	err = a.searchTimeRange()
	if err != nil {
		return err
	}
	week := time.Now().Format("2006-01-02")
	file_path := path.Join(conststat.ROOT_PATH, week+".txt")
	f, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	defer f.Close()

	txt := `一周病毒样本态势情况分析
一、企业运行情况
1.目前累计掌握样本 %v 个，其中病毒样本总数为 %v 个。
2.本周新收集样本 %v 个，其中黑样本 %v 个。
3.本周更新产品病毒库XX次，其中正常更新XX次，临时更新XX次。
二、样本统计情况
1.本周病毒样本TOP %v
2.新增重点病毒样本情况分析
	病毒样本分布情况（国外、国内样本分布统计）（国内各省份样本分布情况统计）
	高危病毒数量 %v
	勒索病毒数量 %v
	挖矿病毒数量 %v
3.本周主要传播的病毒类型
	病毒家族排名 %v
	样本涉及APT组织排名 
三、研判分析
1.本周病毒传播数量较上周（上升、下降、平稳），Virus：%v Worm：%v
2.本周影响较大的病毒事件……
3.本周全球反病毒行业重要情况……`
	str := fmt.Sprintf(txt, a.total, a.virus_total, a.week_total, a.week_virus_total, a.week_virus_type_top, a.high_risk_total, a.extortion_total, a.mining_total, a.week_virus_name_top, a.week_virus_spread_top, a.week_worm_spread_top)
	f.WriteString(str)
	return err
}

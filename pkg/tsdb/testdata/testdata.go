package testdata

import (
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/tsdb"
)

type TestDataExecutor struct {
	*tsdb.DataSourceInfo
	log log.Logger
}

func NewTestDataExecutor(dsInfo *tsdb.DataSourceInfo) tsdb.Executor {
	return &TestDataExecutor{
		DataSourceInfo: dsInfo,
		log:            log.New("tsdb.testdata"),
	}
}

func init() {
	tsdb.RegisterExecutor("grafana-testdata-datasource", NewTestDataExecutor)
}

func (e *TestDataExecutor) Execute(queries tsdb.QuerySlice, context *tsdb.QueryContext) *tsdb.BatchResult {
	result := &tsdb.BatchResult{}
	result.QueryResults = make(map[string]*tsdb.QueryResult)

	for _, query := range queries {
		scenarioId := query.Model.Get("scenarioId").MustString("random_walk")
		if scenario, exist := ScenarioRegistry[scenarioId]; exist {
			result.QueryResults[query.RefId] = scenario.Handler(query, context)
			result.QueryResults[query.RefId].RefId = query.RefId
		} else {
			e.log.Error("Scenario not found", "scenarioId", scenarioId)
		}
	}

	return result
}

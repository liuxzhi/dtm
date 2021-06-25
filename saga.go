package dtm

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/yedf/dtm/common"
)

type Saga struct {
	SagaData
	Server string
}

type SagaData struct {
	Gid           string     `json:"gid"`
	TransType     string     `json:"trans_type"`
	Steps         []SagaStep `json:"steps"`
	QueryPrepared string     `json:"query_prepared"`
}
type SagaStep struct {
	Action     string `json:"action"`
	Compensate string `json:"compensate"`
	Data       string `json:"data"`
}

func SagaNew(server string, gid string) *Saga {
	return &Saga{
		SagaData: SagaData{
			Gid:       gid,
			TransType: "saga",
		},
		Server: server,
	}
}
func (s *Saga) Add(action string, compensate string, postData interface{}) *Saga {
	logrus.Printf("saga %s Add %s %s %v", s.Gid, action, compensate, postData)
	step := SagaStep{
		Action:     action,
		Compensate: compensate,
		Data:       common.MustMarshalString(postData),
	}
	s.Steps = append(s.Steps, step)
	return s
}

func (s *Saga) Commit() error {
	logrus.Printf("committing %s body: %v", s.Gid, &s.SagaData)
	resp, err := common.RestyClient.R().SetBody(&s.SagaData).Post(fmt.Sprintf("%s/commit", s.Server))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("commit failed: %v", resp.Body())
	}
	return nil
}

func (s *Saga) Prepare(queryPrepared string) error {
	s.QueryPrepared = common.OrString(queryPrepared, s.QueryPrepared)
	logrus.Printf("preparing %s body: %v", s.Gid, &s.SagaData)
	resp, err := common.RestyClient.R().SetBody(&s.SagaData).Post(fmt.Sprintf("%s/prepare", s.Server))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("prepare failed: %v", resp.Body())
	}
	return nil
}
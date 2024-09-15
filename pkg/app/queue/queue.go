package queue

import (
	"encoding/json"
	"fmt"

	"github.com/olbrichattila/gofra/pkg/app/db"

	builder "github.com/olbrichattila/gosqlbuilder/pkg"
)

func New() Quer {
	return &Queue{}
}

type Quer interface {
	Construct(db.DBer, builder.Builder)
	Dispatch(string, string, map[string]interface{}) error
	Pull(string) (map[string]interface{}, error)
}

type Queue struct {
	db         db.DBer
	sqlBuilder builder.Builder
}

func (q *Queue) Construct(d db.DBer, b builder.Builder) {
	q.db = d
	q.sqlBuilder = b
}

func (q *Queue) Dispatch(topic, name string, message map[string]interface{}) error {
	strMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	sql, err := q.sqlBuilder.Insert("jobs").Fields("topic", "name", "message", "is_visible").Values(topic, name, string(strMessage), 1).AsSQL()
	if err != nil {
		return err
	}

	_, err = q.db.Execute(sql, q.sqlBuilder.GetParams()...)
	if err != nil {
		return err
	}

	return nil
}

func (q *Queue) Pull(topic string) (map[string]interface{}, error) {
	// todo add set is_visible false
	sql, err := q.sqlBuilder.
		Select("jobs").
		Fields("id", "name", "message").
		Where("topic", "=", topic).
		Where("is_visible", "=", 1).
		OrderBy("id"). // This should be id desc, but the SQL builder quotes it, need to add the ASC, DESC to the order by in the builder
		AsSQL()
	if err != nil {
		return nil, err
	}

	var message map[string]interface{}
	res, err := q.db.QueryOne(sql, q.sqlBuilder.GetParams()...)
	if err != nil {
		return nil, err
	}

	if mess, ok := res["message"]; ok {
		err = json.Unmarshal([]byte(mess.(string)), &message)
		if err != nil {
			return nil, err
		}
	}

	if id, ok := res["id"]; ok {

		sql, err := q.sqlBuilder.Delete("jobs").Where("id", "=", id).AsSQL()
		if err != nil {
			return nil, err
		}
		_, err = q.db.Execute(sql, q.sqlBuilder.GetParams()...)
		if err != nil {
			return nil, err
		}
		return message, nil
	}

	return nil, fmt.Errorf("could not get message from queue")
}

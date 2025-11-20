package school

import (
	"academic-api/internal/domain"
	"fmt"

	"github.com/gocraft/dbr/v2"
	"github.com/sirupsen/logrus"
)

type SchoolReader struct {
	domain.Reader
	State string `json:"state"`
}

func NewSchoolReader(cursors domain.CursorSet, pageSize int, state string) *SchoolReader {
	return &SchoolReader{
		Reader: domain.Reader{
			Cursors:  cursors,
			PageSize: pageSize,
		},
		State: state,
	}
}

func (r *SchoolReader) Validate() error {
	if r.State != "" && len(r.State) != 2 {
		return fmt.Errorf("State code not valid.")
	}
	return nil
}

func (r *SchoolReader) Query(db *dbr.Tx) (*School, error) {
	err := r.Validate()
	if err != nil {
		return nil, err
	}

	// TODO: implement filtering by state

	school := &School{}
	err = db.Select("*").
		From("schools").
		Where("id = ?", r.Id).
		LoadOne(school)
	if err != nil {
		logrus.WithError(err).Error("Failed to query schools table.")
		return nil, err
	}

	return school, nil
}

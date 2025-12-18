package service

import (
	"academic-api/internal/domain/school"
	"encoding/json"
	"io"

	"github.com/gocraft/dbr/v2"
	"github.com/sirupsen/logrus"
)

type ISchoolService interface {
	initRequest(reqBody io.ReadCloser) (*school.SchoolRequest, *dbr.Tx, error)
	initWriter(reqBody io.ReadCloser) (*school.School, *dbr.Tx, error)
	Create(reqBody io.ReadCloser) (*school.School, error)
	Query(reqBody io.ReadCloser) (*school.SchoolResponse, error)
}

type SchoolService struct {
	ISchoolService
	DbSession *dbr.Session
}

func NewSchoolService(session *dbr.Session) *SchoolService {
	return &SchoolService{
		DbSession: session,
	}
}

func (s *SchoolService) initRequest(reqBody io.ReadCloser) (*school.SchoolRequest, *dbr.Tx, error) {
	reader := &school.SchoolRequest{}
	err := json.NewDecoder(reqBody).Decode(reader)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode request body.")
		return nil, nil, err
	}

	tx, err := s.DbSession.Begin()
	if err != nil {
		logrus.WithError(err).Error("Failed to create database transaction.")
		return nil, nil, err
	}

	return reader, tx, nil
}

func (s *SchoolService) initWriter(reqBody io.ReadCloser) (*school.School, *dbr.Tx, error) {
	schoolObj := &school.School{}
	err := json.NewDecoder(reqBody).Decode(schoolObj)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode request body.")
		return nil, nil, err
	}

	tx, err := s.DbSession.Begin()
	if err != nil {
		logrus.WithError(err).Error("Failed to create database transaction.")
		return nil, nil, err
	}

	return schoolObj, tx, nil

}

func (s *SchoolService) Create(reqBody io.ReadCloser) (*school.School, error) {
	schoolObj, tx, err := s.initWriter(reqBody)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	err = schoolObj.Create(tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return schoolObj, err
}

func (s *SchoolService) Query(reqBody io.ReadCloser) (*school.SchoolResponse, error) {
	reader, tx, err := s.initRequest(reqBody)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize read transaction.")
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	schools, err := reader.Query(tx)
	if err != nil {
		logrus.WithError(err).Error("Failed to query schools table.")
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return schools, nil
}

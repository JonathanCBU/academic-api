package service

import (
	"academic-api/internal/domain/school"
	"encoding/json"
	"io"

	"github.com/gocraft/dbr/v2"
	"github.com/sirupsen/logrus"
)

type ISchoolService interface {
	Create(reqBody io.ReadCloser) (*school.School, error)
	Query(reqBody io.ReadCloser) (*school.School, error)
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

func (s *SchoolService) Create(reqBody io.ReadCloser) (*school.School, error) {
	schoolObj := &school.School{}
	err := json.NewDecoder(reqBody).Decode(schoolObj)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode request body.")
		return nil, err
	}

	tx, err := s.DbSession.Begin()
	if err != nil {
		logrus.WithError(err).Error("Failed to create database transaction.")
		return nil, err
	}

	err = schoolObj.Create(tx)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return schoolObj, err
}

func (s *SchoolService) Query(reqBody io.ReadCloser) (*school.School, error) {
	reader := &school.SchoolReader{}
	err := json.NewDecoder(reqBody).Decode(reader)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode request body.")
		return nil, err
	}

	tx, err := s.DbSession.Begin()
	if err != nil {
		logrus.WithError(err).Error("Failed to create database transaction.")
		return nil, err
	}

	schoolObj, err := reader.Query(tx)
	if err != nil {
		logrus.WithError(err).Error("Failed to query schools table.")
		return nil, err
	}

	return schoolObj, nil
}

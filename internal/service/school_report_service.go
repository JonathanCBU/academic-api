package service

import (
	schoolreport "academic-api/internal/domain/school_report"
	"encoding/json"
	"io"

	"github.com/gocraft/dbr/v2"
	"github.com/sirupsen/logrus"
)

type ISchoolReportService interface {
	initRequest(reqBody io.ReadCloser) (*schoolreport.SchoolReportRequest, *dbr.Tx, error)
	initWriter(reqBody io.ReadCloser) (*schoolreport.SchoolReport, *dbr.Tx, error)
	Create(reqBody io.ReadCloser) (*schoolreport.SchoolReport, error)
	Query(reqBody io.ReadCloser) (*schoolreport.SchoolReportResponse, error)
}

type SchoolReportService struct {
	ISchoolReportService
	DbSession *dbr.Session
}

func NewSchoolReportService(session *dbr.Session) *SchoolReportService {
	return &SchoolReportService{
		DbSession: session,
	}
}

func (s *SchoolReportService) initWriter(reqBody io.ReadCloser) (*schoolreport.SchoolReport, *dbr.Tx, error) {
	reportObj := &schoolreport.SchoolReport{}
	err := json.NewDecoder(reqBody).Decode(reportObj)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode request body.")
		return nil, nil, err
	}

	tx, err := s.DbSession.Begin()
	if err != nil {
		logrus.WithError(err).Error("Failed to create database transaction.")
		return nil, nil, err
	}

	return reportObj, tx, nil
}

func (s *SchoolReportService) initRequest(reqBody io.ReadCloser) (*schoolreport.SchoolReportRequest, *dbr.Tx, error) {
	reader := &schoolreport.SchoolReportRequest{}
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

func (s *SchoolReportService) Create(reqBody io.ReadCloser) (*schoolreport.SchoolReport, error) {
	logrus.Info("School report service create")
	reportObj, tx, err := s.initWriter(reqBody)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	err = reportObj.Create(tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return reportObj, err
}

func (s *SchoolReportService) Query(reqBody io.ReadCloser) (*schoolreport.SchoolReportResponse, error) {
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

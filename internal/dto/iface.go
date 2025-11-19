package dto

import "academic-api/internal/domain"

type IDto interface {
	ToDomain() domain.Model
}

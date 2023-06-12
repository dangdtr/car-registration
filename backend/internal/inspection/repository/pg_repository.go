package repository

import (
	"backend/internal/inspection"
	"backend/internal/models"
	"backend/pkg/utils"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type inspectionRepo struct {
	db *gorm.DB
}

func (i *inspectionRepo) CountByQuarterAndYear(ctx context.Context, quarter int, year int) (int, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "inspectionRepo.CountByQuarterAndYear")
	defer span.Finish()
	//var count int64
	//err := i.db.Table("inspections").
	//	Where("QUARTER(inspection_date) = ? AND YEAR(inspection_date) = ?", quarter, year).
	//	Count(&count).Error
	//if err != nil {
	//	return 0, err
	//}
	//
	//return int(count), nil
	db := i.db.Table("inspections")

	startDate, endDate := getQuarterDates(quarter, year)

	var count int64
	err := db.
		Where("inspection_date >= ? AND inspection_date <= ?", startDate, endDate).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func getQuarterDates(quarter int, year int) (time.Time, time.Time) {
	var startDate, endDate time.Time

	switch quarter {
	case 1:
		startDate = time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(year, time.March, 31, 23, 59, 59, 0, time.UTC)
	case 2:
		startDate = time.Date(year, time.April, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(year, time.June, 30, 23, 59, 59, 0, time.UTC)
	case 3:
		startDate = time.Date(year, time.July, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(year, time.September, 30, 23, 59, 59, 0, time.UTC)
	case 4:
		startDate = time.Date(year, time.October, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)
	}

	return startDate, endDate
}

func (i *inspectionRepo) GetByInspectionDate(ctx context.Context, month int, year int, query *utils.PaginationQuery) (*models.InspectionsList, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "inspectionRepo.GetByInspectionDate")
	defer span.Finish()

	result := i.db.First(&models.Inspection{})

	totalCount := int(result.RowsAffected)
	if result := i.db.First(&models.Inspection{}); result.Error != nil {
		return nil, errors.Wrap(result.Error, "inspectionRepo.GetByInspectionDate.Results")
	}

	if totalCount == 0 {
		return &models.InspectionsList{
			TotalCount:  totalCount,
			TotalPages:  utils.GetTotalPages(totalCount, query.GetSize()),
			Page:        query.GetPage(),
			Size:        query.GetSize(),
			HasMore:     utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Inspections: make([]*models.Inspection, 0),
		}, nil
	}

	var inspections = make([]*models.Inspection, 0, query.GetSize())
	clauseWhere := ""
	if year != 0 {
		clauseWhere += fmt.Sprintf("EXTRACT(YEAR FROM expiry_date) = %d", year)
	}
	if month != 0 && year == 0 {
		clauseWhere += fmt.Sprintf("EXTRACT(MONTH FROM expiry_date) = %d", month)
	}
	if month != 0 && year != 0 {
		clauseWhere += fmt.Sprintf(" and EXTRACT(MONTH FROM expiry_date) = %d", month)
	}
	if records := i.db.Where(clauseWhere).Limit(query.GetLimit()).Offset(query.GetOffset()).Order(query.GetOrderBy()).Find(&inspections); records.Error != nil {
		return nil, errors.Wrap(records.Error, "inspectionRepo.GetGetByInspectionDateAll.Query")
	}

	return &models.InspectionsList{
		TotalCount:  totalCount,
		TotalPages:  utils.GetTotalPages(totalCount, query.GetSize()),
		Page:        query.GetPage(),
		Size:        query.GetSize(),
		HasMore:     utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Inspections: inspections,
	}, nil
}

func (i *inspectionRepo) GetByExpiryDate(ctx context.Context, month int, year int, query *utils.PaginationQuery) (*models.InspectionsList, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "inspectionRepo.GetByExpiryDate")
	defer span.Finish()

	result := i.db.First(&models.Inspection{})

	totalCount := int(result.RowsAffected)
	if result := i.db.First(&models.Inspection{}); result.Error != nil {
		return nil, errors.Wrap(result.Error, "inspectionRepo.GetByExpiryDate.Results")
	}

	if totalCount == 0 {
		return &models.InspectionsList{
			TotalCount:  totalCount,
			TotalPages:  utils.GetTotalPages(totalCount, query.GetSize()),
			Page:        query.GetPage(),
			Size:        query.GetSize(),
			HasMore:     utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Inspections: make([]*models.Inspection, 0),
		}, nil
	}

	var inspections = make([]*models.Inspection, 0, query.GetSize())
	clauseWhere := ""
	if year != 0 {
		clauseWhere += fmt.Sprintf("EXTRACT(YEAR FROM expiry_date) = %d", year)
	}
	if month != 0 && year == 0 {
		clauseWhere += fmt.Sprintf("EXTRACT(MONTH FROM expiry_date) = %d", month)
	}
	if month != 0 && year != 0 {
		clauseWhere += fmt.Sprintf(" and EXTRACT(MONTH FROM expiry_date) = %d", month)
	}
	//fmt.Println(clauseWhere)
	if records := i.db.Where(clauseWhere).Limit(query.GetLimit()).Offset(query.GetOffset()).Order(query.GetOrderBy()).Find(&inspections); records.Error != nil {
		return nil, errors.Wrap(records.Error, "inspectionRepo.GetByExpiryDate.Query")
	}

	return &models.InspectionsList{
		TotalCount:  totalCount,
		TotalPages:  utils.GetTotalPages(totalCount, query.GetSize()),
		Page:        query.GetPage(),
		Size:        query.GetSize(),
		HasMore:     utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Inspections: inspections,
	}, nil
}

func (i *inspectionRepo) GetAll(ctx context.Context, query *utils.PaginationQuery) (*models.InspectionsList, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "inspectionRepo.GetAll")
	defer span.Finish()

	result := i.db.First(&models.Inspection{})

	totalCount := int(result.RowsAffected)
	if result := i.db.First(&models.Inspection{}); result.Error != nil {
		return nil, errors.Wrap(result.Error, "inspectionRepo.GetAll.Results")
	}

	if totalCount == 0 {
		return &models.InspectionsList{
			TotalCount:  totalCount,
			TotalPages:  utils.GetTotalPages(totalCount, query.GetSize()),
			Page:        query.GetPage(),
			Size:        query.GetSize(),
			HasMore:     utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Inspections: make([]*models.Inspection, 0),
		}, nil
	}

	var inspections = make([]*models.Inspection, 0, query.GetSize())
	if records := i.db.Limit(query.GetLimit()).Offset(query.GetOffset()).Order(query.GetOrderBy()).Find(&inspections); records.Error != nil {
		return nil, errors.Wrap(records.Error, "inspectionRepo.GetAll.Query")
	}

	return &models.InspectionsList{
		TotalCount:  totalCount,
		TotalPages:  utils.GetTotalPages(totalCount, query.GetSize()),
		Page:        query.GetPage(),
		Size:        query.GetSize(),
		HasMore:     utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Inspections: inspections,
	}, nil
}

func (i *inspectionRepo) GetByID(ctx context.Context, ID int) (*models.Inspection, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "inspectionRepo.GetByID")
	defer span.Finish()
	var inspect models.Inspection
	if result := i.db.Where("inspection_id = ?", ID).First(&inspect); result.Error != nil {
		return nil, errors.Wrap(result.Error, "inspectionRepo.GetByID.Where.First")
	}
	return &inspect, nil
}

func (i *inspectionRepo) GetByStationCode(ctx context.Context, stationCode string, query *utils.PaginationQuery) (*models.InspectionsList, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "inspectionRepo.GetByStationCode")
	defer span.Finish()

	result := i.db.First(&models.Inspection{})

	totalCount := int(result.RowsAffected)
	if result := i.db.First(&models.Inspection{}); result.Error != nil {
		return nil, errors.Wrap(result.Error, "inspectionRepo.GetByStationCode.Results")
	}

	if totalCount == 0 {
		return &models.InspectionsList{
			TotalCount:  totalCount,
			TotalPages:  utils.GetTotalPages(totalCount, query.GetSize()),
			Page:        query.GetPage(),
			Size:        query.GetSize(),
			HasMore:     utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Inspections: make([]*models.Inspection, 0),
		}, nil
	}

	var inspections = make([]*models.Inspection, 0, query.GetSize())
	if records := i.db.Where("station_code = ?", stationCode).Limit(query.GetLimit()).Offset(query.GetOffset()).Order(query.GetOrderBy()).Find(&inspections); records.Error != nil {
		return nil, errors.Wrap(records.Error, "inspectionRepo.GetByStationCode.Query")
	}

	return &models.InspectionsList{
		TotalCount:  totalCount,
		TotalPages:  utils.GetTotalPages(totalCount, query.GetSize()),
		Page:        query.GetPage(),
		Size:        query.GetSize(),
		HasMore:     utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Inspections: inspections,
	}, nil
}

func (i *inspectionRepo) GetByRegistrationID(ctx context.Context, registrationID string, query *utils.PaginationQuery) (*models.InspectionsList, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "inspectionRepo.GetByRegistrationID")
	defer span.Finish()

	result := i.db.First(&models.Inspection{})

	totalCount := int(result.RowsAffected)
	if result := i.db.First(&models.Inspection{}); result.Error != nil {
		return nil, errors.Wrap(result.Error, "inspectionRepo.GetByRegistrationID.Results")
	}

	if totalCount == 0 {
		return &models.InspectionsList{
			TotalCount:  totalCount,
			TotalPages:  utils.GetTotalPages(totalCount, query.GetSize()),
			Page:        query.GetPage(),
			Size:        query.GetSize(),
			HasMore:     utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Inspections: make([]*models.Inspection, 0),
		}, nil
	}

	var inspections = make([]*models.Inspection, 0, query.GetSize())
	if records := i.db.Where("registration_id = ?", registrationID).Limit(query.GetLimit()).Offset(query.GetOffset()).Order(query.GetOrderBy()).Find(&inspections); records.Error != nil {
		return nil, errors.Wrap(records.Error, "inspectionRepo.GetUsers.Query")
	}

	return &models.InspectionsList{
		TotalCount:  totalCount,
		TotalPages:  utils.GetTotalPages(totalCount, query.GetSize()),
		Page:        query.GetPage(),
		Size:        query.GetSize(),
		HasMore:     utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Inspections: inspections,
	}, nil
}

func (i *inspectionRepo) Create(ctx context.Context, inspection *models.Inspection) (*models.Inspection, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "inspectionRepo.Create")
	defer span.Finish()

	if result := i.db.Create(&inspection); result.Error != nil {
		return nil, errors.Wrap(result.Error, "inspectionRepo.Create.Create")
	}

	return inspection, nil
}

func (i *inspectionRepo) Update(ctx context.Context, inspection *models.Inspection) (*models.Inspection, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "inspectionRepo.Update")
	defer span.Finish()

	if result := i.db.Where("inspection_id = ?", inspection.InspectionID).Updates(&inspection); result.Error != nil {
		return nil, errors.Wrap(result.Error, "inspectionRepo.Update.Where.Update")
	}

	return inspection, nil
}

func (i *inspectionRepo) Delete(ctx context.Context, ID int) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "inspectionRepo.Delete")
	defer span.Finish()

	if result := i.db.Where("inspection_id = ?", ID).Delete(&models.User{}); result.Error != nil {
		return errors.Wrap(result.Error, "inspectionRepo.Delete.Where.Update")
	}

	return nil
}

// NewInspectionRepository Ins Repository constructor
func NewInspectionRepository(db *gorm.DB) inspection.Repository {
	return &inspectionRepo{db: db}
}

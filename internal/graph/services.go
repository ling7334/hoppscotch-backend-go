package graph

import (
	"dto"
	"encoding/json"
	"fmt"
	"model"

	"github.com/lucsky/cuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func createTeam(db *gorm.DB, name, creator string) (*model.Team, error) {
	teamID := cuid.New()
	user := model.User{}
	err := db.First(&user, "uid=?", creator).Error
	if err != nil {
		return nil, err
	}
	team := &model.Team{ID: teamID, Name: name, Teammembers: []model.TeamMember{
		{ID: cuid.New(), TeamID: teamID, UserUID: creator, Role: model.OWNER, User: user},
	}}
	err = db.Create(team).Error
	return team, err
}

func removeTeamMember(db *gorm.DB, teamID, userUID string) (string, error) {
	var member []model.TeamMember
	result := db.Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
	).Delete(
		&member,
		`"teamID" = ? AND "userUid" = ?`, teamID, userUID,
	)
	if result.Error != nil {
		return "", result.Error
	}
	return member[0].ID, nil
}

func removeTeamCollection(db *gorm.DB, collectionID string) error {
	colls := []model.TeamCollection{}
	if err := db.Where(`"parentID"=?`, collectionID).Find(&colls).Error; err != nil {
		return err
	}
	for _, coll := range colls {
		if err := removeTeamCollection(db, coll.ID); err != nil {
			return err
		}
	}
	if err := db.Delete(&model.TeamRequest{}, `"collectionID"=?`, collectionID).Error; err != nil {
		return err
	}
	if err := db.Delete(&model.TeamCollection{}, "id=?", collectionID).Error; err != nil {
		return err
	}
	return nil
}

func getUserMaxOrderIndex(db *gorm.DB, model model.Orderable, userUID string, parentID *string) int32 {
	var maxValue int32
	base := db.Model(model).Select(`MAX("orderIndex")`).Where(`"UserUid"=?`, userUID)
	if parentID != nil {
		base.Where(fmt.Sprintf(`"%s" = ?`, model.ParentColName()), *parentID).Row().Scan(&maxValue)
	} else {
		base.Where(fmt.Sprintf(`"%s" IS NULL`, model.ParentColName())).Row().Scan(&maxValue)
	}
	return maxValue
}

func getTeamMaxOrderIndex(db *gorm.DB, model model.Orderable, teamID *string, parentID *string) int32 {
	var maxValue int32
	base := db.Model(model).Select(`MAX("orderIndex")`)
	if teamID != nil {
		base = base.Where(`"teamID"=?`, teamID)
	}
	if parentID != nil {
		base.Where(fmt.Sprintf(`"%s" = ?`, model.ParentColName()), *parentID).Row().Scan(&maxValue)
	} else {
		base.Where(fmt.Sprintf(`"%s" IS NULL`, model.ParentColName())).Row().Scan(&maxValue)
	}
	return maxValue
}

func moveTeamOrderIndex(db *gorm.DB, model model.Orderable, teamID string, parentID *string, upper, lower *int32, up bool) error {
	expr := gorm.Expr(`"orderIndex" + 1`)
	if up {
		expr = gorm.Expr(`"orderIndex" - 1`)
	}
	base := db.Model(model).Where(fmt.Sprintf(`"teamID"=? AND "%s"=?`, model.ParentColName()), teamID, parentID)
	if upper != nil {
		base = base.Where(`"orderIndex">?`, *upper)
	}
	if lower != nil {
		base = base.Where(`"orderIndex"<?`, *lower)
	}
	return base.Update(`"orderIndex"`, expr).Error
}

func moveUserOrderIndex(db *gorm.DB, model model.Orderable, teamID string, parentID *string, upper, lower *int32, up bool) error {
	expr := gorm.Expr(`"orderIndex" + 1`)
	if up {
		expr = gorm.Expr(`"orderIndex" - 1`)
	}
	base := db.Model(model).Where(fmt.Sprintf(`"userUid"=? AND "%s"=?`, model.ParentColName()), teamID, parentID)
	if upper != nil {
		base = base.Where(`"orderIndex">?`, *upper)
	}
	if lower != nil {
		base = base.Where(`"orderIndex"<?`, *lower)
	}
	return base.Update(`"orderIndex"`, expr).Error
}

func createTeamRequest(db *gorm.DB, teamID string, collID string, request dto.TeamRequestExportJSON, order int32) (*model.TeamRequest, error) {
	req := model.ReqDetail{
		V:                json.Number(request.V),
		Auth:             request.Auth,
		Body:             request.Body,
		Name:             request.Name,
		Method:           request.Method,
		Params:           request.Params,
		Headers:          request.Headers,
		Endpoint:         request.Endpoint,
		TestScript:       request.TestScript,
		PreRequestScript: request.PreRequestScript,
	}
	child := &model.TeamRequest{
		ID:           cuid.New(),
		CollectionID: collID,
		TeamID:       teamID,
		Title:        request.Name,
		Request:      req,
		OrderIndex:   order,
	}
	err := db.Create(child).Error
	return child, err
}

func createTeamCollection(db *gorm.DB, teamID string, colls []dto.TeamCollectionImportJSON, parentCollectionID *string) error {
	for i, coll := range colls {
		this := &model.TeamCollection{}
		order := int32(i)
		if parentCollectionID == nil {
			var err error
			data, err := json.Marshal(coll.Data)
			if err != nil {
				return err
			}
			datastr := string(data)
			this, err = createRootTeamCollection(db, coll.Name, &datastr, teamID, &order)
			if err != nil {
				return err
			}
		} else {
			var err error
			if err := db.Where(`"id"=?`, parentCollectionID).First(this).Error; err != nil {
				return err
			}
			data, err := json.Marshal(coll.Data)
			if err != nil {
				return err
			}
			datastr := string(data)
			this, err = createChildTeamCollection(db, coll.Name, &datastr, *this, &order)
			if err != nil {
				return err
			}
		}
		err := createTeamCollection(db, teamID, coll.Folders, &this.ID)
		if err != nil {
			return err
		}
		for i, r := range coll.Requests {
			_, err := createTeamRequest(db, teamID, this.ID, r, int32(i))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createRootTeamCollection(db *gorm.DB, title string, data *string, teamID string, order *int32) (*model.TeamCollection, error) {
	if order == nil {
		*order = getTeamMaxOrderIndex(db, &model.TeamCollection{}, &teamID, nil) + 1
	}
	child := &model.TeamCollection{
		ID:         cuid.New(),
		ParentID:   nil,
		TeamID:     teamID,
		Title:      title,
		Data:       data,
		OrderIndex: *order,
	}
	err := db.Create(child).Error
	return child, err
}

func createChildTeamCollection(db *gorm.DB, title string, data *string, coll model.TeamCollection, order *int32) (*model.TeamCollection, error) {
	if order == nil {
		*order = getTeamMaxOrderIndex(db, &coll, &coll.TeamID, &coll.ID) + 1
	}
	child := &model.TeamCollection{
		ID:         cuid.New(),
		ParentID:   &coll.ID,
		TeamID:     coll.TeamID,
		Title:      title,
		Data:       data,
		OrderIndex: *order,
		Parent:     &coll,
	}
	err := db.Create(child).Error
	return child, err
}

func moveTeamCollection(db *gorm.DB, coll *model.TeamCollection, parentCollectionID *string) error {
	tx := db.Begin()
	if err := moveTeamOrderIndex(tx, coll, coll.TeamID, coll.ParentID, &coll.OrderIndex, nil, true); err != nil {
		tx.Rollback()
		return err
	}
	coll.ParentID = parentCollectionID
	coll.OrderIndex = getTeamMaxOrderIndex(tx, coll, &coll.TeamID, parentCollectionID) + 1
	if err := tx.Save(&coll).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func moveUserCollection(db *gorm.DB, coll *model.UserCollection, parentCollectionID *string) error {
	tx := db.Begin()
	if err := moveUserOrderIndex(tx, coll, coll.UserUID, coll.ParentID, &coll.OrderIndex, nil, true); err != nil {
		tx.Rollback()
		return err
	}
	coll.ParentID = parentCollectionID
	coll.OrderIndex = getUserMaxOrderIndex(tx, coll, coll.UserUID, parentCollectionID) + 1
	if err := tx.Save(&coll).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func moveTeamRequest(db *gorm.DB, req *model.TeamRequest, parentCollectionID string) error {
	tx := db.Begin()
	err := moveTeamOrderIndex(tx, req, req.TeamID, &req.CollectionID, &req.OrderIndex, nil, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	req.OrderIndex = getTeamMaxOrderIndex(tx, req, &req.TeamID, &parentCollectionID) + 1
	err = tx.Save(&req).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func moveUserRequest(db *gorm.DB, uid string, sourceCollectionID string, requestID string, destinationCollectionID string, nextRequestID *string) (*dto.UserRequestReorderData, error) {
	req := &model.UserRequest{}
	nextReq := &model.UserRequest{}
	if err := db.Where(`"userUid" = ?`, uid).First(req, "id=?", requestID).Error; err != nil {
		return nil, err
	}
	if err := db.Where(`"userUid" = ?`, uid).First(nextReq, "id=?", nextRequestID).Error; err != nil {
		return nil, err
	}
	tx := db.Begin()
	if err := tx.Model(req).Where(`"collectionID"=?`, sourceCollectionID).Where(`"orderIndex">?`, req.OrderIndex).Update(`"orderIndex"`, gorm.Expr(`"orderIndex" - 1`)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Model(req).Where(`"collectionID"=?`, destinationCollectionID).Where(`"orderIndex">=?`, nextReq.OrderIndex).Update(`"orderIndex"`, gorm.Expr(`"orderIndex" + 1`)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	req.CollectionID = destinationCollectionID
	req.OrderIndex = nextReq.OrderIndex
	if err := tx.Save(req).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	if err := db.Where(`"userUid" = ?`, uid).First(req, "id=?", requestID).Error; err != nil {
		return nil, err
	}
	if err := db.Where(`"userUid" = ?`, uid).First(nextReq, "id=?", nextRequestID).Error; err != nil {
		return nil, err
	}
	return &dto.UserRequestReorderData{
		Request:     req,
		NextRequest: nextReq,
	}, nil
}

func updateTeamCollectionOrder(db *gorm.DB, coll *model.TeamCollection, destColl *model.TeamCollection) (*dto.CollectionReorderData, error) {
	tx := db.Begin()
	if coll.OrderIndex < destColl.OrderIndex {
		if err := moveTeamOrderIndex(db, destColl, coll.TeamID, coll.ParentID, &coll.OrderIndex, &destColl.OrderIndex, true); err != nil {
			tx.Rollback()
			return nil, err
		}
		coll.OrderIndex = destColl.OrderIndex - 1
		if err := tx.Save(coll).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		upper := destColl.OrderIndex - 1
		if err := moveTeamOrderIndex(db, destColl, coll.TeamID, coll.ParentID, &upper, &coll.OrderIndex, false); err != nil {
			tx.Rollback()
			return nil, err
		}
		coll.OrderIndex = upper
		if err := tx.Save(coll).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &dto.CollectionReorderData{
		Collection:     coll,
		NextCollection: destColl,
	}, nil
}

func updateTeamRequestOrder(db *gorm.DB, req *model.TeamRequest, nextReq *model.TeamRequest) (*dto.RequestReorderData, error) {
	if req.OrderIndex != nextReq.OrderIndex-1 {
		tx := db.Begin()
		if req.OrderIndex < nextReq.OrderIndex {
			if err := moveTeamOrderIndex(tx, nextReq, req.TeamID, &req.CollectionID, &req.OrderIndex, &nextReq.OrderIndex, true); err != nil {
				tx.Rollback()
				return nil, err
			}
			req.OrderIndex = nextReq.OrderIndex - 1
			if err := tx.Save(req).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			upper := nextReq.OrderIndex - 1
			if err := moveTeamOrderIndex(tx, nextReq, req.TeamID, &req.CollectionID, &upper, &req.OrderIndex, false); err != nil {
				tx.Rollback()
				return nil, err
			}
			req.OrderIndex = upper + 1
			if err := tx.Save(req).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		tx.Commit()
	}
	return &dto.RequestReorderData{
		Request:     req,
		NextRequest: nextReq,
	}, nil
}

func UserRequestToDTO(req model.UserRequest) dto.UserRequestExportJSON {
	return dto.UserRequestExportJSON{
		ID:               req.ID,
		Name:             req.Title,
		V:                req.Request.V.String(),
		Auth:             req.Request.Auth,
		Body:             req.Request.Body,
		Method:           req.Request.Method,
		Params:           req.Request.Params,
		Headers:          req.Request.Headers,
		Endpoint:         req.Request.Endpoint,
		TestScript:       req.Request.TestScript,
		PreRequestScript: req.Request.PreRequestScript,
	}
}

func UserCollectionToDTO(collection model.UserCollection) dto.UserCollectionExportJSON {
	sub := []dto.UserCollectionExportJSON{}
	for _, c := range collection.Children {
		sub = append(sub, UserCollectionToDTO(c))
	}
	reqs := []dto.UserRequestExportJSON{}
	for _, r := range collection.Requests {
		reqs = append(reqs, UserRequestToDTO(r))
	}
	return dto.UserCollectionExportJSON{
		ID:       collection.ID,
		Name:     collection.Title,
		Data:     collection.Data,
		Folders:  sub,
		Requests: reqs,
	}
}

func TeamRequestToDTO(req model.TeamRequest) dto.TeamRequestExportJSON {
	return dto.TeamRequestExportJSON{
		Name:             req.Title,
		V:                req.Request.V.String(),
		Auth:             req.Request.Auth,
		Body:             req.Request.Body,
		Method:           req.Request.Method,
		Params:           req.Request.Params,
		Headers:          req.Request.Headers,
		Endpoint:         req.Request.Endpoint,
		TestScript:       req.Request.TestScript,
		PreRequestScript: req.Request.PreRequestScript,
	}
}

func TeamCollectionToDTO(collection model.TeamCollection) dto.TeamCollectionExportJSON {
	sub := []dto.TeamCollectionExportJSON{}
	for _, c := range collection.Children {
		sub = append(sub, TeamCollectionToDTO(c))
	}
	reqs := []dto.TeamRequestExportJSON{}
	for _, r := range collection.Requests {
		reqs = append(reqs, TeamRequestToDTO(r))
	}
	return dto.TeamCollectionExportJSON{
		Name:     collection.Title,
		Data:     collection.Data,
		Folders:  sub,
		Requests: reqs,
	}
}

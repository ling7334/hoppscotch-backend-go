package graph

import (
	"dto"
	"encoding/json"
	"fmt"
	"model"
	"strconv"
	"strings"

	ex "exception"

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

func ShortcodeToDTO(sc *model.Shortcode) *dto.ShortcodeWithUserEmail {
	return &dto.ShortcodeWithUserEmail{
		ID:         sc.ID,
		Request:    sc.Request,
		Properties: sc.EmbedProperties,
		CreatedOn:  sc.CreatedOn,
		Creator: &dto.ShortcodeCreator{
			UID:   sc.Creator.UID,
			Email: *sc.Creator.Email,
		},
	}
}

func AssignUserCollections(colls *[]model.UserCollection, reqType model.ReqType, parentCollectionID *string) {
	for _, coll := range *colls {
		coll.ID = cuid.New()
		coll.ParentID = parentCollectionID
		for _, r := range coll.Requests {
			r.ID = cuid.New()
			r.CollectionID = coll.ID
			r.Type = reqType
		}
		AssignUserCollections(&coll.Children, reqType, &coll.ID)
	}
	return
}

// Following fields are not updatable by `infraConfigs` Mutation. Use dedicated mutations for these fields instead.
// EXCLUDE_FROM_UPDATE_CONFIGS := []string{
// 	"VITE_ALLOWED_AUTH_PROVIDERS",
// 	"ALLOW_ANALYTICS_COLLECTION",
// 	"ANALYTICS_USER_ID",
// 	"IS_FIRST_TIME_INFRA_SETUP",
// 	"MAILER_SMTP_ENABLE",
// 	"USER_HISTORY_STORE_ENABLED",
// }

// Following fields can not be fetched by `infraConfigs` Query. Use dedicated queries for these fields instead.
// EXCLUDE_FROM_FETCH_CONFIGS := []string{
// 	"VITE_ALLOWED_AUTH_PROVIDERS",
// 	"ANALYTICS_USER_ID",
// 	"IS_FIRST_TIME_INFRA_SETUP",
// }

// -------------------------------
// Casting and helpers
// -------------------------------

func cast(dbCfg model.InfraConfig) (*model.InfraConfig, error) {
	var plainValue string
	enable := "ENABLE"
	disable := "DISABLE"
	if dbCfg.IsEncrypted {
		v, err := Decrypt(*dbCfg.Value, nil)
		if err != nil {
			return nil, err
		}
		plainValue = v
	} else {
		plainValue = *dbCfg.Value
	}

	// 模拟 TS 中的转换逻辑
	if plainValue == "true" {
		dbCfg.Value = &enable
	} else {
		dbCfg.Value = &disable
	}

	return &dbCfg, nil
}

// GetInfraConfigsMap 返回 map[name]value（自动解密）
func GetInfraConfigsMap(r *gorm.DB) (map[string]string, error) {
	var cfgs []*model.InfraConfig
	err := r.Find(&cfgs).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]string)
	for _, c := range cfgs {
		v, err := cast(*c)
		if err != nil {
			return nil, err
		}
		out[string(c.Name)] = *v.Value
	}
	return out, nil
}

// -------------------------------
// CRUD-like Operations
// -------------------------------

func UpdateInfraConfig(r *gorm.DB, name dto.InfraConfigEnum, value string, isEnc bool) (res model.InfraConfig, err error) {
	err = r.Model(&res).Where("name = ?", name).Clauses(clause.Returning{}).Updates(model.InfraConfig{Value: &value, IsEncrypted: isEnc}).Error
	return
}

// Update 单个 infra config
func Update(r *gorm.DB, name dto.InfraConfigEnum, value string, restartEnabled bool) (model.InfraConfig, error) {
	// 验证
	if ok, err := validateEnvValues([]model.InfraConfig{{Name: name.String(), Value: &value}}); !ok {
		return model.InfraConfig{}, err
	}

	// 先查 isEncrypted
	var existing model.InfraConfig
	err := r.First(&existing, "name = ?", name).Error
	if err != nil {
		return model.InfraConfig{}, ex.ErrInfraConfigNotFound
	}

	saveValue := value
	isEnc := false
	if existing.IsEncrypted {
		enc, err := Encrypt(value, nil)
		if err != nil {
			return model.InfraConfig{}, err
		}
		saveValue = enc
		isEnc = true
	}

	updated, err := UpdateInfraConfig(r, name, saveValue, isEnc)
	if err != nil {
		return model.InfraConfig{}, ex.ErrInfraConfigUpdateFailed
	}

	if restartEnabled {
		fmt.Println("Restart requested after update (implement actual restart)")
	}

	casted, err := cast(updated)
	if err != nil {
		return model.InfraConfig{}, err
	}
	return *casted, nil
}

func validateEnvValues(configs []model.InfraConfig) (bool, error) {
	for _, cfg := range configs {
		name := cfg.Name
		value := cfg.Value
		fail := func() (bool, error) {
			fmt.Printf("[Infra Validation Failed] Key: %s\n", name)
			return false, ex.ErrInfraConfigInvalidInput
		}

		switch name {
		case "MAILER_SMTP_ENABLE", "MAILER_USE_CUSTOM_CONFIGS", "MAILER_SMTP_SECURE", "MAILER_TLS_REJECT_UNAUTHORIZED":
			if *value != "true" && *value != "false" {
				return fail()
			}
		case "MAILER_SMTP_URL":
			if !validateSMTPUrl(*value) {
				return fail()
			}
		case "MAILER_ADDRESS_FROM":
			if !validateSMTPEmail(*value) {
				return fail()
			}
		case "MAILER_SMTP_HOST", "MAILER_SMTP_PORT", "MAILER_SMTP_USER", "MAILER_SMTP_PASSWORD",
			"GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET", "GOOGLE_SCOPE",
			"GITHUB_CLIENT_ID", "GITHUB_CLIENT_SECRET", "GITHUB_SCOPE",
			"MICROSOFT_CLIENT_ID", "MICROSOFT_CLIENT_SECRET", "MICROSOFT_SCOPE", "MICROSOFT_TENANT":
			if *value == "" {
				return fail()
			}
		case "GOOGLE_CALLBACK_URL", "GITHUB_CALLBACK_URL", "MICROSOFT_CALLBACK_URL":
			if !validateUrl(*value) {
				return fail()
			}
		case "VITE_ALLOWED_AUTH_PROVIDERS":
			parts := strings.Split(*value, ",")
			if len(parts) == 0 {
				return fail()
			}
			for _, p := range parts {
				switch strings.ToUpper(p) {
				case "GOOGLE", "GITHUB", "MICROSOFT", "EMAIL":
					// OK
				default:
					return fail()
				}
			}
		case "TOKEN_SALT_COMPLEXITY", "MAGIC_LINK_TOKEN_VALIDITY", "ACCESS_TOKEN_VALIDITY", "REFRESH_TOKEN_VALIDITY", "RATE_LIMIT_TTL", "RATE_LIMIT_MAX":
			n, err := strconv.Atoi(*value)
			if err != nil || n < 1 {
				return fail()
			}
		default:
			// no-op
		}
	}
	return true, nil
}

// -------------------------------
// Validation helpers
// -------------------------------

func validateSMTPUrl(u string) bool {
	// 简单检查
	return strings.HasPrefix(u, "smtp://") || strings.HasPrefix(u, "smtps://")
}
func validateSMTPEmail(e string) bool {
	// 非严格 email 检查
	return strings.Contains(e, "@")
}
func validateUrl(u string) bool {
	return strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")
}

package controller

//type SnapshotController struct{}
//
//func (this SnapshotController) RegisterRoute(g *echo.Group) {
//	g.POST("/data", this.Post)
//	g.GET("/data", this.Get)
//}
//
//func (this *SnapshotController) Post(ctx echo.Context) error {
//	var ps []model.WishOrginalData
//
//	b, err := ioutil.ReadAll(ctx.Request().Body)
//	if err != nil {
//		return ctx.String(http.StatusBadRequest, err.Error())
//	}
//	err = json.Unmarshal(b, ps)
//	if err != nil {
//		return ctx.String(http.StatusBadRequest, err.Error())
//	}
//
//	for _, p := range ps {
//		if len(p.Data.Contest.ID) <= 0 {
//			continue
//		}
//
//		value, err := json.Marshal(&p)
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//
//		var ps model.ProductSnapshot
//		_, err = ini.AppWish.
//			Where("wish_id = ?", p.Data.Contest.ID).
//			And("created = ?", now.BeginningOfDay()).Get(&ps)
//
//		if err == nil && len(ps.Data) > 0 {
//			//d := ZipBytes(value)
//			//ps.Data = string(d)
//			//if _, err = o.Update(&ps, "data"); err != nil {
//			//	utility.Errorln(0, err)
//			//}
//		} else {
//			d := ZipBytes(value)
//			ps.Data = string(d)
//			ps.Created = now.BeginningOfDay()
//			ps.WishId = p.Data.Contest.ID
//
//			_, err := ini.AppWish.Insert(&ps)
//			if err != nil {
//				util.Errorln(0, err)
//			}
//		}
//	}
//
//	return ctx.String(http.StatusOK, "")
//
//}
//
//func (this *SnapshotController) Get(ctx echo.Context) error {
//	var JSONData WishIdJson
//	JSONData.Code = 200
//	page, err := strconv.Atoi(ctx.QueryParam("page"))
//	if err != nil {
//		page = 0
//	}
//	var datas []string
//	var start = 0
//	var end = 0
//
//	if page*size+size > global.WeekSalesCacheLenght {
//		start = page * size
//		end = global.WeekSalesCacheLenght - page*size
//	} else {
//		start = page * size
//		end = page*size + size
//	}
//
//	if ids, err := ini.RedisClient.
//		LRange(global.WEEK_SALES_GREATER_THAN_ZERO, int64(start), int64(end)).
//		Result(); err == nil {
//		datas = ids
//	} else {
//		util.Errorln(0, err)
//	}
//
//	if len(datas) > 0 {
//		JSONData.Data = datas
//		JSONData.Users = model.GetUsers()
//	}
//
//	return ctx.JSON(http.StatusOK, JSONData)
//}

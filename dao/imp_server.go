package dao

//
//type serverImp struct {
//	*baseDao
//}
//
//func NewServerImp(db *sql.DB) Server {
//	return &serverImp{
//		baseDao: newBaseDao(db, &models.Server{}, ChannelServer, TableServer),
//	}
//}
//
//func (s *serverImp) GetById(id string) (srv *models.Server, err error) {
//	srv = &models.Server{}
//	sqlString := fmt.Sprintf("select %s from %s where id= '%s' ;", s.queryFields, s.table, id)
//	err = s.queryResults(sqlString, srv)
//	return
//}
//
//func (s *serverImp) GetByName(name string) (srv *models.Server, err error) {
//	srv = &models.Server{}
//	sqlString := fmt.Sprintf("select %s from %s where name= '%s' ;", s.queryFields, s.table, name)
//	err = s.queryResults(sqlString, srv)
//	return
//}
//
//func (s *serverImp) GetByIdOrName(idOrName string) (idc *models.Server, err error) {
//	if uid.IsUUID(idOrName) {
//		return s.GetById(idOrName)
//	}
//	return s.GetByName(idOrName)
//}
//
//func (s *serverImp) GetList() (srvs []*models.Server, err error) {
//	err = s.queryAll(&srvs)
//	return
//}
//
//func (s *serverImp) Create(srv *models.Server) error {
//	createBase(&srv.Base)
//	return s.insertAndSendEvent(srv, string(srv.Id))
//}
//
//func (s *serverImp) Update(id string, srv *models.Server) error {
//	updateBase(&srv.Base)
//	return s.updateAndSendEvent(newIdCond(id), srv, id)
//}
//
//func (s *serverImp) Delete(id string) error {
//	return s.deleteAndSendEvent(newIdCond(id), id)
//}
//
//func (s *serverImp) QueryIdcReference(idc_id string) (bool, error) {
//	c, err := s.queryCount(newCond().eq("idc_id", idc_id).String())
//	if err != nil {
//		return false, err
//	}
//	if c > 0 {
//		return true, nil
//	}
//	return false, nil
//}

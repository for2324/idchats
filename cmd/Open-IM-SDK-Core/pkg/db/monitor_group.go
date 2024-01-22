package db

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) CheckAndInsertMonitorGroup(groupInfo *model_struct.LocalGroup) error {
	d.monitorMtx.Lock()
	defer d.monitorMtx.Unlock()
	var g model_struct.LocalGroup
	fmt.Println("CheckAndInsertMonitorGroup,CheckAndInsertMonitorGroup,CheckAndInsertMonitorGroup")
	if err := d.conn.Table(constant.MonitorGroupTableName).Where("group_id = ?", groupInfo.GroupID).Take(&g).Error; err == gorm.ErrRecordNotFound {
		return utils.Wrap(d.conn.Table(constant.MonitorGroupTableName).Create(groupInfo).Error, "CheckAndInsertMonitorGroup failed")
	}
	return nil
}
func (d *DataBase) InsertMonitorGroup(groupInfo *model_struct.LocalGroup) error {
	d.monitorMtx.Lock()
	defer d.monitorMtx.Unlock()
	return utils.Wrap(d.conn.Table(constant.MonitorGroupTableName).Create(groupInfo).Error, "InsertMonitorGroup failed")
}
func (d *DataBase) DeleteMonitorGroup(groupID string) error {
	d.monitorMtx.Lock()
	defer d.monitorMtx.Unlock()
	localGroup := model_struct.LocalGroup{GroupID: groupID}
	return utils.Wrap(d.conn.Table(constant.MonitorGroupTableName).Delete(&localGroup).Error, "DeleteMonitorGroup failed")
}
func (d *DataBase) UpdateMonitorGroup(groupInfo *model_struct.LocalGroup) error {
	d.monitorMtx.Lock()
	defer d.monitorMtx.Unlock()

	t := d.conn.Table(constant.MonitorGroupTableName).Model(groupInfo).Select("*").Updates(*groupInfo)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")

}
func (d *DataBase) GetJoinedMonitorGroupList() ([]*model_struct.LocalGroup, error) {
	d.monitorMtx.Lock()
	defer d.monitorMtx.Unlock()
	var groupList []model_struct.LocalGroup
	err := d.conn.Table(constant.MonitorGroupTableName).Find(&groupList).Error
	var transfer []*model_struct.LocalGroup
	for _, v := range groupList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetJoinedGroupList failed ")
}
func (d *DataBase) GetJoinedMonitorGroupIDList() ([]string, error) {
	groupList, err := d.GetJoinedMonitorGroupList()
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var groupIDList []string
	for _, v := range groupList {
		if v.GroupType == constant.WorkingGroup {
			groupIDList = append(groupIDList, v.GroupID)
		}
	}
	return groupIDList, nil
}

func (d *DataBase) GetJoinedMonitorWorkingGroupGroupList() ([]*model_struct.LocalGroup, error) {
	groupList, err := d.GetJoinedMonitorGroupList()
	var transfer []*model_struct.LocalGroup
	for _, v := range groupList {
		if v.GroupType == constant.WorkingGroup {
			transfer = append(transfer, v)
		}
	}
	return transfer, utils.Wrap(err, "GetJoinedSuperGroupList failed ")
}

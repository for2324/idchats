package main

import (
	"fmt"

	"Open_IM/cmd/init_data/group"
	"Open_IM/cmd/init_data/task"
)

func main() {
	fmt.Println("init data")
	// user
	// task.InitUserHeadInfo()
	// task.InitTaskCheckIsFollowTwitter()
	// group
	// task.InItOfficialGroupMenberJoinTask()
	// task.InitUserCreateSpaceTask()
	// create task

	fmt.Println("init data CreateTaskList start ... ")
	task.CreateTaskList()
	fmt.Println("init data CreateTaskList end ... ")

	fmt.Println("init data AutoDismissGroup start ... ")
	group.AutoDismissGroup()
	fmt.Println("init data AutoDismissGroup end ... ")

	fmt.Println("init data AutoCraeteUserGroup start ... ")
	group.AutoCraeteUserGroup()
	fmt.Println("init data AutoCraeteUserGroup end ... ")

	fmt.Println("init data FollowGroupCreator start ... ")
	group.FollowGroupCreator()
	fmt.Println("init data FollowGroupCreator end ... ")

	fmt.Println("init data success")
}

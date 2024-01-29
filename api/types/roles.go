package types

// TODO: 临时参数定义，后续优化

type Menus struct {
	MenuIDS []int64 `json:"menu_ids"`
}

type Roles struct {
	RoleIds []int64 `json:"role_ids"`
}

type RoleReq struct {
	Memo     string `json:"memo" `      // 备注
	Name     string `json:"name"`       // 名称
	Sequence int    `json:"sequence" `  // 排序值
	ParentID int64  `json:"parent_id" ` // 父级ID
	Status   int8   `json:"status" `    // 0 表示禁用，1 表示启用
}

type UpdateRoleReq struct {
	Memo            string `json:"memo" `      // 备注
	Name            string `json:"name"`       // 名称
	Sequence        int    `json:"sequence" `  // 排序值
	ParentID        int64  `json:"parent_id" ` // 父级ID
	Status          int8   `json:"status" `    // 0 表示禁用，1 表示启用
	ResourceVersion int64  `json:"resource_version"`
}

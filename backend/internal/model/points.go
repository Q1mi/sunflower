package model

// 定义 积分 相关模型

type PointsTransactionInput struct {
	UserId uint64 // 用户ID
	Points int64  // 积分数量
	Desc   string // 描述
	Type   int    // 积分类型
}

// PointsRecordsInput 积分记录查询模型
type PointsRecordsInput struct {
	UserId uint64 // 用户ID
	Limit  int    // 分页
	Offset int    // 分页
}

type PointsRecordsOutput struct {
	List    []*PointsRecordItem
	HasMore bool
	Total   int
}

type PointsRecordItem struct {
	Points          int64  // 积分数量
	TransactionType int    // 积分类型
	Description     string // 描述
	Date            string // 时间
}

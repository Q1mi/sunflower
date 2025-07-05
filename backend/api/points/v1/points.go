package v1

import "github.com/gogf/gf/v2/frame/g"

// SummaryReq 总积分请求结构体
type SummaryReq struct {
	g.Meta `path:"/points/summary" method:"get" sm:"总积分" tags:"积分"`
}

type SummaryRes struct {
	Total int `json:"total"`
}

// RecordsReq 积分明细请求结构体
type RecordsReq struct {
	g.Meta `path:"/points/records" method:"get" sm:"积分明细" tags:"积分"`
	Limit  int `p:"limit" d:"10" dc:"分页大小，默认为10"`
	Offset int `p:"offset" d:"0" dc:"分页偏移"`
}

type RecordsRes struct {
	Total   int       `json:"total"`
	HasMore bool      `json:"hasMore"`
	List    []*Record `json:"list"`
}
type Record struct {
	PointsChange    int64  `json:"pointsChange"`    // 积分变化量
	TransactionType int    `json:"transactionType"` // 交易类型
	Description     string `json:"description"`     // 描述
	TransactionTime string `json:"transactionTime"` // 交易时间
}

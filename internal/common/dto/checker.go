package dto

import (
	"errors"
	"go-tpl/internal/domain/entity"

	"github.com/shopspring/decimal"
)

type RuleDTO struct {
	CustomRule *entity.BiddingCustomRule
	OtherRule  *entity.BiddingOtherRule

	// 扩展字段
	OtherRulePrice []*entity.BiddingRulePrice
}

type BiddingParamsDTO struct {
	Order          *OrderDTO
	CustomRuleList []*entity.BiddingCustomRule
	OtherRuleList  []*RuleDTO
	IsBidding      bool
	IsLimitPrice   bool
}

type BiddingResultDTO struct {
	Status        int                       `json:"status"`        // 0=失败，1=成功
	Price         decimal.Decimal           `json:"price"`         // 计算后价格
	Type          string                    `json:"type"`          // 校验类型
	UID           int64                     `json:"uid"`           // 用户ID
	RuleId        int64                     `json:"ruleId"`        // 规则ID
	RuleName      string                    `json:"ruleName"`      // 规则名称
	IfBottomCover bool                      `json:"ifBottomCover"` // 是否底价覆盖
	RuleInfo      *entity.BiddingCustomRule `json:"ruleInfo"`      // 规则信息（成功时返回）
	UserInfo      *entity.User              `json:"userInfo"`      // 用户信息（成功时返回）
	OtherRuleInfo *entity.BiddingOtherRule  `json:"otherRuleInfo"` // 其他规则信息（成功时返回）
	FailReason    string                    `json:"failReason"`    // 主要失败原因（最后一条失败）
	AllFailReason []string                  `json:"allFailReason"` // 所有失败原因
	AllFailType   []string                  `json:"allFailType"`   // 所有失败类型
}

type CheckSubResultDTO struct {
	Passed     bool   // 是否通过
	FailReason string // 失败原因
	FailType   string // 失败类型
}

type CheckFailCollector struct {
	allReasons []string
	allTypes   []string
	lastReason string
	lastType   string
}

func (c *CheckFailCollector) Add(reason, typ string) {
	c.allReasons = append(c.allReasons, reason)
	c.allTypes = append(c.allTypes, typ)
	c.lastReason = reason
	c.lastType = typ
}

func (c *CheckFailCollector) Get() (allReasons []string, allTypes []string, lastReason, lastType string) {
	return c.allReasons, c.allTypes, c.lastReason, c.lastType
}

func (c *CheckFailCollector) HasErr() bool {
	return len(c.allReasons) > 0
}

func (c *CheckFailCollector) GetErr() error {
	return errors.New(c.lastType + "-" + c.lastReason)
}

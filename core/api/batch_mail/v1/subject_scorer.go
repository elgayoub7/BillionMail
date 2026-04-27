package v1

import "github.com/gogf/gf/v2/frame/g"

type ScoreSubjectReq struct {
	g.Meta   `path:"/batch_mail/score_subject" method:"post" tags:"BatchMail" summary:"Score a subject line for spam likelihood"`
	Subject  string `json:"subject" dc:"Email subject line to score"`
}

type ScoreSubjectRes struct {
	Score    int      `json:"score"`
	Warnings []string `json:"warnings"`
}

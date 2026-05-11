package notification

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sNotification struct{}

func (s *sNotification) Add(ctx context.Context, in model.AddNotificationInput) (out model.AddNotificationOutput, err error) {
	notification := do.Notifications{
		UserId:     in.UserId,
		QuestionId: in.QuestionId,
		AnswerId:   in.AnswerId,
		ReplyToId:  in.ReplyToId,
		Type:       in.Type,
	}
	id, err := dao.Notifications.Ctx(ctx).InsertAndGetId(notification)
	if err != nil {
		return model.AddNotificationOutput{}, err
	}
	out = model.AddNotificationOutput{
		Id: int(id),
	}
	return out, nil
}

// Get 一次 JOIN 拉取所有通知 + 关联的问题/回答/用户信息。
// 原实现 6+ 次 SQL，现在合并为 1 条 JOIN + 1 轮 Go 组装。
func (s *sNotification) Get(ctx context.Context, in model.GetNotificationsInput) (out model.GetNotificationsOutput, err error) {
	// 单条 SQL 拿全部需要的字段
	const query = `
	SELECT
		n.id              AS nid,
		n.type            AS ntype,
		n.is_read         AS is_read,
		n.created_at      AS n_created_at,
		n.question_id     AS qid,
		n.answer_id       AS aid,
		n.reply_to_id     AS rid,
		q.title           AS q_title,
		q.contents        AS q_contents,
		q.dst_user_id     AS q_dst_user_id,
		a.contents        AS a_contents,
		a.user_id         AS a_user_id,
		r.contents        AS r_contents,
		r.user_id         AS r_user_id,
		ua.nickname       AS a_nickname,
		ur.nickname       AS r_nickname
	FROM notifications n
	LEFT JOIN questions q ON q.id = n.question_id
	LEFT JOIN answers   a ON a.id = n.answer_id
	LEFT JOIN answers   r ON r.id = n.reply_to_id
	LEFT JOIN users    ua ON ua.id = a.user_id
	LEFT JOIN users    ur ON ur.id = r.user_id
	WHERE n.user_id = ?
	ORDER BY n.is_read ASC, n.created_at DESC
	`

	rows, err := g.DB().Ctx(ctx).Query(ctx, query, in.UserId)
	if err != nil {
		return model.GetNotificationsOutput{}, err
	}

	out = model.GetNotificationsOutput{
		NewQuestion: make([]model.NotificationNewQuestion, 0),
		NewAnswer:   make([]model.NotificationNewAnswer, 0),
		NewReply:    make([]model.NotificationNewReply, 0),
	}

	for _, row := range rows {
		ntype := row["ntype"].String()
		qid := row["qid"].Int()
		base := model.NotificationBase{
			Id:              int64(row["nid"].Int()),
			QuestionId:      qid,
			QuestionTitle:   row["q_title"].String(),
			QuestionContent: row["q_contents"].String(),
			IsRead:          row["is_read"].Int() == 1,
			CreatedAt:       row["n_created_at"].GTime().TimestampMilli(),
		}

		switch ntype {
		case consts.NewQuestion:
			out.NewQuestion = append(out.NewQuestion, model.NotificationNewQuestion{
				NotificationBase: base,
				UserName:         consts.DefaultUserName,
				UserId:           consts.DefaultUserId,
			})

		case consts.NewAnswer:
			respdName := row["a_nickname"].String()
			respdId := row["a_user_id"].Int()
			out.NewAnswer = append(out.NewAnswer, model.NotificationNewAnswer{
				NotificationBase: base,
				AnswerId:         row["aid"].Int(),
				AnswerContent:    row["a_contents"].String(),
				RespondentName:   respdName,
				RespondentId:     respdId,
			})

		case consts.NewReply:
			respdName := row["r_nickname"].String()
			respdId := row["r_user_id"].Int()
			// 如果问题的目标老师就是当前用户，回复者匿名化
			if row["q_dst_user_id"].Int() == in.UserId {
				respdName = consts.DefaultUserName
				respdId = consts.DefaultUserId
			}
			out.NewReply = append(out.NewReply, model.NotificationNewReply{
				NotificationBase: base,
				AnswerId:         row["aid"].Int(),
				AnswerContent:    row["a_contents"].String(),
				ReplyToId:        row["rid"].Int(),
				ReplyToContent:   row["r_contents"].String(),
				RespondentName:   respdName,
				RespondentId:     respdId,
			})
		}
	}
	return out, nil
}

func (s *sNotification) Update(ctx context.Context, in model.UpdateNotificationInput) (out model.UpdateNotificationOutput, err error) {
	_, err = dao.Notifications.Ctx(ctx).Where(dao.Notifications.Columns().Id, in.Id).
		Update(do.Notifications{IsRead: true})
	if err != nil {
		return model.UpdateNotificationOutput{}, err
	}
	out.IsRead = true
	out.Id = in.Id
	return out, nil
}

func (s *sNotification) UpdateAoQ(ctx context.Context, in model.UpdateAoQInput) (out model.UpdateAoQOutput, err error) {
	_, err = dao.Notifications.Ctx(ctx).Where(dao.Notifications.Columns().UserId, in.UserID).Where(dao.Notifications.Columns().QuestionId, in.QuestionID).Update(do.Notifications{IsRead: true})
	if err != nil {
		return model.UpdateAoQOutput{}, err
	}
	out.QuestionID = in.QuestionID
	out.IsRead = true
	return out, nil
}

func (s *sNotification) Delete(ctx context.Context, in model.DeleteNotificationInput) (out model.DeleteNotificationOutput, err error) {
	_, err = dao.Notifications.Ctx(ctx).Where(dao.Notifications.Columns().Id, in.Id).Delete()
	if err != nil {
		return model.DeleteNotificationOutput{}, err
	}
	out = model.DeleteNotificationOutput{}
	return out, nil
}

func (s *sNotification) NewNotificationCount(ctx context.Context, in model.NewNotificationCountInput) (out model.NewNotificationCountOutput, err error) {
	// 优化：原来 3 次 COUNT 查询合并为 1 次 GROUP BY
	type typeCount struct {
		Type string `json:"type"`
		Cnt  int    `json:"cnt"`
	}
	var counts []typeCount
	err = g.DB().Ctx(ctx).Model("notifications").
		Fields("type, COUNT(*) AS cnt").
		Where("user_id = ? AND is_read = 0", in.UserId).
		Group("type").
		Scan(&counts)
	if err != nil {
		return model.NewNotificationCountOutput{}, err
	}
	for _, c := range counts {
		switch c.Type {
		case consts.NewQuestion:
			out.NewQuestionCount = c.Cnt
		case consts.NewAnswer:
			out.NewAnswerCount = c.Cnt
		case consts.NewReply:
			out.NewReplyCount = c.Cnt
		}
	}
	return out, nil
}

func init() {
	service.RegisterNotification(New())
}

func New() *sNotification {
	return &sNotification{}
}

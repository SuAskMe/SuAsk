package notification

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	qutil "suask/internal/logic/questions_util"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sNotification struct{}

const notificationQuestionQuery = `
	SELECT
		n.id AS nid,
		n.is_read AS is_read,
		n.created_at AS n_created_at,
		n.question_id AS qid,
		CASE WHEN q.deleted_at IS NULL THEN q.title ELSE '' END AS q_title,
		CASE WHEN q.deleted_at IS NULL THEN q.contents ELSE '' END AS q_contents,
		q.dst_user_id AS q_dst_user_id,
		CASE WHEN q.deleted_at IS NULL THEN q.src_user_id ELSE 0 END AS q_src_user_id,
		CASE WHEN q.deleted_at IS NULL THEN uq.nickname ELSE '' END AS q_src_nickname,
		CASE WHEN q.deleted_at IS NULL THEN uq.avatar_file_id ELSE 0 END AS q_src_avatar_file_id
	FROM notifications n
	LEFT JOIN questions q ON q.id = n.question_id
	LEFT JOIN users uq ON uq.id = q.src_user_id AND uq.deleted_at IS NULL
	WHERE n.user_id = ? AND n.type = ? AND n.deleted_at IS NULL
	ORDER BY n.is_read ASC, n.created_at DESC`

const notificationAnswerQuery = `
	SELECT
		n.id AS nid,
		n.is_read AS is_read,
		n.created_at AS n_created_at,
		n.question_id AS qid,
		n.answer_id AS aid,
		CASE WHEN q.deleted_at IS NULL THEN q.title ELSE '' END AS q_title,
		CASE WHEN q.deleted_at IS NULL THEN q.contents ELSE '' END AS q_contents,
		q.dst_user_id AS q_dst_user_id,
		CASE WHEN a.deleted_at IS NULL THEN a.contents ELSE '' END AS a_contents,
		CASE WHEN a.deleted_at IS NULL THEN a.user_id ELSE 0 END AS a_user_id,
		CASE WHEN a.deleted_at IS NULL THEN ua.nickname ELSE '' END AS a_nickname,
		CASE WHEN a.deleted_at IS NULL THEN ua.avatar_file_id ELSE 0 END AS a_avatar_file_id
	FROM notifications n
	LEFT JOIN questions q ON q.id = n.question_id
	LEFT JOIN answers a ON a.id = n.answer_id
	LEFT JOIN users ua ON ua.id = a.user_id AND ua.deleted_at IS NULL
	WHERE n.user_id = ? AND n.type = ? AND n.deleted_at IS NULL
	ORDER BY n.is_read ASC, n.created_at DESC`

const notificationReplyQuery = `
	SELECT
		n.id AS nid,
		n.is_read AS is_read,
		n.created_at AS n_created_at,
		n.question_id AS qid,
		n.answer_id AS aid,
		n.reply_to_id AS rid,
		CASE WHEN q.deleted_at IS NULL THEN q.title ELSE '' END AS q_title,
		CASE WHEN q.deleted_at IS NULL THEN q.contents ELSE '' END AS q_contents,
		q.dst_user_id AS q_dst_user_id,
		CASE WHEN a.deleted_at IS NULL THEN a.contents ELSE '' END AS a_contents,
		CASE WHEN r.deleted_at IS NULL THEN r.contents ELSE '' END AS r_contents,
		CASE WHEN r.deleted_at IS NULL THEN r.user_id ELSE 0 END AS r_user_id,
		CASE WHEN r.deleted_at IS NULL THEN ur.nickname ELSE '' END AS r_nickname,
		CASE WHEN r.deleted_at IS NULL THEN ur.avatar_file_id ELSE 0 END AS r_avatar_file_id
	FROM notifications n
	LEFT JOIN questions q ON q.id = n.question_id
	LEFT JOIN answers a ON a.id = n.answer_id
	LEFT JOIN answers r ON r.id = n.reply_to_id
	LEFT JOIN users ur ON ur.id = r.user_id AND ur.deleted_at IS NULL
	WHERE n.user_id = ? AND n.type = ? AND n.deleted_at IS NULL
	ORDER BY n.is_read ASC, n.created_at DESC`

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

func (s *sNotification) Get(ctx context.Context, in model.GetNotificationsInput) (out model.GetNotificationsOutput, err error) {
	questionRows, err := queryNotificationRows(ctx, notificationQuestionQuery, in.UserId, consts.NewQuestion)
	if err != nil {
		return model.GetNotificationsOutput{}, err
	}
	answerRows, err := queryNotificationRows(ctx, notificationAnswerQuery, in.UserId, consts.NewAnswer)
	if err != nil {
		return model.GetNotificationsOutput{}, err
	}
	replyRows, err := queryNotificationRows(ctx, notificationReplyQuery, in.UserId, consts.NewReply)
	if err != nil {
		return model.GetNotificationsOutput{}, err
	}

	avatarFileIDs := make([]int, 0, len(questionRows)+len(answerRows)+len(replyRows))
	for _, row := range questionRows {
		avatarFileIDs = append(avatarFileIDs, row["q_src_avatar_file_id"].Int())
	}
	for _, row := range answerRows {
		avatarFileIDs = append(avatarFileIDs, row["a_avatar_file_id"].Int())
	}
	for _, row := range replyRows {
		avatarFileIDs = append(avatarFileIDs, row["r_avatar_file_id"].Int())
	}
	avatarURLMap, err := qutil.BatchGetFileURLs(ctx, qutil.CollectUniqueFileIDs(avatarFileIDs))
	if err != nil {
		return model.GetNotificationsOutput{}, gerror.Wrap(err, "resolve notification avatars")
	}

	out = model.GetNotificationsOutput{
		NewQuestion: make([]model.NotificationNewQuestion, 0, len(questionRows)),
		NewAnswer:   make([]model.NotificationNewAnswer, 0, len(answerRows)),
		NewReply:    make([]model.NotificationNewReply, 0, len(replyRows)),
	}

	defaultActor := buildNotificationActor(0, "", 0, avatarURLMap)
	for _, row := range questionRows {
		questioner := buildNotificationActor(
			row["q_src_user_id"].Int(),
			row["q_src_nickname"].String(),
			row["q_src_avatar_file_id"].Int(),
			avatarURLMap,
		)
		out.NewQuestion = append(out.NewQuestion, model.NotificationNewQuestion{
			NotificationBase: buildNotificationBase(row),
			UserAvatar:       questioner.Avatar,
			UserName:         questioner.Name,
			UserId:           questioner.Id,
		})
	}
	for _, row := range answerRows {
		respondent := buildNotificationActor(
			row["a_user_id"].Int(),
			row["a_nickname"].String(),
			row["a_avatar_file_id"].Int(),
			avatarURLMap,
		)
		out.NewAnswer = append(out.NewAnswer, model.NotificationNewAnswer{
			NotificationBase: buildNotificationBase(row),
			AnswerId:         row["aid"].Int(),
			AnswerContent:    row["a_contents"].String(),
			RespondentAvatar: respondent.Avatar,
			RespondentName:   respondent.Name,
			RespondentId:     respondent.Id,
		})
	}
	for _, row := range replyRows {
		respondent := buildNotificationActor(
			row["r_user_id"].Int(),
			row["r_nickname"].String(),
			row["r_avatar_file_id"].Int(),
			avatarURLMap,
		)
		if row["q_dst_user_id"].Int() == in.UserId {
			respondent = defaultActor
		}
		out.NewReply = append(out.NewReply, model.NotificationNewReply{
			NotificationBase: buildNotificationBase(row),
			ReplyToId:        row["rid"].Int(),
			ReplyToContent:   row["r_contents"].String(),
			RespondentAvatar: respondent.Avatar,
			RespondentName:   respondent.Name,
			RespondentId:     respondent.Id,
			AnswerId:         row["aid"].Int(),
			AnswerContent:    row["a_contents"].String(),
		})
	}
	return out, nil
}

func queryNotificationRows(ctx context.Context, query string, userID int, notificationType string) (gdb.Result, error) {
	rows, err := g.DB().Ctx(ctx).Query(ctx, query, userID, notificationType)
	if err != nil {
		return nil, gerror.Wrap(err, "query notifications")
	}
	return rows, nil
}

type notificationActor struct {
	Id     int
	Name   string
	Avatar string
}

func buildNotificationBase(row gdb.Record) model.NotificationBase {
	var createdAt int64
	if t := row["n_created_at"].GTime(); t != nil {
		createdAt = t.TimestampMilli()
	}
	return model.NotificationBase{
		Id:              int64(row["nid"].Int()),
		QuestionId:      row["qid"].Int(),
		QuestionTitle:   row["q_title"].String(),
		QuestionContent: row["q_contents"].String(),
		IsRead:          row["is_read"].Int() == 1,
		CreatedAt:       createdAt,
	}
}

func buildNotificationActor(userID int, nickname string, avatarFileID int, avatarURLMap map[int]string) notificationActor {
	actor := notificationActor{
		Id:     userID,
		Name:   nickname,
		Avatar: resolveNotificationAvatar(avatarFileID, avatarURLMap),
	}
	if actor.Id == 0 {
		actor.Id = consts.DefaultUserId
	}
	if actor.Name == "" {
		actor.Name = consts.DefaultUserName
	}
	return actor
}

func resolveNotificationAvatar(avatarFileID int, avatarURLMap map[int]string) string {
	if avatarFileID == 0 {
		return consts.DefaultAvatarURL
	}
	if url, ok := avatarURLMap[avatarFileID]; ok {
		return url
	}
	return consts.DefaultAvatarURL
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
	type typeCount struct {
		Type string `json:"type"`
		Cnt  int    `json:"cnt"`
	}
	var counts []typeCount
	err = g.DB().Ctx(ctx).Model("notifications").
		Fields("type, COUNT(*) AS cnt").
		Where("user_id = ? AND is_read = 0 AND deleted_at IS NULL", in.UserId).
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

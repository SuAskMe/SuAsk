package notification

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"
	"suask/utility"
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

func (s *sNotification) Get(ctx context.Context, in model.GetNotificationsInput) (out model.GetNotificationsOutput, err error) {
	md := dao.Notifications.Ctx(ctx).Where(dao.Notifications.Columns().UserId, in.UserId).OrderAsc(dao.Notifications.Columns().IsRead).OrderDesc(dao.Notifications.Columns().CreatedAt)
	var newQuestion []*entity.Notifications
	var newAnswer []*entity.Notifications
	var newReply []*entity.Notifications
	var nqc int
	var nac int
	var nrc int
	// 分别获取new_question,new_answer,new_reply
	err = md.Where(dao.Notifications.Columns().Type, "new_question").ScanAndCount(&newQuestion, &nqc, false)
	if err != nil {
		return model.GetNotificationsOutput{}, err
	}
	err = md.Where(dao.Notifications.Columns().Type, "new_answer").ScanAndCount(&newAnswer, &nac, false)
	if err != nil {
		return model.GetNotificationsOutput{}, err
	}
	err = md.Where(dao.Notifications.Columns().Type, "new_reply").ScanAndCount(&newReply, &nrc, false)
	if err != nil {
		return model.GetNotificationsOutput{}, err
	}
	//fmt.Println("newQuestion:", newQuestion, nqc)
	//fmt.Println("newAnswer:", newAnswer, nac)
	//fmt.Println("newReply:", newReply, nrc)

	// 塞入 Set 里面
	var qIDList []int
	var aIDList []int
	for _, n := range newQuestion {
		qIDList = utility.AddUnique(qIDList, n.QuestionId)
	}
	for _, n := range newAnswer {
		qIDList = utility.AddUnique(qIDList, n.QuestionId)
		aIDList = utility.AddUnique(aIDList, n.AnswerId)
	}
	for _, n := range newReply {
		qIDList = utility.AddUnique(qIDList, n.QuestionId)
		aIDList = utility.AddUnique(aIDList, n.AnswerId)
		aIDList = utility.AddUnique(aIDList, n.ReplyToId)
	}
	//fmt.Println("qIDSet:", qIDList)
	//fmt.Println("aIDSet:", aIDList)
	var q []*entity.Questions
	var a []*entity.Answers

	// 根据 Set 查询
	err = dao.Questions.Ctx(ctx).WhereIn(dao.Questions.Columns().Id, qIDList).Scan(&q)
	if err != nil {
		return model.GetNotificationsOutput{}, err
	}
	err = dao.Answers.Ctx(ctx).WhereIn(dao.Answers.Columns().Id, aIDList).Scan(&a)
	if err != nil {
		return model.GetNotificationsOutput{}, err
	}
	// 查出的结果存入 Map
	qMap := make(map[int]*entity.Questions)
	aMap := make(map[int]*entity.Answers)

	for _, q := range q {
		qMap[q.Id] = q
	}
	for _, a := range a {
		aMap[a.Id] = a
	}

	// 获取用户信息
	var userIDs []int
	var u []entity.Users
	userMap := make(map[int]*entity.Users)
	for _, n := range qMap {
		userIDs = utility.AddUnique(userIDs, n.DstUserId)
	}
	for _, n := range aMap {
		userIDs = utility.AddUnique(userIDs, n.UserId)
	}
	err = dao.Users.Ctx(ctx).WhereIn(dao.Users.Columns().Id, userIDs).Scan(&u)
	if err != nil {
		return model.GetNotificationsOutput{}, err
	}
	for _, u := range u {
		userMap[u.Id] = &u
	}

	// 开辟 out 的内存
	out = model.GetNotificationsOutput{
		NewQuestion: make([]model.NotificationNewQuestion, nqc),
		NewAnswer:   make([]model.NotificationNewAnswer, nac),
		NewReply:    make([]model.NotificationNewReply, nrc),
	}
	// 塞入内容
	for i, n := range newQuestion {
		out.NewQuestion[i] = model.NotificationNewQuestion{
			NotificationBase: model.NotificationBase{
				Id:              n.Id,
				QuestionId:      n.QuestionId,
				QuestionTitle:   qMap[n.QuestionId].Title,
				QuestionContent: qMap[n.QuestionId].Contents,
				IsRead:          n.IsRead,
				CreatedAt:       n.CreatedAt.TimestampMilli(),
			},
			UserName: userMap[qMap[n.QuestionId].DstUserId].Name,
			UserId:   userMap[qMap[n.QuestionId].DstUserId].Id,
		}
	}
	for i, n := range newAnswer {
		out.NewAnswer[i] = model.NotificationNewAnswer{
			NotificationBase: model.NotificationBase{
				Id:              n.Id,
				QuestionId:      n.QuestionId,
				QuestionTitle:   qMap[n.QuestionId].Title,
				QuestionContent: qMap[n.QuestionId].Contents,
				IsRead:          n.IsRead,
				CreatedAt:       n.CreatedAt.TimestampMilli(),
			},
			AnswerId:       n.AnswerId,
			AnswerContent:  aMap[n.AnswerId].Contents,
			RespondentName: userMap[aMap[n.UserId].UserId].Nickname,
			RespondentId:   userMap[aMap[n.UserId].UserId].Id,
		}
	}
	for i, n := range newReply {
		out.NewReply[i] = model.NotificationNewReply{
			NotificationBase: model.NotificationBase{
				Id:              n.Id,
				QuestionId:      n.QuestionId,
				QuestionTitle:   qMap[n.QuestionId].Title,
				QuestionContent: qMap[n.QuestionId].Contents,
				IsRead:          n.IsRead,
				CreatedAt:       n.CreatedAt.TimestampMilli(),
			},
			AnswerId:       n.AnswerId,
			AnswerContent:  aMap[n.AnswerId].Contents,
			ReplyToId:      n.ReplyToId,
			ReplyToContent: aMap[n.ReplyToId].Contents,
			RespondentName: userMap[aMap[n.ReplyToId].UserId].Nickname,
			RespondentId:   userMap[aMap[n.ReplyToId].UserId].Id,
		}
	}
	return out, nil
}

func (s *sNotification) Update(ctx context.Context, in model.UpdateNotificationInput) (out model.UpdateNotificationOutput, err error) {
	out = model.UpdateNotificationOutput{}
	_, err = dao.Notifications.Ctx(ctx).WherePri(in.Id).Update(do.Notifications{IsRead: true})
	if err != nil {
		return model.UpdateNotificationOutput{}, err
	}
	out.IsRead = true
	out.Id = in.Id
	return out, nil
}

func (s *sNotification) Delete(ctx context.Context, in model.DeleteNotificationInput) (out model.DeleteNotificationOutput, err error) {
	_, err = dao.Notifications.Ctx(ctx).WherePri(in.Id).Delete()
	if err != nil {
		return model.DeleteNotificationOutput{}, err
	}
	out = model.DeleteNotificationOutput{}
	return out, nil
}

func (s *sNotification) NewNotificationCount(ctx context.Context, in model.NewNotificationCountInput) (out model.NewNotificationCountOutput, err error) {
	md := dao.Notifications.Ctx(ctx).Where(dao.Notifications.Columns().UserId, in.UserId).Where(dao.Notifications.Columns().IsRead, false)
	newQuestionCount, err := md.Where(dao.Notifications.Columns().Type, consts.NewQuestion).Count()
	if err != nil {
		return model.NewNotificationCountOutput{}, err
	}
	newAnswerCount, err := md.Where(dao.Notifications.Columns().Type, consts.NewAnswer).Count()
	if err != nil {
		return model.NewNotificationCountOutput{}, err
	}
	newReplyCount, err := md.Where(dao.Notifications.Columns().Type, consts.NewReply).Count()
	if err != nil {
		return model.NewNotificationCountOutput{}, err
	}
	out = model.NewNotificationCountOutput{
		NewQuestionCount: newQuestionCount,
		NewAnswerCount:   newAnswerCount,
		NewReplyCount:    newReplyCount,
	}
	return out, nil
}

func init() {
	service.RegisterNotification(New())
}

func New() *sNotification {
	return &sNotification{}
}

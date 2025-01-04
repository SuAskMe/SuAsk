package notification

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"
)

type sNotification struct{}

type pad struct{}

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
	qIDSet := make(map[int]pad)
	aIDSet := make(map[int]pad)
	for _, n := range newQuestion {
		if _, ok := qIDSet[n.QuestionId]; !ok {
			qIDSet[n.QuestionId] = pad{}
		}
	}
	for _, n := range newAnswer {
		if _, ok := qIDSet[n.QuestionId]; !ok {
			qIDSet[n.QuestionId] = pad{}
		}
		if _, ok := aIDSet[n.AnswerId]; !ok {
			aIDSet[n.AnswerId] = pad{}
		}
	}
	for _, n := range newReply {
		if _, ok := qIDSet[n.QuestionId]; !ok {
			aIDSet[n.AnswerId] = pad{}
		}
		if _, ok := aIDSet[n.AnswerId]; !ok {
			aIDSet[n.AnswerId] = pad{}
		}
		if _, ok := aIDSet[n.ReplyToId]; !ok {
			aIDSet[n.ReplyToId] = pad{}
		}
	}

	// 转换成 List
	qIDList := make([]int, 0, len(qIDSet))
	aIDList := make([]int, 0, len(aIDSet))
	for k := range qIDSet {
		qIDList = append(qIDList, k)
	}
	for k := range aIDSet {
		aIDList = append(aIDList, k)
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
	userIDs := make([]int, 0, len(qIDSet)+len(aIDSet))
	userIDSet := make(map[int]pad)
	var u []entity.Users
	userMap := make(map[int]*entity.Users)
	for _, n := range qMap {
		if _, ok := qMap[n.DstUserId]; !ok {
			userIDSet[n.DstUserId] = pad{}
		}
	}
	for _, n := range aMap {
		if _, ok := userIDSet[n.UserId]; !ok {
			userIDSet[n.UserId] = pad{}
		}
	}
	for k := range userIDSet {
		userIDs = append(userIDs, k)
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
		qid := n.QuestionId
		out.NewQuestion[i] = model.NotificationNewQuestion{
			NotificationBase: model.NotificationBase{
				Id:              n.Id,
				QuestionId:      qid,
				QuestionTitle:   qMap[qid].Title,
				QuestionContent: qMap[qid].Contents,
				IsRead:          n.IsRead,
				CreatedAt:       n.CreatedAt.TimestampMilli(),
			},
			// 只有老师会有这个提醒，统一设置成匿名用户
			UserName: consts.DefaultUserName,
			UserId:   consts.DefaultUserId,
		}
	}
	for i, n := range newAnswer {
		qid := n.QuestionId
		out.NewAnswer[i] = model.NotificationNewAnswer{
			NotificationBase: model.NotificationBase{
				Id:              n.Id,
				QuestionId:      qid,
				QuestionTitle:   qMap[qid].Title,
				QuestionContent: qMap[qid].Contents,
				IsRead:          n.IsRead,
				CreatedAt:       n.CreatedAt.TimestampMilli(),
			},
			AnswerId:      n.AnswerId,
			AnswerContent: aMap[n.AnswerId].Contents,
			// 自己问题的回复，可以看到所有回复者信息
			RespondentName: userMap[aMap[n.AnswerId].UserId].Nickname,
			RespondentId:   userMap[aMap[n.AnswerId].UserId].Id,
		}
	}
	for i, n := range newReply {
		qid := n.QuestionId
		respdName := userMap[aMap[n.ReplyToId].UserId].Nickname
		respdId := userMap[aMap[n.ReplyToId].UserId].Id
		// 在问老师的问题里，如果问题的目标是自己，回复者是匿名用户
		// 依赖其他的逻辑保证，只有回答拥有者才能获得该通知
		if qMap[qid].DstUserId == in.UserId {
			respdName = consts.DefaultUserName
			respdId = consts.DefaultUserId
		}
		out.NewReply[i] = model.NotificationNewReply{
			NotificationBase: model.NotificationBase{
				Id:              n.Id,
				QuestionId:      qid,
				QuestionTitle:   qMap[qid].Title,
				QuestionContent: qMap[qid].Contents,
				IsRead:          n.IsRead,
				CreatedAt:       n.CreatedAt.TimestampMilli(),
			},
			AnswerId:       n.AnswerId,
			AnswerContent:  aMap[n.AnswerId].Contents,
			ReplyToId:      n.ReplyToId,
			ReplyToContent: aMap[n.ReplyToId].Contents,
			RespondentName: respdName,
			RespondentId:   respdId,
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

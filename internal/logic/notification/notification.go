package notification

import (
	"context"
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
	md := dao.Notifications.Ctx(ctx).Where(dao.Notifications.Columns().UserId, in.UserId).OrderDesc(dao.Notifications.Columns().CreatedAt)
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
	// 开辟 out 的内存
	out = model.GetNotificationsOutput{
		NewQuestion: make([]model.Notification, nqc),
		NewAnswer:   make([]model.Notification, nac),
		NewReply:    make([]model.Notification, nrc),
	}
	// 塞入内容
	for i, n := range newQuestion {
		out.NewQuestion[i] = model.Notification{
			Id:              n.Id,
			QuestionId:      n.QuestionId,
			QuestionTitle:   qMap[n.QuestionId].Title,
			QuestionContent: qMap[n.QuestionId].Contents,
			IsRead:          n.IsRead,
			CreatedAt:       n.CreatedAt.TimestampMilli(),
		}
	}
	for i, n := range newAnswer {
		out.NewAnswer[i] = model.Notification{
			Id:              n.Id,
			QuestionId:      n.QuestionId,
			QuestionTitle:   qMap[n.QuestionId].Title,
			QuestionContent: qMap[n.QuestionId].Contents,
			AnswerId:        n.AnswerId,
			AnswerContent:   aMap[n.AnswerId].Contents,
			IsRead:          n.IsRead,
			CreatedAt:       n.CreatedAt.TimestampMilli(),
		}
	}
	for i, n := range newReply {
		out.NewReply[i] = model.Notification{
			Id:              n.Id,
			QuestionId:      n.QuestionId,
			QuestionTitle:   qMap[n.QuestionId].Title,
			QuestionContent: qMap[n.QuestionId].Contents,
			AnswerId:        n.AnswerId,
			AnswerContent:   aMap[n.AnswerId].Contents,
			IsRead:          n.IsRead,
			CreatedAt:       n.CreatedAt.TimestampMilli(),
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
	_, err = dao.Notifications.Ctx(ctx).Delete(in.Id)
	if err != nil {
		return model.DeleteNotificationOutput{}, err
	}
	out = model.DeleteNotificationOutput{}
	return out, nil
}

func init() {
	service.RegisterNotification(New())
}

func New() *sNotification {
	return &sNotification{}
}

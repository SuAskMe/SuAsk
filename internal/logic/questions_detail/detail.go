package questions_detail

import (
	"context"
	"fmt"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"
	"suask/module/validation"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sQuestionDetail struct{}

func AddResponseCnt(ctx context.Context, teacherId int) error {
	md := dao.Teachers.Ctx(ctx).Where(dao.Teachers.Columns().Id, teacherId)
	_, err := md.Increment(dao.Teachers.Columns().Responses, 1)
	if err != nil {
		return err
	}
	return nil
}

func (sQuestionDetail) GetQuestionBase(ctx context.Context, in *model.GetQuestionBaseInput) (*model.GetQuestionBaseOutput, error) {
	md := dao.Questions.Ctx(ctx).Unscoped().Where(dao.Questions.Columns().Id, in.QuestionId)
	var question entity.Questions
	err := md.Scan(&question)
	if err != nil {
		return nil, err
	}
	// 权限验证
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	// 不是老师本人时，要先过老师提问箱权限（私有/受保护/公开）
	if question.DstUserId != UserId {
		if err = validation.TeacherPerm(ctx, question.DstUserId); err != nil {
			return nil, err
		}
	}
	err = validation.QuestionPerm(ctx, &question)
	if err != nil {
		return nil, err
	}
	canReply := true
	err = validation.AnswerPerm(ctx, &question)
	if err != nil {
		canReply = false
	}
	// 获取问题详情
	var imgList []custom.Image
	var count int
	err = dao.Attachments.Ctx(ctx).Where(dao.Attachments.Columns().QuestionId, question.Id).ScanAndCount(&imgList, &count, false)
	if err != nil {
		return nil, err
	}
	imgIdList := extractFileIDs(imgList)
	isFavorite := false
	if in.UserId != 0 {
		one, err := dao.Favorites.Ctx(ctx).Where(dao.Favorites.Columns().QuestionId, in.QuestionId).Where(dao.Favorites.Columns().UserId, in.UserId).One()
		if !one.IsEmpty() {
			isFavorite = true
		}
		if err != nil {
			return nil, err
		}
	}
	output := model.GetQuestionBaseOutput{
		ID:         question.Id,
		Title:      question.Title,
		Content:    question.Contents,
		Views:      question.Views,
		CreatedAt:  question.CreatedAt.TimestampMilli(),
		CanReply:   canReply,
		ImageList:  imgIdList,
		IsFavorite: isFavorite,
		DstUserId:  question.DstUserId,
	}
	return &output, nil
}

func (sQuestionDetail) GetAnswers(ctx context.Context, in *model.GetAnswerDetailInput) (*model.GetAnswerDetailOutput, error) {
	md := dao.Answers.Ctx(ctx).Where(dao.Answers.Columns().QuestionId, in.QuestionId).Where("deleted_at IS NULL")
	var answers []entity.Answers
	err := md.Scan(&answers)
	if err != nil {
		return nil, err
	}
	// 获取回答详情

	answerList := make([]model.AnswerWithDetails, len(answers))
	IdList := make([]int, len(answers)) // 回答的ID列表
	IdMap := make(map[int]int)          // 回答的ID映射
	UserIdMap := make(map[int][]int)    // 用户ID所对应的回答ID列表

	for i, ans := range answers {
		IdList[i] = ans.Id
		IdMap[ans.Id] = i
		UserId := ans.UserId

		// "问大家"已下线，每条问题都必定有 DstUserId；这里只判断"不是老师本人的回答"，
		// 若是学生回答一律匿名化（显示默认用户头像）。
		if in.DstUserId != UserId {
			UserId = consts.DefaultUserId
		}

		if _, ok := UserIdMap[UserId]; !ok {
			UserIdMap[UserId] = []int{ans.Id}
		} else {
			UserIdMap[UserId] = append(UserIdMap[UserId], ans.Id)
		}

		answerList[i].Id = ans.Id
		answerList[i].UserId = UserId
		answerList[i].InReplyTo = ans.InReplyTo
		answerList[i].Contents = ans.Contents
		answerList[i].CreatedAt = ans.CreatedAt.TimestampMilli()
		answerList[i].Upvotes = ans.Upvotes
	}

	// 获取用户点赞信息
	//UserId := 2
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	md = dao.Upvotes.Ctx(ctx).WhereIn(dao.Upvotes.Columns().AnswerId, IdList).Where(dao.Upvotes.Columns().UserId, UserId)
	var upvotes []custom.MyUpvotes
	err = md.Scan(&upvotes)
	if err != nil {
		return nil, err
	}
	for _, upvote := range upvotes {
		answerList[IdMap[upvote.AnswerId]].IsUpvoted = true
	}
	// 获取回答者的信息
	UserIdList := make([]int, 0, len(answers)) // 用户ID列表
	for k := range UserIdMap {
		UserIdList = append(UserIdList, k)
	}
	md = dao.Users.Ctx(ctx).WhereIn(dao.Users.Columns().Id, UserIdList)
	var userInfo []custom.UserInfo // 用户信息
	err = md.Scan(&userInfo)
	if err != nil {
		return nil, err
	}
	AvatarMap := make(map[int][]int) // 头像ID对应的回答ID列表
	seenUserIDs := make(map[int]struct{}, len(userInfo))
	for _, info := range userInfo {
		seenUserIDs[info.UserId] = struct{}{}
		for _, v := range UserIdMap[info.UserId] {
			answerList[IdMap[v]].NickName = info.NickName
		}
		if info.Role == consts.TEACHER { // 如果是老师，则显示用户名
			for _, v := range UserIdMap[info.UserId] {
				answerList[IdMap[v]].TeacherName = info.Name
			}
		}
		if info.AvatarFileId == 0 {
			for _, v := range UserIdMap[info.UserId] {
				answerList[IdMap[v]].UserAvatar = consts.DefaultAvatarURL
			}
		} else {
			AvatarMap[info.AvatarFileId] = UserIdMap[info.UserId]
		}

	}
	for userId, answerIDs := range UserIdMap {
		if _, ok := seenUserIDs[userId]; ok {
			continue
		}
		nickname := "未知用户"
		for _, answerId := range answerIDs {
			answerList[IdMap[answerId]].NickName = nickname
			answerList[IdMap[answerId]].UserAvatar = consts.DefaultAvatarURL
		}
		AvatarMap[0] = append(AvatarMap[0], answerIDs...)
	}
	// 获取回答的图片
	md = dao.Attachments.Ctx(ctx).WhereIn(dao.Attachments.Columns().AnswerId, IdList)
	var imgList []custom.AnswerImage
	err = md.Scan(&imgList)
	if err != nil {
		return nil, err
	}
	ImgMap := make(map[int][]int)
	for _, img := range imgList {
		if _, ok := ImgMap[img.AnswerId]; !ok {
			ImgMap[img.AnswerId] = make([]int, 0, 8)
		}
		ImgMap[img.AnswerId] = append(ImgMap[img.AnswerId], img.FileID)
	}
	return &model.GetAnswerDetailOutput{
		IdMap:      IdMap,
		Answers:    answerList,
		AvatarsMap: AvatarMap,
		ImageMap:   ImgMap,
	}, nil
}

func (sQuestionDetail) AddQuestionView(ctx context.Context, in *model.AddViewInput) (*model.AddViewOutput, error) {
	db := g.DB()
	// ViewsLock.Lock()
	_, err := db.Exec(ctx, "UPDATE questions SET views = views + 1 WHERE id = ?", in.QuestionId)
	// ViewsLock.Unlock()
	if err != nil {
		return nil, err
	}
	return &model.AddViewOutput{}, nil
}

func (sQuestionDetail) AddAnswerUpvote(ctx context.Context, in *model.UpvoteInput) (*model.UpvoteOutput, error) {
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	if UserId == consts.DefaultUserId {
		return nil, fmt.Errorf("you are not allowed to upvote")
	}
	db := g.DB()
	// 事务包裹：保证 upvotes 表和 answers.upvotes 计数一致
	var result *model.UpvoteOutput
	err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cnt, err := tx.Model("upvotes").
			Where("answer_id = ? AND user_id = ?", in.AnswerId, UserId).
			Count()
		if err != nil {
			return err
		}
		if cnt > 0 {
			// 取消点赞
			_, err = tx.Model("upvotes").
				Where("answer_id = ? AND user_id = ?", in.AnswerId, UserId).
				Delete()
			if err != nil {
				return err
			}
			row, err := tx.Query(
				"UPDATE answers SET upvotes = upvotes - 1 WHERE id = ? RETURNING upvotes", in.AnswerId)
			if err != nil {
				return err
			}
			newCnt := 0
			if len(row) > 0 {
				newCnt = row[0]["upvotes"].Int()
			}
			result = &model.UpvoteOutput{IsUpvoted: false, Upvotes: newCnt}
		} else {
			// 点赞
			_, err = tx.Model("upvotes").Data(do.Upvotes{AnswerId: in.AnswerId, UserId: UserId}).Insert()
			if err != nil {
				return err
			}
			row, err := tx.Query(
				"UPDATE answers SET upvotes = upvotes + 1 WHERE id = ? RETURNING upvotes", in.AnswerId)
			if err != nil {
				return err
			}
			newCnt := 0
			if len(row) > 0 {
				newCnt = row[0]["upvotes"].Int()
			}
			result = &model.UpvoteOutput{IsUpvoted: true, Upvotes: newCnt}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sQuestionDetail) ReplyQuestion(ctx context.Context, in *model.AddAnswerInput) (*model.AddAnswerOutput, error) {
	md := dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().Id, in.QuestionId).Fields("id, src_user_id, dst_user_id, reply_cnt")
	var question entity.Questions
	err := md.Scan(&question)
	if err != nil {
		return nil, err
	}
	// 权限验证
	// "问大家"已下线：所有问题都指向某个老师。非老师本人时要过提问箱权限。
	if question.DstUserId != in.UserId {
		if err = validation.TeacherPerm(ctx, question.DstUserId); err != nil {
			return nil, err
		}
	}
	err = validation.AnswerPerm(ctx, &question)
	if err != nil {
		return nil, fmt.Errorf("you are not allowed to access this question")
	}

	out := &model.AddAnswerOutput{}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model("answers").Data(do.Answers{
			QuestionId: in.QuestionId,
			UserId:     in.UserId,
			Contents:   in.Content,
			InReplyTo:  in.InReplyTo,
		}).InsertAndGetId()
		if err != nil {
			return err
		}

		row, err := tx.Query("UPDATE questions SET reply_cnt = reply_cnt + 1 WHERE id = ? RETURNING reply_cnt", in.QuestionId)
		if err != nil {
			return err
		}
		replyCnt := 0
		if len(row) > 0 {
			replyCnt = row[0]["reply_cnt"].Int()
		}
		if replyCnt == 1 {
			if _, err = tx.Exec("UPDATE teachers SET responses = responses + 1 WHERE id = ?", question.DstUserId); err != nil {
				return err
			}
		}
		if replyCnt <= consts.MaxAvatarsPerQuestion {
			count, err := tx.Model("user_relation").
				Where("question_id = ? AND user_id = ?", in.QuestionId, in.UserId).
				Count()
			if err != nil {
				return err
			}
			if count == 0 {
				if _, err = tx.Model("user_relation").Data(do.UserRelation{QuestionId: in.QuestionId, UserId: in.UserId}).Insert(); err != nil {
					return err
				}
			}
		}

		out.Id = int(id)
		out.ReplyCnt = replyCnt
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (sQuestionDetail) AddReplyCnt(ctx context.Context, in *model.AddReplyCntInput) (*model.AddReplyCntOutput, error) {
	row, err := g.DB().Query(ctx,
		"UPDATE questions SET reply_cnt = reply_cnt + 1 WHERE id = ? RETURNING reply_cnt",
		in.QuestionId)
	if err != nil {
		return nil, err
	}
	cnt := 0
	if len(row) > 0 {
		cnt = row[0]["reply_cnt"].Int()
	}
	return &model.AddReplyCntOutput{ReplyCnt: cnt}, nil
}

func (s *sQuestionDetail) BuildRelation(ctx context.Context, in *model.BuildRelationInput) (*model.BuildRelationOutput, error) {
	// 保存关系
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	md := dao.UserRelation.Ctx(ctx)
	count, err := md.Where(dao.UserRelation.Columns().QuestionId, in.QuestionId).
		Where(dao.UserRelation.Columns().UserId, UserId).
		Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		_, err = md.Insert(do.UserRelation{
			QuestionId: in.QuestionId,
			UserId:     UserId,
		})
		if err != nil {
			return nil, err
		}
	}
	return &model.BuildRelationOutput{}, nil
}

// DeleteQuestion 软删问题。允许：提问者本人 / 目标老师 / admin。
func (sQuestionDetail) DeleteQuestion(ctx context.Context, questionId, userId int) error {
	var question entity.Questions
	err := dao.Questions.Ctx(ctx).
		Where(dao.Questions.Columns().Id, questionId).
		Where("deleted_at IS NULL").
		Scan(&question)
	if err != nil {
		return fmt.Errorf("问题不存在")
	}
	// 权限：本人 / 目标老师 / admin
	if question.SrcUserId != userId && question.DstUserId != userId {
		// 检查是否 admin
		var user entity.Users
		if err := dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, userId).Fields("role").Scan(&user); err != nil {
			return fmt.Errorf("无权删除")
		}
		if user.Role != consts.ADMIN {
			return fmt.Errorf("无权删除该问题")
		}
	}
	_, err = g.DB().Exec(ctx, "UPDATE questions SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", questionId)
	return err
}

// RestoreQuestion 恢复软删的问题。允许：提问者本人 / 目标老师 / admin。
func (sQuestionDetail) RestoreQuestion(ctx context.Context, questionId, userId int) error {
	var question entity.Questions
	err := dao.Questions.Ctx(ctx).Unscoped().
		Where(dao.Questions.Columns().Id, questionId).
		Scan(&question)
	if err != nil {
		return fmt.Errorf("问题不存在")
	}
	if question.DeletedAt == nil {
		return fmt.Errorf("该问题尚未删除")
	}
	// 权限：本人 / 目标老师 / admin
	if question.SrcUserId != userId && question.DstUserId != userId {
		// 检查是否 admin
		var user entity.Users
		if err := dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, userId).Fields("role").Scan(&user); err != nil {
			return fmt.Errorf("无权恢复")
		}
		if user.Role != consts.ADMIN {
			return fmt.Errorf("无权恢复该问题")
		}
	}
	_, err = g.DB().Exec(ctx, "UPDATE questions SET deleted_at = NULL WHERE id = ?", questionId)
	return err
}

// DeleteAnswer 软删回答。允许：回答者本人 / admin。
func (sQuestionDetail) DeleteAnswer(ctx context.Context, answerId, userId int) error {
	var answer entity.Answers
	err := dao.Answers.Ctx(ctx).
		Where(dao.Answers.Columns().Id, answerId).
		Where("deleted_at IS NULL").
		Scan(&answer)
	if err != nil {
		return fmt.Errorf("回答不存在")
	}
	if answer.UserId != userId {
		var user entity.Users
		if err := dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, userId).Fields("role").Scan(&user); err != nil {
			return fmt.Errorf("无权删除")
		}
		if user.Role != consts.ADMIN {
			return fmt.Errorf("无权删除该回答")
		}
	}
	_, err = g.DB().Exec(ctx, "UPDATE answers SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", answerId)
	return err
}

func init() {
	service.RegisterQuestionDetail(New())
}

func New() *sQuestionDetail {
	return &sQuestionDetail{}
}

package questions

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type sQuestionUtil struct{}

func (sQuestionUtil) GetImages(ctx context.Context, input *model.GetImagesInput) (*model.GetImagesOutput, error) {
	idList := input.QuestionIDs
	if len(idList) == 0 {
		return &model.GetImagesOutput{
			ImageMap: make(map[int][]int),
		}, nil
	}
	md := dao.Attachments.Ctx(ctx).WhereIn(dao.Attachments.Columns().QuestionId, idList)
	var Images []*custom.Image
	err := md.Scan(&Images)
	if err != nil {
		return nil, err
	}
	imageMap := make(map[int][]int)
	for _, img := range Images {
		if _, ok := imageMap[img.QuestionId]; !ok {
			imageMap[img.QuestionId] = make([]int, 0, 8)
		}
		imageMap[img.QuestionId] = append(imageMap[img.QuestionId], img.FileID)
	}
	output := model.GetImagesOutput{
		ImageMap: imageMap,
	}
	return &output, nil
}

func (s sQuestionUtil) GetQuestionListAssets(ctx context.Context, input *model.GetQuestionListAssetsInput) (*model.GetQuestionListAssetsOutput, error) {
	out := &model.GetQuestionListAssetsOutput{
		ImageURLMap:     make(map[int][]string),
		AnswerAvatarMap: make(map[int][]string),
		AnswerUserMap:   make(map[int][]model.AnswerUserSummary),
	}
	if input == nil || len(input.QuestionIDs) == 0 {
		return out, nil
	}

	imagesOutput, err := s.GetImages(ctx, &model.GetImagesInput{QuestionIDs: input.QuestionIDs})
	if err != nil {
		return nil, err
	}
	imageURLMap, err := BatchGetFileURLs(ctx, CollectFileIDs(imagesOutput.ImageMap, nil))
	if err != nil {
		return nil, err
	}
	for questionId, ids := range imagesOutput.ImageMap {
		urls := make([]string, 0, len(ids))
		for _, id := range ids {
			if url, ok := imageURLMap[id]; ok {
				urls = append(urls, url)
			}
		}
		out.ImageURLMap[questionId] = urls
	}

	answersOutput, err := s.GetAnswers(ctx, &model.GetAnswersInput{QuestionIDs: input.QuestionIDs})
	if err != nil {
		return nil, err
	}

	var (
		avatarsMap    map[int][]int
		answerUserMap map[int][]model.AnswerUserAsset
	)
	if answersOutput != nil {
		avatarsMap = answersOutput.AvatarsMap
		answerUserMap = answersOutput.AnswerUserMap
	}
	if avatarsMap == nil {
		avatarsMap = make(map[int][]int)
	}
	if answerUserMap == nil {
		answerUserMap = make(map[int][]model.AnswerUserAsset)
	}

	avatarFileIds := make([]int, 0)
	for questionId, ids := range avatarsMap {
		if input.DstUserIDMap != nil && input.DstUserIDMap[questionId] != 0 {
			continue
		}
		for _, id := range ids {
			if id != 0 {
				avatarFileIds = append(avatarFileIds, id)
			}
		}
	}

	teacherUserIDs := make([]int, 0)
	teacherUserIDSet := make(map[int]struct{})
	if input.DstUserIDMap != nil {
		for _, qid := range input.QuestionIDs {
			if dstUserID, ok := input.DstUserIDMap[qid]; ok && dstUserID != 0 {
				if _, exists := teacherUserIDSet[dstUserID]; !exists {
					teacherUserIDSet[dstUserID] = struct{}{}
					teacherUserIDs = append(teacherUserIDs, dstUserID)
				}
			}
		}
	}

	teacherUserMap := make(map[int]model.AnswerUserAsset)
	if len(teacherUserIDs) > 0 {
		var teacherUsers []struct {
			Id           int    `json:"id"`
			Name         string `json:"name"`
			Nickname     string `json:"nickname"`
			AvatarFileId int    `json:"avatar_file_id"`
		}
		err = dao.Users.Ctx(ctx).Fields("id", "name", "nickname", "avatar_file_id").WhereIn(dao.Users.Columns().Id, teacherUserIDs).Scan(&teacherUsers)
		if err != nil {
			return nil, err
		}
		for _, u := range teacherUsers {
			teacherUserMap[u.Id] = model.AnswerUserAsset{
				AvatarFileID: u.AvatarFileId,
				Nickname:     answerUserDisplayName(u.Nickname, u.Name),
			}
			if u.AvatarFileId != 0 {
				avatarFileIds = append(avatarFileIds, u.AvatarFileId)
			}
		}
	}

	avatarURLMap, err := BatchGetFileURLs(ctx, avatarFileIds)
	if err != nil {
		return nil, err
	}

	for _, questionId := range input.QuestionIDs {
		if dstUserID, ok := input.DstUserIDMap[questionId]; ok && dstUserID != 0 {
			teacherUser, ok := teacherUserMap[dstUserID]
			if !ok {
				out.AnswerAvatarMap[questionId] = []string{}
				out.AnswerUserMap[questionId] = []model.AnswerUserSummary{}
				continue
			}
			avatar := resolveAvatarURL(avatarURLMap, teacherUser.AvatarFileID)
			out.AnswerAvatarMap[questionId] = []string{avatar}
			out.AnswerUserMap[questionId] = []model.AnswerUserSummary{{
				Avatar:   avatar,
				Nickname: teacherUser.Nickname,
			}}
			continue
		}

		ids, ok := avatarsMap[questionId]
		if !ok || len(ids) == 0 {
			out.AnswerAvatarMap[questionId] = []string{}
			out.AnswerUserMap[questionId] = []model.AnswerUserSummary{}
			continue
		}
		urls := make([]string, 0, len(ids))
		for _, id := range ids {
			urls = append(urls, resolveAvatarURL(avatarURLMap, id))
		}
		out.AnswerAvatarMap[questionId] = urls

		answerUsers := answerUserMap[questionId]
		users := make([]model.AnswerUserSummary, 0, len(answerUsers))
		for _, user := range answerUsers {
			users = append(users, model.AnswerUserSummary{
				Avatar:   resolveAvatarURL(avatarURLMap, user.AvatarFileID),
				Nickname: user.Nickname,
			})
		}
		out.AnswerUserMap[questionId] = users
	}
	return out, nil
}

func resolveAvatarURL(avatarURLMap map[int]string, avatarFileID int) string {
	if avatarFileID == 0 {
		return consts.DefaultAvatarURL
	}
	if url, ok := avatarURLMap[avatarFileID]; ok {
		return url
	}
	return consts.DefaultAvatarURL
}

func answerUserDisplayName(nickname, name string) string {
	if nickname != "" {
		return nickname
	}
	return name
}

func (sQuestionUtil) Favorite(ctx context.Context, in *model.FavoriteInput) (out *model.FavoriteOutput, err error) {
	md := dao.Favorites.Ctx(ctx)
	//UserId := 1
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	if UserId == consts.DefaultUserId {
		return nil, gerror.New("user not login")
	}
	md = md.Where(dao.Favorites.Columns().UserId, UserId).Where(dao.Favorites.Columns().QuestionId, in.QuestionID)
	md = md.WhereNot(dao.Favorites.Columns().Package, consts.OnTop)
	cnt, err := md.Count()
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		_, err = md.Delete()
		if err != nil {
			return nil, err
		}
		return &model.FavoriteOutput{
			IsFavorite: false,
		}, nil
	} else {
		_, err = md.Insert(do.Favorites{
			UserId:     UserId,
			QuestionId: in.QuestionID,
			Package:    "default",
		})
		if err != nil {
			return nil, err
		}
		return &model.FavoriteOutput{
			IsFavorite: true,
		}, nil
	}
}

func (sQuestionUtil) GetQuestionSrcUserId(ctx context.Context, questionID int) (out int, err error) {
	var res *entity.Questions
	err = dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().Id, questionID).
		Fields("src_user_id").Scan(&res)
	if err != nil {
		return 0, err
	}
	return res.SrcUserId, nil
}

// AddQuestion 创建一个问题。
// "问大家"模块已下线，dst_user_id 必填；controller 层会在调用前做校验。
func (sQuestionUtil) AddQuestion(ctx context.Context, in *model.AddQuestionInput) (out *model.AddQuestionOutput, err error) {
	question := do.Questions{
		SrcUserId: in.SrcUserID,
		DstUserId: in.DstUserID,
		Title:     in.Title,
		Contents:  in.Content,
	}
	out = &model.AddQuestionOutput{}
	id, err := dao.Questions.Ctx(ctx).InsertAndGetId(question)
	if err != nil {
		if gstr.Contains(err.Error(), "FOREIGN KEY (`src_user_id`)") {
			return nil, gerror.New("找不到发送者")
		} else if gstr.Contains(err.Error(), "FOREIGN KEY (`dst_user_id`)") {
			return nil, gerror.New("找不到老师")
		}
		return nil, err
	}
	out.ID = int(id)
	return out, nil
}

// GetAnswers 批量拿"问题 → 回答者头像/昵称"的映射。之前在 sPublicQuestion 里实现，
// 供列表页拼接回答者小头像使用；和是否"公开提问"无关，因此挪到 util。
func (sQuestionUtil) GetAnswers(ctx context.Context, input *model.GetAnswersInput) (*model.GetAnswersOutput, error) {
	if len(input.QuestionIDs) == 0 {
		return &model.GetAnswersOutput{
			AvatarsMap:    map[int][]int{},
			AnswerUserMap: map[int][]model.AnswerUserAsset{},
		}, nil
	}
	db := g.DB()
	sqlStr := `
	SELECT question_id, avatar_file_id, nickname, name
	FROM (
		SELECT
			ur.question_id,
			u.avatar_file_id,
			u.nickname,
			u.name,
			ROW_NUMBER() OVER (
				PARTITION BY ur.question_id
				ORDER BY ur.user_id
			) AS rn
		FROM user_relation ur
		JOIN users u ON u.id = ur.user_id
		WHERE ur.question_id IN (?)
	) ranked
	WHERE rn <= ?
	ORDER BY question_id, rn`
	res, err := db.Query(ctx, sqlStr, input.QuestionIDs, consts.MaxAvatarsPerQuestion)
	if err != nil {
		return nil, err
	}
	avatarsMap := make(map[int][]int)
	answerUserMap := make(map[int][]model.AnswerUserAsset)
	for _, row := range res {
		id := row["question_id"].Int()
		if _, ok := avatarsMap[id]; !ok {
			avatarsMap[id] = make([]int, 0, consts.MaxAvatarsPerQuestion)
			answerUserMap[id] = make([]model.AnswerUserAsset, 0, consts.MaxAvatarsPerQuestion)
		}
		if len(avatarsMap[id]) >= consts.MaxAvatarsPerQuestion {
			continue
		}
		avatarFileID := row["avatar_file_id"].Int()
		avatarsMap[id] = append(avatarsMap[id], avatarFileID)
		answerUserMap[id] = append(answerUserMap[id], model.AnswerUserAsset{
			AvatarFileID: avatarFileID,
			Nickname:     answerUserDisplayName(row["nickname"].String(), row["name"].String()),
		})
	}
	return &model.GetAnswersOutput{
		AvatarsMap:    avatarsMap,
		AnswerUserMap: answerUserMap,
	}, nil
}

func init() {
	service.RegisterQuestionUtil(New())
}

func New() *sQuestionUtil {
	return &sQuestionUtil{}
}

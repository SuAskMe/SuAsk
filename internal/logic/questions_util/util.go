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
	}
	if input == nil || len(input.QuestionIDs) == 0 {
		return out, nil
	}

	imagesOutput, err := s.GetImages(ctx, &model.GetImagesInput{QuestionIDs: input.QuestionIDs})
	if err != nil {
		return nil, err
	}
	imageFileIds := make([]int, 0)
	for _, ids := range imagesOutput.ImageMap {
		imageFileIds = append(imageFileIds, ids...)
	}
	imageURLMap, err := getFileURLMap(ctx, imageFileIds)
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

	var avatarsMap map[int][]int
	if answersOutput != nil {
		avatarsMap = answersOutput.AvatarsMap
	} else {
		avatarsMap = make(map[int][]int)
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

	teacherAvatarFileIDMap := make(map[int]int)
	if len(teacherUserIDs) > 0 {
		var teacherUsers []struct {
			Id           int `json:"id"`
			AvatarFileId int `json:"avatar_file_id"`
		}
		err = dao.Users.Ctx(ctx).Fields("id", "avatar_file_id").WhereIn(dao.Users.Columns().Id, teacherUserIDs).Scan(&teacherUsers)
		if err != nil {
			return nil, err
		}
		for _, u := range teacherUsers {
			teacherAvatarFileIDMap[u.Id] = u.AvatarFileId
			if u.AvatarFileId != 0 {
				avatarFileIds = append(avatarFileIds, u.AvatarFileId)
			}
		}
	}

	avatarURLMap, err := getFileURLMap(ctx, avatarFileIds)
	if err != nil {
		return nil, err
	}

	for _, questionId := range input.QuestionIDs {
		if dstUserID, ok := input.DstUserIDMap[questionId]; ok && dstUserID != 0 {
			avatarFileID := teacherAvatarFileIDMap[dstUserID]
			if avatarFileID == 0 {
				out.AnswerAvatarMap[questionId] = []string{consts.DefaultAvatarURL}
			} else if url, ok := avatarURLMap[avatarFileID]; ok {
				out.AnswerAvatarMap[questionId] = []string{url}
			} else {
				out.AnswerAvatarMap[questionId] = []string{consts.DefaultAvatarURL}
			}
			continue
		}

		ids, ok := avatarsMap[questionId]
		if !ok || len(ids) == 0 {
			out.AnswerAvatarMap[questionId] = []string{}
			continue
		}
		urls := make([]string, 0, len(ids))
		for _, id := range ids {
			if id == 0 {
				urls = append(urls, consts.DefaultAvatarURL)
				continue
			}
			if url, ok := avatarURLMap[id]; ok {
				urls = append(urls, url)
			}
		}
		out.AnswerAvatarMap[questionId] = urls
	}
	return out, nil
}

func getFileURLMap(ctx context.Context, idList []int) (map[int]string, error) {
	urlMap := make(map[int]string)
	if len(idList) == 0 {
		return urlMap, nil
	}
	files, err := service.File().GetList(ctx, model.FileListGetInput{IdList: idList})
	if err != nil {
		return nil, err
	}
	for i, id := range files.FileId {
		if i < len(files.URL) {
			urlMap[id] = files.URL[i]
		}
	}
	return urlMap, nil
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

// GetAnswers 批量拿"问题 → 回答者头像"的映射。之前在 sPublicQuestion 里实现，
// 供列表页拼接回答者小头像使用；和是否"公开提问"无关，因此挪到 util。
func (sQuestionUtil) GetAnswers(ctx context.Context, input *model.GetAnswersInput) (*model.GetAnswersOutput, error) {
	if len(input.QuestionIDs) == 0 {
		return &model.GetAnswersOutput{AvatarsMap: map[int][]int{}}, nil
	}
	db := g.DB()
	sqlStr := `
	SELECT ur.question_id, u.avatar_file_id FROM user_relation ur, users u
	WHERE ur.question_id IN (?) AND u.id = ur.user_id;`
	res, err := db.Query(ctx, sqlStr, input.QuestionIDs)
	if err != nil {
		return nil, err
	}
	avatarsMap := make(map[int][]int)
	for _, row := range res {
		id := row["question_id"].Int()
		if _, ok := avatarsMap[id]; !ok {
			avatarsMap[id] = make([]int, 0, consts.MaxAvatarsPerQuestion)
		}
		if len(avatarsMap[id]) >= consts.MaxAvatarsPerQuestion {
			continue
		}
		avatarsMap[id] = append(avatarsMap[id], row["avatar_file_id"].Int())
	}
	return &model.GetAnswersOutput{AvatarsMap: avatarsMap}, nil
}

func init() {
	service.RegisterQuestionUtil(New())
}

func New() *sQuestionUtil {
	return &sQuestionUtil{}
}

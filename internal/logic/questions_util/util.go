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
			Package:    "默认收藏夹",
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
	err = dao.Questions.Ctx(ctx).WherePri(questionID).Fields("src_user_id").Scan(&res)
	if err != nil {
		return 0, err
	}
	return res.SrcUserId, nil
}

func init() {
	service.RegisterQuestionUtil(New())
}

func New() *sQuestionUtil {
	return &sQuestionUtil{}
}

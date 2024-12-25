package history

import (
	"context"
	"math"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"
)

type sHistoryOperation struct{}

// 查找历史提问模块需要的信息
func (sHistoryOperation) LoadHistoryInfo(ctx context.Context, in *model.GetHistoryInput) (out *model.GetHistoryOutput, err error) {

	// 得到用户提问的所有提问记录
	md := dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().SrcUserId, in.UserId)

	// 按照时间顺序降序，问题id降序
	md.Order("CreatedAt desc, Id desc")

	// 将排序之后的结果分页并得到对应的查询结果
	// historyQuestionAll := md.Page(in.Page, consts.NumOfQuestionsPerPage)
	historyQuestionAll := md.Page(in.Page, 10)
	var mqq []*model.MultiQueryQuestions

	err = historyQuestionAll.Scan(&mqq)
	if err != nil {
		return nil, err
	}
	// 初始化myhistoryquestion
	mhq := make([]model.MyHistoryQuestion, len(mqq))

	// key是问题的id  value是对应附件所有的url
	imageUrlMap := make(map[int][]string)

	for _, m := range mqq {
		for i := range m.Images {
			tempFileID := m.Images[i].FileID

			// 得到结构体sFile的  用tempFileID做出来的
			var tempFileGetInput = model.FileGetInput{
				Id: tempFileID, // 确保 tempFileID 是 int 类型
			}
			var tempFileGetOutput = model.FileGetOutput{}
			var tempUrl string
			tempFileGetOutput, err = service.IFile.Get(service.File(), ctx, tempFileGetInput) // 这步放到controller
			tempUrl = tempFileGetOutput.URL
			if err != nil {
				return nil, err
			}
			imageUrlMap[m.Id] = append(imageUrlMap[m.Id], tempUrl)
		}
	}

	// MyQuestionList切片完成
	for i := range mqq {
		mhq[i] = model.MyHistoryQuestion{
			Id:        mqq[i].Id,                         //int
			Title:     mqq[i].Title,                      //string
			Contents:  mqq[i].Contents,                   //string
			CreatedAt: mqq[i].CreatedAt.TimestampMilli(), //int64
			Views:     mqq[i].Views,                      //int
			ImageURLs: imageUrlMap[mqq[i].Id],            //[]string
		}
	}

	limit := consts.NumOfQuestionsPerPage
	limit = 10
	total, err := dao.Questions.Ctx(ctx).Where(do.Questions{SrcUserId: in.UserId}).Count()
	if err != nil {
		return nil, err
	}
	pageNum := math.Ceil(float64(total) / float64(limit))
	remain := int(pageNum) - in.Page

	ultimate_out := model.GetHistoryOutput{
		Question:   mhq,
		Total:      total,
		Size:       limit,
		PageNum:    int(pageNum),
		RemainPage: remain,
	}

	return &ultimate_out, nil
}

func init() {
	service.RegisterHistoryOperation(New())
}

func New() *sHistoryOperation {
	return &sHistoryOperation{}
}

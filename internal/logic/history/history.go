package history

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/service"
)

type HistoryOperation struct{}

// 查找历史提问模块需要的信息
func (h HistoryOperation) LoadHistoryInfo(ctx context.Context, in *model.GetHistoryInput) (out *model.GetHistoryOutput, err error) {

	// 得到用户提问的所有提问记录
	md := dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().SrcUserId, in.UserId)

	// 按照时间顺序降序，问题id降序
	md.Order("CreatedAt desc").Order("Id desc")

	// 将排序之后的结果分页并得到对应的查询结果
	historyQuestionAll := md.Page(in.Page, consts.NumOfQuestionsPerPage)

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
			Id:        mqq[i].Id,              //int
			Contents:  mqq[i].Contents,        //string
			CreatedAt: mqq[i].CreatedAt,       //*gtime.Time
			Views:     mqq[i].Views,           //int
			ImageURLs: imageUrlMap[mqq[i].Id], //[]string
		}
	}
	remain, err := md.Count()
	if err != nil {
		return nil, err
	}
	remainNum := remain - consts.NumOfQuestionsPerPage*in.Page
	remain = remainNum / consts.NumOfQuestionsPerPage
	if remainNum%consts.NumOfQuestionsPerPage > 0 {
		remain += 1
	}
	ultimate_out := model.GetHistoryOutput{
		Question:   mhq,
		RemainPage: remain,
	}

	return &ultimate_out, nil
}

// // 辅助函数仅通过文件id得到图片url
// func (h HistoryOperation) GetUrlUseFileId(ctx context.Context, id int) (out string, err error) {
// 	file := entity.Files{}
// 	err = dao.Files.Ctx(ctx).Where(dao.Files.Columns().Id, id).Scan(&file)
// 	if err != nil {
// 		return "", fmt.Errorf("无法查询文件记录：%w", err)
// 	}
// 	URL, err := files.GetURL(file.Hash, file.Name)
// 	if err != nil {
// 		return "", fmt.Errorf("生成文件 URL 失败：%w", err)
// 	}

// 	return URL, nil
// }

func init() {
	service.RegisterHistory(New())
}

func New() *HistoryOperation {
	return &HistoryOperation{}
}

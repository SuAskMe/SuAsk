package history

import (
	"context"
	"fmt"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/entity"
	"suask/utility/files"
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

	var mqq *[]model.MultiQueryQuestions

	err = historyQuestionAll.Scan(&mqq)

	return
}

func New() *HistoryOperation {
	return &HistoryOperation{}
}

func (h HistoryOperation) GetUrlUseFileId(ctx context.Context, id string) (out string, err error) {
	file := entity.Files{}
	err = dao.Files.Ctx(ctx).Where(dao.Files.Columns().Id, id).Scan(&file)
	if err != nil {
		return "", fmt.Errorf("无法查询文件记录：%w", err)
	}
	URL, err := files.GetURL(file.Hash, file.Name)
	if err != nil {
		return "", fmt.Errorf("生成文件 URL 失败：%w", err)
	}

	return URL, nil
}

package errcode

var (
	ErrorUpdateRepositoriesFail = NewError(20010001, "更新helm repo列表失败")
	ErrorListRepoChartsFail     = NewError(20010002, "查询helm repo列表失败")

	ErrorShowChartInfoFail       = NewError(20020001, "查询chart信息失败")
	ErrorDListUploadedChartsFail = NewError(20020002, "获取上传的chart列表失败")

	ErrorShowReleaseInfoFail = NewError(20030001, "获取release信息失败")
	ErrorInstallReleaseFail  = NewError(20030002, "部署release失败")
	ErrorCreateArticleFail   = NewError(20020003, "创建文章失败")
	ErrorUpdateArticleFail   = NewError(20020004, "更新文章失败")
	ErrorDeleteArticleFail   = NewError(20020005, "删除文章失败")

	ErrorUploadFileFail = NewError(20040001, "上传文件失败")
)

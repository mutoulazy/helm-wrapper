package errcode

var (
	ErrorUpdateRepositoriesFail = NewError(20010001, "更新helm repo列表失败")
	ErrorListRepoChartsFail     = NewError(20010002, "查询helm repo列表失败")

	ErrorShowChartInfoFail       = NewError(20020001, "查询chart信息失败")
	ErrorDListUploadedChartsFail = NewError(20020002, "获取上传的chart列表失败")

	ErrorShowReleaseInfoFail      = NewError(20030001, "获取release信息失败")
	ErrorInstallReleaseFail       = NewError(20030002, "部署release失败")
	ErrorUninstallReleaseFail     = NewError(20030003, "卸载release失败")
	ErrorRollbackReleaseFail      = NewError(20030004, "回滚release失败")
	ErrorUpgradeReleaseFail       = NewError(20030005, "更新release失败")
	ErrorListReleasesFail         = NewError(20030006, "获取release列表信息失败")
	ErrorGetReleaseStatusFail     = NewError(20030007, "获取release状态信息失败")
	ErrorListReleaseHistoriesFail = NewError(20030008, "获取release历史版本信息失败")

	ErrorUploadFileFail = NewError(20040001, "上传文件失败")
)

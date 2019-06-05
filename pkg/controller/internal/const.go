package internal

const (
	LabelNameKey      = "dlkit.njuics.cn/name"
	LabelMemberKey    = "dlkit.njuics.cn/member"
	LabelTypeKey      = "dlkit.njuics.cn/type"
	LabelWorkspaceKey = "dlkit.njuics.cn/workspace"
	LabelUserKey      = "dlkit.njuics.cn/user"

	LabelJobPod = "group_name"
	ValueJobPod = "kubeflow.org"
)

const (
	NvidiaGpuResource = "nvidia.com/gpu"
	NvidiaGpuEnv      = "CUDA_VISIBLE_DEVICES"
)

const (
	WorkspaceVolumeMountPath = "/workspace"
	DatasetVolumeMountPath   = "/workspace"
)

package sparkcluster

// Const Variable
const (
	Master           = "master"
	MasterPvc        = "namenode-pvc"
	MasterImage      = "registry.njuics.cn/wdongyu/spark_master_on_kube:1.0.2"
	SlaveImage       = "registry.njuics.cn/wdongyu/spark_slave_on_kube:1.0.2"
	Slave            = "slave"
	SlavePvc         = "datanode-pvc"
	UIService        = "hadoop-ui-service"
	ShareServer      = "114.212.189.141"
	StorageClassName = "cephfs"
)

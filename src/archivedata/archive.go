package archivedata

type IArchive interface {
	
	PushData(data string)
	Exit()
	
}

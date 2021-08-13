package worlddata

func Init() {
	svrts_init()
	seq_init()

	g_resetter.Open()
}

func Close() {
	seq_save()
}

package system

func SecurePanic(err error) {
	if err != nil {
		panic(err)
	}
}

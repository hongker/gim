package errgroup

func Do(fn ...func() error) error {
	for _, f := range fn {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

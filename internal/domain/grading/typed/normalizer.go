package typed

type normalizer func(string) string

func applyNormalizers(s string, normalizers ...normalizer) string {
	for i := range normalizers {
		if normalizers[i] != nil {
			s = normalizers[i](s)
		}
	}

	return s
}

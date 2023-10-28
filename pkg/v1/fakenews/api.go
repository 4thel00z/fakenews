package fakenews

const (
	StructTagName = "fakenews"
)

func Fake[T any](val T, rules ...Rule[T]) (T, error) {
	for _, rule := range rules {
		val, err := rule.Run(&val)
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

package fakenews

const (
	StructTagName = "fakenews"
)

func Fake(val interface{}, rules ...Rule) (interface{}, error) {
	for _, rule := range rules {
		val, err := rule.Run(val)
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

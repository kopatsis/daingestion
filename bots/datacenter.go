package bots

func FromDataCenter(asn uint, list map[uint]string) string {
	if center, ok := list[asn]; ok {
		return center
	}
	return ""
}

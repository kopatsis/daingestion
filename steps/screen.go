package steps

type ScreenBuckets struct {
	InnerWidthBucket   string
	InnerHeightBucket  string
	ScreenWidthBucket  string
	ScreenHeightBucket string
}

func bucketWidth(w int) string {
	switch {
	case w < 480:
		return "ultra_narrow"
	case w < 768:
		return "narrow"
	case w < 1280:
		return "standard"
	case w < 1920:
		return "wide"
	default:
		return "ultra_wide"
	}
}

func bucketHeight(h int) string {
	switch {
	case h < 600:
		return "short"
	case h < 900:
		return "medium"
	case h < 1200:
		return "tall"
	default:
		return "extra_tall"
	}
}

func BucketScreenSizes(innerW, innerH, screenW, screenH int) ScreenBuckets {
	return ScreenBuckets{
		InnerWidthBucket:   bucketWidth(innerW),
		InnerHeightBucket:  bucketHeight(innerH),
		ScreenWidthBucket:  bucketWidth(screenW),
		ScreenHeightBucket: bucketHeight(screenH),
	}
}

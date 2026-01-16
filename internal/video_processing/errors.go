package video_processing

type ErrFFmpegNotFound struct {
	Hint string
}

func (e ErrFFmpegNotFound) Error() string {
	if e.Hint == "" {
		return "ffmpeg not found"
	}
	return e.Hint
}

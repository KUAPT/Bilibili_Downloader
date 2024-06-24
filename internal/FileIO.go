package internal

/*
import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// 将 Response 结构体转换为 Markdown 格式
func responseToMarkdown(response *Response) string {
	var sb strings.Builder

	// 写入数据列表
	for _, video := range response.Data.List {
		sb.WriteString(fmt.Sprintf("# %s\n", video.Title))
		sb.WriteString(fmt.Sprintf("![%s](%s)\n", video.Title, video.Pic))
		sb.WriteString(fmt.Sprintf("- **发布地点**: %s\n", video.PubLocation))
		sb.WriteString(fmt.Sprintf("- **发布日期**: %s\n", time.Unix(video.Pubdate, 0).Format("2006-01-02")))
		sb.WriteString(fmt.Sprintf("- **分类**: %s\n", video.Tname))
		sb.WriteString(fmt.Sprintf("- **简介**: %s\n", video.Desc))
		sb.WriteString(fmt.Sprintf("- **作者**: %s\n", video.Owner.Name))
		sb.WriteString(fmt.Sprintf("- **观看数**: %d\n", video.Stat.View))
		sb.WriteString(fmt.Sprintf("- **弹幕数**: %d\n", video.Stat.Danmaku))
		sb.WriteString(fmt.Sprintf("- **收藏数**: %d\n", video.Stat.Favorite))
		sb.WriteString(fmt.Sprintf("- **硬币数**: %d\n", video.Stat.Coin))
		sb.WriteString(fmt.Sprintf("- **分享数**: %d\n", video.Stat.Share))
		sb.WriteString(fmt.Sprintf("- **点赞数**: %d\n", video.Stat.Like))
		sb.WriteString(fmt.Sprintf("- **推荐理由**: %s\n", video.RcmdReason.Content))
		sb.WriteString(fmt.Sprintf("- [视频链接](%s)\n", video.ShortLinkV2))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func FileStore(response *Response) {
	// 将 Response 转换为 Markdown 格式
	markdown := responseToMarkdown(response)
	// 将 Markdown 数据保存为文件
	err := os.WriteFile("当前热榜.md", []byte(markdown), 0644)
	if err != nil {
		log.Fatalf("保存文件错误: %v\n", err)
		return
	}
}
*/

package toolkit

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func CheckAndCreateCacheDir() error {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	cacheDirPath := filepath.Join(currentDir, "download_cache")

	// 检查目录是否存在
	if _, err := os.Stat(cacheDirPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err = os.MkdirAll(cacheDirPath, 0755)
		if err != nil {
			return fmt.Errorf("临时下载目录创建失败: %v", err)
		}
		fmt.Println("成功创建缓存目录:", cacheDirPath)
		log.Println("建立缓存目录正常:", cacheDirPath)
	} else if err != nil {
		// 如果 os.Stat 返回了错误，但不是 os.IsNotExist
		return fmt.Errorf("检查目录时发生错误: %v", err)
	} else {
		// 目录存在
		fmt.Println("缓存目录已经存在，继续使用:", cacheDirPath)
		log.Println("缓存目录已经存在，继续使用:", cacheDirPath)
	}

	return nil
}

func RemoveCacheDir() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	cacheDirPath := filepath.Join(currentDir, "download_cache")

	if err := os.RemoveAll(cacheDirPath); err != nil {
		return fmt.Errorf("移除临时cache目录失败: %v", err)
	}
	log.Println("成功移除cache目录:", cacheDirPath)
	return nil
}

func ClearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Println("清屏命令执行失败：", err)
		}
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Println("清屏命令执行失败：", err)
		}
	default:
		log.Println("无法清屏，不支持的平台！")
	}
}

func CheckAndCreateDir(dir string) error {
	configDir := dir

	// 检查目录是否存在
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.Mkdir(configDir, 0755)
		if err != nil {
			return fmt.Errorf("无法创建目录 %s: %v", configDir, err)
		}
		fmt.Println("目录已创建:", configDir)
	} else if err != nil {
		// 其他错误
		return fmt.Errorf("检查目录 %s 时出错: %v", configDir, err)
	} else {
		// 目录已存在
		fmt.Println("目录已存在:", configDir)
	}
	return nil
}

// CheckAndCleanFileName 检查文件名是否包含不允许的字符，并进行清理
func CheckAndCleanFileName(fileName string) string {
	disallowedChars := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	// 检查文件名中的每个字符
	for _, char := range disallowedChars {
		if strings.Contains(fileName, char) {
			// 替换不允许的字符为下划线
			fileName = strings.ReplaceAll(fileName, char, "_")
		}
	}
	return fileName
}

func DownloadAndTrackProgress(src io.Reader, dst io.Writer, bar *pb.ProgressBar) error {
	reader := bar.NewProxyReader(io.TeeReader(src, dst))
	_, err := io.Copy(io.Discard, reader)
	return err
}

func CatchAndCheckBVid() string {
	//正则对BV号进行基本检查
	BVCheak := regexp.MustCompile(`^BV[1-9A-HJ-NP-Za-km-z]{10}$`)
	var BVid string
	for {
		fmt.Printf("请输入需要下载视频的BV号：")
		if _, err := fmt.Scanln(&BVid); err != nil {
			ClearScreen()
			fmt.Println("输入读取错误，请重试！")
			log.Println("读取输入错误：", err)
			continue
		}
		if BVCheak.MatchString(BVid) {
			break
		} else {
			ClearScreen()
			fmt.Println("BV号格式错误，请检查格式后重试！")
		}
	}
	return BVid
}

func ObtainUserResolutionSelection(downloadInfoResponse *DownloadInfoResponse) (int, int, string) {
	definition := make(map[int]string, 10)
	for i := 0; i < len(downloadInfoResponse.Data.AcceptDescription); i++ {
		definition[downloadInfoResponse.Data.AcceptQuality[i]] = downloadInfoResponse.Data.AcceptDescription[i]
	}

	effectiveDefinition := make([]int, 0, 10)
	effectiveDefinitionMap := make(map[int]bool)
	for i := range downloadInfoResponse.Data.Dash.Video {
		id := downloadInfoResponse.Data.Dash.Video[i].ID
		if !effectiveDefinitionMap[id] {
			effectiveDefinition = append(effectiveDefinition, id)
			effectiveDefinitionMap[id] = true
		}
	}

	var choose int
	for true {
		fmt.Println("\n请选择想要下载的分辨率：(ps:此处仅显示当前登录账号有权获取的所有分辨率选项)")
		for i := range effectiveDefinition {
			fmt.Println(i+1, definition[effectiveDefinition[i]])
		}
		fmt.Printf("请输入分辨率前的序号(单个数字)：")
		if _, err := fmt.Scanln(&choose); err != nil {
			ClearScreen()
			log.Println("读取输入发生错误")
			fmt.Println("读取输入发生错误,请检查输入格式后重试，若问题依旧，请携带日志log文件向开发者反馈！")
			continue
		}
		if choose < 1 || choose > len(effectiveDefinition) {
			ClearScreen()
			fmt.Println("输入错误，请检查输入后重试！")
			continue
		}
		choose -= 1
		break
	}
	var videoIndex int
	for i := range downloadInfoResponse.Data.Dash.Video {
		if downloadInfoResponse.Data.Dash.Video[i].ID == effectiveDefinition[choose] {
			videoIndex = i
		}
	}

	resolutions := map[int]string{
		6:   "240P",
		16:  "360P",
		32:  "480P",
		64:  "720P",
		74:  "720P60",
		80:  "1080P",
		112: "1080P+",
		116: "1080P60",
		120: "4K",
		125: "HDR",
		126: "杜比视界",
		127: "8K超高清",
	}
	videoCode := effectiveDefinition[choose]
	resolutionDescription := resolutions[videoCode]

	return videoIndex, videoCode, resolutionDescription
}

func IsWantToContinueDownloading() bool {
	reader := bufio.NewReader(os.Stdin)
	var isContinue rune
	for true {
		fmt.Printf("是否继续下载其他视频？(y/n):")
		if _, err := fmt.Scanf("%c", &isContinue); err != nil {
			_, _ = reader.ReadString('\n')
			log.Println("读取输入发生错误", err)
			fmt.Println("读取输入发生错误,请检查输入格式后重试，若问题依旧，请携带日志log文件向开发者反馈！")
			continue
		}
		if isContinue != 'y' && isContinue != 'Y' && isContinue != 'n' && isContinue != 'N' {
			_, _ = reader.ReadString('\n')
			ClearScreen()
			fmt.Println("输入错误，请检查输入(y/n)[不区分大小写]！")
			continue
		}
		_, _ = reader.ReadString('\n')
		break
	}
	if isContinue == 'y' || isContinue == 'Y' {
		ClearScreen()
		return true
	} else {
		return false
	}
}

package recoverTool

import (
	"fmt"
	"html"
	"strings"
)

func FormatStack(stack string) string {
	lines := strings.Split(stack, "\n")
	var formatted []string
	for i := 0; i < len(lines); i += 2 {
		if i+1 < len(lines) {
			file := strings.TrimSpace(lines[i])
			function := strings.TrimSpace(lines[i+1])
			if strings.HasPrefix(file, "goroutine ") {
				formatted = append(formatted, file)
			} else {
				formatted = append(formatted, fmt.Sprintf("%s\n    at %s", function, file))
			}
		}
	}
	return strings.Join(formatted, "\n")
}

func formatHtml(service, errMsg string, stack string) string {
	// 创建一个简单的 HTML 结构
	htmlContent := "<html><head>\n"
	htmlContent += "<style>\n"
	htmlContent += "body { font-family: monospace; font-size: 16px; margin: 20px; }\n"
	htmlContent += ".error { color: #ff0000; font-weight: bold; margin-bottom: 10px; }\n"
	htmlContent += ".stack { white-space: pre-wrap; }\n"
	htmlContent += "</style>\n"
	htmlContent += "</head><body>\n"

	// 添加服务名称和错误信息
	htmlContent += fmt.Sprintf("<div class='error'>service: %s</div>\n", service)
	htmlContent += fmt.Sprintf("<div class='error'>panic: %s</div>\n", html.EscapeString(errMsg))

	// 添加堆栈信息，保持原始格式但转义特殊字符
	htmlContent += "<div class='stack'>\n"
	htmlContent += html.EscapeString(stack)
	htmlContent += "</div>\n"

	htmlContent += "</body></html>"

	return htmlContent
}

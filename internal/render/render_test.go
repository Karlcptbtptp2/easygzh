package render

import (
	"strings"
	"testing"
)

const testMD = `# 欢迎使用 easyGZH

这是一段**正文**，含[外链](https://example.com)。

## 二级标题

- 项一
- 项二
`

const testCSS = `
#easygzh-root { font-size: 15px; color: #333; }
#easygzh-root h1 { color: #1a73e8; }
#easygzh-root strong { color: #1a73e8; }
`

func TestRender_FullPipeline(t *testing.T) {
	out, err := Render(testMD, PipelineOptions{ThemeCSS: testCSS, LinkFootnotes: true})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("Render returned empty output")
	}
	t.Logf("output:\n%s", out)
	if !strings.Contains(out, `<h1`) {
		t.Error("output missing <h1")
	}
	if !strings.Contains(out, "1a73e8") {
		t.Error("theme color not inlined")
	}
	if !strings.Contains(out, "[1]") {
		t.Error("footnote [1] not present")
	}
}

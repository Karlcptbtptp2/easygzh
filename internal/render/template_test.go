package render

import (
	"strings"
	"testing"
)

const testTemplateHTML = `<!-- template:
  name: test-template
  description: A test template
-->
<style data-template-css>
.test-brand { font-size: 13px; color: #7A7A7A; }
.test-cta { background: #E8923C; color: #fff; }
</style>

<!-- slot: brand -->
<section class="test-brand">{{BRAND_LABEL}} {{BRAND_DATE}}</section>

<!-- slot: body -->

<!-- slot: cta -->
<section class="test-cta">{{CTA_TEXT}}</section>
`

func TestParseTemplate_ExtractsMeta(t *testing.T) {
	tmpl, err := ParseTemplate(testTemplateHTML)
	if err != nil {
		t.Fatalf("ParseTemplate error: %v", err)
	}
	if tmpl.Name != "test-template" {
		t.Errorf("Name = %q, want test-template", tmpl.Name)
	}
	if tmpl.Description != "A test template" {
		t.Errorf("Description = %q, want 'A test template'", tmpl.Description)
	}
}

func TestParseTemplate_ExtractsCSS(t *testing.T) {
	tmpl, _ := ParseTemplate(testTemplateHTML)
	if !strings.Contains(tmpl.CSS, ".test-brand") {
		t.Errorf("CSS should contain .test-brand, got: %s", tmpl.CSS)
	}
	if !strings.Contains(tmpl.CSS, ".test-cta") {
		t.Errorf("CSS should contain .test-cta, got: %s", tmpl.CSS)
	}
}

func TestParseTemplate_RemovesStyleBlock(t *testing.T) {
	tmpl, _ := ParseTemplate(testTemplateHTML)
	if strings.Contains(tmpl.HTML, "data-template-css") {
		t.Error("HTML should not contain the style block")
	}
	if strings.Contains(tmpl.HTML, ".test-brand") {
		t.Error("HTML should not contain CSS rules")
	}
}

func TestApplyStructureTemplate_BodyInjection(t *testing.T) {
	tmpl, _ := ParseTemplate(testTemplateHTML)
	bodyHTML := `<h1>Test Title</h1><p>Test paragraph</p>`

	result, err := ApplyStructureTemplate(bodyHTML, tmpl, SlotData{})
	if err != nil {
		t.Fatalf("ApplyStructureTemplate error: %v", err)
	}

	// Body content should be present.
	if !strings.Contains(result, "Test paragraph") {
		t.Error("Result should contain the body HTML")
	}

	// The body slot comment should be replaced.
	if strings.Contains(result, "<!-- slot: body -->") {
		t.Error("Body slot comment should be replaced with actual content")
	}
}

func TestApplyStructureTemplate_PlaceholderFilling(t *testing.T) {
	tmpl, _ := ParseTemplate(testTemplateHTML)

	result, _ := ApplyStructureTemplate("<p>body</p>", tmpl, SlotData{
		BrandLabel:  "( Test Brand )",
		CTAText:     "Join now",
		BrandFooter: "TEST CORP",
	})

	if !strings.Contains(result, "( Test Brand )") {
		t.Error("Brand label should be filled in")
	}
	if !strings.Contains(result, "Join now") {
		t.Error("CTA text should be filled in")
	}
	// {{BRAND_DATE}} should auto-fill with current month.
	if strings.Contains(result, "{{BRAND_DATE}}") {
		t.Error("Brand date should be auto-filled, not left as placeholder")
	}
}

func TestApplyStructureTemplate_TitleFromH1(t *testing.T) {
	tmpl, _ := ParseTemplate(testTemplateHTML)

	result, _ := ApplyStructureTemplate("<h1>My Article Title</h1><p>body</p>", tmpl, SlotData{})

	if !strings.Contains(result, "My Article Title") {
		t.Error("Title should be extracted from h1 and filled into {{TITLE}}")
	}
}

func TestApplyStructureTemplate_NilTemplate(t *testing.T) {
	bodyHTML := "<p>just a paragraph</p>"
	result, err := ApplyStructureTemplate(bodyHTML, nil, SlotData{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != bodyHTML {
		t.Error("Nil template should return input unchanged (backward compatibility)")
	}
}

func TestApplyStructureTemplate_EmptyHTML(t *testing.T) {
	tmpl := &StructureTemplate{Name: "empty", HTML: ""}
	result, err := ApplyStructureTemplate("<p>test</p>", tmpl, SlotData{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "<p>test</p>" {
		t.Error("Empty template HTML should return input unchanged")
	}
}

func TestApplyStructureTemplate_NoUnfilledPlaceholders(t *testing.T) {
	tmpl, _ := ParseTemplate(testTemplateHTML)

	result, _ := ApplyStructureTemplate("<h1>Title</h1><p>body</p>", tmpl, SlotData{})

	// All known placeholders should be filled.
	knownPlaceholders := []string{
		"{{TITLE}}", "{{BRAND_LABEL}}", "{{BRAND_DATE}}",
		"{{SUBTITLE}}", "{{CTA_TEXT}}", "{{BRAND_FOOTER}}",
	}
	for _, ph := range knownPlaceholders {
		if strings.Contains(result, ph) {
			t.Errorf("Placeholder %s should have been filled", ph)
		}
	}
}

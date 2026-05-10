package files

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/net/ghttp"
)

// makeUploadFile 构造一个可交给 ghttp 的上传文件壳子：走一遍真实 multipart 解析，
// 这样 fh.Open() / fh.Size / fh.Filename 都是标准库提供的，与生产一致。
func makeUploadFile(t *testing.T, filename string, payload []byte) *ghttp.UploadFile {
	t.Helper()

	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	header := make(textproto.MIMEHeader)
	header.Set(
		"Content-Disposition",
		`form-data; name="file"; filename="`+filename+`"`,
	)
	header.Set("Content-Type", "application/octet-stream")
	part, err := w.CreatePart(header)
	if err != nil {
		t.Fatalf("create part: %v", err)
	}
	if _, err := part.Write(payload); err != nil {
		t.Fatalf("write part: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	r := multipart.NewReader(&body, w.Boundary())
	form, err := r.ReadForm(int64(body.Len()) + 1024)
	if err != nil {
		t.Fatalf("read form: %v", err)
	}
	files := form.File["file"]
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	return &ghttp.UploadFile{FileHeader: files[0]}
}

// 真实 PNG 魔术字节前缀（8 字节）+ 几个垃圾字节，能通过 DetectContentType 识别为 image/png
func pngBytes(extra int) []byte {
	head := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}
	return append(head, bytes.Repeat([]byte{0x00}, extra)...)
}

// 真实 JPEG 魔术字节（\xFF\xD8\xFF）
func jpegBytes(extra int) []byte {
	head := []byte{0xFF, 0xD8, 0xFF, 0xE0}
	return append(head, bytes.Repeat([]byte{0x00}, extra)...)
}

func defaultPolicy() UploadPolicy {
	return UploadPolicy{
		MaxBytes:     defaultMaxUploadBytes,
		AllowedMimes: append([]string(nil), defaultAllowedMimes...),
	}
}

// ---- 放行 ----

func TestValidateUpload_AllowsPNG(t *testing.T) {
	fh := makeUploadFile(t, "cat.png", pngBytes(100))
	if err := ValidateUpload(fh, defaultPolicy()); err != nil {
		t.Fatalf("valid PNG should pass, got err: %v", err)
	}
}

func TestValidateUpload_AllowsJPEG(t *testing.T) {
	fh := makeUploadFile(t, "photo.jpg", jpegBytes(200))
	if err := ValidateUpload(fh, defaultPolicy()); err != nil {
		t.Fatalf("valid JPEG should pass, got err: %v", err)
	}
}

// ---- 拒绝：MIME 不在白名单 ----

func TestValidateUpload_RejectsPlainText(t *testing.T) {
	fh := makeUploadFile(t, "notes.txt",
		[]byte(strings.Repeat("hello world\n", 40)))
	err := ValidateUpload(fh, defaultPolicy())
	if err == nil {
		t.Fatal("plain text must be rejected")
	}
	if !strings.Contains(err.Error(), "不支持") {
		t.Logf("(提示) 错误信息未包含 '不支持'，当前: %v", err)
	}
}

// 伪造：扩展名是 png，但内容实际是 HTML
func TestValidateUpload_RejectsMimeMismatchByContent(t *testing.T) {
	payload := []byte(
		"<!DOCTYPE html><html><body>hi</body></html>" +
			strings.Repeat(" ", 500),
	)
	fh := makeUploadFile(t, "fake.png", payload)
	err := ValidateUpload(fh, defaultPolicy())
	if err == nil {
		t.Fatal("HTML 内容即使扩展名是 png 也必须拒绝")
	}
}

// ---- 拒绝：扩展名黑名单 ----

func TestValidateUpload_RejectsDangerousExtension(t *testing.T) {
	// 用 PNG 魔术字节但扩展名是 .exe —— 内容看上去是 image/png，
	// 扩展名层面仍然要拦下
	fh := makeUploadFile(t, "payload.exe", pngBytes(200))
	err := ValidateUpload(fh, defaultPolicy())
	if err == nil {
		t.Fatal(".exe 必须被拒")
	}
}

func TestValidateUpload_RejectsSVG(t *testing.T) {
	// SVG 是 XSS 高风险载体，必须在扩展名层面挡掉
	fh := makeUploadFile(t, "shiny.svg", pngBytes(200))
	err := ValidateUpload(fh, defaultPolicy())
	if err == nil {
		t.Fatal(".svg 必须被拒")
	}
}

// ---- 拒绝：尺寸 ----

func TestValidateUpload_RejectsOversize(t *testing.T) {
	p := defaultPolicy()
	p.MaxBytes = 1024 // 1 KiB
	fh := makeUploadFile(t, "big.png", pngBytes(4096))
	err := ValidateUpload(fh, p)
	if err == nil {
		t.Fatal("超过 MaxBytes 必须被拒")
	}
	if !strings.Contains(err.Error(), "过大") {
		t.Logf("(提示) 错误信息未包含 '过大'，当前: %v", err)
	}
}

// ---- 边界：空文件 ----

func TestValidateUpload_RejectsEmpty(t *testing.T) {
	// 空内容的 detect 会返回 application/octet-stream，不在白名单，应被拒
	fh := makeUploadFile(t, "empty.png", nil)
	if err := ValidateUpload(fh, defaultPolicy()); err == nil {
		t.Fatal("空文件（非图片内容）应被拒")
	}
}

// ---- 扩展名白名单：无扩展名的图片也允许 ----

func TestValidateUpload_AllowsImageWithoutExtension(t *testing.T) {
	fh := makeUploadFile(t, "screenshot", pngBytes(100))
	if err := ValidateUpload(fh, defaultPolicy()); err != nil {
		t.Fatalf("无扩展名的合法图片应放行，got err: %v", err)
	}
}

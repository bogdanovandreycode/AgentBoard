package service

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bogdanovandreycode/agentboard/internal/core"
	"github.com/bogdanovandreycode/agentboard/internal/persistence"
)

func screenshotPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 24, 20))
	for y := 0; y < 20; y++ {
		for x := 0; x < 24; x++ {
			img.SetNRGBA(x, y, color.NRGBA{R: uint8(x * 8), G: uint8(y * 9), B: 80, A: 127})
		}
	}
	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}

func imageMarkdown(url string) string { return "[Screenshot](" + url + ")" }

func TestImageMarkdownReferenceSemantics(t *testing.T) {
	first := "11111111-1111-1111-1111-111111111111"
	second := "22222222-2222-2222-2222-222222222222"
	link := imageMarkdown("/api/images/" + first)
	for name, tc := range map[string]struct {
		content string
		want    []string
	}{
		"link":                  {link, []string{first}},
		"adjacent":              {link + "![second](/api/images/" + second + ")", []string{first, second}},
		"deduplicated":          {link + " " + link, []string{first}},
		"angle and title":       {`[Screen \]1](</api/images/` + first + `> "Bug")`, []string{first}},
		"bare path":             {"/api/images/" + first, nil},
		"external URL":          {imageMarkdown("https://example.com/api/images/" + first), nil},
		"inline code":           {"`" + link + "`", nil},
		"long inline delimiter": {"`` ` " + link + " ` ``", nil},
		"fenced code":           {"```markdown\n" + link + "\n```", nil},
		"tilde fenced code":     {"~~~markdown\n" + link + "\n~~~", nil},
		"longer closing fence":  {"```markdown\n" + link + "\n````\n" + link, []string{first}},
		"unclosed fence":        {"```markdown\n" + link, nil},
		"code then attachment":  {"`" + link + "` " + link, []string{first}},
	} {
		t.Run(name, func(t *testing.T) {
			got := imageReferenceIDs(tc.content)
			if strings.Join(got, ",") != strings.Join(tc.want, ",") {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestImagesStoredAsWebPWithAlphaAndSharedDatabaseDirectory(t *testing.T) {
	f := setup(t)
	img, err := f.svc.HumanUploadImage(f.ctx, f.project.ID, `C:\screenshots\bug.png`, screenshotPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	if img.Name != "bug.webp" || img.Width != 24 || img.Height != 20 || !filepath.IsAbs(img.Path) || img.URL != "/api/images/"+img.ID {
		t.Fatalf("unexpected image metadata: %+v", img)
	}
	if filepath.Dir(img.Path) != f.store.ImageDirectory() || filepath.Base(img.Path) != img.ID+".webp" {
		t.Fatalf("wrong storage path: %s", img.Path)
	}
	var filename, databasePath string
	if err := f.store.DB.QueryRow(`SELECT path FROM images WHERE id=?`, img.ID).Scan(&filename); err != nil {
		t.Fatal(err)
	}
	if filename != img.ID+".webp" {
		t.Fatalf("database must contain only relative file reference, got %q", filename)
	}
	if err := f.store.DB.QueryRow(`SELECT file FROM pragma_database_list WHERE name='main'`).Scan(&databasePath); err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(f.store.ImageDirectory()) != filepath.Dir(databasePath) {
		t.Fatal("images must follow the selected database")
	}
	// A separate HTTP/MCP process opening the same database sees the same file.
	other, err := persistence.Open(databasePath)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	metadata, data, err := New(other).HumanImage(f.ctx, img.ID)
	if err != nil || metadata.Path != img.Path {
		t.Fatalf("read through second store: %+v %v", metadata, err)
	}
	if len(data) < 12 || string(data[:4]) != "RIFF" || string(data[8:12]) != "WEBP" {
		t.Fatal("stored file is not WebP")
	}
	decoded, format, err := image.Decode(bytes.NewReader(data))
	if err != nil || format != "webp" {
		t.Fatalf("decode stored image: %s %v", format, err)
	}
	want := color.NRGBA{R: 56, G: 81, B: 80, A: 127}
	if got := color.NRGBAModel.Convert(decoded.At(7, 9)); got != want {
		t.Fatalf("pixels or alpha changed: got %v, want %v", got, want)
	}
	if files, err := os.ReadDir(f.store.ImageDirectory()); err != nil || len(files) != 1 {
		t.Fatalf("upload left temporary files: %v %v", files, err)
	}
}

func TestImagesNormalizeJPEGAndWebP(t *testing.T) {
	f := setup(t)
	var jpg bytes.Buffer
	if err := jpeg.Encode(&jpg, image.NewRGBA(image.Rect(0, 0, 12, 8)), nil); err != nil {
		t.Fatal(err)
	}
	img, err := f.svc.HumanUploadImage(f.ctx, f.project.ID, "photo.jpg", jpg.Bytes())
	if err != nil || img.Name != "photo.webp" {
		t.Fatalf("JPEG upload: %+v %v", img, err)
	}
	_, data, err := f.svc.HumanImage(f.ctx, img.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.svc.HumanUploadImage(f.ctx, f.project.ID, "already.webp", data); err != nil {
		t.Fatalf("WebP upload: %v", err)
	}
}

func TestImagesRejectInvalidOversizedAndForeignReferences(t *testing.T) {
	f := setup(t)
	largeDimensions := append([]byte(nil), screenshotPNG(t)...)
	binary.BigEndian.PutUint32(largeDimensions[16:20], 16385)
	binary.BigEndian.PutUint32(largeDimensions[29:33], crc32.ChecksumIEEE(largeDimensions[12:29]))
	largePixels := append([]byte(nil), screenshotPNG(t)...)
	binary.BigEndian.PutUint32(largePixels[16:20], 5000)
	binary.BigEndian.PutUint32(largePixels[20:24], 5000)
	binary.BigEndian.PutUint32(largePixels[29:33], crc32.ChecksumIEEE(largePixels[12:29]))
	for name, input := range map[string][]byte{
		"empty": nil, "invalid": []byte("not an image"), "svg": []byte(`<svg xmlns="http://www.w3.org/2000/svg"/>`),
		"too many bytes": make([]byte, MaxImageBytes+1), "too wide": largeDimensions, "too many pixels": largePixels,
		"truncated": screenshotPNG(t)[:40],
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := f.svc.HumanUploadImage(f.ctx, f.project.ID, name, input); !errors.Is(err, core.ErrInvalidInput) {
				t.Fatalf("invalid upload accepted: %v", err)
			}
		})
	}
	img, err := f.svc.HumanUploadImage(f.ctx, f.project.ID, "screen.png", screenshotPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	refs, err := f.svc.ImageReferences(f.ctx, f.project.ID, "[Screen 1]("+img.URL+") and ![Screen 2]("+img.URL+")")
	if err != nil || len(refs) != 1 || refs[0].ID != img.ID {
		t.Fatalf("references not deduplicated: %v %v", refs, err)
	}
	other, err := f.store.EnsureProject(f.ctx, "other", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.svc.ImageReferences(f.ctx, other.ID, imageMarkdown(img.URL)); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("cross-project reference accepted: %v", err)
	}
	if _, err = f.svc.ImageReferences(f.ctx, f.project.ID, imageMarkdown("/api/images/00000000-0000-0000-0000-000000000000")); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("missing reference accepted: %v", err)
	}
	if _, err := f.svc.ImageReferences(f.ctx, f.project.ID, imageMarkdown(img.URL+"extra")); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("image URL prefix matched: %v", err)
	}
}

func TestAgentImagePermissionsAndReferenceOwnership(t *testing.T) {
	f := setup(t)
	task := createTask(t, f, f.a, "features", "ai")
	otherTask := createTask(t, f, f.a, "features", "ai")
	hiddenTask := createTask(t, f, f.b, "features", "ai")
	backlog := createTask(t, f, f.a, "backlog", "ai")
	for _, target := range []core.Task{hiddenTask, backlog} {
		if _, err := f.svc.AgentUploadImage(f.ctx, f.ac, target.ID, "screen.png", screenshotPNG(t)); !errors.Is(err, core.ErrTaskNotAccessible) {
			t.Fatalf("upload to inaccessible task accepted: %v", err)
		}
	}
	img, err := f.svc.AgentUploadImage(f.ctx, f.ac, task.ID, "screen.png", screenshotPNG(t))
	if err != nil || img.UploadedTaskID != task.ID {
		t.Fatalf("agent upload: %+v %v", img, err)
	}
	if err := f.svc.AgentValidateImageReferences(f.ctx, f.ac, task.ID, imageMarkdown(img.URL)); err != nil {
		t.Fatalf("cannot attach own new screenshot: %v", err)
	}
	if err := f.svc.AgentValidateImageReferences(f.ctx, f.ac, otherTask.ID, imageMarkdown(img.URL)); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("unreferenced screenshot shared across tasks: %v", err)
	}
	if _, _, err := f.svc.AgentImage(f.ctx, f.ac, task.ID, img.ID); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("unattached image readable: %v", err)
	}
	if _, err = f.svc.AgentAddHistory(f.ctx, f.ac, task.ID, "[Screen 1]("+img.URL+")"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := f.svc.AgentImage(f.ctx, f.ac, task.ID, img.ID); err != nil {
		t.Fatalf("attached screenshot unreadable: %v", err)
	}
	if _, _, err := f.svc.AgentImage(f.ctx, f.bc, task.ID, img.ID); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("other worker read screenshot: %v", err)
	}
	if err := f.svc.AgentValidateImageReferences(f.ctx, f.ac, "", imageMarkdown(img.URL)); err != nil {
		t.Fatalf("readable screenshot not permitted on new feature: %v", err)
	}
	hiddenImage, err := f.svc.AgentUploadImage(f.ctx, f.bc, hiddenTask.ID, "hidden.png", screenshotPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := f.svc.AgentValidateImageReferences(f.ctx, f.ac, task.ID, imageMarkdown(hiddenImage.URL)); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("guessed hidden screenshot can be attached: %v", err)
	}
	if err := f.svc.AgentValidateImageReferences(f.ctx, f.ac, "", imageMarkdown(hiddenImage.URL)); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("guessed hidden screenshot can be used in new feature: %v", err)
	}
	if _, err := f.svc.HumanMove(f.ctx, task.ID, "complete"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := f.svc.AgentImage(f.ctx, f.ac, task.ID, img.ID); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("screenshot from completed task readable: %v", err)
	}
}

func TestHumanDescriptionImageIsVisibleToAssignedAgent(t *testing.T) {
	f := setup(t)
	img, err := f.svc.HumanUploadImage(f.ctx, f.project.ID, "bug.png", screenshotPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	task, err := f.svc.HumanCreateTask(f.ctx, f.project.ID, core.TaskInput{Title: "Fix screenshot", Description: "[Screen 1](" + img.URL + ")", State: "features", TestingMode: "ai", Assignee: workerAssignee(f.a.ID)})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := f.svc.AgentImage(f.ctx, f.ac, task.ID, img.ID); err != nil {
		t.Fatalf("human description screenshot unavailable to worker: %v", err)
	}
	if err := f.svc.AgentValidateImageReferences(f.ctx, f.ac, task.ID, "[Screen 1]("+img.URL+")"); err != nil {
		t.Fatalf("agent cannot mention existing human screenshot: %v", err)
	}
}

func TestTestingInstructionImagesAreVisibleToAssignedAgent(t *testing.T) {
	f := setup(t)
	img, err := f.svc.HumanUploadImage(f.ctx, f.project.ID, "test.png", screenshotPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"ai", "human"} {
		t.Run(field, func(t *testing.T) {
			in := core.TaskInput{Title: "Check screenshot", Description: "```unclosed code fence", State: "features", TestingMode: "ai", Assignee: workerAssignee(f.a.ID)}
			if field == "ai" {
				in.AITestInstructions = imageMarkdown(img.URL)
			} else {
				in.HumanTestInstructions = imageMarkdown(img.URL)
			}
			task, err := f.svc.HumanCreateTask(f.ctx, f.project.ID, in)
			if err != nil {
				t.Fatal(err)
			}
			if _, _, err := f.svc.AgentImage(f.ctx, f.ac, task.ID, img.ID); err != nil {
				t.Fatalf("testing instruction screenshot unavailable: %v", err)
			}
			for _, read := range []func() (core.TaskDetails, error){
				func() (core.TaskDetails, error) { return f.svc.AgentTask(f.ctx, f.ac, task.ID) },
				func() (core.TaskDetails, error) { return f.svc.HumanTask(f.ctx, task.ID) },
			} {
				got, err := read()
				if err != nil || len(got.Images) != 1 || got.Images[0].ID != img.ID {
					t.Fatalf("code fence in description hid testing instruction image: %+v %v", got.Images, err)
				}
			}
		})
	}
}

func TestImageReferencesKeepMarkdownFieldsIndependent(t *testing.T) {
	f := setup(t)
	img, err := f.svc.HumanUploadImage(f.ctx, f.project.ID, "screen.png", screenshotPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	refs, err := f.svc.ImageReferencesInTexts(f.ctx, f.project.ID, "```unclosed\n", imageMarkdown(img.URL), imageMarkdown(img.URL))
	if err != nil || len(refs) != 1 || refs[0].ID != img.ID {
		t.Fatalf("fields not independent/deduplicated: %+v %v", refs, err)
	}
	other, err := f.store.EnsureProject(f.ctx, "other", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.ImageReferencesInTexts(f.ctx, other.ID, "```unclosed\n", imageMarkdown(img.URL)); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("code fence hid foreign image in another field: %v", err)
	}
}

func TestSystemEventsCannotGrantAgentImageAccess(t *testing.T) {
	f := setup(t)
	task := createTask(t, f, f.a, "features", "ai")
	hiddenTask := createTask(t, f, f.b, "features", "ai")
	img, err := f.svc.AgentUploadImage(f.ctx, f.bc, hiddenTask.ID, "hidden.png", screenshotPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	link := imageMarkdown(img.URL)
	if _, err := f.svc.AgentAddArtifact(f.ctx, f.ac, task.ID, core.ArtifactInput{Name: link, Path: "output/screenshot.txt"}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := f.svc.AgentImage(f.ctx, f.ac, task.ID, img.ID); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("artifact name granted hidden image access: %v", err)
	}
	if err := f.svc.AgentValidateImageReferences(f.ctx, f.ac, task.ID, link); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("artifact event allowed attaching a hidden screenshot: %v", err)
	}
	if _, err := f.svc.CreateFeature(f.ctx, f.ac, core.FeatureInput{Title: "Feature with hidden screenshot", Reason: link, Assignee: workerAssignee(f.a.ID)}); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("feature reason accepted hidden screenshot: %v", err)
	}
	// Historical system messages must remain harmless even if they were created
	// before reason validation was introduced.
	if _, err := f.store.DB.Exec(`INSERT INTO history_entries(id,task_id,actor_type,entry_type,content,created_at) VALUES(?,?,?,?,?,?)`, id(), task.ID, "system", "task_created", "Reason: "+link, core.Now()); err != nil {
		t.Fatal(err)
	}
	if _, _, err := f.svc.AgentImage(f.ctx, f.ac, task.ID, img.ID); !errors.Is(err, core.ErrTaskNotAccessible) {
		t.Fatalf("legacy feature reason granted hidden image access: %v", err)
	}
}

func TestImageRejectsTamperedDatabasePathAndMissingFile(t *testing.T) {
	f := setup(t)
	img, err := f.svc.HumanUploadImage(f.ctx, f.project.ID, "screen.png", screenshotPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(img.Path); err != nil {
		t.Fatal(err)
	}
	if _, _, err := f.svc.HumanImage(f.ctx, img.ID); !errors.Is(err, core.ErrNotFound) {
		t.Fatalf("missing image file not reported: %v", err)
	}
	if _, err := f.store.DB.Exec(`UPDATE images SET path=? WHERE id=?`, "../secret.webp", img.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, err := f.svc.HumanImage(f.ctx, img.ID); !errors.Is(err, core.ErrInvalidInput) || !strings.Contains(err.Error(), "path") {
		t.Fatalf("tampered path not rejected: %v", err)
	}
}

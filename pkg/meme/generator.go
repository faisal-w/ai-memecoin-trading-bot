package meme

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// Generator handles meme generation
type Generator struct {
	TemplatesDir string
	OutputDir    string
	font         *truetype.Font
}

// NewGenerator creates a new meme generator
func NewGenerator(templatesDir, outputDir string) (*Generator, error) {
	// Create directories if they don't exist
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create templates directory: %w", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	return &Generator{
		TemplatesDir: templatesDir,
		OutputDir:    outputDir,
		font:         font,
	}, nil
}

// CreateMeme generates a meme from a template with text overlay
func (g *Generator) CreateMeme(templateName, topText, bottomText string) (string, error) {
	templatePath := filepath.Join(g.TemplatesDir, templateName)
	
	// Open template image
	file, err := os.Open(templatePath)
	if err != nil {
		return "", fmt.Errorf("template not found: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Create a new image to draw on
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Add text to image
	if err := g.addText(rgba, topText, true); err != nil {
		return "", err
	}
	if err := g.addText(rgba, bottomText, false); err != nil {
		return "", err
	}

	// Save the image
	outputPath := filepath.Join(g.OutputDir, "meme_"+templateName)
	return outputPath, g.saveImage(rgba, outputPath)
}

// CreateMemeFromURL generates a meme from a URL with text overlay
func (g *Generator) CreateMemeFromURL(imageURL, topText, bottomText string) (string, error) {
	// Download image
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Create a new image to draw on
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Add text to image
	if err := g.addText(rgba, topText, true); err != nil {
		return "", err
	}
	if err := g.addText(rgba, bottomText, false); err != nil {
		return "", err
	}

	// Save the image
	outputPath := filepath.Join(g.OutputDir, "meme_from_url.png")
	return outputPath, g.saveImage(rgba, outputPath)
}

// ListTemplates returns a list of available templates
func (g *Generator) ListTemplates() ([]string, error) {
	entries, err := os.ReadDir(g.TemplatesDir)
	if err != nil {
		return nil, err
	}

	var templates []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" {
			templates = append(templates, name)
		}
	}
	return templates, nil
}

func (g *Generator) addText(img *image.RGBA, text string, isTop bool) error {
	if text == "" {
		return nil
	}

	text = strings.ToUpper(text)
	bounds := img.Bounds()
	
	// Calculate font size based on image height
	fontSize := float64(bounds.Dy()) / 10.0
	
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(g.font)
	c.SetFontSize(fontSize)
	c.SetClip(bounds)
	c.SetDst(img)

	// Calculate text position
	pt := freetype.Pt(0, 0)
	if isTop {
		pt.Y = c.PointToFixed(fontSize * 1.5)
	} else {
		pt.Y = c.PointToFixed(float64(bounds.Dy()) - fontSize)
	}

	// Center text horizontally (approximate)
	textWidth := len(text) * int(fontSize) / 2
	pt.X = c.PointToFixed(float64(bounds.Dx()-textWidth) / 2.0)

	// Draw text outline (black)
	for dx := -2; dx <= 2; dx++ {
		for dy := -2; dy <= 2; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			c.SetSrc(image.NewUniform(color.Black))
			outlinePt := pt
			outlinePt.X += c.PointToFixed(float64(dx))
			outlinePt.Y += c.PointToFixed(float64(dy))
			c.DrawString(text, outlinePt)
		}
	}

	// Draw main text (white)
	c.SetSrc(image.NewUniform(color.White))
	_, err := c.DrawString(text, pt)
	return err
}

func (g *Generator) saveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".jpg" || ext == ".jpeg" {
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
	}
	return png.Encode(file, img)
}

// SaveImageTo saves an image to a writer
func (g *Generator) SaveImageTo(img image.Image, w io.Writer) error {
	return png.Encode(w, img)
}

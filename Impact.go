package Impact

import (
	"bufio"
	"fmt"
	"image"
	"strings"

	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func Impact(imagePath string, fontSize float64, toptext, bottomtext string) string {

	fontBytes, err := ioutil.ReadFile("resources/impact.ttf")
	if err != nil {
		panic(err)
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var img_bg image.Image

	switch imagePath[len(imagePath)-4:] {
	case ".png":
		img_bg, err = png.Decode(file)
	case ".jpg":
		img_bg, err = jpeg.Decode(file)
	}

	if err != nil {
		log.Fatal(err)
	}

	rgba := image.NewRGBA(img_bg.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img_bg, image.Point{}, draw.Src)
	c := freetype.NewContext()

	py := img_bg.Bounds().Dy() / 100

	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetHinting(font.HintingFull)

	face := truetype.NewFace(f, &truetype.Options{Size: fontSize})

	toptextArr := strings.Split(toptext, " ")
	toplines := []string{toptextArr[0]}
	toptextArr = toptextArr[1:]
	topline := 0
	toplenght := 0

	for _, v := range toptextArr {
		if toplenght+font.MeasureString(face, v).Round() < img_bg.Bounds().Dx() {
			toplines[topline] += " " + v
			toplenght = font.MeasureString(face, toplines[topline]).Round()
		} else {
			topline += 1
			toplenght = 0
			toplines = append(toplines, v)
		}
	}

	n := 3 // stroke

	for i, v := range toplines {
		textWidth := font.MeasureString(face, v).Round()
		c.SetSrc(image.Black)
		for dx := -n; dx <= n; dx++ {
			for dy := -n; dy < n; dy++ {
				if (dx*dx+dy*dy >= n*n) || ((dx == 0) && (dy == 0)) {
					continue
				}
				c.DrawString(v, freetype.Pt(img_bg.Bounds().Dx()/2-textWidth/2+(dx*2), int(fontSize)*(i+1)+(dy*2)))
			}
		}
		c.SetSrc(image.White)
		c.DrawString(v, freetype.Pt(img_bg.Bounds().Dx()/2-textWidth/2, int(fontSize)*(i+1)))
	}

	bottomtextArr := strings.Split(bottomtext, " ")
	bottomlines := []string{bottomtextArr[0]}
	bottomtextArr = bottomtextArr[1:]
	bottomline := 0
	bottomlenght := 0

	for _, v := range bottomtextArr {
		if bottomlenght+font.MeasureString(face, v).Round() < img_bg.Bounds().Dx() {
			bottomlines[bottomline] += " " + v
			bottomlenght = font.MeasureString(face, bottomlines[bottomline]).Round()
		} else {
			bottomline += 1
			bottomlenght = 0
			bottomlines = append(bottomlines, v)
		}
	}

	for i, v := range bottomlines {
		textWidth := font.MeasureString(face, v).Round()
		c.SetSrc(image.Black)
		for dx := -n; dx <= n; dx++ {
			for dy := -n; dy < n; dy++ {
				if (dx*dx+dy*dy >= n*n) || ((dx == 0) && (dy == 0)) {
					continue
				}
				c.DrawString(v, freetype.Pt(img_bg.Bounds().Dx()/2-textWidth/2+(dx*2), img_bg.Bounds().Dy()-(int(fontSize)*(len(bottomlines)-(i+1))+(dy*2)+ +py*2)))
			}
		}
		c.SetSrc(image.White)
		c.DrawString(v, freetype.Pt(img_bg.Bounds().Dx()/2-textWidth/2, img_bg.Bounds().Dy()-(int(fontSize)*(len(bottomlines)-(i+1))+py*2)))
	}

	a := strings.Split(imagePath, "/")

	outFile, err := os.Create(a[0] + "/" + "impact_" + a[1])
	if err != nil {
		log.Println(err)

	}
	defer outFile.Close()

	fmt.Println(1)

	b := bufio.NewWriter(outFile)

	switch imagePath[len(imagePath)-4:] {
	case ".png":
		err = png.Encode(b, rgba)
	case ".jpg":
		err = jpeg.Encode(b, rgba, &jpeg.Options{Quality: 100})
	}
	fmt.Println(2)
	if err != nil {
		panic(err)
	}

	fmt.Println(2)

	err = b.Flush()
	if err != nil {
		panic(err)
	}

	return a[0] + "/" + "impact_" + a[1]
}


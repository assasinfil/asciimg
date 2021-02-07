package main

import (
	"flag"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func scale(img image.Image, w int, h int) image.Image {
	dstImg := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.NearestNeighbor.Scale(dstImg, dstImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dstImg
}

func decodeImageFile(imgName string) (image.Image, error) {
	imgFile, err := os.Open(imgName)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(imgFile)

	return img, err
}

func processPixel(c color.Color) rune {
	gc := color.GrayModel.Convert(c)
	r, _, _, _ := gc.RGBA()
	r = r >> 8
	//fmt.Println(r)
	symbols := []rune("@80GCLft1i;:,. ")
	index := int(r) * len(symbols) / 256
	return symbols[index]
}

func convertToAscii(img image.Image) [][]rune {
	textImg := make([][]rune, img.Bounds().Dy())
	for i := range textImg {
		textImg[i] = make([]rune, img.Bounds().Dx())
	}

	for i := range textImg {
		for j := range textImg[i] {
			textImg[i][j] = processPixel(img.At(j, i))
		}
	}
	return textImg
}

var (
	output  = flag.String("o", "", "-o <out.txt>")
	weigth  = flag.Int("w", 200, "")
	heigth  = flag.Int("h", 40, "")
	noscale = flag.Bool("noscale", false, "")
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: asciimg <image.jpg>")
		os.Exit(0)
	}
	imgName := flag.Arg(0)

	img, err := decodeImageFile(imgName)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	if !*noscale && *output == "" {
		img = scale(img, *weigth, *heigth)
	}

	textImg := convertToAscii(img)

	if *output == "" {
		for i := range textImg {
			for j := range textImg[i] {
				fmt.Printf("%c", textImg[i][j])
			}
			fmt.Println()
		}
	} else {
		file, _ := os.Create(*output)
		defer file.Close()
		for i := range textImg {
			for j := range textImg[i] {
				fmt.Fprintf(file, "%c", textImg[i][j])
			}
			fmt.Fprintf(file, "\n")
		}
	}

}

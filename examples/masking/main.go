// Copyright 2016 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/common"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

var (
	gophersImage   *ebiten.Image
	fiveyearsImage *ebiten.Image
	maskImage      *ebiten.Image
	spotLightImage *ebiten.Image
	spotLightX     = 0
	spotLightY     = 0
	spotLightVX    = 1
	spotLightVY    = 1
)

func update(screen *ebiten.Image) error {
	spotLightX += spotLightVX
	spotLightY += spotLightVY
	if spotLightX < 0 {
		spotLightX = -spotLightX
		spotLightVX = -spotLightVX
	}
	if spotLightY < 0 {
		spotLightY = -spotLightY
		spotLightVY = -spotLightVY
	}
	w, h := spotLightImage.Size()
	maxX, maxY := screenWidth-w, screenHeight-h
	if maxX < spotLightX {
		spotLightX = -spotLightX + 2*maxX
		spotLightVX = -spotLightVX
	}
	if maxY < spotLightY {
		spotLightY = -spotLightY + 2*maxY
		spotLightVY = -spotLightVY
	}

	if err := maskImage.Clear(); err != nil {
		return err
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(spotLightX), float64(spotLightY))
	if err := maskImage.DrawImage(spotLightImage, op); err != nil {
		return err
	}

	op = &ebiten.DrawImageOptions{}
	op.CompositeMode = ebiten.CompositeModeSourceOut
	if err := maskImage.DrawImage(fiveyearsImage, op); err != nil {
		return err
	}

	if err := screen.Fill(color.RGBA{0x00, 0x00, 0x80, 0xff}); err != nil {
		return err
	}
	if err := screen.DrawImage(gophersImage, &ebiten.DrawImageOptions{}); err != nil {
		return err
	}
	if err := screen.DrawImage(maskImage, &ebiten.DrawImageOptions{}); err != nil {
		return err
	}

	return nil
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	var err error
	gophersImage, _, err = common.AssetImage("gophers.jpg", ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	fiveyearsImage, _, err = common.AssetImage("fiveyears.jpg", ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	maskImage, err = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}

	as := image.Point{128, 128}
	a := image.NewAlpha(image.Rectangle{image.ZP, as})
	for j := 0; j < as.Y; j++ {
		for i := 0; i < as.X; i++ {
			r := as.X / 2
			d := math.Sqrt(float64((i-r)*(i-r) + (j-r)*(j-r)))
			b := uint8(max(0, min(0xff, 3*0xff-int(d*3*0xff)/r)))
			a.SetAlpha(i, j, color.Alpha{b})
		}
	}
	spotLightImage, err = ebiten.NewImageFromImage(a, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Masking (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}

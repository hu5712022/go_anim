package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	. "github.com/andlabs/ui"
)

var (
	startPX          = 400.0
	startPY          = 250.0
	mArea            *Area
	mAreaDp          *AreaDrawParams
	handler          MyHandler
	radiu            = 100.0
	leftX, leftY     = startPX, startPY
	centerX, centerY = startPX, startPY
	rightX, rightY   = startPX, startPY

	states = 1 // 1 切割 2 旋转 3 闪烁

	v = 5.0 // 切割时 贝塞尔曲线的弧度变化速率参数

	angle   = 0.0   // 旋转时  旋转的角度
	radiu_2 = 200.0 // 旋转时 两圆之间得到半径

	largerCount = 0

	r, g, b, al = 1.0, 1.0, 0.0, 1.0

	rate1, rate2, rate3 time.Duration = 50, 100, 10 //动画的频率 每次更新次数
	rateColor           time.Duration = 10          // 颜色变化速率
)

func main() {
	showWindow()
}

func color() {
	for {
		time.Sleep(time.Second / rateColor)
		i := randInt(8)
		//红色 #FF0000 1  0   0  1
		//橙色 #FF7F00 1  0.5 0  1
		//黄色 #FFFF00 1  1   0  1
		//绿色 #00FF00 0  1   0  1
		//青色 #00FFFF 0  1   1  1
		//蓝色 #0000FF 0  0   1  1
		//紫色 #8B00FF 0.5 0  1  1
		switch i {
		case 0:
			r, g, b, al = 0.9, 0.1, 0.1, 0.8
		case 1:
			r, g, b, al = 0.2, 0.7, 0.1, 0.8
		case 2:
			r, g, b, al = 0.9, 1.0, 0.0, 0.8
		case 3:
			r, g, b, al = 0.9, 0.5, 0.0, 0.8
		case 4:
			r, g, b, al = 0.1, 1.0, 0.9, 0.8
		case 5:
			r, g, b, al = 0.1, 0.1, 1.0, 0.8
		case 6:
			r, g, b, al = 0.5, 0.1, 1.0, 0.9
		case 7:
			r = randFloat()
			g = randFloat()
			b = randFloat()
		}
	}
}

func update() {
	for {
		if states == 1 {
			time.Sleep(time.Second / rate1)
			radiu -= 0.3
			v -= 0.03
			leftX -= 1.5
			rightX += 1.5
			if rightX-centerX > radiu_2 {
				states = 2
				v = 5.0
			}
		} else if states == 2 {
			time.Sleep(time.Second / rate2)
			radiu_2 -= 0.2
			if radiu_2 <= 0 {
				states = 3
				radiu = 100
			}
			v -= 0.01
			if v <= 0 {
				v = 5.0
			}

			rateColor = (time.Duration)(10 / v)
			angle += 5
			lenY := radiu_2 * math.Sin(math.Pi*angle/v/180.0)
			lenX := radiu_2 * math.Cos(math.Pi*angle/v/180.0)

			leftX = centerX - lenX
			leftY = centerY - lenY

			rightX = centerX + lenX
			rightY = centerY + lenY
		} else if states == 3 {
			time.Sleep(time.Second / rate3)
			//			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			//			radiu = float64(r.Intn(200))
			if largerCount < 4 {
				radiu += 20
				if radiu > 400 {
					largerCount++
					if largerCount < 3 {
						radiu = 100
					}

				}
			} else {
				radiu -= 20
			}

			if radiu <= 0 {
				radiu = 100
				states = 1
				v = 5.0
				radiu_2 = 200.0
				angle = 0.0
				largerCount = 0
			}
		}

		if radiu == 0 {
			return
		}

		//fmt.Println(radiu)
		fmt.Println("states:", states)
		mArea.QueueRedrawAll()

	}
}

func showWindow() {

	err := Main(func() {

		handler = MyHandler{}
		mArea = NewScrollingArea(handler, 100, 1000)
		//		mArea := NewArea(handler)
		bt := NewButton("走你")

		bt.OnClicked(func(b *Button) {
			b.Hide()
			go update()
			go color()
		})

		box := NewVerticalBox()
		box.Append(bt, false)
		box.Append(mArea, true)

		//		创建window窗口。并设置长宽。
		window := NewWindow("第一个应用程序。", 800, 500, true)
		window.Margined()
		window.SetChild(box)
		window.OnClosing(func(*Window) bool {
			Quit()
			return true
		})
		window.Show()

	})
	if err != nil {

		panic(err)
	}
}

func Text(dc *DrawContext, str string) {
	//font family 一定要填正确，百度下window字体  simsun 宋体
	fd := &FontDescriptor{"SimSun", 20, TextWeightNormal, TextItalicNormal, TextStretchNormal}
	fmt.Println(fd)
	f := LoadClosestFont(fd)
	dc.Text(300, 20, NewTextLayout(str, f, -1))

}

func (h MyHandler) Draw(a *Area, dp *AreaDrawParams) {
	fmt.Println("Draw")

	Text(dp.Context, "hia hia hia")

	p := NewPath(Winding)

	p.NewFigure(100, 100) // or NewFigureWithArc
	p.NewFigureWithArc(leftX, leftY, radiu, 0, math.Pi*2, true)
	p.NewFigureWithArc(rightX, rightY, radiu, 0, math.Pi*2, true)
	p.NewFigure(leftX, leftY+radiu)

	if states == 1 {
		if radiu-(rightX-rightY)/v > 0 && radiu+(rightX-leftX)/v > 0 {
			p.BezierTo(leftX, leftY+radiu, startPX, startPY+radiu-(rightX-leftX)/v, rightX, rightY+radiu)
			p.LineTo(rightX, rightY-radiu)

			p.BezierTo(rightX, rightY-radiu, startPX, startPY-radiu+(rightX-leftX)/v, leftX, leftY-radiu)
			p.LineTo(leftX, leftY+radiu)
		}
	}

	p.CloseFigure()
	p.End()

	brush := &Brush{Solid, r, g, b, al, 0, 0, 0, 0, 0, []GradientStop{{1, 0.2, 0.8, 0.5, 0.9}}}

	//	brush := &Brush{Solid, 1, 1, 0, 0.5, 0, 0, 0, 0, 0, []GradientStop{{1, 0.2, 0.8, 0.5, 0.9}}}
	dp.Context.Fill(p, brush)
	p.Free()
}
func (h MyHandler) MouseEvent(a *Area, me *AreaMouseEvent) {
	fmt.Println("MouseEvent")

}
func (h MyHandler) MouseCrossed(a *Area, left bool) {
	fmt.Println("MouseCrossed")
}
func (h MyHandler) DragBroken(a *Area) {
	fmt.Println("DragBroken")
}
func (h MyHandler) KeyEvent(a *Area, ke *AreaKeyEvent) bool {
	fmt.Println("KeyEvent")
	return true
}

type MyHandler struct {
}

func randFloat() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return float64(r.Intn(100)) / 100
}

func randInt(m int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(m)
}

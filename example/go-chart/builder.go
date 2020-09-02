package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/PreciselyData/compose-chart-api/pic"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

type builder struct {
	chart.RendererProvider
	*pic.Config
	pic.NumberFormat
	chartType     string
	width, height int
	dpi           float64
	data          *pic.Data
	title         string
	titleFont     pic.Font
	axisFont      pic.Font
	bgColor       pic.Color
}

type chartRenderer interface {
	Render(rp chart.RendererProvider, w io.Writer) error
}

func newBuilder(c *pic.Config) *builder {
	return &builder{
		Config:       c,
		NumberFormat: c.NumberFormat(),
		chartType:    c.Name(),
		data:         c.Data(),
		title:        c.Value("title").Text(),
		titleFont:    c.Font("titleFont"),
		axisFont:     c.Font("axisFont"),
		bgColor:      c.Color("bgColor"),
	}
}

// SetFormat is part of the pic.Builder interface. The given format and
// colorSpace represent the required image format. The builder can change
// these values to those supported by the renderer, but doing so will mean
// that Designer/Generate will need to convert the image on output.
func (b *builder) SetFormat(format *pic.ImageFormat, colorSpace *pic.ColorSpace) {
	switch *format {
	case pic.SVG:
		b.RendererProvider = chart.SVG
	default:
		b.RendererProvider = chart.PNG
		*format = pic.PNG
	}
	*colorSpace = pic.RGB
}

// SetSize is part of the pic.Builder interface.
func (b *builder) SetSize(width, height pic.Twiplet, dpi int32) {
	b.width = width.Pixels(dpi)
	b.height = height.Pixels(dpi)
	b.dpi = float64(dpi)
}

// Render is part of the pic.Builder interface.
func (b *builder) Render() (*bytes.Buffer, error) {
	switch b.chartType {
	case "pie":
		return b.renderChart(b.newPieChart())
	case "donut":
		return b.renderChart(b.newDonutChart())
	case "line":
		return b.renderChart(b.newLineChart())
	default:
		return nil, fmt.Errorf("unknown chart type '%s'", b.chartType)
	}
}

func (b *builder) renderChart(renderer chartRenderer) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := renderer.Render(b.RendererProvider, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (b *builder) newPieChart() *chart.PieChart {
	return &chart.PieChart{
		Title:        b.title,
		TitleStyle:   b.titleStyle(),
		ColorPalette: b,
		Width:        b.width,
		Height:       b.height,
		DPI:          b.dpi,
		Values:       b.buildSingleSeriesValues(),
	}
}

func (b *builder) newDonutChart() *chart.DonutChart {
	return &chart.DonutChart{
		Title:        b.title,
		TitleStyle:   b.titleStyle(),
		ColorPalette: b,
		Width:        b.width,
		Height:       b.height,
		DPI:          b.dpi,
		Values:       b.buildSingleSeriesValues(),
	}
}

func (b *builder) newLineChart() *chart.Chart {
	font := b.ResolveFont(b.axisFont)
	graph := &chart.Chart{
		Title:        b.title,
		TitleStyle:   b.titleStyle(),
		ColorPalette: b,
		Width:        b.width,
		Height:       b.height,
		DPI:          b.dpi,
		XAxis: chart.XAxis{
			Style: b.fontStyle(font),
			Ticks: b.buildXTicks(),
		},
		YAxis: chart.YAxis{
			Style:          b.fontStyle(font),
			ValueFormatter: b.valueFormatter,
		},
		Series: b.buildMultiSeries(),
	}
	b.addLegend(graph)
	return graph
}

// BackgroundColor is part of the chart.ColorPalette interface.
func (b *builder) BackgroundColor() drawing.Color {
	return drawing.Color{
		R: b.bgColor.R,
		G: b.bgColor.G,
		B: b.bgColor.B,
		A: 255,
	}
}

// BackgroundStrokeColor is part of the chart.ColorPalette interface.
func (b *builder) BackgroundStrokeColor() drawing.Color {
	return drawing.ColorBlack
}

// CanvasColor is part of the chart.ColorPalette interface.
func (b *builder) CanvasColor() drawing.Color {
	return drawing.Color{
		R: b.bgColor.R,
		G: b.bgColor.G,
		B: b.bgColor.B,
		A: 255,
	}
}

// CanvasStrokeColor is part of the chart.ColorPalette interface.
func (b *builder) CanvasStrokeColor() drawing.Color {
	return drawing.ColorBlack
}

// AxisStrokeColor is part of the chart.ColorPalette interface.
func (b *builder) AxisStrokeColor() drawing.Color {
	return drawing.ColorBlack
}

// TextColor is part of the chart.ColorPalette interface.
func (b *builder) TextColor() drawing.Color {
	return drawing.ColorBlack
}

// GetSeriesColor is part of the chart.ColorPalette interface.
func (b *builder) GetSeriesColor(index int) drawing.Color {
	if index < len(b.data.Colors) {
		color := b.data.Colors[index]
		return drawing.Color{
			R: color.R,
			G: color.G,
			B: color.B,
			A: 255,
		}
	}
	return drawing.ColorRed
}

func (b *builder) titleStyle() chart.Style {
	if b.title == "" {
		return chart.Style{Hidden: true}
	}
	font := b.ResolveFont(b.titleFont)
	return b.fontStyle(font)
}

func (b *builder) addLegend(graph *chart.Chart) {
	if b.Value("legend").True() {
		if b.Value("legendPos") == "left" {
			color := b.Color("legendColor")
			opacity := b.Integer("legendOpacity")
			legendStyle := chart.Style{
				FillColor: drawing.Color{
					R: color.R,
					G: color.G,
					B: color.B,
					A: uint8(255 * opacity / 100),
				},
			}
			graph.Elements = []chart.Renderable{
				chart.LegendLeft(graph, legendStyle),
			}
		} else {
			graph.Background = chart.Style{
				Padding: chart.Box{
					Top: b.Twiplet("legendOffset").Pixels(int32(b.dpi)),
				},
			}
			graph.Elements = []chart.Renderable{
				chart.LegendThin(graph),
			}
		}
	}
}

func (b *builder) buildSingleSeriesValues() (values []chart.Value) {
	for i, val := range b.data.Values[0] {
		num := b.ResolveNumber(val)
		values = append(
			values,
			chart.Value{
				Value: num,
				Label: b.singleSeriesLabel(i, num),
				Style: b.fontStyle(b.singleSeriesFont(i)),
			},
		)
	}
	return
}

func (b *builder) singleSeriesLabel(i int, val float64) string {
	label := ""
	if i < len(b.data.Labels) {
		label = b.data.Labels[i].Text()
	}
	format := b.data.Formats.CustomFormat(0, i).Text()
	if format != "" {
		label = strings.ReplaceAll(format, "{label}", label)
		label = strings.ReplaceAll(label, "{value}", b.valueFormatter(val))
	}
	return label
}

func (b *builder) singleSeriesFont(i int) *pic.FontStyle {
	if i < len(b.data.Fonts) {
		return b.ResolveFont(b.data.Fonts[i])
	}
	return b.ResolveFont(pic.DefaultFont)
}

func (b *builder) buildMultiSeries() (series []chart.Series) {
	for i, values := range b.data.Values {
		series = append(
			series,
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeWidth:     b.multiSeriesLineWidth(i),
					StrokeDashArray: b.multiSeriesLineDash(i),
				},
				Name:    b.multiSeriesName(i),
				XValues: b.buildXValues(values),
				YValues: b.buildYValues(values),
			},
		)
	}
	return
}

func (b *builder) multiSeriesLineWidth(i int) float64 {
	val := b.data.Styles.Setting(i, 0, "lineWidth")
	pixels := b.ResolveTwiplet(val).Pixels(int32(b.dpi))
	return float64(pixels)
}

func (b *builder) multiSeriesLineDash(i int) []float64 {
	val := b.data.Styles.Setting(i, 0, "lineStyle")
	if val != "dash" {
		return nil
	}
	return []float64{10, 5}
}

func (b *builder) multiSeriesName(i int) string {
	if i < len(b.data.Titles) {
		return b.data.Titles[i].Text()
	}
	return ""
}

func (b *builder) buildXTicks() (ticks []chart.Tick) {
	for i, label := range b.data.Labels {
		ticks = append(
			ticks,
			chart.Tick{
				Value: float64(i),
				Label: label.Text(),
			},
		)
	}
	return
}

func (b *builder) buildXValues(values []pic.Value) []float64 {
	xv := make([]float64, len(values))
	for i := range values {
		xv[i] = float64(i)
	}
	return xv
}

func (b *builder) buildYValues(values []pic.Value) []float64 {
	yv := make([]float64, len(values))
	for i, val := range values {
		yv[i] = b.ResolveNumber(val)
	}
	return yv
}

func (b *builder) valueFormatter(v interface{}) string {
	s := chart.FloatValueFormatterWithFormat(v, chart.DefaultFloatFormat)
	if b.DecimalPoint != '.' {
		s = strings.Replace(s, ".", string(b.DecimalPoint), 1)
	}
	return s
}

func (b *builder) fontStyle(fs *pic.FontStyle) chart.Style {
	return chart.Style{
		Font:     fs.TruetypeFont,
		FontSize: fs.PointSize,
		FontColor: drawing.Color{
			R: fs.R,
			G: fs.G,
			B: fs.B,
			A: 255,
		},
	}
}

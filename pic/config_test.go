package pic

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
)

type mockCallback struct {
	fontResources map[GUID]int
	fontStyles    map[GUID]int
}

func newMockCallback() *mockCallback {
	return &mockCallback{
		fontResources: make(map[GUID]int),
		fontStyles:    make(map[GUID]int),
	}
}

func (mockCallback) integer(s string) (int32, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	return int32(i), err
}

func (mockCallback) number(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func (mockCallback) numberFormat() NumberFormat {
	return NumberFormat{
		ThousandsSeparator: ',',
		DecimalPoint:       '.',
	}
}

func (mc *mockCallback) fontResource(guid GUID) (*FontResource, error) {
	if v, ok := mc.fontResources[guid]; ok {
		mc.fontResources[guid] = v + 1
	} else {
		mc.fontResources[guid] = 1
	}
	fr := &FontResource{
		Typeface:   "mockfont",
		Attributes: Bold | Italic | EmulateBold | EmulateItalic | EmulateTypeface,
	}
	return fr, nil
}

func (mc *mockCallback) fontStyle(guid GUID) (*FontStyle, error) {
	if v, ok := mc.fontStyles[guid]; ok {
		mc.fontStyles[guid] = v + 1
	} else {
		mc.fontStyles[guid] = 1
	}
	fs := &FontStyle{
		FontResource: &FontResource{
			Typeface: "mockstyle",
		},
	}
	return fs, nil
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	code := m.Run()
	os.Exit(code)
}

func TestTextValue(t *testing.T) {
	c := newConfig(newMockCallback(), "dog=canine\ncat=feline\r\nempty=\n", "")
	assertEqual(t, c.Value("dog").Text(), "canine")
	assertEqual(t, c.Value("cat").Text(), "feline")
	assertEqual(t, c.Value("empty").Text(), "")
}

func TestTrueValue(t *testing.T) {
	c := newConfig(newMockCallback(), "yes=true\nno=false\nnum=1", "")
	assertEqual(t, c.Value("yes").True(), true)
	assertEqual(t, c.Value("no").True(), false)
	assertEqual(t, c.Value("num").True(), false)
}

func TestSymbolValue(t *testing.T) {
	p := fmt.Sprintf("dog=%[1]csym1\ncat=%[1]csym2\ncurrency=%[1]csym3", ascDLE)
	s := fmt.Sprintf("sym1=canine\nsym3=%c$£8.99", ascESC)
	c := newConfig(newMockCallback(), p, s)
	assertEqual(t, c.Value("dog").Text(), "canine")
	assertEqual(t, c.Value("cat").Text(), "")
	assertEqual(t, c.Value("currency").Text(), "£8.99")
}

func TestIntegerValue(t *testing.T) {
	c := newConfig(newMockCallback(), "int=12345\ninvalid=foo", "")
	assertEqual(t, c.Integer("int"), int32(12345))
	assertEqual(t, c.Integer("invalid"), int32(0))
	assertEqual(t, c.Integer("missing"), int32(0))
}

func TestTwipletValue(t *testing.T) {
	c := newConfig(newMockCallback(), "x=14400\ny=-288000", "")
	x := c.Twiplet("x")
	y := c.Twiplet("y")
	assertEqual(t, x, Twiplet(14400))
	assertEqual(t, x.Pixels(300), 30)
	assertEqual(t, y, Twiplet(-288000))
	assertEqual(t, y.Pixels(300), -600)
}

func TestNumberValue(t *testing.T) {
	c := newConfig(newMockCallback(), "num=123.45\ninvalid=foo", "")
	assertEqual(t, c.Number("num"), 123.45)
	assertEqual(t, c.Number("invalid"), 0.0)
	assertEqual(t, c.Number("missing"), 0.0)
}

func TestNumberFormat(t *testing.T) {
	c := newConfig(newMockCallback(), "", "")
	nf := c.NumberFormat()
	assertEqual(t, nf.ThousandsSeparator, ',')
	assertEqual(t, nf.DecimalPoint, '.')
}

func TestValueType(t *testing.T) {
	p := fmt.Sprintf("int=%[1]ci42\nnum=%[1]cn3.14\ndate=%[1]cd31/12/1999\n"+
		"time=%[1]ct12:34:56\ncurr=%[1]c$9.99\nstr=txt\nfoo=%[1]c?bar", ascESC)
	c := newConfig(newMockCallback(), p, "")
	assertEqual(t, c.Value("int").Type(), Integer)
	assertEqual(t, c.Value("num").Type(), Number)
	assertEqual(t, c.Value("date").Type(), Date)
	assertEqual(t, c.Value("time").Type(), Time)
	assertEqual(t, c.Value("curr").Type(), Currency)
	assertEqual(t, c.Value("str").Type(), Neutral)
	assertEqual(t, c.Value("foo").Type(), NotSet)
}

func TestColorValue(t *testing.T) {
	c := newConfig(newMockCallback(), "color1=15\ncolor2=0,4,16711935,6553600\ninvalid1=foo\ninvalid2=99\ninvalid3=1,2,3", "")
	color1 := c.Color("color1")
	color2 := c.Color("color2")
	invalid1 := c.Color("invalid1")
	invalid2 := c.Color("invalid2")
	invalid3 := c.Color("invalid3")
	missing := c.Color("missing")
	assertEqual(t, color1.R, uint8(255))
	assertEqual(t, color1.G, uint8(255))
	assertEqual(t, color1.B, uint8(255))
	assertEqual(t, color2.R, uint8(255))
	assertEqual(t, color2.G, uint8(0))
	assertEqual(t, color2.B, uint8(255))
	assertEqual(t, color2.M, uint8(100))
	assertEqual(t, invalid1.K, uint8(100))
	assertEqual(t, invalid2.K, uint8(100))
	assertEqual(t, invalid3.K, uint8(100))
	assertEqual(t, missing.K, uint8(100))
}

func TestFontValue(t *testing.T) {
	p := fmt.Sprintf("font1=%cfCAFE000000000000000000000000F00D|0,0,0,100|0\nfont2=invalid", ascESC)
	c := newConfig(newMockCallback(), p, "")
	f := c.Font("font1")
	invalid := c.Font("font2")
	missing := c.Font("missing")
	assertEqual(t, f.IsStyle, false)
	assertEqual(t, f.GUID[0], byte(0xCA))
	assertEqual(t, f.GUID[1], byte(0xFE))
	assertEqual(t, f.GUID[14], byte(0xF0))
	assertEqual(t, f.GUID[15], byte(0x0D))
	assertEqual(t, invalid.GUID.IsZero(), true)
	assertEqual(t, invalid.GUID.IsZero(), true)
	assertEqual(t, missing.GUID.IsZero(), true)
}

func TestResolveFont(t *testing.T) {
	p := fmt.Sprintf("font=%cfCAFE000000000000000000000000F00D|0,0,0,100|0", ascESC)
	c := newConfig(newMockCallback(), p, "")
	fs := c.ResolveFont(c.Font("font"))
	assertEqual(t, fs.Typeface, "mockfont")
	assertEqual(t, fs.Bold(), true)
	assertEqual(t, fs.Italic(), true)
	assertEqual(t, fs.EmulateBold(), true)
	assertEqual(t, fs.EmulateItalic(), true)
	assertEqual(t, fs.EmulateTypeface(), true)
}

func TestFontResourceCache(t *testing.T) {
	p := fmt.Sprintf("font1=%[1]cfCAFE000000000000000000000000F00D|0,0,0,100|0\n"+
		"font2=%[1]cfCAFE000000000000000000000000F00D|0,1,0,100|0\n"+
		"font3=%[1]cfdead000000000000000000000000d0d0|0,0,0,100|0", ascESC)
	mc := newMockCallback()
	c := newConfig(mc, p, "")
	c.ResolveFont(c.Font("font1"))
	assertEqual(t, len(mc.fontResources), 1)
	c.ResolveFont(c.Font("font2"))
	assertEqual(t, len(mc.fontResources), 1)
	c.ResolveFont(c.Font("font3"))
	assertEqual(t, len(mc.fontResources), 2)
}

func TestFontStyleValue(t *testing.T) {
	p := fmt.Sprintf("style=%cf$CAFE000000000000000000000000F00D", ascESC)
	c := newConfig(newMockCallback(), p, "")
	f := c.Font("style")
	assertEqual(t, f.IsStyle, true)
	assertEqual(t, f.GUID[0], byte(0xCA))
	assertEqual(t, f.GUID[1], byte(0xFE))
	assertEqual(t, f.GUID[14], byte(0xF0))
	assertEqual(t, f.GUID[15], byte(0x0D))
}

func TestResolveFontStyle(t *testing.T) {
	p := fmt.Sprintf("style=%cf$CAFE000000000000000000000000F00D", ascESC)
	c := newConfig(newMockCallback(), p, "")
	f := c.Font("style")
	fs := c.ResolveFont(f)
	assertEqual(t, fs.FontResource.Typeface, "mockstyle")
}

func TestFontStyleCache(t *testing.T) {
	p := fmt.Sprintf("style1=%[1]cf$CAFE000000000000000000000000F00D\n"+
		"style2=%[1]cf$CAFE000000000000000000000000F00D\n"+
		"style3=%[1]cf$dead000000000000000000000000d0d0", ascESC)
	mc := newMockCallback()
	c := newConfig(mc, p, "")
	c.ResolveFont(c.Font("style1"))
	assertEqual(t, len(mc.fontStyles), 1)
	c.ResolveFont(c.Font("style2"))
	assertEqual(t, len(mc.fontStyles), 1)
	c.ResolveFont(c.Font("style3"))
	assertEqual(t, len(mc.fontStyles), 2)
}

func TestName(t *testing.T) {
	c := newConfig(newMockCallback(), "config=test", "")
	assertEqual(t, c.Name(), "test")
}

func TestData(t *testing.T) {
	p := `
data.values=4,2,3,4|2,4,1,3|8,5,4,5
data.titles=Series 1|Series 2|Series 3
data.colors=
data.styles=line:|line:|line:
data.labels=Category 1,Category 2,Category 3,Category 4
data.fonts=
data.formats=
`
	c := newConfig(newMockCallback(), p, "")
	d := c.Data()
	assertEqual(t, d.Values[0][0].Text(), "4")
	assertEqual(t, d.Values[0][1].Text(), "2")
	assertEqual(t, d.Values[0][2].Text(), "3")
	assertEqual(t, d.Values[0][3].Text(), "4")
}

func TestEmptyDataset(t *testing.T) {
	c := newConfig(newMockCallback(), "data=", "")
	d := c.Dataset("data")
	assertEqual(t, len(d), 1)
	assertEqual(t, len(d[0]), 1)
	assertEqual(t, d[0][0].Text(), "")
}

func TestSingleValueDataset(t *testing.T) {
	c := newConfig(newMockCallback(), "data=val", "")
	d := c.Dataset("data")
	assertEqual(t, len(d), 1)
	assertEqual(t, len(d[0]), 1)
	assertEqual(t, d[0][0].Text(), "val")
}

func TestSingleRowConstantDataset(t *testing.T) {
	c := newConfig(newMockCallback(), "data=val1,val2,val3,val4", "")
	d := c.Dataset("data")
	assertEqual(t, len(d), 1)
	assertEqual(t, len(d[0]), 4)
	assertEqual(t, d[0][0].Text(), "val1")
	assertEqual(t, d[0][1].Text(), "val2")
	assertEqual(t, d[0][2].Text(), "val3")
	assertEqual(t, d[0][3].Text(), "val4")
}

func TestSingleRowSymbolDataset(t *testing.T) {
	p := fmt.Sprintf("data=%c%csym1", ascSOH, ascDLE)
	s := fmt.Sprintf("sym1=val1%[1]cval2%[1]cval3%[1]cval4", ascUS)
	c := newConfig(newMockCallback(), p, s)
	d := c.Dataset("data")
	assertEqual(t, len(d), 1)
	assertEqual(t, len(d[0]), 4)
	assertEqual(t, d[0][0].Text(), "val1")
	assertEqual(t, d[0][1].Text(), "val2")
	assertEqual(t, d[0][2].Text(), "val3")
	assertEqual(t, d[0][3].Text(), "val4")
}

func TestSingleRowMixedDataset(t *testing.T) {
	p := fmt.Sprintf("data=%c%csym1%cval5", ascSOH, ascDLE, ascUS)
	s := fmt.Sprintf("sym1=val1%[1]cval2%[1]cval3%[1]cval4", ascUS)
	c := newConfig(newMockCallback(), p, s)
	d := c.Dataset("data")
	assertEqual(t, len(d), 1)
	assertEqual(t, len(d[0]), 5)
	assertEqual(t, d[0][0].Text(), "val1")
	assertEqual(t, d[0][1].Text(), "val2")
	assertEqual(t, d[0][2].Text(), "val3")
	assertEqual(t, d[0][3].Text(), "val4")
	assertEqual(t, d[0][4].Text(), "val5")
}

func TestDataTitles(t *testing.T) {
	p := fmt.Sprintf("data.titles=%[1]cSeries 1%[2]cSeries 2%[2]cSeries 3", ascSOH, ascRS)
	c := newConfig(newMockCallback(), p, "")
	titles := c.DataTitles()
	assertEqual(t, len(titles), 3)
	assertEqual(t, titles[0].Text(), "Series 1")
	assertEqual(t, titles[1].Text(), "Series 2")
	assertEqual(t, titles[2].Text(), "Series 3")
}

func TestDataLabels(t *testing.T) {
	p := fmt.Sprintf("data.labels=%[1]cCategory 1%[2]cCategory 2%[2]cCategory 3%[2]cCategory 4", ascSOH, ascUS)
	c := newConfig(newMockCallback(), p, "")
	labels := c.DataLabels()
	assertEqual(t, len(labels), 4)
	assertEqual(t, labels[0].Text(), "Category 1")
	assertEqual(t, labels[1].Text(), "Category 2")
	assertEqual(t, labels[2].Text(), "Category 3")
	assertEqual(t, labels[3].Text(), "Category 4")
}

func TestDataStyles(t *testing.T) {
	p := "data.styles=line:+style=solid+width=7200|line:+style=dash+width=3600"
	c := newConfig(newMockCallback(), p, "")
	ds := c.DataStyles()
	assertEqual(t, len(ds), 2)
	assertEqual(t, ds[0][0].Type, "line")
	assertEqual(t, ds.Setting(0, 0, "style").Text(), "solid")
	assertEqual(t, ds.Setting(0, 0, "width").Text(), "7200")
	assertEqual(t, ds[1][0].Type, "line")
	assertEqual(t, ds.Setting(1, 0, "style").Text(), "dash")
	assertEqual(t, ds.Setting(1, 0, "width").Text(), "3600")
}

func TestDataFormats(t *testing.T) {
	p := fmt.Sprintf("data.formats=default:%[1]ccustomFmt={value},custom:%[1]ccustomFmt={value}", ascSTX)
	c := newConfig(newMockCallback(), p, "")
	ds := c.DataFormats()
	assertEqual(t, len(ds), 1)
	assertEqual(t, len(ds[0]), 2)
	assertEqual(t, ds[0][0].Type, "default")
	assertEqual(t, ds.CustomFormat(0, 0).Text(), "")
	assertEqual(t, ds[0][1].Type, "custom")
	assertEqual(t, ds.CustomFormat(0, 1).Text(), "{value}")
}

func TestSingleSeriesDataColors(t *testing.T) {
	testDataColors(t, ascUS)
}

func TestMultiSeriesDataColors(t *testing.T) {
	testDataColors(t, ascRS)
}

func testDataColors(t *testing.T, sep byte) {
	p := fmt.Sprintf(
		"data.colors=%[1]c1,5,16724787,5263360%[2]c1,3,3407667,1342197760%[2]c1,11,6711039,1010565120%[2]c0,7,16776960,25600",
		ascSOH, sep,
	)
	c := newConfig(newMockCallback(), p, "")
	colors := c.DataColors()
	assertEqual(t, len(colors), 4)
	assertEqual(t, colors[0].R, uint8(255))
	assertEqual(t, colors[0].G, uint8(51))
	assertEqual(t, colors[0].B, uint8(51))
	assertEqual(t, colors[1].R, uint8(51))
	assertEqual(t, colors[1].G, uint8(255))
	assertEqual(t, colors[1].B, uint8(51))
	assertEqual(t, colors[2].R, uint8(102))
	assertEqual(t, colors[2].G, uint8(102))
	assertEqual(t, colors[2].B, uint8(255))
	assertEqual(t, colors[3].R, uint8(255))
	assertEqual(t, colors[3].G, uint8(255))
	assertEqual(t, colors[3].B, uint8(0))
}

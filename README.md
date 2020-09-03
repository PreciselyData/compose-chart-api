# Compose Chart API

[![LICENSE](https://img.shields.io/badge/license-MIT-green)](https://github.com/PreciselyData/compose-chart-api/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/PreciselyData/compose-chart-api/pic)](https://goreportcard.com/report/github.com/PreciselyData/compose-chart-api/pic)
[![GoDoc](https://godoc.org/github.com/PreciselyData/compose-chart-api/pic?status.svg)](https://godoc.org/github.com/PreciselyData/compose-chart-api/pic)

EngageOne Designer/Generate Plug-in Chart API for [Go](https://golang.org).

This repository contains code to help you create your own plug-in chart engine to use with Designer and Generate. Package [pic](https://github.com/PreciselyData/compose-chart-api/tree/master/pic) defines the API, and [example/go-chart](https://github.com/PreciselyData/compose-chart-api/tree/master/example/go-chart) provides an example implementation.

## Table of Contents

- [Building the example](#building-the-example)
  - [Install GCC for Windows](#install-gcc-for-windows)
  - [Go Build](#go-build)
- [Installing the example](#installing-the-example)
- [Installing on EngageOne Server](#installing-on-engageone-server)
  - [Interactive Editor setup file](#interactive-editor-setup-file)
  - [ActiveX cab file](#activex-cab-file)
- [How it works](#how-it-works)
  - [Configuration](#configuration)
  - [Property attributes](#property-attributes)
    - [Indent](#indent)
    - [Enable](#enable)
    - [Type](#type)
  - [Other elements](#other-elements)
    - [Data set](#data-set)
    - [Property group](#property-group)
- [Compatibility](#compatibility)

## Building the example

To use the example implementation in Designer you will need to build it as a 32-bit DLL. If you run 64-bit Generate on Windows, you will also need to build it as a 64-bit DLL. If your Generate platform is Linux then you will need to build it as a shared object (.so file).

Although Go supports cross-compilation for building executables (for example, building a Linux executable on a Windows machine), this is not practical when building shared libraries (.dll or .so files) as the C libraries for the target platform must be available on the build machine. In other words, you will need to build your `.so` file on Linux and your `.dll` file(s) on Windows. Note that you do not need to install Go to use the shared libraries, you just need to install Go to build them.

To build the example application and the `pic` package, you must have at least version 1.12 of Go installed. Visit https://golang.org/dl for download links and installation instructions.

To build a Go module as a shared library you also need to have GCC installed. If you need to install GCC on Windows, follow the instructions below. For Linux, GCC binaries are typically included as part of the distribution but may need to be installed using the package manager. Instructions to install GCC on Linux are specific to the distribution and are not covered here.

### Install GCC for Windows

GCC binary releases for Windows are available from [various websites](https://gcc.gnu.org/install/binaries.html), however only the `Cygwin` and `mingw-w64` projects offer both 32-bit and 64-bit installations. The following steps are for the `mingw-w64` installer which can be downloaded from [here](http://mingw-w64.org/doku.php/download/mingw-builds) by clicking on the `Sourceforge` link. This will redirect you to the sourceforge website where the download should start automatically.

To install both the 32-bit and 64-bit compilers you will need to run the installer twice.

For 32-bit, select the i686 architecture option.
For 64-bit, select the x86_64 architecture option.

Each installation will go into a separate folder - `C:\Program Files (x86)\mingw-w64\i686...`
for 32-bit and `C:\Program Files\mingw-w64\x86_64...` for 64-bit. You can change these directories
during installation if you wish.

To build a shared library or DLL, the Go compiler (via `cgo`) needs to know the location of GCC, so GCC must
be in your path. The `mingw-w64.bat` file in the 32-bit and 64-bit installation folders can be used to open a
command shell with the correct PATH set. Run the 32-bit version first to build the DLL for Designer, as described below.

### Go Build

From a command prompt, run the following to download the example application and `pic` package code.
```
go get -d -u github.com/PreciselyData/compose-chart-api/example/go-chart
```
The code should now be in the `%GOPATH%\src\github.com\PreciselyData\compose-chart-api` directory on Windows, or the `$(go env GOPATH)/src/github.com/PreciselyData/compose-chart-api` directory on Linux. Change to this directory in your command prompt.

To build a shared library, the environment variable `CGO_ENABLED` must be set to `1`. To build the 32-bit DLL, the environment variable `GOARCH` must be set to `386`. The default value for `CGO_ENABLED` is `1`, but changing the value of `GOARCH` also changes `CGO_ENABLED` to `0`, so make sure you set `CGO_ENABLED` back to `1` after setting `GOARCH` to build the 32-bit DLL on 64-bit Windows. You shouldn't need to set either of these environment variables otherwise, but you can run `go env` to check they are correct before building.

Change directory to `example/go-chart` and run the following to build the example:
```
go build -buildmode=c-shared -o go-chart.dll
```
Replace `.dll` with `.so` when building on Linux. If you omit the `-o` option, the output will be created without an extension.

## Installing the example

To install the example for use in Designer, copy `go-chart.dll` to the Designer client folder (where `cockpit.exe` resides).

You must also copy the contents of the `config` folder to the Designer client subfolder `propertytemplates\charts\en` (where `chtdir.cfg` resides).

If your installed Designer language is not English then you will also need to copy the template files to the appropriate `propertytemplates\charts\<language-id>` subfolder and change the `locale` attribute in `go-chart.xml` to match. In this case, you may also wish to localise the values of the `name` and `description` attributes in the `go-chart.xml` file, and the default values in the `go-chart.cfg` file.

To install the example for Generate, copy `go-chart.dll` (or `go-chart.so` on Linux) into the same folder as `doc1gen`.

## Installing on EngageOne Server

For Generate, copy your 64-bit `.dll` or `.so` file into the same folder as `doc1gen` in the batch and composition bundles. For the Interactive Editor, follow the steps below depending on your installation type.

### Interactive Editor setup file

If you install the Interactive Editor Application or the ActiveX Editor on your client machines using one of the `setup.exe` files, you will need to copy your 32-bit `.dll` to the installation folder after the installation is complete. See the EngageOne Server Installation Guide where it talks about deploying an `ApplicationStyle.xaml` file to change the visual style of the Interactive Editor Application. You will need to follow the same steps to deploy your `.dll`.

### ActiveX cab file

If you allow your clients to install the ActiveX Editor via Internet Explorer you will need to add your 32-bit `.dll` to the `eoeditor.cab` file. This means that you will also need to subsequently sign the cab file with your own certificate.

## How it works

Once installed, the example implementation can be seen working inside Designer from the Plug-in Chart dialog. You will see a new option in the Engine drop-down list called `Go-chart example`. When Designer loads the Plug-in Chart dialog it scans the `propertytemplates\charts\<language-id>` directory for XML files with a `propertyTemplate` root element. The value of the `name` attribute in this element is added to the engine list, and the value of the `id` element is the file name (minus extension) of the DLL to load in order to create the chart image. You can see these values in the [go-chart.xml](https://github.com/PreciselyData/compose-chart-api/blob/master/example/go-chart/config/go-chart.xml) file.

The DLL must export the functions `EnchCreateImage` and `EnchDestroyImage` otherwise an error will be shown in the dialog where the chart image is usually displayed. The DLL (or shared object on Linux) must also export the `EnchTerminate` function in order to work in Generate. These functions are exported by the `pic` package which takes full care of `EnchDestroyImage` and `EnchTerminate`. When `EnchCreateImage` is called, `pic` calls the interface method `Client.NewBuilder()` to create the image via the `Builder` interface. You can see how the example implements `pic.Builder` [here](https://github.com/PreciselyData/compose-chart-api/blob/master/example/go-chart/builder.go).

### Configuration

The chart properties are defined by an xml file (go-chart.xml for example), and the default values for the properties are supplied by a collection of cfg files. The xml and cfg files are tied together by the value of the `propertyTemplate.id` attribute in the xml file and the value of the `engine` property in the cfg files. Each chart engine can have multiple configurations, and each configuration has its own cfg file containing the default values. The format of a cfg file name is `<engine>-<id>.cfg`, for exmaple go-chart-pie.cfg. The `<engine>.cfg` file contains the default values common to each configuration. Each configuration is displayed in the Type drop-down list in the Plug-in Chart dialog.

The layout of the options in the dialog are also determined by the xml file. Each `category` element in the xml is represented on the left-hand side of the dialog, and each `property` element in the category is represented on the right-hand side. The categories are grouped together inside each `configuration` element. Categories that are shared by more than one configuration can be defined before the configuration definitions and referred to by the configurations using the `categoryRef` element. Alternatively, a category definition can be embedded inside a configuration definition.

When Designer/Generate calls the `EnchCreateImage` function, the configuration is supplied as a list of `property=value` settings where `property` is the `id` attribute of a `property` element in the xml. The `pic.Config` struct provides methods to fetch the values from the configuration in order to build the chart image.

### Property attributes

Each property can have a number of attributes which define its type, how it is displayed and when it should be enabled. These attributes are described below.

#### Indent

The `indent` attribute specifies the level of indentation displayed in the dialog. This attribute is used in conjunction with the `enable` attribute.

#### Enable

The `enable` attribute determines whether or not the property should be enabled in the dialog. Its value is a condition to determine whether a parent property is set to a particular value, for example `enable="legend=true"`. The `!` character can be used to negate the condition, for example `enable="legend=!true"`. A parent property is a property with a lower or no `indent` value.

#### Type

The `type` attribute of a `property` element in the xml defines which type of value can be entered into the field on the dialog. The available types are:

Type | Summary | Description
---- | ------- | -----------
`vp` | Value Picker | Allows the selection of a variable, data field, etc.
`fp` | Font Picker | Allows the selection of a font (including text color).
`cp` | Color Picker | Allows the selection of a color.
`mu` | Measurement Unit | Allows the configuration of a measurement.
`bool` | Boolean | Provides a drop-down list containing the options Yes and No.
`int` | Integer | Allows the entry of a number. The `min` and `max` attributes can be used to limit the range of the value.
`opt` | Option | Provides a drop-down list containing the options defined by each `option` child element.
`optSort` | Sorted Option | Same as `opt` except that the list will be sorted alphabetically.

### Other elements

#### Data set

You must define the `dataSet` element in your xml file. It can be empty (apart from the `id` and `name` attributes), in which case the data style options will not be available in the dialog. You must also reference the data using the `dataSetRef` element in each configuration.

Data styles are made available in the dialog by defining `property` elements with the`id="dataStyle"` attribute. This element must also contain `option` elements with `id` attributes that match the configuration to which they apply. Other properties within the data style can be enabled using the `enable="dataStyle=<option-id>"` attribute. A number of options can be listed in the `enable` value by separating each option with the `|` character.

To tie a configuration to a data style, include a `variant` element with an `id` attribute matching that of an `option` element in the `dataStyle` definition.

#### Property group

The `propertyGroup` element is used to group `property` elements together. Each `id` of a property group must be unique for all chart engines. A property group is referenced using the `propertyGroupRef` element. The `prefix` attribute is used to create a unique name for each property in the referenced group when saved to the configuration for `EnchCreateImage`. If a property is not required for a particular reference it can be removed with the `remove=<property-id>` attribute. To remove more than one property, separate each `id` with a comma.

## Compatibility

This API was published to coincide with the release of Designer/Generate 6.6 SP10 and is therefore compatible with version 6.6 SP10 and later. The API should also be compatible with previous releases of Designer/Generate version 6, but this has not been tested. The following known issues exist with versions of Designer/Generate prior to 6.6 SP10.
1. If the chart engine implementation does not export the required functions, Designer may crash or report a meaningless error. Designer 6.6 SP10 reports a meaningful error.
2. Modifying an existing chart and changing the engine setting results in errors on the data page. This option is disabled in Designer 6.6 SP10 when modifying a chart.
3. 64-bit Generate for Windows does not support PNG images returned by the chart engine. This issue has been resolved in Generate 6.6 SP10.
// Copyright 2019 Pitney Bowes Inc.

// Package pic defines the Plug-in Chart API for EngageOne® Designer/Generate.
// It provides the export functions called by Designer/Generate to create a
// chart image, and functions to help parse the chart configuration settings.
//
// To use this package you must implement the pic.Client and pic.Builder
// interfaces. Tell pic about your Client implementation by calling the
// pic.SetClient() method in the init() function of your main.go file.
// When Designer/Generate needs to create a chart image, pic will call
// your Client.NewBuilder() method to create the image via your Builder
// implementation.
//
// You must build your application with the -buildmode=c-shared option in
// order to create a .dll file that can be loaded by Designer and Generate.
// If you run Generate on Linux you will also need to build a .so file.
package pic

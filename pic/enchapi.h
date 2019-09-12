//============================================================================
// File:        enchapi.h
// Synopsis:    Enhanced Chart API definitions.
//============================================================================
#ifndef fEnchApiIncluded
#define fEnchApiIncluded

#include <wchar.h>

// Enhanced Chart API (ENCHAPI) is supported on the following platforms -
//   Windows
//   Linux (Intel 64-bit)
//   Z/Linux
//   Solaris (Sparc)
//   AIX (RS/6000, pSeries)
//   HP-UX IA-64 (Itanium), but not HP-UX PA-RISC
//   Tru64 UNIX
#if defined(WIN32) || \
   defined(LNX64) || \
   defined(LNXZLINUX) || \
   defined(SPARC) || \
   defined(RS6000) || \
   (defined(HPUX11) && defined(__ia64)) || \
   defined(TRU64)
#define ENCHAPI_SUPPORTED
#endif

// Set alignment consistent with DOC1.
#if defined(RS6000)
#pragma options align=packed
#endif

// API data types.
typedef unsigned char UTF8CHAR;    // UTF-8 encoded character

// API return codes.
typedef enum
{
   ENCHRC_OK                   = 0,
   ENCHRC_Failed               = 1,
   ENCHRC_NotImplemented       = 2,
   ENCHRC_InvalidFilePath      = 3,
   ENCHRC_MissingProperty      = 4,
   ENCHRC_InvalidValue         = 5,
   ENCHRC_UnresolvedFont       = 6,
   ENCHRC_InvalidDataString    = 7,
   ENCHRC_EmptyDataString      = 8,
   ENCHRC_JavaException        = 9,
   ENCHRC_FailedToCreateJavaVM = 10
   // Please update enchrcStrings in commonlib when adding new error codes.
} ENCHRC;

// Character buffer sizes.
#define ENCH_FileNameSize        300
#define ENCH_TypefaceSize        128

// Font attribute flags.
#define ENCH_FontBold            0x0001
#define ENCH_FontItalic          0x0002
#define ENCH_FontEmulateBold     0x0004
#define ENCH_FontEmulateItalic   0x0008
#define ENCH_FontEmulateTypeface 0x0010

// Character style flags.
#define ENCH_StyleUnderline      0x0001

// Colour types.
#define ENCH_ColorNamed          0
#define ENCH_ColorRgb            1
#define ENCH_ColorCmyk           2

// Data types.
#define ENCH_DataNotSet          (-1)
#define ENCH_DataNeutral         0
#define ENCH_DataInteger         1
#define ENCH_DataNumber          2
#define ENCH_DataDate            3
#define ENCH_DataTime            4
#define ENCH_DataCurrency        5

// Image formats.
#define ENCH_ImageBmp            0
#define ENCH_ImagePng            1
#define ENCH_ImageJpg            2
#define ENCH_ImageSvg            3

// Colour structure.
typedef struct tagEnchColor
{
   int master;                          // Master colour definition; one of ENCH_Color*
   int named;                           // Named colour value 0-15
   int red, green, blue;                // RGB value
   int cyan, magenta, yellow, keyBlack; // CMYK value
} EnchColor;

// Font resource structure.
typedef struct tagEnchFontResource
{
   wchar_t szFileName[ENCH_FileNameSize];
   wchar_t szTypeface[ENCH_TypefaceSize];
   int nDeciPointSize;
   unsigned short fsFlags; // combination of ENCH_Font* flags
} EnchFontResource;

// Style resource structure.
typedef struct tagEnchStyleResource
{
   EnchFontResource fontResource;
   EnchColor color;
   unsigned short fsFlags; // combination of ENCH_Style* flags
} EnchStyleResource;

// Date value structure.
typedef struct tagEnchDate
{
   int nDay;
   int nMonth;
   int nYear;
} EnchDate;

// Time value structure.
typedef struct tagEnchTime
{
   int nHour;
   int nMinute;
   int nSecond;
} EnchTime;

// Data value structure.
typedef struct tagEnchDataValue
{
   int nType; // Data type; one of ENCH_Data*
   union
   {
      int iValue;          // ENCH_DataInteger
      double dValue;       // ENCH_DataNumber or ENCH_DataCurrency
      EnchDate dateValue;  // ENCH_DataDate
      EnchTime timeValue;  // ENCH_DataTime
   } data;
} EnchDataValue;

// Number format structure.
typedef struct tagEnchNumberFormat
{
   wchar_t chThousandsSeparator;
   wchar_t chDecimalPoint;
} EnchNumberFormat;

// Date and time format structure.
typedef struct tagEnchDateTimeFormat
{
   wchar_t** ppszMonthNames;
   int cMonthNames;
   wchar_t** ppszWeekDayNames;
   int cWeekDayNames;
   wchar_t* pszAm;
   wchar_t* pszPm;
   wchar_t* pszShortDateFormat;
   wchar_t* pszLongDateFormat;
   wchar_t* pszTimeFormat;
} EnchDateTimeFormat;

// Structure used to identify the chart when calling back into DOC1.
typedef struct tagENCHSHELLDATA *PENCHSHELLDATA;

// Callback structure defining the functions for calling back into DOC1.
typedef struct tagEnchCallback EnchCallback;
struct tagEnchCallback
{
   PENCHSHELLDATA pEnchShellData;
   int (*pfnGetFont)(EnchCallback* pCallback, const unsigned char* pGuid, EnchFontResource* pFontResource);
   int (*pfnGetStyle)(EnchCallback* pCallback, const unsigned char* pGuid, EnchStyleResource* pStyleResource);
   ENCHRC (*pfnGetDataValue)(EnchCallback* pCallback, const wchar_t* pszValue, EnchDataValue* pDataValue, int nExpectedType);
   int (*pfnGetNumberFormat)(EnchCallback* pCallback, EnchNumberFormat* pNumberFormat);
   int (*pfnGetDateTimeFormat)(EnchCallback* pCallback, EnchDateTimeFormat* pDateTimeFormat);
   void (*pfnFreeDateTimeFormat)(EnchDateTimeFormat* pDateTimeFormat);
};

// Structure identifying the chart image.
typedef struct tagEnchImage
{
   int nFormat;                // Preferred image format; one of ENCH_Image*
   int nColorSpace;            // Preferred image colour space; ENCH_ColorRgb or ENCH_ColorCmyk.
   int nLogWidth;              // Image width in logical units (1/100ths of a twip).
   int nLogHeight;             // Image height in logical units (1/100ths of a twip).
   int nRes;                   // Image resolution.

   // The following fields are set by the interface when EnchCreateImage is called and
   // are used by the interface to destroy the data when EnchDestroyImage is called.

   void* pvImageHandle;        // Interface-specific image handle.
   unsigned char* pbImageData; // Physical image data.
   unsigned int cbImageData;   // Number of bytes in the image data.
} EnchImage;

//============================================================================
// Exported function definitions.
//============================================================================

#if defined (__cplusplus)
extern "C"
{
#endif

// EnchCreateImage is called by Designer/Generate to create a chart image
// from a list of configuration settings (pszProperties). These settings
// appear in the form name=value where value may be a constant, or refer
// to a field in the symbol table (pszSymbolTable). Image information is
// conveyed via pImage. The pCallback parameter contains pointers to
// functions defined by Designer/Generate to resolve data values, fonts
// and locale information.
#ifdef ENCH_IMPORT
typedef ENCHRC (*PFNEnchCreateImage)
#else
ENCHRC EnchCreateImage
#endif
(
   EnchCallback* pCallback,        // Pointer to DOC1 callback structure
   const UTF8CHAR* pszProperties,  // List of property=value configuration settings
   const UTF8CHAR* pszSymbolTable, // List of variable data symbols and their values
   EnchImage* pImage               // Pointer to the chart image details
);

// EnchDestroyImage is called by Designer/Generate to destroy the chart image
// data created by EnchCreateImage.
#ifdef ENCH_IMPORT
typedef ENCHRC (*PFNEnchDestroyImage)
#else
ENCHRC EnchDestroyImage
#endif
(
   EnchImage* pImage               // Pointer to the chart image details
);

// EnchTerminate is called by Generate to tidy up before the program exits.
// Return ENCHRC_OK for Generate to unload the module, or ENCHRC_Failed to
// keep the module loaded.
#ifdef ENCH_IMPORT
typedef ENCHRC (*PFNEnchTerminate)
#else
ENCHRC EnchTerminate
#endif
(
);

#if defined(__cplusplus)
}
#endif

#endif

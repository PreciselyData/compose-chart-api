#include <stdlib.h>
#include <string.h>
#include "enchapi.h"

wchar_t* Utf8ToWideChar(const char* pszUtf8);
char* Utf8FromWideChar(const wchar_t* pwsz);

typedef struct tagEnchFontResourceUtf8
{
   char* pszFileName;
   char* pszTypeface;
   int nDeciPointSize;
   unsigned short fsFlags; // combination of ENCH_Font* flags
} EnchFontResourceUtf8;

typedef struct tagEnchStyleResourceUtf8
{
   EnchFontResourceUtf8 fontResource;
   EnchColor color;
   unsigned short fsFlags; // combination of ENCH_Style* flags
} EnchStyleResourceUtf8;

ENCHRC EnchGetDataValue(void* pvCallback, const char* pszValue, EnchDataValue* pDataValue, int nExpectedType)
{
   EnchCallback* pCallback = (EnchCallback*)pvCallback;
   wchar_t* pwsz = Utf8ToWideChar(pszValue);
   ENCHRC rc = pCallback->pfnGetDataValue(pCallback, pwsz, pDataValue, nExpectedType);
   free(pwsz);
   return rc;
}

ENCHRC EnchGetInteger(void* pvCallback, const char* pszValue, int* pValue)
{
   EnchDataValue dataValue;
   ENCHRC rc = EnchGetDataValue(pvCallback, pszValue, &dataValue, ENCH_DataInteger);
   if (rc != ENCHRC_OK)
   {
      return rc;
   }
   switch (dataValue.nType)
   {
   case ENCH_DataInteger:
      *pValue = dataValue.data.iValue;
      break;
   case ENCH_DataNumber:
   case ENCH_DataCurrency:
      *pValue = (int)dataValue.data.dValue;
      break;
   default:
      rc = ENCHRC_InvalidValue;
      break;
   }
   return rc;
}

ENCHRC EnchGetNumber(void* pvCallback, const char* pszValue, double* pValue)
{
   EnchDataValue dataValue;
   ENCHRC rc = EnchGetDataValue(pvCallback, pszValue, &dataValue, ENCH_DataNumber);
   if (rc != ENCHRC_OK)
   {
      return rc;
   }
   switch (dataValue.nType)
   {
   case ENCH_DataInteger:
      *pValue = dataValue.data.iValue;
      break;
   case ENCH_DataNumber:
   case ENCH_DataCurrency:
      *pValue = dataValue.data.dValue;
      break;
   default:
      rc = ENCHRC_InvalidValue;
      break;
   }
   return rc;
}

int EnchGetNumberFormat(void* pvCallback, EnchNumberFormat* pNumberFormat)
{
   EnchCallback* pCallback = (EnchCallback*)pvCallback;
   return pCallback->pfnGetNumberFormat(pCallback, pNumberFormat);
}

int EnchGetFont(void* pvCallback, const unsigned char* pGuid, EnchFontResourceUtf8* pFontResource)
{
   EnchCallback* pCallback = (EnchCallback*)pvCallback;
   EnchFontResource fr;
   int rc = pCallback->pfnGetFont(pCallback, pGuid, &fr);
   if (rc)
   {
      pFontResource->pszFileName = Utf8FromWideChar(fr.szFileName);
      pFontResource->pszTypeface = Utf8FromWideChar(fr.szTypeface);
      pFontResource->nDeciPointSize = fr.nDeciPointSize;
      pFontResource->fsFlags = fr.fsFlags;
   }
   return rc;
}

void EnchFreeFont(EnchFontResourceUtf8* pFontResource)
{
   free(pFontResource->pszFileName);
   free(pFontResource->pszTypeface);
}

int EnchGetStyle(void* pvCallback, const unsigned char* pGuid, EnchStyleResourceUtf8* pStyleResource)
{
   EnchCallback* pCallback = (EnchCallback*)pvCallback;
   EnchStyleResource sr;
   int rc = pCallback->pfnGetStyle(pCallback, pGuid, &sr);
   if (rc)
   {
      pStyleResource->fontResource.pszFileName = Utf8FromWideChar(sr.fontResource.szFileName);
      pStyleResource->fontResource.pszTypeface = Utf8FromWideChar(sr.fontResource.szTypeface);
      pStyleResource->fontResource.nDeciPointSize = sr.fontResource.nDeciPointSize;
      pStyleResource->fontResource.fsFlags = sr.fontResource.fsFlags;
      pStyleResource->color = sr.color;
      pStyleResource->fsFlags = sr.fsFlags;
   }
   return rc;
}

wchar_t* Utf8ToWideChar(const char* pszUtf8)
{
   // by looking at the first byte determine how many bytes need reading
   // by using this table
   static UTF8CHAR abFirstByte[256] =
   {
      0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
      0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
      0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
      0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
      0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
      0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
      1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1, 1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,
      2,2,2,2,2,2,2,2,2,2,2,2,2,2,2,2, 3,3,3,3,3,3,3,3,4,4,4,4,5,5,5,5
   };

   static unsigned int auiOffsetsFromUTF8[6] =
   {
      0x00000000, 0x00003080, 0x000E2080,
      0x03C82080, 0xFA082080, 0x82082080
   };

   const UTF8CHAR* pch = (const UTF8CHAR*)pszUtf8;
   int cch = strlen(pszUtf8);
   wchar_t* pwc = (wchar_t*)malloc(((2 * cch) + 1) * sizeof(wchar_t));
   int iCountOut = 0;
   int cBytesToRead;
   unsigned int uiUnicode;

   // For each char in the input
   for (int iCounter = 0; iCounter < cch; iCounter++)
   {
      if ( ( *pch == 0xEF ) && ( *(pch+1) == 0xBB )  && ( *(pch+2) == 0xBF ) ) // ignore byte order marks
      {
         pch += 3;       // Skip the BOM
         iCounter += 2;  // Increment the count of chars processes by 2
                         // It's already been incremented by 1 in the "for " statement above
      }
      else
      {
         uiUnicode = 0;

         // According to the number of "1"s at the start of the byte
         cBytesToRead = (int)abFirstByte[*pch];

         switch (cBytesToRead)
         {
         case 5:
            uiUnicode += *pch++; uiUnicode <<= 6;
         case 4:
            uiUnicode += *pch++; uiUnicode <<= 6;
         case 3:
            uiUnicode += *pch++; uiUnicode <<= 6;
         case 2:
            uiUnicode += *pch++; uiUnicode <<= 6;
         case 1:
            uiUnicode += *pch++; uiUnicode <<= 6;
         case 0:
            uiUnicode += *pch++;
         };
         uiUnicode -= auiOffsetsFromUTF8[cBytesToRead];

         // Skip on th enumber of EXTRA bytes read
         iCounter += cBytesToRead;

         if (uiUnicode <= 0x0000FFFF)
         {
            // unicode will fit in one TCHAR
            pwc[iCountOut++] = (wchar_t)uiUnicode;
         }
         else if (uiUnicode > 0x0010FFFF)
         {
            // This is an error as unicode cant be more than this
            pwc[iCountOut++] = 0xFFFD; // put in the undefined
         }
         else
         {
            // put the char in two bytes
            uiUnicode -= 0x0010000;
            pwc[iCountOut++] = (wchar_t)((uiUnicode >> 10) + 0xD800);
            pwc[iCountOut++] = (wchar_t)((uiUnicode & 0x3FF) + 0xDC00);
         }
      }
   }

   pwc[iCountOut] = 0;
   return pwc;
}

char* Utf8FromWideChar(const wchar_t* pwsz)
{
   const wchar_t* pwch = pwsz;
   int            cwch = wcslen(pwsz);
   char*          pszUtf8 = (char*)malloc((6 * cwch) + 1);
   UTF8CHAR*      pbOut = (UTF8CHAR*)pszUtf8;
   int            iCountOut = 0;
   int            cBytesToWrite = 0;
   unsigned int   uiUnicode;
   unsigned short usUnicodeNext;
   UTF8CHAR       abByteMark[7] = {0x00, 0x00, 0xC0, 0xE0, 0xF0, 0xF8, 0xFC};

   for (int iCounter = 0; iCounter < cwch; iCounter++)
   {
      uiUnicode = (*pwch & 0xFFFF);

      // If this is a surrogate pair
      if ( ( uiUnicode >= 0xD800 ) && ( uiUnicode <= 0xDBFF ) )
      {
         // Get the next char and make sure its the second byte of a pair
         usUnicodeNext = (*(pwch+1) & 0xFFFF);

         if ( ( usUnicodeNext >= 0xDC00 ) && ( usUnicodeNext <= 0xDFFF ) )
         {
            // make up the actual value
            // move the char map on...
            iCounter++;
            pwch++;

            // Some bit magic
            uiUnicode = ((uiUnicode - 0xD800) << 10)
               + (usUnicodeNext - 0xDC00) + 0x0010000;
         }
      }

      // Convert the unicode to UTF8, determine mow many bytes to write
      if (uiUnicode < 0x80)
      {
         cBytesToWrite = 1;
      }
      else if (uiUnicode < 0x800)
      {
         cBytesToWrite = 2;
      }
      else if (uiUnicode < 0x10000)
      {
         cBytesToWrite = 3;
      }
      else if (uiUnicode < 0x200000)
      {
         cBytesToWrite = 4;
      }
      else if (uiUnicode < 0x4000000)
      {
         cBytesToWrite = 5;
      }
      else if (uiUnicode <= 0x7FFFFFFF)  // Largest UCS-4 value
      {
         cBytesToWrite = 6;
      }
      else
      {
         cBytesToWrite = 2;
         uiUnicode = 0xFFFD;
      };

      // Work backwards through the output bytes, see RFC2279
      pbOut += cBytesToWrite;

      switch (cBytesToWrite)
      {
      case 6:
         *--pbOut = (UTF8CHAR)((uiUnicode | 0x80) & 0xBF); uiUnicode >>= 6;
      case 5:
         *--pbOut = (UTF8CHAR)((uiUnicode | 0x80) & 0xBF); uiUnicode >>= 6;
      case 4:
         *--pbOut = (UTF8CHAR)((uiUnicode | 0x80) & 0xBF); uiUnicode >>= 6;
      case 3:
         *--pbOut = (UTF8CHAR)((uiUnicode | 0x80) & 0xBF); uiUnicode >>= 6;
      case 2:
         *--pbOut = (UTF8CHAR)((uiUnicode | 0x80) & 0xBF); uiUnicode >>= 6;
      case 1:
         *--pbOut = (UTF8CHAR)(uiUnicode | abByteMark[cBytesToWrite]);
      };

      pbOut += cBytesToWrite;
      iCountOut += cBytesToWrite;
      pwch++;
   }

   pszUtf8[iCountOut] = 0;
   return pszUtf8;
}

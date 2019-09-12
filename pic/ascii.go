package pic

// ASCII control codes used in the configuration settings.
const (
	ascSOH byte = 0x01 // Start Of Heading used to begin a dataset
	ascSTX byte = 0x02 // Start of Text used to delimit data style properties
	ascDLE byte = 0x10 // Data Link Escape used to denote a data field register
	ascESC byte = 0x1b // Escape used to denote a special property type
	ascRS  byte = 0x1e // Record Separator used to delimit a dataset
	ascUS  byte = 0x1f // Unit Separator used to delimit values in a dataset
)

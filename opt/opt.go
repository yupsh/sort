package opt

// Custom types for parameters
type Field int
type Delimiter string

// Boolean flag types with constants
type ReverseFlag bool

const (
	Reverse   ReverseFlag = true
	NoReverse ReverseFlag = false
)

type NumericFlag bool

const (
	Numeric   NumericFlag = true
	NoNumeric NumericFlag = false
)

type UniqueFlag bool

const (
	Unique   UniqueFlag = true
	NoUnique UniqueFlag = false
)

type IgnoreCaseFlag bool

const (
	IgnoreCase    IgnoreCaseFlag = true
	CaseSensitive IgnoreCaseFlag = false
)

type RandomFlag bool

const (
	Random   RandomFlag = true
	NoRandom RandomFlag = false
)

type IgnoreLeadingBlanksFlag bool

const (
	IgnoreLeadingBlanks   IgnoreLeadingBlanksFlag = true
	NoIgnoreLeadingBlanks IgnoreLeadingBlanksFlag = false
)

type VersionSortFlag bool

const (
	VersionSort   VersionSortFlag = true
	NoVersionSort VersionSortFlag = false
)

type HumanNumericFlag bool

const (
	HumanNumeric   HumanNumericFlag = true
	NoHumanNumeric HumanNumericFlag = false
)

type MonthSortFlag bool

const (
	MonthSort   MonthSortFlag = true
	NoMonthSort MonthSortFlag = false
)

type StableSortFlag bool

const (
	StableSort   StableSortFlag = true
	NoStableSort StableSortFlag = false
)

// Flags represents the configuration options for the sort command
type Flags struct {
	Reverse             ReverseFlag             // Reverse sort order (-r)
	Numeric             NumericFlag             // Numeric sort (-n)
	Unique              UniqueFlag              // Remove duplicates (-u)
	IgnoreCase          IgnoreCaseFlag          // Case insensitive sort (-f)
	Field               Field                   // Sort by specific field (-k)
	Delimiter           Delimiter               // Field delimiter (-t)
	Random              RandomFlag              // Random shuffle (-R)
	IgnoreLeadingBlanks IgnoreLeadingBlanksFlag // Ignore leading blanks (-b)
	VersionSort         VersionSortFlag         // Version sort (-V)
	HumanNumeric        HumanNumericFlag        // Human numeric sort (-h)
	MonthSort           MonthSortFlag           // Month sort (-M)
	StableSort          StableSortFlag          // Stable sort (-s)
}

// Configure methods for the opt system
func (r ReverseFlag) Configure(flags *Flags)             { flags.Reverse = r }
func (n NumericFlag) Configure(flags *Flags)             { flags.Numeric = n }
func (u UniqueFlag) Configure(flags *Flags)              { flags.Unique = u }
func (i IgnoreCaseFlag) Configure(flags *Flags)          { flags.IgnoreCase = i }
func (f Field) Configure(flags *Flags)                   { flags.Field = f }
func (d Delimiter) Configure(flags *Flags)               { flags.Delimiter = d }
func (r RandomFlag) Configure(flags *Flags)              { flags.Random = r }
func (i IgnoreLeadingBlanksFlag) Configure(flags *Flags) { flags.IgnoreLeadingBlanks = i }
func (v VersionSortFlag) Configure(flags *Flags)         { flags.VersionSort = v }
func (h HumanNumericFlag) Configure(flags *Flags)        { flags.HumanNumeric = h }
func (m MonthSortFlag) Configure(flags *Flags)           { flags.MonthSort = m }
func (s StableSortFlag) Configure(flags *Flags)          { flags.StableSort = s }

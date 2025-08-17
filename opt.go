package command

type Field int
type Delimiter string

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

type flags struct {
	Reverse             ReverseFlag
	Numeric             NumericFlag
	Unique              UniqueFlag
	IgnoreCase          IgnoreCaseFlag
	Field               Field
	Delimiter           Delimiter
	Random              RandomFlag
	IgnoreLeadingBlanks IgnoreLeadingBlanksFlag
	VersionSort         VersionSortFlag
	HumanNumeric        HumanNumericFlag
	MonthSort           MonthSortFlag
	StableSort          StableSortFlag
}

func (r ReverseFlag) Configure(flags *flags)             { flags.Reverse = r }
func (n NumericFlag) Configure(flags *flags)             { flags.Numeric = n }
func (u UniqueFlag) Configure(flags *flags)              { flags.Unique = u }
func (i IgnoreCaseFlag) Configure(flags *flags)          { flags.IgnoreCase = i }
func (f Field) Configure(flags *flags)                   { flags.Field = f }
func (d Delimiter) Configure(flags *flags)               { flags.Delimiter = d }
func (r RandomFlag) Configure(flags *flags)              { flags.Random = r }
func (i IgnoreLeadingBlanksFlag) Configure(flags *flags) { flags.IgnoreLeadingBlanks = i }
func (v VersionSortFlag) Configure(flags *flags)         { flags.VersionSort = v }
func (h HumanNumericFlag) Configure(flags *flags)        { flags.HumanNumeric = h }
func (m MonthSortFlag) Configure(flags *flags)           { flags.MonthSort = m }
func (s StableSortFlag) Configure(flags *flags)          { flags.StableSort = s }

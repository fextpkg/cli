package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	compareVersionTrue = [][3]string{
		{"1.0.2", ">", "1.0.0"},
		{"1.0.0", "==", "1.0.0"},
		{"1.1.0", "==", "1.1"},
		{"1.1", "==", "1.1"},
		{"1", "==", "1"},
		{"2.0.1", ">=", "2.0.0"},
		{"2.0.0", ">=", "2.0.0"},
		{"2.0.0", "<", "2.0.1"},
		{"2.0.0", "<=", "2.0.0"},
		{"2.0.0", "!=", "2.0.1"},
		{"2.0.1", "!=", "2.0.0"},
		{"2", "!=", "3"},
		{"2", "<", "3"},
		{"3", ">", "2"},
		{"3", ">", "2"},
		{"2.0.1a1", "!=", "2.0.1a2"},
		{"2.0.1a1", "<", "2.0.1a2"},
		{"2.0.1a1", "<=", "2.0.1a1"},
		{"2.0.1b4", ">", "2.0.1a1"},
		{"2.0.1b4", "==", "2.0.1b4"},
		{"2b4", "==", "2b4"},
		{"2.1b4", ">", "2.1b2"},
		{"2b4", "==", "2b4"},
		{"2b4", ">", "2b2"},
		{"2b0", "<", "2b2"},
		{"2.1", "~=", "2.0"},
		{"2.0.1", "~=", "2.0.0"},
		{"2.2.1", "~=", "2.0.0"},
		{"2.2.1a3", "~=", "2.0.0"},
		{"2.2a3", "~=", "2.0.0"},
		{"2a3", "~=", "2a0"},
	}
	compareVersionFalse = [][3]string{
		{"1.0.0", "!=", "1.0.0"},
		{"1.0.0", ">", "1.0.0"},
		{"1.0.0", "<", "1.0.0"},
		{"1.0.0", ">=", "1.0.2"},
		{"1.0.0", ">", "1.0.2"},
		{"1.0.2", "<=", "1.0.0"},
		{"1.0.2", "<", "1.0.0"},
		{"1.0", ">", "1.0.0"},
		{"1.0", "<", "1.0.0"},
		{"1.0", "==", "1.0.2"},
		{"1.0", "==", "1.0.2"},
		{"1", "==", "1.0.2"},
		{"1", ">", "1.0.2"},
		{"1", "<", "1.0.0"},
		{"1.0.0b4", "==", "1.0.2"},
		{"1.0.0b4", "==", "1.0.2b3"},
		{"1.0.0b4", "<", "1.0.0b3"},
		{"1.0.0b3", ">", "1.0.2b3"},
		{"1b3", ">", "1.0.2b3"},
		{"2.0", "~=", "2.0.1"},
		{"3.0", "~=", "2.0.1"},
		{"1.0", "~=", "2.0.1"},
		{"2.0.0", "~=", "2.0.1"},
		{"2.45.0", "~=", "2.46.1"},
		{"2a3", "~=", "2a4"},
	}
	compareVersionInvalid = [][3]string{
		{"test", "==", "test2"},
		{"test", "==", "2.2"},
		{"1.1", "==", "test"},
		{"1.1", "asd", "2.2"},
	}
)

func TestParePreVersion(t *testing.T) {
	patchVersion, preVersion, err := parsePreVersion("1a2")
	assert.Nil(t, err)
	assert.Equal(t, preVersion, 147)
	assert.Equal(t, patchVersion, 1)

	patchVersion, preVersion, err = parsePreVersion("1a")
	assert.Nil(t, err)
	assert.Equal(t, preVersion, 97)
	assert.Equal(t, patchVersion, 1)

	patchVersion, preVersion, err = parsePreVersion("1")
	assert.Nil(t, err)
	assert.Equal(t, preVersion, 0)
	assert.Equal(t, patchVersion, 0)

	patchVersion, preVersion, err = parsePreVersion("")
	assert.Nil(t, err)
	assert.Equal(t, preVersion, 0)
	assert.Equal(t, patchVersion, 0)

	patchVersion, preVersion, err = parsePreVersion("a1")
	assert.NotNil(t, err)
}

func TestGetStrIndexSum(t *testing.T) {
	sum := getStringIndexSum("a")
	assert.Equal(t, sum, 97)

	sum = getStringIndexSum("aa")
	assert.Equal(t, sum, 97*2)

	sum = getStringIndexSum("a1")
	assert.Equal(t, sum, 146)
}

func TestParseVersion(t *testing.T) {
	version, preVersion, err := parseVersion("1.2.3")
	assert.Nil(t, err)
	assert.Zero(t, preVersion)
	assert.Equal(t, version, [3]int{1, 2, 3})

	version, preVersion, err = parseVersion("1.2.3a2")
	assert.Nil(t, err)
	// 147 = 97 (a) + 50 (2) = 147 (a2)
	assert.Equal(t, preVersion, 147)
	assert.Equal(t, version, [3]int{1, 2, 3})

	version, preVersion, err = parseVersion("1.2a3")
	assert.Nil(t, err)
	// 148 = 97 (a) + 51 (3) = 148 (a3)
	assert.Equal(t, preVersion, 148)
	assert.Equal(t, version, [3]int{1, 2, 0})

	version, preVersion, err = parseVersion("1")
	assert.Nil(t, err)
	assert.Zero(t, preVersion)
	assert.Equal(t, version, [3]int{1, 0, 0})

	version, preVersion, err = parseVersion("")
	assert.Nil(t, err)
	assert.Zero(t, preVersion)
	assert.Equal(t, version, [3]int{0, 0, 0})

	version, preVersion, err = parseVersion("0.1.0")
	assert.Nil(t, err)
	assert.Zero(t, preVersion)
	assert.Equal(t, version, [3]int{0, 1, 0})
}

func TestSplitConditions(t *testing.T) {
	cond := splitConditions(">=1.2.3<2.3.4")
	assert.Len(t, cond, 2)
	assert.Equal(t, cond[0].Operator, ">=")
	assert.Equal(t, cond[0].Value, "1.2.3")
	assert.Equal(t, cond[1].Operator, "<")
	assert.Equal(t, cond[1].Value, "2.3.4")

	cond = splitConditions(">=1.2.3")
	assert.Len(t, cond, 1)
	assert.Equal(t, cond[0].Operator, ">=")
	assert.Equal(t, cond[0].Value, "1.2.3")

	cond = splitConditions(">=")
	assert.Len(t, cond, 0)

	cond = splitConditions("1.2.3")
	assert.Len(t, cond, 1)
	assert.Equal(t, cond[0].Operator, "")
	assert.Equal(t, cond[0].Value, "1.2.3")
}

func TestParseConditions(t *testing.T) {
	name, cond := ParseConditions("name>=1.2.3")
	assert.Equal(t, name, "name")
	assert.Len(t, cond, 1)
	assert.Equal(t, cond[0].Operator, ">=")
	assert.Equal(t, cond[0].Value, "1.2.3")

	name, cond = ParseConditions("name>=1.2.3!=2.3.4")
	assert.Equal(t, name, "name")
	assert.Len(t, cond, 2)
	assert.Equal(t, cond[0].Operator, ">=")
	assert.Equal(t, cond[0].Value, "1.2.3")
	assert.Equal(t, cond[1].Operator, "!=")
	assert.Equal(t, cond[1].Value, "2.3.4")

	name, cond = ParseConditions("name")
	assert.Equal(t, name, "name")
	assert.Len(t, cond, 0)

	name, cond = ParseConditions("name>=")
	assert.Equal(t, name, "name")
	assert.Len(t, cond, 0)

	name, cond = ParseConditions("name1.2.3")
	assert.Equal(t, name, "name1.2.3")
	assert.Len(t, cond, 0)
}

func TestIsOperator(t *testing.T) {
	for _, op := range []rune{'>', '<', '=', '!'} {
		assert.True(t, isOperator(op))
	}

	for _, char := range []rune{'a', ')', '0', '-', '+'} {
		assert.False(t, isOperator(char))
	}
}

func TestCompareVersion(t *testing.T) {
	for _, v := range compareVersionTrue {
		result, err := CompareVersion(v[0], v[1], v[2])
		assert.Nil(t, err)
		assert.True(t, result)
	}

	for _, v := range compareVersionFalse {
		result, err := CompareVersion(v[0], v[1], v[2])
		assert.Nil(t, err)
		assert.False(t, result)
	}
}

func TestCompareVersionInvalid(t *testing.T) {
	for _, v := range compareVersionInvalid {
		_, err := CompareVersion(v[0], v[1], v[2])
		assert.NotNil(t, err)
	}
}

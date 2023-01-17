package plugins

import (
	"testing"
)

type GetDiscordAndRuneScapeNameTest struct {
	segments              []string
	expectedDiscordName   string
	expectedRunescapeName string
}

func TestGetDiscordAndRuneScapeName(t *testing.T) {
	t.Parallel()

	getDiscordAndRuneScapeNameTests := []GetDiscordAndRuneScapeNameTest{
		{[]string{"lord", "ex#1234", "i", "ex", "i"}, "lord ex#1234", "i ex i"},
		{[]string{"joey#1337", "bender", "life"}, "joey#1337", "bender life"},
	}

	for _, test := range getDiscordAndRuneScapeNameTests {
		discordNameAndDiscriminator, runescapeName, err := getDiscordAndRuneScapeName(test.segments)
		if err != nil {
			t.Error(err)
		}
		if discordNameAndDiscriminator != test.expectedDiscordName {
			t.Errorf("Expected discord name and discriminator to be %s, got %s", test.expectedDiscordName, discordNameAndDiscriminator)
		}
		if runescapeName != test.expectedRunescapeName {
			t.Errorf("Expected RuneScape name to be %s, got %s", test.expectedRunescapeName, runescapeName)
		}
	}
}

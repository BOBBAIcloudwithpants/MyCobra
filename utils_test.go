package bobra

import (
	"reflect"
	"testing"
)



func Test_StripFlags(t *testing.T) {
	cmd := &Command{
		Short: "test",
		Long: "test",
		Example: "test",
	}
	cmd.Flags().StringP("aaaa", "a", "YOUR NAME", "author name for copyright attribution")
	cmd.Flags().StringP("ddd", "d", "YOUR NAME", "author name for copyright attribution")
	cmd.Flags().StringP("c", "c", "YOUR NAME", "author name for copyright attribution")

	input := []string{"-a", "subcmd1", "subcmd2", "-d123", "-c=14", "subcmd3"}
	r := stripFlags(input, cmd)
	expected := []string{ "subcmd2", "subcmd3"}

	if !reflect.DeepEqual(r, expected) {
		t.Errorf("expected '%q' but got '%q'", expected, r)
	}
}

func Test_RemoveFirstMatchStr(t *testing.T) {
	array := []string {"aasdf", "adggdd", "移除", "ddddd"}
	r := removeFirstMatchStr(array, "移除")
	expected := []string{"aasdf", "adggdd",  "ddddd"}
	if !reflect.DeepEqual(r, expected) {
		t.Errorf("expected '%q' but got '%q'", expected, r)
	}
}

package fail

import "fmt"

func Error() error {
	return fmt.Errorf("All requests failed")
}

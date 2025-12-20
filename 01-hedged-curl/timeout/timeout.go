package timeout

import "fmt"

var RequestsTimeout = 228

func Error() error {
	return fmt.Errorf("All requests timed out")
}

package efi

import (
	"fmt"
	"strings"
)

func GetStubPath() string {
	return fmt.Sprintf("/usr/lib/systemd/boot/efi/linux%s.efi.stub", strings.ToLower(Arch))
}

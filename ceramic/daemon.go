package ceramic

import (
	"context"
	"fmt"
	"log"
	"os/exec"
)

func Daemon(ctx context.Context) {
	cmd := exec.CommandContext(ctx, "ceramic", "daemon")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		log.Fatalln("ceramic panic error")
	}
}

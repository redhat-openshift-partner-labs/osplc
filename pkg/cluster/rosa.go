package cluster

import (
	"log"
	"os/exec"
)

func StartRosaCluster(name string) error {
	cmd := exec.Command("ocm", "resume", "cluster", name)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	otp, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	print(string(otp))
	return nil
}

func StopRosaCluster(name string) error {
	cmd := exec.Command("ocm", "hibernate", "cluster", name)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	otp, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	print(string(otp))
	return nil
}

package daemon

// Check root rights to use system service
func checkPrivileges() bool {
	cmd := exec.Command("id", "-g")
	if output, err := cmd.Output(); err == nil {
		gid := strings.TrimSpace(string(output))
		return gid == "0"
	} else {
		log.Println(err)
	}
	return false
}

package utils

import "testing"

func TestServerConfig_ReloadConfig(t *testing.T) {
	GlobalConfig.ReloadConfig()
}

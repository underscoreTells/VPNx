package config

import (
	"fmt"
	"testing"

	"github.com/underscoreTells/vpn-exit-node/internal/config/app"
	"github.com/underscoreTells/vpn-exit-node/internal/config/common"
)

var (
	testSchemaVersion  = app.DEFAULT_CONFIG_VERSION
	testGluetunVersion = app.GLUETUN_TARGET_VERSION
	testVPNProvider    = "custom"
	testVPNProtocol    = common.ProtocolWireguard
	testUsernameFrom   = "env"
	testUsernameName   = "VPN_USERNAME"
	testPasswordFrom   = "env"
	testPasswordName   = "VPN_PASSWORD"
	testLogLevel       = app.LOG_LEVEL_INFO
	testLogDestination = "stdout"
	testLogFilename    = "vpnx.log"
	testSchema         = fmt.Appendf(nil, `{
	"schema_version": %d,
	"gluetun_version": "%s",
	"vpn_config": {
		"provider": "%s",
		"protocol": "%s",
		"credentials": {
			"username": {
				"from": "%s",
				"name": "%s"
			},
			"password": {
				"from": "%s",
				"name": "%s"
			}
		}
	},
	"log": {
		"level": "%s",
		"destination": "%s",
		"filename": "%s"
	}
}`,
		testSchemaVersion,
		testGluetunVersion,
		testVPNProvider,
		testVPNProtocol,
		testUsernameFrom,
		testUsernameName,
		testPasswordFrom,
		testPasswordName,
		testLogLevel,
		testLogDestination,
		testLogFilename,
	)
	testDefaultsSchema = fmt.Appendf(nil, `{
	"gluetun_version": "%s",
	"vpn_config": {
		"provider": "%s",
		"credentials": {
			"username": {
				"from": "%s",
				"name": "%s"
			},
			"password": {
				"from": "%s",
				"name": "%s"
			}
		}
	},
	"log": {
		"destination": "%s",
		"filename": "%s"
	}
}`,
		testGluetunVersion,
		testVPNProvider,
		testUsernameFrom,
		testUsernameName,
		testPasswordFrom,
		testPasswordName,
		testLogDestination,
		testLogFilename,
	)
)

func TestGetConfigVersion(t *testing.T) {
	version, err := GetAppConfigVersion(testSchema)
	if err != nil {
		t.Errorf("GetConfigVersion error: %v", err)
	}

	if version != testSchemaVersion {
		t.Errorf("GetConfigVersion returned wrong version: %d, expected %d", version, testSchemaVersion)
	}
}

func TestLoadFromBytes(t *testing.T) {
	config, err := LoadAppConfigFromBytes(testSchema)
	if err != nil {
		t.Errorf("LoadFromBytes error: %v", err)
	}

	if config.VPNConfig.Provider != testVPNProvider {
		t.Errorf("LoadFromBytes returned wrong VPN provider: %s, expected %s", config.VPNConfig.Provider, testVPNProvider)
	}

	if config.VPNConfig.Protocol != testVPNProtocol {
		t.Errorf("LoadFromBytes returned wrong VPN protocol: %s, expected %s", config.VPNConfig.Protocol, testVPNProtocol)
	}

	if config.VPNConfig.Credentials.Username.From != testUsernameFrom {
		t.Errorf("LoadFromBytes returned wrong username from: %s, expected %s", config.VPNConfig.Credentials.Username.From, testUsernameFrom)
	}

	if config.VPNConfig.Credentials.Username.Name != testUsernameName {
		t.Errorf("LoadFromBytes returned wrong username name: %s, expected %s", config.VPNConfig.Credentials.Username.Name, testUsernameName)
	}

	if config.VPNConfig.Credentials.Password.From != testPasswordFrom {
		t.Errorf("LoadFromBytes returned wrong password from: %s, expected %s", config.VPNConfig.Credentials.Password.From, testPasswordFrom)
	}

	if config.VPNConfig.Credentials.Password.Name != testPasswordName {
		t.Errorf("LoadFromBytes returned wrong password name: %s, expected %s", config.VPNConfig.Credentials.Password.Name, testPasswordName)
	}

	if config.Log.Level != testLogLevel {
		t.Errorf("LoadFromBytes returned wrong log level: %s, expected %s", config.Log.Level, testLogLevel)
	}

	if config.Log.Destination != testLogDestination {
		t.Errorf("LoadFromBytes returned wrong log destination: %s, expected %s", config.Log.Destination, testLogDestination)
	}

	if config.Log.Filename != testLogFilename {
		t.Errorf("LoadFromBytes returned wrong log file name: %s, expected %s", config.Log.Filename, testLogFilename)
	}
}

func TestDefaults(t *testing.T) {
	config, err := LoadAppConfigFromBytes(testDefaultsSchema)
	if err != nil {
		t.Errorf("LoadFromBytes error: %v", err)
	}

	if config.SchemaVersion != app.DEFAULT_CONFIG_VERSION {
		t.Errorf("LoadFromBytes returned wrong schema version: %d, expected %d", config.SchemaVersion, app.DEFAULT_CONFIG_VERSION)
	}

	if config.VPNConfig.Protocol != app.DEFAULT_VPN_PROTOCOL {
		t.Errorf("LoadFromBytes returned wrong VPN protocol: %s, expected %s", config.VPNConfig.Protocol, app.DEFAULT_VPN_PROTOCOL)
	}

	if config.Log.Level != app.DEFAULT_LOG_LEVEL {
		t.Errorf("LoadFromBytes returned wrong log level: %s, expected %s", config.Log.Level, app.DEFAULT_LOG_LEVEL)
	}
}

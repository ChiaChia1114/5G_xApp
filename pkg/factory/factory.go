/*
 * AUSF Configuration Factory
 */

package factory

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"xApp/internal/logger"
)

var XAppConfig Config

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	if content, err := ioutil.ReadFile(f); err != nil {
		return err
	} else {
		XAppConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &XAppConfig); yamlErr != nil {
			return yamlErr
		}
	}

	return nil
}

func CheckConfigVersion() error {
	currentVersion := XAppConfig.GetVersion()

	if currentVersion != XApp_EXPECTED_CONFIG_VERSION {
		return fmt.Errorf("config version is [%s], but expected is [%s].",
			currentVersion, XApp_EXPECTED_CONFIG_VERSION)
	}
	logger.CfgLog.Infof("config version [%s]", currentVersion)

	return nil
}

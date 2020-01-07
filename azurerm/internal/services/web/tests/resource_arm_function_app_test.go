package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFunctionApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					testCheckAzureRMFunctionAppHasNoContentShare(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "~1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "inbound_ip_addresses"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "possible_inbound_ip_addresses"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "outbound_ip_addresses"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "possible_outbound_ip_addresses"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					testCheckAzureRMFunctionAppHasNoContentShare(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMFunctionApp_requiresImport),
		},
	})
}

func TestAccAzureRMFunctionApp_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_tagsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_tagsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.hello", "Berlin"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.%", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_credential.#", "1"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_appSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.hello", "world"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_siteConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_alwaysOn(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.always_on", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_linuxFxVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_linuxFxVersion(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "functionapp,linux,container"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.linux_fx_version", "DOCKER|(golang:latest)"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_connectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_connectionStrings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.name", "Example"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.type", "PostgreSQL"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_siteConfigMulti(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.%", "0"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_appSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.hello", "world"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_appSettingsAlwaysOn(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.hello", "world"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.always_on", "true"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_appSettingsAlwaysOnLinuxFxVersion(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "functionapp,linux,container"),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.hello", "world"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.always_on", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.linux_fx_version", "DOCKER|(golang:latest)"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_appSettingsAlwaysOnLinuxFxVersionConnectionStrings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "functionapp,linux,container"),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.hello", "world"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.always_on", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.linux_fx_version", "DOCKER|(golang:latest)"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.name", "Example"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.type", "PostgreSQL"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_updateVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_version(data, "~1"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "~1"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_version(data, "~2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "~2"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_3264bit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_64bit(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.use_32_bit_worker_process", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_httpsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_httpsOnly(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "https_only", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_consumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_consumptionPlan(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					testCheckAzureRMFunctionAppHasContentShare(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_consumptionPlanUppercaseName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_consumptionPlanUppercaseName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					testCheckAzureRMFunctionAppHasContentShare(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.use_32_bit_worker_process", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_createIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_basicIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_updateIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "0"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_basicIdentity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.principal_id", validate.UUIDRegExp),
					resource.TestMatchResourceAttr(data.ResourceName, "identity.0.tenant_id", validate.UUIDRegExp),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_loggingDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_loggingDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					testCheckAzureRMFunctionAppHasNoContentShare(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_builtin_logging", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_updateLogging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_builtin_logging", "true"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_loggingDisabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_builtin_logging", "false"),
				),
			},
			{
				Config: testAccAzureRMFunctionApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_builtin_logging", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMFunctionApp_authSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	tenantID := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_authSettings(data, tenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.issuer", fmt.Sprintf("https://sts.windows.net/%s", tenantID)),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.runtime_version", "1.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.unauthenticated_client_action", "RedirectToLoginPage"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.token_refresh_extension_hours", "75"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.token_store_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.additional_login_params.test_key", "test_value"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.allowed_external_redirect_urls.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.allowed_external_redirect_urls.0", "https://terra.form"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_id", "aadclientid"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.client_secret", "aadsecret"),
					resource.TestCheckResourceAttr(data.ResourceName, "auth_settings.0.active_directory.0.allowed_audiences.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_corsSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_corsSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.cors.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.cors.0.support_credentials", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.cors.0.allowed_origins.#", "4"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_vnetName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	vnetName := strings.ToLower(acctest.RandString(11))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_vnetName(data, vnetName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.virtual_network_name", vnetName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_enableHttp2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_enableHttp2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.http2_enabled", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_minTlsVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_minTlsVersion(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.min_tls_version", "1.2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFunctionApp_ftpsState(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFunctionApp_ftpsState(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ftps_state", "AllAllowed"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMFunctionAppDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_function_app" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMFunctionAppExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		functionAppName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Function App: %s", functionAppName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, functionAppName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Function App %q (resource group: %q) does not exist", functionAppName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMFunctionAppHasContentShare(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		functionAppName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Function App: %s", functionAppName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		appSettingsResp, err := client.ListApplicationSettings(ctx, resourceGroup, functionAppName)
		if err != nil {
			return fmt.Errorf("Error making Read request on AzureRM Function App AppSettings %q: %+v", functionAppName, err)
		}

		for k := range appSettingsResp.Properties {
			if strings.EqualFold("WEBSITE_CONTENTSHARE", k) {
				return nil
			}
		}

		return fmt.Errorf("Function App %q does not contain the Website Content Share!", functionAppName)
	}
}

func testCheckAzureRMFunctionAppHasNoContentShare(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		functionAppName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Function App: %s", functionAppName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		appSettingsResp, err := client.ListApplicationSettings(ctx, resourceGroup, functionAppName)
		if err != nil {
			return fmt.Errorf("Error making Read request on AzureRM Function App AppSettings %q: %+v", functionAppName, err)
		}

		for k, v := range appSettingsResp.Properties {
			if strings.EqualFold("WEBSITE_CONTENTSHARE", k) && v != nil && *v != "" {
				return fmt.Errorf("Function App %q contains the Website Content Share!", functionAppName)
			}
		}

		return nil
	}
}

func testAccAzureRMFunctionApp_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMFunctionApp_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_function_app" "import" {
  name                      = "${azurerm_function_app.test.name}"
  location                  = "${azurerm_function_app.test.location}"
  resource_group_name       = "${azurerm_function_app.test.resource_group_name}"
  app_service_plan_id       = "${azurerm_function_app.test.app_service_plan_id}"
  storage_connection_string = "${azurerm_function_app.test.storage_connection_string}"
}
`, template)
}

func testAccAzureRMFunctionApp_tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_tagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  tags = {
    environment = "production"
    hello       = "Berlin"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_version(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  version                   = "%[4]s"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, version)
}

func testAccAzureRMFunctionApp_appSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  app_settings = {
    "hello" = "world"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_alwaysOn(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    always_on = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_linuxFxVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "Linux"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }

  reserved = true
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  version                   = "~2"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    linux_fx_version = "DOCKER|(golang:latest)"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_connectionStrings(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_appSettingsAlwaysOn(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  app_settings = {
    "hello" = "world"
  }

  site_config {
    always_on = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_appSettingsAlwaysOnLinuxFxVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "Linux"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }

  reserved = true
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  version                   = "~2"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  app_settings = {
    "hello" = "world"
  }

  site_config {
    always_on        = true
    linux_fx_version = "DOCKER|(golang:latest)"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_appSettingsAlwaysOnLinuxFxVersionConnectionStrings(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "Linux"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }

  reserved = true
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  version                   = "~2"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  app_settings = {
    "hello" = "world"
  }

  site_config {
    always_on        = true
    linux_fx_version = "DOCKER|(golang:latest)"
  }

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_64bit(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    use_32_bit_worker_process = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionApp_httpsOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
  https_only                = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionApp_consumptionPlan(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionApp_consumptionPlanUppercaseName(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-FuncWithUppercase"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionApp_basicIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_loggingDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"
  enable_builtin_logging    = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_authSettings(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  auth_settings {
    enabled                       = true
    issuer                        = "https://sts.windows.net/%[4]s"
    runtime_version               = "1.0"
    unauthenticated_client_action = "RedirectToLoginPage"
    token_refresh_extension_hours = 75
    token_store_enabled           = true

    additional_login_params = {
      test_key = "test_value"
    }

    allowed_external_redirect_urls = [
      "https://terra.form",
    ]

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, tenantID)
}

func testAccAzureRMFunctionApp_corsSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    cors {
      allowed_origins = [
        "http://www.contoso.com",
        "www.contoso.com",
        "contoso.com",
        "http://localhost:4201",
      ]

      support_credentials = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMFunctionApp_vnetName(data acceptance.TestData, vnetName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%[1]d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    virtual_network_name = "%[4]s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, vnetName)
}

func testAccAzureRMFunctionApp_enableHttp2(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    http2_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionApp_minTlsVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    min_tls_version = "1.2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFunctionApp_ftpsState(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                      = "acctest-%d-func"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  app_service_plan_id       = "${azurerm_app_service_plan.test.id}"
  storage_connection_string = "${azurerm_storage_account.test.primary_connection_string}"

  site_config {
    ftps_state = "AllAllowed"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

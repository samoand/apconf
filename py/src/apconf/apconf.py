"""Global config object and config processing."""
import os.path
from typing import Any, Callable, Optional

from . import process_yaml
from . import util

# Type aliases for clearer signatures
ConfigDict = dict[str, Any]
Preprocessor = Callable[[ConfigDict], None]
Deployer = Callable[[ConfigDict, ConfigDict, util.ConfigDiffResult], None]
Validator = Callable[[ConfigDict, ConfigDict, util.ConfigDiffResult], bool]


class ApconfException(Exception):
    """Module-level exception."""


# pylint:disable=too-many-positional-arguments
class Config:
    """
    Config holder and orchestrator.

    Loads configuration from YAML files, applies preprocessors,
    validates, and deploys configuration changes.

    Attributes:
        config: The current configuration dictionary
        config_root: Root directory containing config subdirectories
        config_basenames: List of config directory names to load

    Example:
        >>> config = Config(
        ...     config_root="/app/config",
        ...     config_basenames=["base", "production"],
        ...     template_params={"env": "prod"},
        ...     config_deployers=[my_deployer],
        ...     config_preprocessors=None,
        ...     config_validators=None,
        ... )
        >>> print(config.config["my_setting"])
    """

    def __init__(
            self,
            config_root: str,
            config_basenames: list[str],
            template_params: Optional[dict[str, Any]],
            config_preprocessors: Optional[list[Preprocessor]],
            config_deployers: Optional[list[Deployer]],
            config_validators: Optional[list[Validator]],
            *,
            raise_on_deployer_error: bool = False,
    ) -> None:
        """
        Initialize the Config holder.

        Args:
            config_root: Root directory of the config dirs
            config_basenames: List of subdirectory names to load configs from
            template_params: Parameters for Jinja2 template rendering in YAML
            config_preprocessors: List of functions that mutate config
                before deployment (in-place mutations)
            config_deployers: List of functions that deploy the config.
                Called with (config_new, config_old, config_diff_result)
            config_validators: List of functions that validate the config.
                Must return True for deployment to proceed.
            raise_on_deployer_error: If True, raise exceptions from deployers
                immediately. If False (default), collect exceptions and return
                them from apply(). Default False for backward compatibility.
        """
        self.config_root = config_root
        self.config_basenames = config_basenames
        self.template_params = template_params
        self.config_preprocessors = config_preprocessors
        self.config_deployers = config_deployers
        self.config_validators = config_validators
        self.raise_on_deployer_error = raise_on_deployer_error
        self.config: ConfigDict = {}
        self._init()

    def _init(self) -> None:
        """Load config from YAML files and apply it."""
        config_dirs = [os.path.join(self.config_root, config_dir) for
                       config_dir in self.config_basenames]

        for config_dir in config_dirs:
            if not os.path.isdir(config_dir):
                raise ApconfException(
                    f"Config directory {config_dir} does not exist.")

        config = process_yaml.process_yaml_dirs(
            config_dirs, self.template_params)["Config"]
        self.apply(config)

    def _preprocess(self, config: ConfigDict) -> None:
        """
        In-place preprocessing of config after it was
        posted or read from the config files but before
        it's deployed.
        """
        if self.config_preprocessors is None:
            return
        for preprocessor in self.config_preprocessors:
            preprocessor(config)

    def _validate(
            self,
            config_new: ConfigDict,
            config_diff_result: util.ConfigDiffResult
    ) -> bool:
        """
        Validate the config using registered validators.

        Returns:
            True if all validators pass (or no validators registered),
            False otherwise.
        """
        if self.config_validators is None:
            return True
        return all(
            validator(config_new, self.config, config_diff_result)
            for validator in self.config_validators
        )

    def _deploy(
            self,
            config_new: ConfigDict,
            config_diff_result: util.ConfigDiffResult
    ) -> Optional[list[Exception]]:
        """
        Deploy the config by calling all registered deployers.

        Args:
            config_new: The new configuration to deploy
            config_diff_result: Diff between new and old config

        Returns:
            List of exceptions raised by deployers, or None if no errors.
            Only returns exceptions if raise_on_deployer_error is False.

        Raises:
            Exception: If raise_on_deployer_error is True and a deployer fails.
        """
        exs: Optional[list[Exception]] = None
        if self.config_deployers is None:
            return exs
        for deployer in self.config_deployers:
            try:
                deployer(config_new, self.config, config_diff_result)
            # pylint:disable=broad-exception-caught
            except Exception as e:
                if self.raise_on_deployer_error:
                    raise
                if exs is None:
                    exs = []
                exs.append(e)
        return exs

    def apply(self, config: ConfigDict) -> Optional[list[Exception]]:
        """
        Apply a new configuration.

        Preprocesses, validates, and deploys the config. On initial load,
        all keys appear in config_diff_result.diff_added. On updates,
        keys appear in diff_changed, diff_added, or diff_removed as
        appropriate.

        Args:
            config: The new configuration dictionary to apply.

        Returns:
            List of exceptions from deployers (if raise_on_deployer_error
            is False), or None if no errors occurred.

        Raises:
            Exception: If raise_on_deployer_error is True and a deployer fails.
        """
        exs: Optional[list[Exception]] = None
        config_diff_result = util.config_diff(
            config_new=config, config_old=self.config)
        self._preprocess(config)
        if self._validate(config, config_diff_result):
            exs = self._deploy(config, config_diff_result)
            self.config = config
        return exs

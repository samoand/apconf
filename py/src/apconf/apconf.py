"""Global config object and config processing."""
import os.path
from typing import Callable, Optional

from . import process_yaml
from . import util

class ApconfException(Exception):
    """Module-level exception."""

# pylint:disable=too-many-positional-arguments
class Config:
    """ Config holder."""
    def __init__(self,
                 config_root:str,
                 config_basenames:list[str],
                 template_params:Optional[dict],
                 config_preprocessors:Optional[list[Callable[[dict], None]]],
                 config_deployers:Optional[
                     list[Callable[[dict, dict, util.ConfigDiffResult], None]]],
                 config_validators:Optional[list[
                     Callable[[dict, dict, util.ConfigDiffResult], bool]]]):
        """
        Init.
        Inputs:
        - config_root: str: root directory of the config dirs
        - config_preprocessors: list[Callable[[dict], None]]: list of mutators
            that will be applied to the config before it's deployed
        - config_deployers:
          list[Callable[[dict, dict, util.ConfigDiffResult], None]]:
            list of functions that deploy the config.
        - config_validators:
          list[Callable[[dict, dict, util.ConfigDiffResult], None]]:
            list of functions that validate the config.
            These function validate the config which has already been
            mutated by the preprocessors.
        """
        self.config_root = config_root
        self.config_basenames = config_basenames
        self.template_params = template_params
        self.config_preprocessors = config_preprocessors
        self.config_deployers = config_deployers
        self.config_validators = config_validators
        self.config = {}
        self._init()

    def _init(self):
        config_dirs = [os.path.join(self.config_root, config_dir) for
                       config_dir in self.config_basenames]

        for config_dir in config_dirs:
            if not os.path.isdir(config_dir):
                raise ApconfException(
                    f"Config directory {config_dir} does not exist.")

        config = process_yaml.process_yaml_dirs(
            config_dirs, self.template_params)["Config"]
        self.apply(config)

    def _preprocess(self, config:dict):
        """
        In-place preprocessing of config after it was
        posted or read from the config files but before
        it's deployed.
        """
        if self.config_preprocessors is None:
            return
        for preprocessor in self.config_preprocessors:
            preprocessor(config)

    def _validate(self,
                  config_new:dict,
                  config_diff_result:util.ConfigDiffResult) -> bool:
        """Deploy the config."""
        return True if self.config_validators is  None else all(
            validator(config_new,
                      self.config,
                      config_diff_result) for
            validator in self.config_validators)

    def _deploy(self,
                config_new:dict,
                config_diff_result:util.ConfigDiffResult) -> Optional[
                    list[Exception]]:
        """Deploy the config."""
        exs:Optional[list[Exception]] = None
        if self.config_deployers is None:
            return exs
        for deployer in self.config_deployers:
            try:
                deployer(config_new,
                         self.config,
                         config_diff_result)
            # pylint:disable=broad-exception-caught
            except Exception as e:
                if exs is None:
                    exs = []
                exs.append(e)
        return exs

    def apply(self, config:dict) -> Optional[list[Exception]]:
        """
        Post the config.
        Inputs:
        - config: dict: config to post.
        """
        exs:Optional[list[Exception]] = None
        config_diff_result = util.config_diff(
            config_new=config, config_old=self.config)
        self._preprocess(config)
        if self._validate(config, config_diff_result):
            exs = self._deploy(config, config_diff_result)
            self.config = config
        return exs

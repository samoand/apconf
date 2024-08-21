"""Illustrate use of the library."""
from pathlib import Path
from typing import Callable

import copy
import logging.config
import re
import os.path

from apconf import apconf, util, preprocess_config

def _path_from_desc(path_desc):
    return os.path.normpath(
        os.path.join(
            *re.split(r'\s+', path_desc)))

# In-place post processing of config descriptors after they were
#   read from the config files.
# pylint:disable=unnecessary-lambda
_config_preprocessors:list[Callable[[dict], None]] = [
    preprocess_config.preprocessor(
        key_filter=lambda key: key.endswith("dir_desc"),
        key_transformer=lambda old_key: old_key[:-len("_desc")],
        value_transformer=lambda old_value: _path_from_desc(old_value),
        keep_filt_keys=False),
    preprocess_config.preprocessor(
        key_filter=lambda key: key.endswith("root_desc"),
        key_transformer=lambda old_key: old_key[:-len("_desc")],
        value_transformer=lambda old_value: _path_from_desc(old_value),
        keep_filt_keys=False),
    preprocess_config.preprocessor(
        key_filter=lambda key: key == "log_path_desc",
        key_transformer=lambda _: "filename",
        value_transformer=lambda old_value: _path_from_desc(old_value),
        keep_filt_keys=False),
    ]

def __apply_log(updated_config:dict,
                _:dict,
                config_diff_result:util.ConfigDiffResult) -> None:
    if not any("logging_config" in d for d in
           [config_diff_result.diff_changed,
            config_diff_result.diff_added,
            config_diff_result.diff_removed]):
        return
    logging_config = updated_config["logging_config"]["spec"]
    for _, handler in logging_config["handlers"].items():
        if "filename" in handler:
            Path(handler["filename"]).parent.mkdir(
                parents=True, exist_ok=True)
    logging.config.dictConfig(logging_config)

_config_deployers:list[Callable[[dict, dict, util.ConfigDiffResult], None]] = [
    __apply_log
]

if __name__ == '__main__':
    logger = logging.getLogger(__name__)
    project_root = os.path.normpath(
        util.find_git_root_or_error(os.path.realpath(__file__)))
    config_root = os.path.join(project_root, "test_config")
    config = apconf.Config(
        config_root=config_root,
        config_basenames=["crawl", "logging-py"],
        template_params={
            "project_root": project_root,
            "proc_id": os.getpid()
        },
        config_preprocessors=_config_preprocessors,
        config_deployers=_config_deployers,
        config_validators=None)
    # this should go to log but not to console
    logger.warning("This is a warning")
    # this should go to both log and console
    logger.error("This is an error")

    # update logging level, see that now warning
    # goes to console
    new_config = copy.deepcopy(config.config)
    new_config["logging_config"][
        "spec"]["handlers"]["console_handler"]["level"] = "WARNING"
    config.apply(new_config)
    logger.warning("This is another warning")

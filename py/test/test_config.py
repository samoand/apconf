""" Test configurations. """

import copy
import logging.config
import os
from pathlib import Path
import re
from typing import Callable
import unittest

from apconf import apconf, preprocess_config, util

logger = logging.getLogger(__name__)

def _path_from_desc(path_desc):
    return os.path.normpath(
        os.path.join(
            *re.split(r'\s+', path_desc)))

def _apply_log(new_config:dict,
               _:dict,
               config_diff_result:util.ConfigDiffResult) -> None:
    if not any("logging_config" in d for d in
           [config_diff_result.diff_changed,
            config_diff_result.diff_added,
            config_diff_result.diff_removed]):
        return
    logging_config = new_config["logging_config"]["spec"]
    for _, handler in logging_config["handlers"].items():
        if "filename" in handler:
            Path(handler["filename"]).parent.mkdir(
                parents=True, exist_ok=True)

    logging.config.dictConfig(logging_config)

class TestConfig(unittest.TestCase):
    """
    Test cases for app configuration created from config files.
    """
    # setup
    # pylint:disable=unnecessary-lambda
    def setUp(self):
        config_preprocessors:list[Callable[[dict], None]] = [
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
        config_deployers:list[
            Callable[[dict, dict, util.ConfigDiffResult], None]] = [
            _apply_log
        ]
        project_root = os.path.normpath(
            util.find_git_root_or_error(os.path.realpath(__file__)))
        config_root = os.path.join(project_root, "test_config")
        self.proc_id = os.getpid()
        self.config = apconf.Config(
            config_root=config_root,
            config_basenames=["crawl", "logging-py"],
            template_params={
                "project_root": project_root,
                "proc_id": self.proc_id
            },
            config_preprocessors=config_preprocessors,
            config_deployers=config_deployers,
            config_validators=None)
        self.configured_filename = self.config.config[
            "logging_config"]["spec"]["handlers"][
                "file_handler"]["filename"]

    def test_init_config(self):
        """Test initial configuration declared in config files."""

        # See that "crawler" config is there.
        self.assertEqual(
            self.config.config["crawler_config"]["spec"]["num_workers"],
            30)

        # See that "logging" config is properly set and enforced
        log_file_end = os.path.join(
            "artifacts", "test", "log",
            "myapp." + str(self.proc_id) + ".log")
        self.assertTrue(
            self.configured_filename.endswith(log_file_end))
        warning_msg = "This is a warning"
        logger.warning(warning_msg)
        # read the log file, verify the warning message
        with open(self.configured_filename, "r", encoding='utf-8') as log_file:
            log_content = log_file.read()
            self.assertIn(warning_msg, log_content)
        info_msg = "This is an info message"
        logger.info(info_msg)
        # Read the log file, verify that info message
        # hasn't been logged
        with open(self.configured_filename, "r", encoding='utf-8') as log_file:
            log_content = log_file.read()
            self.assertNotIn(info_msg, log_content)

    def test_modified_log_config(self):
        """Create and apply new config with INFO logging level."""
        new_config = copy.deepcopy(self.config.config)
        new_config["logging_config"]["spec"]["handlers"][
            "file_handler"]["level"] = "INFO"
        new_config["logging_config"]["spec"]["loggers"][
            "root"]["level"] = "INFO"
        self.config.apply(new_config)
        info_msg = "This is another info message"
        logger.info(info_msg)
        # read the log file, verify that info message
        # has been logged this time
        with open(self.configured_filename, "r", encoding='utf-8') as log_file:
            log_content = log_file.read()
            self.assertIn(info_msg, log_content)

if __name__ == '__main__':
    unittest.main()

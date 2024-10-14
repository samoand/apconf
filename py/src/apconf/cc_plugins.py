"""
Cookie-cutter plugins (cc-plugins): most common plugins.
"""
from pathlib import Path
import logging.config

from . import util

def apply_log(
        updated_config:dict,
        _:dict,
        config_diff_result:util.ConfigDiffResult) -> None:
    """Apply log expressed as config."""
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

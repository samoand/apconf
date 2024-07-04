"""Post-processing functions for config."""
from typing import Callable

def _process_list(lst,
                  key_filter:Callable[[str], bool],
                  key_transformer:Callable[[str], str],
                  value_transformer:Callable[[object], object],
                  processed_ids:set[int],
                  keep_filt_keys:bool=True) -> None:
    for item in lst:
        if isinstance(item, dict):
            _process_dict(item,
                          key_filter,
                          key_transformer,
                          value_transformer,
                          processed_ids,
                          keep_filt_keys)
        elif isinstance(item, list):
            _process_list(item,
                          key_filter,
                          key_transformer,
                          value_transformer,
                          processed_ids,
                          keep_filt_keys)

def _process_dict(d,
                  key_filter:Callable[[str], bool],
                  key_transformer:Callable[[str], str],
                  value_transformer:Callable[[object], object],
                  processed_ids:set[int],
                  keep_filt_keys:bool=True) -> None:
    keys_to_add = {}
    filt_keys = []

    for key, value in list(d.items()):
        if key_filter(key):
            new_key = key_transformer(key)
            new_value = value_transformer(value)

            # Avoid processing the same element more than once
            if id(new_value) not in processed_ids:
                processed_ids.add(id(new_value))

                # Recursively process new_value if it's a dict or list
                if isinstance(new_value, dict):
                    _process_dict(new_value,
                                  key_filter,
                                  key_transformer,
                                  value_transformer,
                                  processed_ids)
                elif isinstance(new_value, list):
                    _process_list(new_value,
                                  key_filter,
                                  key_transformer,
                                  value_transformer,
                                  processed_ids)

                keys_to_add[new_key] = new_value
                filt_keys.append(key)
        else:
            if isinstance(value, dict):
                _process_dict(value,
                              key_filter,
                              key_transformer,
                              value_transformer,
                              processed_ids,
                              keep_filt_keys)
            elif isinstance(value, list):
                _process_list(value,
                              key_filter,
                              key_transformer,
                              value_transformer,
                              processed_ids,
                              keep_filt_keys)

    # Update the dictionary with newly transformed keys and values
    for new_key, new_value in keys_to_add.items():
        d[new_key] = new_value
    if not keep_filt_keys:
        for key in filt_keys:
            del d[key]

def preprocessor(
        key_filter:Callable[[str], bool],
        key_transformer:Callable[[str], str],
        value_transformer:Callable[[object], object],
        keep_filt_keys:bool=True) -> Callable[[dict], None]:
    """
    Create a convenience cookiecutter preprocessor for
      a dictionary with config data.

    Args:
        key_filter: A function that takes a key and returns True if the key
            should be transformed, False otherwise.
        key_transformer: A function that takes a key and returns the new key.
        value_transformer: A function that takes a value and returns the new
            value.
        keep_filt_keys: If True, keep the original keys in the dictionary after
            transformation. If False, remove the original keys.
    """
    def inner(config:dict) -> None:
        """
        Args:
            config: The config dictionary to process.
        """
        _process_dict(
            config,
            key_filter,
            key_transformer,
            value_transformer,
            set(),
            keep_filt_keys)
    return inner

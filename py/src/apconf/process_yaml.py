"""Read config data from yaml files in given dirs. Support params."""
from collections import defaultdict
from pathlib import Path
import yaml
from jinja2 import Template

def process_yaml_dirs(dirs, template_params):
    """Process dirs."""
    final_dict = defaultdict(dict, {})

    def merge_dicts(d1, d2):
        for key in d2:
            if key in d1:
                if isinstance(d1[key], dict) and isinstance(d2[key], dict):
                    merge_dicts(d1[key], d2[key])
                elif isinstance(d1[key], list) and isinstance(d2[key], list):
                    d1[key] += d2[key]
                elif d1[key] != d2[key]:
                    # pylint:disable=broad-exception-raised
                    raise Exception(
                        f"Conflict in merging dictionaries at key: {key}")
            else:
                d1[key] = d2[key]

    def process_yaml_content(content):
        loaded_content = yaml.safe_load_all(content)
        for doc in loaded_content:
            if 'kind' not in doc or 'metadata' not in doc or 'name' not in doc['metadata']:
                continue
            kind = doc['kind']
            name = doc['metadata']['name']
            wrapped = {kind: {name: doc}}
            merge_dicts(final_dict, wrapped)

    for dir_path in dirs:
        for file_path in Path(dir_path).rglob('*.yaml'):
            if file_path.name.startswith ('.'):
                continue
            with open(file_path, 'r', encoding='utf-8') as file:
                process_yaml_content(
                    Template(file.read()).render(template_params))
    return dict(final_dict)

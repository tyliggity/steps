#!/usr/bin/env python3

# import glob
import os
import yaml
from pathlib import Path
import glob

PATH_TO_VENDORS = './steps'

def get_all_manifests_paths():
    return list(Path(PATH_TO_VENDORS).rglob('*/manifest.yaml'))

def get_vendor_manifests_paths(vendor):
    all_paths = get_all_manifests_paths()

    return [ path for path in all_paths if get_vendor(path) == vendor]

def get_vendors():
    return listdirs(PATH_TO_VENDORS)

def get_step(path):
    return path.parent.name

def get_vendor(path):
    parents = list(path.parents)[::-1]
    return parents[2].name

def get_vendor_steps(paths, vendor):
    steps = []
    for path in paths:
        if get_vendor(path) == vendor:
            steps.append(get_step(path))

    return set(steps)


def listdirs(path):
    return [d for d in os.listdir(path) if os.path.isdir(os.path.join(path, d))]

class HealthReport:
    def __init__(self, isDetailed):
        self.isDetailed = isDetailed
        self.invalid_yamls = []
        self.other_errors = []
        self.vendors = 0
        self.manifests = 0
        self.manifests_with_envs = 0
        self.manifests_with_args = 0
        self.manifests_with_both = 0

    def collect_error(self, err, path):
        if type(err) == yaml.scanner.ScannerError:
            self.invalid_yamls.append({'path': path, 'error': err})
        else:
            self.other_errors.append({'path': path, 'error': err})

    def incr_var_types_on_manifests(self, manifest):
        if "envs" in manifest:
            self.manifests_with_envs += 1
        if "args" in manifest:
            self.manifests_with_args += 1
        if "args" in manifest and "envs" in manifest:
            self.manifests_with_both += 1


    def print(self):
        print('\n * Manifests Report *')
        if self.isDetailed: print ('(Detailed)') 
        print('total vendors: ', self.vendors)
        print('total manifests: ', self.manifests)
        print(f'manifests containing (args: {self.manifests_with_args}), (envs: {self.manifests_with_envs}) , (both: {self.manifests_with_both})')
        self.printErrors('invalid yamls', self.invalid_yamls)
        self.printErrors('other errors', self.other_errors)


    def printErrors(self, error_name, errors):
        print(f'{error_name}: {len(errors)}')

        if self.isDetailed:
            for error in errors:
                print(error['error'], error['path'])
            print('')

def report_manifests_health(isDetailed):
    manifest_paths = get_all_manifests_paths()

    report = HealthReport(isDetailed)
    report.vendors = len(get_vendors())
    report.manifests = len(manifest_paths)

    for manifest_path in manifest_paths:
        try:
            with open(manifest_path, 'r') as stream:
                manifest = yaml.safe_load(stream)

                report.incr_var_types_on_manifests(manifest)
        except Exception as e:
            report.collect_error(e, manifest_path)

    report.print()

def validate_envs_consitancy(isDetailed, vendor = None):
    inconsistent_fields = []
    vendors = []

    if vendor:
        vendors = [vendor]
    else:
        vendors = get_vendors()

    for vendor in vendors:
        manifests_paths = get_vendor_manifests_paths(vendor)
        vendor_envs = get_envs_by_step(manifests_paths)
        
        inconsistent_fields = inconsistent_fields + assert_equal_vars(vendor_envs, vendor)

    print("\n * validate envs inconsistancy *")
    if vendor: print(f'Vendor: {vendor}')
    print('Total inconsistant fields: ', len(inconsistent_fields))
    if isDetailed: print(*inconsistent_fields)


def get_envs_by_step(manifest_paths):
    envs = {}

    for m_path in manifest_paths:
        envs[get_step(m_path)] = get_env_vars(m_path)
    
    return envs


def get_env_vars(manifest_path):
    try:
        manifest = yaml.safe_load(open(manifest_path, 'r'))
        if 'envs' in manifest:
            return manifest['envs']
    except (FileNotFoundError, yaml.scanner.ScannerError):
        print("FileNotFoundError", manifest_path)
        pass
    return None


def assert_equal_vars(vendor_envs, vendor):
    first_vars = {}
    inconsistent_fields = []

    for step, env_vars in vendor_envs.items():
        for var in env_vars:
            var_name = var['name']

            if var_name in first_vars:
                inconsistent_fields = inconsistent_fields + assert_each_field(first_vars[var_name], var, {'vendor': vendor, 'step': step})
            else:
                first_vars[var_name] = var

    return inconsistent_fields
            
 
def assert_each_field(first_var, curr_var, error_info):
    inconsistent_fields = []
    for key, expected_value in first_var.items():
        if key in curr_var and curr_var[key] != expected_value:
            inconsistent_fields.append(f'\n ({error_info["vendor"]}/{error_info["step"]}/envs/{first_var["name"]}) {key}\n is: {curr_var[key]} \n expected to be: {expected_value}\n')

    return inconsistent_fields


def validate_examples_exist(isDetailed, vendor = None):
    vars_without_examples = []
    total_vars = 0

    if vendor:
        manifest_paths = get_vendor_manifests_paths(vendor)
    else:
        manifest_paths = get_all_manifests_paths()

    for manifest_path in manifest_paths:
        env_vars = get_env_vars(manifest_path)
        if env_vars:
            total_vars += len(env_vars)
            vars_without_examples = vars_without_examples + assert_example_exist(env_vars, manifest_path)
    
    print("\n * validate examples exists *")
    if vendor: print(f'Vendor: {vendor}')
    print('Total envs: ', total_vars)
    print('Lacking examples: ', len(vars_without_examples))
    if vendor and isDetailed:
        print('details: ', vars_without_examples)


def assert_example_exist(env_vars, manifest_path):
    vars_without_examples = []

    for var in env_vars:
        if 'example' not in var:
            vars_without_examples.append({'vendor': get_vendor(manifest_path), 'step': get_step(manifest_path), 'var_name': var["name"]})

    return vars_without_examples



## print count or explantion
detailed = False

## can be vendor (dir name) or None
vendor = 'slack'

print(f'vendor: {vendor}, detailed: {detailed}')
report_manifests_health(detailed)
validate_envs_consitancy(detailed, vendor)
validate_examples_exist(detailed, vendor)

## README
# so anyone who adds examples or envs can know:
# * which steps have examples, by vendor/all
# * which envs has inconsistent values in the same vendor - for example that all KUBECONFIG_CONTENT have the same examples, descriptions, required, etc. its not asserting, just prints reports to inform you so you can decide what to do with it (helped me find a lot of mistakes, typos etc)
# * general report - counts + errors in opening yamls

# pip3 install -> python3 validate-env-consistency.py

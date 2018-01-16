# In The Name Of God
# ========================================
# [] File Name : setup.py
#
# [] Creation Date : 16-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
# Always prefer setuptools over distutils
from setuptools import setup
# To use a consistent encoding
from os import path
# parse requirements.txt
from pip.req import parse_requirements

here = path.abspath(path.dirname(__file__))

# parse_requirements() returns generator of pip.req.InstallRequirement objects
install_reqs = parse_requirements(path.join(here, 'requirements.txt'),
                                  session='hack')

# reqs is a list of requirement
# e.g. ['django==1.5.1', 'mezzanine==1.4.6']
reqs = [str(ir.req) for ir in install_reqs]


setup(
        name='runtime.py',

        # Versions should comply with PEP440.
        # For a discussion on single-sourcing
        # the version across setup.py and the project code, see
        # https://packaging.python.org/en/latest/single_source_version.html
        version='0.2.0',


        # Author details
        author='Parham Alvani',
        author_email='parham.alvani@gmail.com',

        py_modules=['codec', 'main'],

        # List run-time dependencies here.  These will be installed by pip when
        # your project is installed.
        # For an analysis of "install_requires" vs pip's
        # requirements files see:
        # https://packaging.python.org/en/latest/requirements.html
        install_requires=reqs,

        # To provide executable scripts, use entry points in preference to the
        # "scripts" keyword. Entry points provide cross-platform
        # support and allow
        # pip to create the appropriate form of executable
        # for the target platform.
        entry_points={
            'console_scripts': [
                'runtime.py=main:main',
            ],
        },

)

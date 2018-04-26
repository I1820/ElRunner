# In The Name Of God
# ========================================
# [] File Name : setup.py
#
# [] Creation Date : 16-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
'''
setup module for runtime.py
'''

# Always prefer setuptools over distutils
from setuptools import setup


setup(
    name='runtime.py',

    # Versions should comply with PEP440.
    # For a discussion on single-sourcing
    # the version across setup.py and the project code, see
    # https://packaging.python.org/en/latest/single_source_version.html
    version='0.3.2',


    # Author details
    author='Parham Alvani',
    author_email='parham.alvani@gmail.com',

    package_dir={'': 'src'},
    packages=['codec', 'scenario'],
    py_modules=['main'],

    # List run-time dependencies here.  These will be installed by pip when
    # your project is installed.
    # For an analysis of "install_requires" vs pip's
    # requirements files see:
    # https://packaging.python.org/en/latest/requirements.html
    install_requires=[
        # Main
        "wheel",
        # runtime.py
        "click",
        # User Packages
        "cbor",
        # Analyze
        "pylint",
        # Scenario
        "aiohttp",
        "pymongo",
        "geopy",
        "redis",
        "kavenegar",
    ],

    setup_requires=[
        'pytest-runner',
    ],

    tests_require=[
        'json-rpc',
        'werkzeug',
        'pytest-asyncio',
    ],

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

from analyze.pylint_runner import run_pylint, Status

status = run_pylint()
assert status in [Status.no_error, Status.convention, Status.warning, Status.error, Status.fatal]

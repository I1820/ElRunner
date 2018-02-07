from analyze.pylint_runner import run_pylint, Status


def test_analyze():
    status = run_pylint('./scenarios')
    if status not in [Status.no_error, Status.error, Status.fatal]:
        raise AssertionError()

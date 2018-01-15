import codecs
import json
from enum import Enum, unique

import os
from pylint import epylint as lint

config_file_path = "'" + os.path.dirname(os.path.realpath(__file__)) + os.sep + 'pylintrc' + "'"
default_log_file_path = 'log.txt'
message_type_set = {'convention', 'warning', 'error', 'fatal'}


class OutputFormat(Enum):
    json = 'json'
    parseable = 'parseable'


@unique
class Status(Enum):
    no_error = 0
    convention = 1
    warning = 2
    error = 3
    fatal = 4


def pylint_check(path, output_format, ignore=None, msg_template=None):
    """
    Uses pylint function with passing options as a string
    pylint docs can be found here: `Pylint Docs <https://docs.pylint.org/en/1.6.0/run.html>`_
    :param path: code files path
    :param output_format: output messages format that should be OutputFormat class member
    :param ignore: files or directory(<file[,file...]>) to ignore
    :param msg_template:
    :return:
    """
    pylint_options = '{} -j 10 --rcfile={} --output-format={}'.format(path, config_file_path, output_format.value)
    if ignore is not None:
        pylint_options += ' --ignore={}'.format(ignore)
    if msg_template is not None:
        pylint_options += ' --msg-template={}'.format(msg_template)
    (pylint_stdout, pylint_stderr) = lint.py_run(pylint_options, return_std=True)
    return pylint_stdout.getvalue(), pylint_stderr.getvalue()


def log(content, log_file_path):
    """
    Prints and logs given content
    :param content: content to print and log
    :param log_file_path: log file path to store content
    :return: nothing
    """
    print(content)
    log_file = codecs.open(log_file_path, "w+", "utf-8")
    log_file.write(str(content))
    log_file.close()


def run_pylint(path=None, ignore=None, msg_template=None, log_file_path=None):
    """
    Main function to analyze code which returns status code and logs messages in the given file
    :param path: code files path
    :param ignore: files or directory(<file[,file...]>) to ignore
    :param msg_template: output messages template
    :param log_file_path: log file path to store results
    :return: status code which is a member of Status class.
    """
    # set default parameters if not set
    if path is None:
        path = '../scenarios/'
    if msg_template is None:
        msg_template = "'{path}:{line}: [{msg_id}({symbol}), {obj}] {msg}'"
    if log_file_path is None:
        log_file_path = default_log_file_path

    standard_output, err_output = pylint_check(path=path, output_format=OutputFormat.parseable, ignore=ignore,
                                               msg_template=msg_template)
    # log outputs in file
    log(err_output + '\n' + standard_output, log_file_path)

    standard_output, err_output = pylint_check(path=path, output_format=OutputFormat.json, ignore=ignore,
                                               msg_template=msg_template)
    status = get_status(standard_output)
    return status


def get_status(messages_json):
    """
    Computes status code from messages json input.
    :param messages_json: A json array including messages
    :return: status code which is a member of Status class.
    From best to worst: Status.no_error, Status.convention, Status.warning, Status.error, Status.fatal
    """
    status = Status.no_error
    messages = json.loads(messages_json)
    for message in messages:
        message_type = message['type']
        # check if message type is worse than current status
        if message_type in message_type_set and Status[message_type].value > status.value:
            status = Status[message_type]
    return status

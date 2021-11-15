from exceptions import IncorrectCommand
import subprocess

class NMap:
  nmapCommand = 'nmap'

  def __init__(self, command):
    splitCommand = command.split()
    self.command = splitCommand

    if splitCommand[0] != self.nmapCommand:
      raise IncorrectCommand('Expected nmap, found {0}'.format(splitCommand[0]))
    
    for idx, argument in enumerate(splitCommand):
      if splitCommand[idx - 1] == '-oX':
        self.outputFile = argument
        self.outputFileType = 'xml'

  def executeCommand(self):
    nmapCommand = subprocess.run(self.command, stdout=subprocess.DEVNULL)
    print("The exit code is %d" % nmapCommand.returncode)


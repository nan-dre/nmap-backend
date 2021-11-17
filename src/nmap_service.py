from exceptions import IncorrectCommand, MissingBodyArgument
import subprocess
import xmltodict
import json
from typing import Dict
import string
import random

class NMap:
  nmapCommand = "nmap"

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
    subprocess.run(self.command, stdout=subprocess.DEVNULL)

  def getJSONResult(self):
    f = open(self.outputFile)
    xmlContent = f.read()
    f.close()
    return json.dumps(xmltodict.parse(xmlContent), indent=4, sort_keys=True)

class NMapBuilder:
  command = "nmap -A -Pn -oX {0} -p {1}-{2} {3}"

  def __init__(self, requestDictionary: Dict):
    if not requestDictionary.get("host"):
      raise MissingBodyArgument("Missing Argument: host")
    if not requestDictionary.get("start-port"):
      raise MissingBodyArgument("Missing Argument: start-port")
    if not requestDictionary.get("end-port"):
      raise MissingBodyArgument("Missing Argument: end-port")
    
    self.host = requestDictionary["host"][0]
    self.startPort = int(requestDictionary["start-port"][0])
    self.endPort = int(requestDictionary["end-port"][0])
    self.outputFile = "temp/{0}.xml".format(''.join(random.choices(string.ascii_letters + string.digits, k=10)))

  def getCommand(self) -> str:
    self.command = self.command.format(self.outputFile, self.startPort, self.endPort, self.host)
    return self.command

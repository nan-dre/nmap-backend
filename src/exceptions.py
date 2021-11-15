class IncorrectCommand(Exception):
  def __init__(self, *args: object):
      if args:
        self.message = args[0]
      else:
        self.message = "Incorrect Command Provided!"
  
  def __str__(self):
    if self.message:
      return "IncorrectCommand {0}".format(self.message)
    else:
      return "IncorrectCommand has been raised"
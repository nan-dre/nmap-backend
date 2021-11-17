from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib import parse
from exceptions import MissingBodyArgument
from nmap_service import NMap, NMapBuilder

hostName = 'localhost'
serverPort = 8000

class MyServer(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-Type", "text/html")
        self.end_headers()
        mainViewFile = open("views/main.html")
        mainViewTemplate = mainViewFile.read()
        mainViewFile.close()
        self.wfile.write(bytes(mainViewTemplate, "utf-8"))

    def do_POST(self):
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        contentLength = int(self.headers['Content-Length'])
        requestBody = self.rfile.read(contentLength).decode("utf-8")
        try:
            nmapCommand = NMapBuilder(parse.parse_qs(requestBody)).getCommand()
        except MissingBodyArgument as e:
            self.wfile.write(bytes("{\"error\": {0} }".format(e.message), "utf-8"))
        # nmap = NMap("nmap -A -Pn -oX temp/test.xml -p 20000-20005 89.136.180.122")
        # nmap.executeCommand()
        # self.wfile.write(bytes(nmap.getJSONResult(), "utf-8"))
        print(nmapCommand)
        nmap = NMap(nmapCommand)
        nmap.executeCommand()
        self.wfile.write(bytes(nmap.getJSONResult(), "utf-8"))

if __name__ == '__main__':
    webServer = HTTPServer((hostName, serverPort), MyServer)
    print("Server started http://%s:%s" % (hostName, serverPort))

    try:
        webServer.serve_forever()
    except KeyboardInterrupt:
        pass

    webServer.server_close()
    print("Server stopped.")
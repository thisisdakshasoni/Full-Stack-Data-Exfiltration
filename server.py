from flask import Flask, request

app = Flask(__name__)

@app.route('/', methods=['POST'])
def receive():
    data = request.data.decode()
    print(f"[+] Received data:\n{data[:200]}...")
    return "OK", 200

if __name__ == "__main__":
    context = ('/etc/ssl/exfiltration/server.crt', '/etc/ssl/exfiltration/server.key')
    app.run(host='0.0.0.0', port=8443, ssl_context=context)


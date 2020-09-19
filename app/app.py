import logging
from flask import Flask, request, jsonify

app = Flask(__name__)


@app.route('/validate', methods=['POST'])


def deployment_webhook():

    app.logger.info('testing info log')

    request_info = request.get_json()
    uid = request_info["request"].get("uid")

    app.logger.info(uid)

    try:
        if request_info["request"]["object"]["metadata"]["labels"].get("billing"):
            return response(True, uid, "Billing label exists")
    except:
        return response(False, uid, "No labels exist. A Billing label is required")
    
    return response(False, uid, "Not allowed without a billing label")


def response(allowed, uid, message):
     return jsonify({"apiVersion": "admission.k8s.io/v1", "kind": "AdmissionReview", "response": {"allowed": allowed, "uid": uid, "status": {"message": message}}})


if __name__ == '__main__':
    app.run(ssl_context=('certs/mifommcrt.pem', 'certs/mifommkey.pem'),debug=True, host='0.0.0.0')    
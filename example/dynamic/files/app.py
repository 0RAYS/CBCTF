import os
from base64 import b64decode
from flask import Flask, request, jsonify, Response
from generator import generate_attachment

app = Flask(__name__)


@app.route('/gen', methods=['GET'])
def generate():
    data = request.args.to_dict()

    # 身份验证，限制访问
    pwd = data.get('pwd')
    if pwd != os.getenv("generator_pwd"):
        return jsonify({'error': 'Invalid password'}), 403

    id_ = data.get('id')
    # str: base64(base64(flag1),base64(flag2),...)
    flags = data.get('flag')

    if not id_ or not flags:
        return jsonify({'error': 'Invalid parameters'}), 400
    try:
        # bytes: base64(flag1),base64(flag2),...
        flags = b64decode(flags)
    except Exception as e:
        return jsonify({'error': f'Invalid flag encoding: {e}'}), 400

    # list: [flag1, flag2, ...]
    # flags = [b64decode(i) for i in flags.split(b",")]

    file = generate_attachment(flags)
    return Response(file, mimetype='application/zip', headers={
        'Content-Disposition': f'attachment; filename={id_}.zip'
    })


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000, debug=False)

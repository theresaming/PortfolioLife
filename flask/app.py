from flask import Flask, flash, render_template, request, session, json, redirect, make_response
import requests
import os
from werkzeug import secure_filename

from config import MYSQL_URL, OAUTH_CREDENTIALS, allowed_file
from registration import RegistrationForm
from api_strings import *

app = Flask(__name__)
app.secret_key = os.urandom(12)

app.config['SQLALCHEMY_DATABASE_URI'] = MYSQL_URL
app.config['OAUTH_CREDENTIALS'] = OAUTH_CREDENTIALS


@app.route('/')
def login():
    if not session.get('logged_in'):
        return redirect("/login", code=302)
    else:
        return load_home()


@app.route('/login', methods=['GET', 'POST'])
def do_admin_login():
    if request.method == 'POST':
        email = request.form['email']
        print email
        password = request.form['password']
        print password
        data = {
            "email": email,
            "password": password
        }
        jsonStr = json.dumps(data)
        r = requests.post(api_login, jsonStr)
        jsonDict = json.loads(r.text)
        if jsonDict['success']:
            session['logged_in'] = True
            resp = make_response(redirect('/'))
            resp.set_cookie('token', jsonDict['token'])
            return resp
        else:
            flash(jsonDict['message'])
    return render_template('login.html')


@app.route('/registration', methods=['GET', 'POST'])
def register():
    form = RegistrationForm(request.form)
    if request.method == 'POST':
        name = form.username.data
        password = form.password.data
        email = form.email.data
        data = {
            'email': email,
            'name': name,
            'password': password,
            }
        jsonStr = json.dumps(data)
        r = requests.post(api_registration, jsonStr)
        # TODO API does not return error for mismatching passwords/invalid email/etc.
        # TODO Flask does not handle error, including a missing message in current implementation.
        jsonDict = json.loads(r.text)
        print "jsonDict['message']", jsonDict['message']
        if jsonDict['success']:
            flash('Thanks for registering')
            return render_template('login.html')
        else:
            flash(jsonDict['message'])
    return render_template('registration.html', form=form)


@app.route("/logout")
def logout():
    session['logged_in'] = False
    requests.post(api_logout, headers={'token': request.cookies.get('token')})
    resp = make_response(redirect('/'))
    resp.set_cookie('token', '', expires=0)
    return resp


@app.route("/home")
def load_home():
    if session.get('logged_in'):
        token = request.cookies.get('token')

        # Get photos from API
        getPhotos = requests.get(api_get_photos, headers={'token': token})
        jsonDict = json.loads(getPhotos.text)

        # Add photos to array
        if jsonDict['success']:
            pictureUrlArr = [picture['url'] for picture in jsonDict['pictures']]
        else:
            flash(jsonDict['message'])
            pictureUrlArr = []
        return render_template('home.html', imageArr=pictureUrlArr)
    else:
        return login()


@app.route("/upload")
def load_upload():
    return render_template('uploadPhotos.html')


@app.route('/uploader', methods=['GET', 'POST'])
def upload_file():
    if request.method == 'POST':
        f = request.files['file']
        # saves files locally
        # f.save(secure_filename(f.filename))
        if f.filename == '':
            flash('No selected file')
            return redirect(request.url)
        if f and allowed_file(f.filename):
            filename = secure_filename(f.filename)
            print request.cookies.get('token')
            req = requests.post(api_upload,
                                headers={'token': request.cookies.get('token')},
                                files={'files': (filename, f, None, None)})
            jsonDict = json.loads(req.text)
            if jsonDict['success']: # if upload successful
                return load_home()
            else:
                return str(req.status_code) + ': ' + jsonDict['message']


@app.route('/deleter', methods=['POST'])
def delete_file():
    if request.method == 'POST':
        pic_id = ""
        # TODO fix this super-hack to get the id off the picture?
        for x in request.form:
            pic_id = x[:-2]

        token = request.cookies.get('token')
        print token
        api_delete_pic = api_delete + pic_id
        req = requests.delete(api_delete_pic,
                              headers={'token': token})
        jsonDict = json.loads(req.text)
        if jsonDict['success']: # if upload successful
            return load_home()
        else:
            return str(req.status_code) + ': ' + jsonDict['message']


@app.route("/delete", methods=['GET'])
def load_delete():
    req = requests.get(api_get_photos, headers={'token': request.cookies.get('token')})
    jsonDict = json.loads(req.text)
    imageArr = [{'url': picture['url'], 'id': picture['pictureID']} for picture in jsonDict['pictures']]
    return render_template('deletePhotos.html', imageArr=imageArr)


if __name__ == "__main__":
    app.run(debug=True,host='0.0.0.0', port=4000)
    # max upload thingy
    app.config['MAX_CONTENT_LENGTH'] = 16 * 1024 * 1024 # 16 MB
    # TODO: throw a 404 page if filesize is too large

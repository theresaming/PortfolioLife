from flask import Flask, flash, render_template, request, session, json, redirect, make_response, url_for
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
        # print email
        password = request.form['password']
        # print password
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
            pictureArr = [(picture['url'], picture['pictureID']) for picture in jsonDict['pictures']]
            pictureUrlArr = [picture['url'] for picture in jsonDict['pictures']]
            pictureIDArr = [picture['pictureID'] for picture in jsonDict['pictures']]
        else:
            flash(jsonDict['message'])
            pictureUrlArr = []

        return render_template('home.html', pictureArr=pictureArr, pictureUrlArr=pictureUrlArr, pictureIDArr=pictureIDArr)
    else:
        return login()


@app.route("/process-audio", methods=['POST'])
def process_audio():
    if request.method == 'POST':
        transcript = request.form['transcript'];
        if "upload" in transcript:
            return redirect(url_for('load_upload'))
        elif "album" in transcript:
            return redirect(url_for('load_home_albums'))
        elif "home" in transcript or "main" in transcript:
            return redirect(url_for('load_home'))
        else:
            return render_template('process-audio.html', transcript=transcript)
    return render_template(not_found_error(404))


@app.route("/albums")
def load_home_albums():
    if session.get('logged_in'):
        token = request.cookies.get('token')

        # Get photos from API
        getAlbums = requests.get(api_album, headers={'token': token})
        jsonDict = json.loads(getAlbums.text)

        if jsonDict['success']:
            albumIdTitle = [(a['title'], a['albumID']) for a in jsonDict['albums']]
        else:
            flash(jsonDict['message'])
            albumIdTitle = []
        print albumIdTitle
        return render_template('home-albums.html', albumData=albumIdTitle)
    else:
        return login()


@app.route("/upload")
def load_upload():
    return render_template('uploadPhotos.html')


@app.route('/uploader', methods=['GET', 'POST'])
def upload_file():
    if request.method == 'POST':
        f = request.files['file']
        # print f
        allFiles = request.files.getlist('file')
        imgList = []
        for f in allFiles:
            if f.filename == '':
                flash('No selected file')
                return redirect(request.url)
            if f and allowed_file(f.filename):
                #TODO: Change to mass upload
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


@app.route("/delete", methods = ['GET'])
def load_delete():
    req = requests.get(api_get_photos, headers={'token': request.cookies.get('token')})
    jsonDict = json.loads(req.text)
    imageUrlArr = [picture['url'] for picture in jsonDict['pictures']]
    return render_template('deletePhotos.html', imageArr=imageUrlArr)

@app.route("/image")
@app.route("/image/<image_id>", methods=['GET'])
def view_image(image_id):
    # print image_id
    req = requests.get(api_photo_view + image_id, headers={'token': request.cookies.get('token')})
    jsonDict = json.loads(req.text)
    if jsonDict['success']:
        image_url = jsonDict['url']

        try:
            return render_template("viewImage.html", imageID=image_id, imageURL = image_url)
        except Exception, e:
        	return(str(e))

@app.route("/delete/<image_id>")
def delete_image(image_id):
    req = requests.delete(api_photo_view + image_id, headers={'token': request.cookies.get('token')})
    jsonDict = json.loads(req.text)
    if jsonDict['success']:
        return load_home()

# this adds the tags
@app.route("/home")
def like_image(image_id):
    return load_home()

@app.route('/image/<image_id>/tagged', methods = ['POST', 'GET'])
def get_post_javascript_data(image_id):
    if request.method == "POST":
        jsdata = request.data
        jsdata = jsdata[1:len(jsdata) - 1]
        jsdata = jsdata.split(",")
        imageurl = jsdata[0]
        imageid = jsdata[1]
        print jsdata
        data = {
            "tags": jsdata[2:]
        }

        print "jsdata[2:]", jsdata[2:]
        print "jsdata[2:][0]", jsdata[2:][0]
        req = requests.post(api_photo_view + "/" + imageid + "/tags", headers={'token': request.cookies.get('token')},
            data={'tags': [jsdata[2:]]})
        print request.cookies.get('token')
        jsonDict = json.loads(req.text)
        print jsonDict
        return render_template("viewImage.html", imageID=jsdata[1], imageURL=jsdata[0], tags=jsdata[1:])

        # jsonDict = json.loads(req.text)
        # if jsonDict['success']:
        #     imageURL = jsonDict['url']
        #     try:
        #         return render_template("viewImage.html", imageID=image_id, imageURL = image_url)
        #     except Exception, e:
        #     	return(str(e))


@app.route("/create-album")
def load_create_album():
    if session.get('logged_in'):
        token = request.cookies.get('token')

        # Get photos from API
        getPhotos = requests.get(api_get_photos, headers={'token': token})
        jsonDict = json.loads(getPhotos.text)

        # Add photos to array
        if jsonDict['success']:
            pictureArr = [(picture['url'], picture['pictureID']) for picture in jsonDict['pictures']]
            pictureUrlArr = [picture['url'] for picture in jsonDict['pictures']]
            pictureIDArr = [picture['pictureID'] for picture in jsonDict['pictures']]
        else:
            flash(jsonDict['message'])
            pictureUrlArr = []
        return render_template('addAlbum.html', pictureArr=pictureArr, pictureUrlArr=pictureUrlArr, pictureIDArr=pictureIDArr)
    else:
        return login()


@app.route("/submit-album", methods=["POST"])
def submit_album():
    if request.method == "POST":
        data = request.form.to_dict()
        pictureIDs = []
        title = "My New Album"
        for key, value in data.items():
            if key == 'title':
                title = value
            else:
                pictureIDs.append(key)

        asDict = {'title': title,
                  'pictureIDs': pictureIDs}
        asJSON = json.dumps(asDict)
        req = requests.post(api_album_create,
                            headers={'token': request.cookies.get('token')},
                            data=asJSON)
        jsonDict = json.loads(req.text)
        if jsonDict['success']:
            return render_template('albumCreated.html', title=jsonDict['title'], albumId=jsonDict['albumID'])
        else:
            return str(req.status_code) + ': ' + jsonDict['message']
    return "Something went wrong!"


@app.route("/album/<title>/<albumId>")
def album_view(title, albumId):
    reqStr = api_album + "/" + albumId
    getAlbum = requests.get(reqStr, headers={'token': request.cookies.get('token')})
    jsonDict = json.loads(getAlbum.text)

    # Add photos to array
    if jsonDict['success']:
        pictureArr = [(picture['url'], picture['pictureID']) for picture in jsonDict['pictures']]
        pictureUrlArr = [picture['url'] for picture in jsonDict['pictures']]
    else:
        flash(jsonDict['message'])
        pictureArr = []
        pictureUrlArr = []
    return render_template('albumView.html', pictureArr=pictureArr, pictureUrlArr=pictureUrlArr, title=title, albumId=albumId)


@app.route("/stretch/albumEdit")
def edit_album():
    return render_template('editAlbum.html')


@app.route("/album/delete/<albumId>")
def delete_album(albumId):
    reqStr = api_album + "/" + albumId
    delAlbum = requests.delete(reqStr, headers={'token': request.cookies.get('token')})
    jsonDict = json.loads(delAlbum.text)

    if jsonDict['success']:
        flash("Album deleted!")
    else:
        flash(jsonDict['message'])
    return load_home_albums()


@app.errorhandler(404)
def not_found_error(error):
    return render_template('404.html'), 404


if __name__ == "__main__":
    app.run(debug=False,host='0.0.0.0', port=5000)
    # max upload thingy
    app.config['MAX_CONTENT_LENGTH'] = 16 * 1024 * 1024 # 16 MB
    # TODO: throw a 404 page if filesize is too large

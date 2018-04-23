# PortfolioLife

The Curators project, Fall17/Spring18

Designed by:

- [Melissa Essue](https://github.com/messue)
- [Paul Hamill](https://github.com/paul-io)
- [Joshua Koh](https://github.com/JoshuaKoh)
- [Theresa Ming](https://github.com/theresaming)
- [David Nelson Taylor](https://github.com/LargeSlurpee)

-----

# Install Guide

Note that some steps in this install guide are dependent on others. For best results, perform the steps of this guide in order!

This guide assumes a basic competence with a command line tool like Terminal. For more information, [click here](https://www.davidbaumgold.com/tutorials/command-line/).

## Project Code

* Install [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git).
* Clone the repository to your machine. Use `git clone https://github.com/theresaming/PortfolioLife`.

## Back End

### Image Storage (Object Storage)

* Create an "S3 Compatible" object storage, such as Amazon S3, DigitalOcean Spaces, or Minio.
* Fill out the details for the service in `server/config.json`.
* Fill out your API key/secret in `server/run-server.sh`.

### MySQL Database

* Install [MySQL](https://www.mysql.com/downloads/).
* Create a MySQL user and database, and fill out the required details in `server/config.json`.

### REST API

* Install [Go](https://golang.org/dl/).
* Open a terminal in `server` and type in `go test -run TestMigration` to setup the database with the preconfigured tables.
* Run `server/run-server.sh` to start the API.

## Front End

### Python Libraries

* Install [pip](https://pypi.org/project/pip/). On OSX, use `sudo easy_install pip` in Terminal.
* Navigate to the flask directory in the project. On OSX, use `cd PortfolioLife/flask`.
* Enter `pip install -r requirements.txt` in your console to download all Python dependencies.

### Running PortfolioLife locally

* Navigate to the flask directory in the project. On OSX, use `cd PortfolioLife/flask`.
* Enter `FLASK_APP=app.py flask run` in your console to run the app. 
* Your console should output a local address (e.g. 123.0.0.1:8080) that you can paste as a URI in your web browser to view PortfolioLife.

-----

# Release Notes

## Version 1.0.0, April 23, 2018

### Release Features

- **Register** for an account and **login** to the service. Users who have not registered are unable to use the service.
- **View** photos and albums on an account by using the sidebar in the top left corner of the homepage.
- **Upload** photos to an account, by clicking on the "Upload Photos" button on the homepage.
- **Delete** photos from an account, by clicking on the trash can button over a photo on the homepage. Deleted photos cannot be recovered!
- **Create** albums on an account, by clicking on the "Create an Album" button on the homepage. Albums cannot be created without a title and must contain at least one picture.
- **Navigate** to the photos homepage, the albums homepage, or the upload photos page using **voice control**. To activate voice control, click the button in the bottom left or press the spacebar. Wait for your spoken text to appear on screen before pressing the button again to submit your request.

### Bug Fixes

- Clicking the Delete button when viewing an album no longer redirects to the homepage.
- Clicking the Edit button when viewing an album no longer redirects to the homepage.
- Attempting to submit a voice request by clicking on the voice button no longer causes the feature to become unresponsive. 

### Known Bugs and Defects

- When viewing pictures on the album creation screen, a thin square border appears inside each photo.
- Share-able links were not implemented.
- Tag searching was not implemented.
- Album editing was not implemented.
- Mobile phone support was not implemented.
- Voice control for photo deleting, photo tagging, search, album creation, and album deletion were not implemented.
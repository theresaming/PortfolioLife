from sqlalchemy import *
from sqlalchemy import create_engine, ForeignKey
from sqlalchemy import Column, Date, Integer, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import relationship, backref
from app import db
import hashlib, random, string

Base = declarative_base()

########################################################################
class User(Base):
    """"""
    __tablename__ = "users"

    id = Column(Integer, primary_key=True)
    name = Column(String(4096)) # This & email lengths are arbitrary max lengths. TEXT type has max length of 2^16 (65536)
    email = Column(String(4096))
    password = Column(String(64)) # sha256 = 256bits (64 bytes)
    salt = Column(String(8)) # we'll generate salts that are 8 bytes long
    oauth = Column(Integer)
    token = Column(String(32)) # also arbitrary but we should limit it later!

    #----------------------------------------------------------------------
    def __init__(self, name, email, password, salt):
        """"""
        self.username = name
        self.password = sha256(password + "+" + salt)
        self.email = email
        self.salt = salt

    def checkPassword(self, rawPassword, salt):
        return sha256(rawPassword + "+" + salt) == self.password

def sha256(raw):
    return hashlib.sha256(raw).hexdigest()

def generateSalt():
    return ''.join(random.SystemRandom().choice(string.ascii_uppercase + string.digits) for _ in range(8))



def registerUser(user):
    db.session.add(user)
    db.session.commit()


# create tables
# Base.metadata.create_all(engine)
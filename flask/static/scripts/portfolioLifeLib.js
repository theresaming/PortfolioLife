function checkImageSelect(clicked_id) {
    var imageId = "img" + clicked_id;
    var image = document.getElementById(imageId);
    if (document.getElementById(clicked_id).checked) {
        image.style.border = "3px solid #00FF00";
    } else {
        image.style.border = "3px solid #FFFFFF";
    }
}
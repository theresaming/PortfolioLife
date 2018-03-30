function checkImageSelect(clicked_id) {
    var imageId = "img" + clicked_id;
    var image = document.getElementById(imageId);
    if (document.getElementById(clicked_id).checked) {
        image.style.border = "3px solid #00FF00";
    } else {
        image.style.border = "3px solid #FFFFFF";
    }
}

function albumValidation() {
    var title = document.getElementById("title-text");
    if ("" == title.value) {
        title.style.border = "1px solid red";
        alert("Please enter an album title.");
        return false;
    }

    var checkBoxes = document.getElementsByClassName( 'imgCheckbox' );
    console.log(checkBoxes.length);
    for (var i = 0; i < checkBoxes.length; i++) {
        if ( checkBoxes[i].checked ) {
            return true;
        }
    }
    alert( 'Please select at least one picture!' );
    return false;

}

function loadImage (el) {
    var img = new Image(), src = el.getAttribute('data-src');
    img.onload = function() {
      console.log("Image finished Loading");
      if (!! el.parent)
        el.parent.replaceChild(img, el)
      else
        el.src = src;

      finishLoadingImage(el)
      //fn? fn() : null;
    }
    img.onerror = failedLoadingImage;
    img.src = src;
}

function loadImages() {
  console.log("Loading the next bunch")
  document.getElementById("loading-box").style.display="inline";
  var no_images_load = 20; // if this number is too high it can cause memory issues on ther server side

  var images = $( ".unloaded" );
  for (i = 0; i < no_images_load ; i++) {
    loading_images += 1;
    loadImage(images[i])
  }
}

function failedLoadingImage (el) {
  el.innerHTML="";
  loading_images -= 1;
  if (loading_images < 1) {
    finishLoadingImages();
  }
}

function finishLoadingImage (el) {
  //console.log(el);
  el.classList.remove("unloaded");
  loading_images -= 1;
  if (loading_images < 1) {
    finishLoadingImages();
  }
}

function finishLoadingImages () {
  document.getElementById("loading-box").style.display="none"; // removes loading placeholder
  console.log("Finished Loading Images");
  loading_images_started = false;
}

function checkScroll () {
	var body = document.body, html = document.documentElement;
	var documentHeight = Math.max( body.scrollHeight, body.offsetHeight,
                       html.clientHeight, html.scrollHeight, html.offsetHeight );

	console.log( "window.pageYOffset: " + window.pageYOffset + " documentHeight " + documentHeight + " window.innerHeight: " + window.innerHeight );
	if ( window.pageYOffset + window.innerHeight > (documentHeight - 200) && !loading_images_started ) {
    loading_images_started = true
    loadImages();
	}
}

var loading_images_started = false;
var loading_images = 0;

function init () {
   loadImages();
   window.onscroll = checkScroll; // only check after the page has loaded
}
window.onload = init;

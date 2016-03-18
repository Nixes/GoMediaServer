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
    img.src = src;
}

function loadImages() {
  console.log("Loading the next bunch")
  $( ".dimmer" ).addClass( "active" );

  var images = $( ".unloaded" );
  for (i = 0; i < 10 ; i++) {
    loading_images += 1;
    loadImage(images[i])
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
  $( ".dimmer" ).removeClass( "active" ); // removes loading placeholder
  console.log("Finished Loading Images");
}

function checkScroll () {
	var body = document.body, html = document.documentElement;
	var documentHeight = Math.max( body.scrollHeight, body.offsetHeight,
                       html.clientHeight, html.scrollHeight, html.offsetHeight );

	console.log( "window.pageYOffset: " + window.pageYOffset + " documentHeight " + documentHeight + " window.innerHeight: " + window.innerHeight );
	if ( window.pageYOffset + window.innerHeight > (documentHeight - 100) ) {
    loadImages();
	}
}

var loading_images = 0;

function init () {
   loadImages();
   window.onscroll = checkScroll; // only check after the page has loaded
}
window.onload = init;

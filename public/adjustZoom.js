function adjustZoom(flag) {
  let pageWidth = 1400;
  let screenWidth = window.innerWidth;

  console.log(pageWidth, screenWidth, window);
  if (!flag && screenWidth < pageWidth) {
    let scaleFactor = screenWidth / pageWidth;
    document.body.style.zoom = scaleFactor;
    console.log(scaleFactor);
    // document.body.style.transformOrigin = 'top left';
  } else {
    document.body.style.zoom = 1;
  }
}

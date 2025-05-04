function adjustZoom(flag) {
  let pageWidth = 1400;
  let screenWidth = window.innerWidth;

  console.log(pageWidth, screenWidth, window);
  const element = document.getElementById("resultSection");
  if (!element) {
    console.error("Element with ID 'resultSection' not found.");
    return;
  }
  if (!flag && screenWidth < pageWidth) {
    let scaleFactor = screenWidth / pageWidth;
    element.style.zoom = scaleFactor;
    console.log(scaleFactor);
    // document.body.style.transformOrigin = 'top left';
  } else {
    element.style.zoom = 1;
  }
}

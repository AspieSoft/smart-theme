onReady(async function(){
  function onInterval(){
    document.querySelectorAll('header .header-img, footer').forEach(function(elm){
      if(elm.clientHeight < 300){
        elm.classList.remove('background-size-w', 'background-size-h');
        return;
      }

      if(elm.classList.contains('background-size-w') || elm.classList.contains('background-size-h')){
        return;
      }

      const css = window.getComputedStyle(elm);
      if(typeof css.getPropertyValue('--img-size') !== 'string' || typeof css.getPropertyValue('--img') !== 'string'){
        return;
      }

      if(css.getPropertyValue('--img-size').includes('cover')){
        css.getPropertyValue('--img').replace(/url\s*\(\s*(["'`])(.*?)\1\s*\)/i, function(_, _, url){
          const img = new Image();
          img.src = url;
          img.onload = function(){
            if(img.width !== 0 && img.height !== 0){
              elm.setAttribute('img-width', img.width);
              elm.setAttribute('img-height', img.height);

              if(calculateImgSizeRatioDiff(elm, img.width, img.height)){
                elm.classList.add('background-size-w');
                elm.classList.remove('background-size-h');
              }else{
                elm.classList.add('background-size-h');
                elm.classList.remove('background-size-w');
              }
            }
            img.remove();
          }
        });
      }
    });

    document.querySelectorAll('header .header-top').forEach(function(elm) {
      const css = window.getComputedStyle(elm);
      const img = css.getPropertyValue('--img');
      if(img !== '' && img !== 'none' && img.match(/^\s*([\w_-]+-gradient|url)/)){
        elm.style['background-color'] = 'transparent';
        elm.style['animation-name'] = 'scroll-header-top-shadow';
      }else{
        elm.style['background-color'] = '';
        elm.style['animation-name'] = '';
      }
    });
  }
  onInterval();
  setInterval(onInterval, 1000);


  function onResize(){
    document.querySelectorAll('header .header-img, footer').forEach(function(elm){
      if(elm.clientHeight < 300){
        elm.classList.remove('background-size-w', 'background-size-h');
        return;
      }
      
      if(!elm.classList.contains('background-size-w') && !elm.classList.contains('background-size-h')){
        return;
      }
      
      const imgWidth = elm.getAttribute('img-width');
      const imgHeight = elm.getAttribute('img-height');
      if(!imgWidth || !imgHeight || imgWidth === 0 || imgHeight === 0){
        return;
      }

      if(calculateImgSizeRatioDiff(elm, imgWidth, imgHeight)){
        elm.classList.add('background-size-w');
        elm.classList.remove('background-size-h');
      }else{
        elm.classList.add('background-size-h');
        elm.classList.remove('background-size-w');
      }
    });
  }
  onResize();
  window.addEventListener('resize', onResize, {passive: true});


  // Im not sure what to name this funnction, but it does some complex math
  // to calculate whether a background image should use
  // `background-size: 100% auto`, or `background-size: auto 100%`
  // to simulate a `background-size: cover` while allowing an animation
  // to adjust the 100% to a 125%
  //
  // the changes to `background-size` are handled by a css class
  function calculateImgSizeRatioDiff(elm, imgWidth, imgHeight){
    if(elm.clientWidth / elm.clientHeight > imgWidth / imgHeight){
      return true;
    }
    return false;

    /* let w = elm.clientWidth - imgWidth;
    let h = elm.clientHeight - imgHeight;

    if(w < h){
      if(elm.clientWidth / imgWidth >= elm.clientHeight / imgHeight){
        return true;
      }
      return false;
    }

    if(elm.clientWidth / imgWidth >= elm.clientHeight / imgHeight && Math.sqrt((w * w) + (h * h)) >= Math.sqrt((imgWidth * imgWidth) + (imgHeight * imgHeight))){
      return true;
    }
    return false; */
  }

});
